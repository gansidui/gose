package main

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 配置文件
type Config struct {
	EntryPath         string
	UrlQueuePath      string
	FilterExtPath     string
	LimitedToThisPath string
	DownWebpagePath   string
	DownUrlDataPath   string
	MaxNum            int32
	IntervalTime      string
}

// 网页信息
type UrlInfo struct {
	md5     string // 修剪之后的url的md5值，如，md5(www.baidu.com)
	url     string // 完整url, http://www.baidu.com/
	path    string // 网页本地保存路径
	content string // 网页的数据
}

// 爬虫
type Spider struct {
	conf          *Config         // 配置信息
	filterExts    map[string]bool // 扩展名过滤
	limitedToThis []string        // 仅仅抓取包含了其中至少一个子串的url，若limitedToThis为空，则该规则无效
	doneUrls      map[string]bool // 已经爬下的url，用md5标记
	exceptionUrls map[string]bool // 异常url以及被过滤的url等，用md5标记
	chUrlsInfo    chan *UrlInfo   // 爬取到的所有url的Info
	chUrl         chan string     // 存储url,供多个gorountine去处理
	chHttp        chan bool       // 控制同时下载url的gorountine数量
	chStopIO      chan bool       // 主线程通知结束磁盘IO gorountine
	chExit        chan bool       // 磁盘IO结束后再通知主线程结束
	wg            sync.WaitGroup  // 等待所有gorountine结束
	pageNum       int32           // 当前爬取的网页数量
	intervalTime  time.Duration   // 间歇时间,如 "5ms"
}

// 返回一个初始化了的Spider实例
func NewSpider() *Spider {
	// 磁盘处理只需一个goroutine，网络可以适当多几个goroutine，在未达到到带宽限制的情况下有利于抢占网络资源，
	// 磁盘IO阻塞后也可以保证有goroutine正在下载资源保存到内存中，另外正则匹配需要开多个goroutine处理
	runtime.GOMAXPROCS(runtime.NumCPU())

	var sp Spider
	sp.conf = NewConfig()
	sp.filterExts = make(map[string]bool)
	sp.limitedToThis = make([]string, 0)
	sp.doneUrls = make(map[string]bool)
	sp.exceptionUrls = make(map[string]bool)
	sp.chUrlsInfo = make(chan *UrlInfo, 100)
	sp.chUrl = make(chan string, 1000000)
	sp.chHttp = make(chan bool, 5)
	sp.chStopIO = make(chan bool)
	sp.chExit = make(chan bool)
	sp.pageNum = 0
	intervalTime, err := time.ParseDuration(sp.conf.IntervalTime)
	if err != nil {
		sp.intervalTime = 500 * time.Millisecond
	} else {
		sp.intervalTime = intervalTime
	}

	sp.readDownUrlData()
	sp.readEntry()
	sp.readUrlQueue()
	sp.readFilterExt()
	sp.readLimitedToThis()

	return &sp
}

// 读取以前爬过的网页数据
func (this *Spider) readDownUrlData() {
	// 创建父目录
	err := os.MkdirAll(filepath.Dir(this.conf.DownUrlDataPath), os.ModePerm)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", this.conf.DownUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 先建表，若表已经存在则建表失败
	db.Exec("create table data(md5 varchar(32), url varchar(256), path varchar(256))")

	// 读取数据
	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var md5, url, path string
	for rows.Next() {
		rows.Scan(&md5, &url, &path)
		this.doneUrls[md5] = true
	}
}

// 读取入口地址
func (this *Spider) readEntry() {
	file, err := os.Open(this.conf.EntryPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer file.Close()

	re := bufio.NewReader(file)
	for {
		urlbyte, _, err := re.ReadLine()
		if err != nil {
			break
		}
		if string(urlbyte) != "" && !this.doneUrls[getMd5FromUrl(string(urlbyte))] {
			this.chUrl <- string(urlbyte)
		}
	}
}

// 读取爬虫上次爬行过程中Url队列
func (this *Spider) readUrlQueue() {
	// 打开数据库
	db, err := sql.Open("sqlite3", this.conf.UrlQueuePath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 先建表，若表已经存在则建表失败
	db.Exec("create table data(md5 varchar(32), url varchar(256))")

	// 读取数据
	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var md5, url string
	for rows.Next() {
		rows.Scan(&md5, &url)
		if !this.doneUrls[md5] {
			this.chUrl <- url
			this.doneUrls[md5] = true
		}
	}
}

// 读取过滤的扩展名
func (this *Spider) readFilterExt() {
	file, err := os.OpenFile(this.conf.FilterExtPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Printf("%v\r\n", err)
		return
	}
	defer file.Close()

	re := bufio.NewReader(file)
	for {
		extbyte, _, err := re.ReadLine()
		if err != nil {
			break
		}
		if string(extbyte) != "" {
			this.filterExts[string(extbyte)] = true
		}
	}
}

// 读取指定抓取的url信息
func (this *Spider) readLimitedToThis() {
	file, err := os.OpenFile(this.conf.LimitedToThisPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Printf("%v\r\n", err)
		return
	}
	defer file.Close()

	re := bufio.NewReader(file)
	for {
		limbyte, _, err := re.ReadLine()
		if err != nil {
			break
		}
		if string(limbyte) != "" {
			this.limitedToThis = append(this.limitedToThis, string(limbyte))
		}
	}
}

// 将这次爬下的UrlInfo写入到数据库中
func (this *Spider) writeUrlInfo() {
	// IO结束后运行主线程退出
	defer func() {
		this.chExit <- true
	}()

	// md5, url, path 写入 this.conf.DownUrlDataPath数据库中
	// 打开数据库
	db, err := sql.Open("sqlite3", this.conf.DownUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()
	// 先建表，若表已经存在则建表失败
	db.Exec("create table data(md5 varchar(32), url varchar(256), path varchar(256))")

	// content 写入 this.conf.DownWebpagePath/xxx.html
	// 创建父目录
	err = os.MkdirAll(this.conf.DownWebpagePath, os.ModePerm)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	// 收到 停止IO 命令后不能马上退出，因为还需要等待chUrlsInfo中数据处理完才可
	var urlInfo *UrlInfo
	var canStop bool = false
	for {
		select {
		case <-this.chStopIO: //收到结束IO通知
			canStop = true
			if len(this.chUrlsInfo) == 0 {
				return
			}
		case urlInfo = <-this.chUrlsInfo:
			fmt.Printf("[%s]---正在写入文件...\n", urlInfo.url)
			// 保存网页
			ioutil.WriteFile(urlInfo.path, []byte(urlInfo.content), os.ModePerm)
			// 将网页信息插入到数据库中，忽略错误
			db.Exec("insert into data(md5, url, path) values(?, ?, ?)", urlInfo.md5, urlInfo.url, urlInfo.path)
			this.pageNum = atomic.AddInt32(&this.pageNum, 1)
			fmt.Printf("[%s]---写入完成.\n", urlInfo.url)
			if canStop && len(this.chUrlsInfo) == 0 {
				return
			}
		}
	}
}

// 将url队列写入到数据库中用于下次从这里开始爬行
func (this *Spider) writeUrlQueue(urls []string) {
	// 打开数据库
	db, err := sql.Open("sqlite3", this.conf.UrlQueuePath)
	if err != nil {
		log.Printf("%v\r\n", err)
		return
	}
	defer db.Close()

	tx, err := db.Begin() // 启动一个事务
	if err != nil {
		log.Printf("%v\r\n", err)
	} else {
		for _, vv := range urls {
			vv = trimUrl(vv)
			md5 := getMd5FromUrl(vv)
			if !this.doneUrls[md5] && !this.exceptionUrls[md5] && !this.beFiltered(vv) {
				_, err = tx.Exec("insert into data(md5, url) values(?, ?)", md5, vv)
				if err != nil {
					log.Printf("%v\r\n", err)
					break
				}
			}
			if len(urls) == 0 {
				break
			}
		}
	}
	tx.Commit() // 提交事务
}

// 判断爬取的网页数已经达到预期
func (this *Spider) isFinished() bool {
	// 这里的atomic并不能保证该函数同时只能被一个goroutine调用
	if atomic.LoadInt32(&this.pageNum) >= this.conf.MaxNum {
		log.Printf("%v\r\n", "爬取的网页数已达预期！！！")
		return true
	}
	return false
}

// 主线程，开始爬取
func (this *Spider) Fetch() {
	if len(this.chUrl) == 0 {
		log.Fatal("entry url is empty.\r\n")
		return
	}

	go this.writeUrlInfo()
	this.work()

	this.chStopIO <- true //通知结束writeUrlInfo()
	<-this.chExit         // 等待writeUrlInfo结束
}

// 工作线程
func (this *Spider) work() {
	for url := range this.chUrl {
		this.chHttp <- true //控制下载网页的goroutine数量

		go func(url string) {
			this.wg.Add(1)
			log.Printf("%v\r\n", "线程下载开始")
			defer func() {
				<-this.chHttp
				this.wg.Done()
				log.Printf("%v\r\n", "线程下载完成")
			}()
			this.do(url)
		}(url)

		log.Printf("len(chUrlsInfo)==%d --- len(chUrl)==%d --- len(chHttp)==%d\r\n", len(this.chUrlsInfo), len(this.chUrl), len(this.chHttp))
		time.Sleep(this.intervalTime) // 慢点爬，怕被网站封IP

		if this.isFinished() {
			log.Printf("%v\r\n", "正在等待各线程结束......")
			this.wg.Wait()
			log.Printf("%v\r\n", "各线程已经结束！！！")
			if len(this.chUrl) == 0 {
				break
			}

			// 保存this.chUrl中剩余的urls
			urls := make([]string, 0)
			for v := range this.chUrl {
				urls = append(urls, v)
				if len(this.chUrl) == 0 {
					break
				}
			}
			this.writeUrlQueue(urls)
			break
		}
	}
}

// 处理url
func (this *Spider) do(url string) {
	client := &http.Client{
		CheckRedirect: doRedirect,
	}
	// 若url重定向，则client.Get(url)里面调用自定义的doRedirect函数处理，
	// 然后将doRedirect的error返回给这里
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("%s\r\n", err)
		this.exceptionUrls[getMd5FromUrl(url)] = true
		return
	}
	defer resp.Body.Close()

	// 不为OK就返回，因为有些可能是500等错误
	if resp.StatusCode != http.StatusOK {
		log.Printf("[%s] resp.StatusCode == [%d]\r\n", url, resp.StatusCode)
		this.exceptionUrls[getMd5FromUrl(url)] = true
		return
	}

	fmt.Printf("[%s]---正在下载...\n", url)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s\r\n", err)
		this.exceptionUrls[getMd5FromUrl(url)] = true
		fmt.Printf("[%s]---异常，结束下载.\n", url)
		return
	}
	fmt.Printf("[%s]---下载完成.\n", url)

	// 保存UrlInfo
	md5 := getMd5FromUrl(url)
	path := this.conf.DownWebpagePath + md5 + ".html"
	this.chUrlsInfo <- &UrlInfo{md5: md5, url: url, path: path, content: string(content)}

	fmt.Printf("[%s]---正在分析...\n", url)
	// 得到新的url
	urls := getURLs(content)
	for i, v := range urls {
		// 已经完成任务，将url队列写入到数据库中用于下次从这里开始爬行
		if this.isFinished() {
			this.writeUrlQueue(urls[i:])
			break
		}

		// 还未到达预期，继续爬取
		v = trimUrl(v)
		md5 := getMd5FromUrl(v)
		if !this.doneUrls[md5] && !this.exceptionUrls[md5] {
			if this.beFiltered(v) {
				this.exceptionUrls[md5] = true
			} else {
				this.chUrl <- v
				this.doneUrls[md5] = true
			}
		}
	}
	fmt.Printf("[%s]---分析完成.\n", url)
}

// 过滤
func (this *Spider) beFiltered(url string) bool {
	b1 := this.filterExts[filepath.Ext(url)] // 后缀过滤
	b2 := true
	// 若limitedToThis不为空，且limitedToThis中的一个串是url的一个子串，就不过滤，否则过滤掉url
	if len(this.limitedToThis) > 0 {
		for _, v := range this.limitedToThis {
			if strings.Contains(url, v) {
				b2 = false
				break
			}
		}
	} else {
		b2 = false
	}

	return b1 || b2
}

// 重定向处理, StatusCode == 302
func doRedirect(req *http.Request, via []*http.Request) error {
	return errors.New(req.URL.String() + " was as an exception url to do.")
}

// 初始化
func init() {
	setLogOutput()
}

// 读取配置文件
func NewConfig() *Config {
	file, err := ioutil.ReadFile("./spider.conf")
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	var conf Config
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	return &conf
}

// 设置log输出
func setLogOutput() {
	// 为log添加短文件名，方便查看行数
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	logfile, err := os.OpenFile("./spider.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	// 注意这里不能关闭logfile
	if err != nil {
		log.Printf("%v\r\n", err)
	}
	log.SetOutput(logfile)
}

// 修剪url, 把 # 后面的字符去掉
func trimUrl(url string) string {
	p := strings.Index(url, "#")
	if p != -1 {
		url = url[:p]
	}
	return url
}

// 修剪之后的url，另外去掉最后的斜杠和scheme 如 md5(www.baidu.com)
func getMd5FromUrl(url string) string {
	url = strings.TrimRight(url, "/")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	m := md5.New()
	io.WriteString(m, url)
	str := fmt.Sprintf("%x", m.Sum(nil)) // 将md5值格式化成字符串
	return str
}

// 从html页面中提取所有的url
func getURLs(content []byte) (urls []string) {
	re := regexp.MustCompile("href\\s*=\\s*['\"]?\\s*(https?://[^'\"\\s]+)\\s*['\"]?")
	allsubmatch := re.FindAllSubmatch([]byte(content), -1)
	for _, v2 := range allsubmatch {
		for k, v := range v2 {
			// k == 0 是表示匹配的全部元素
			if k > 0 {
				urls = append(urls, string(v))
			}
		}
	}
	return urls
}

func main() {
	log.Printf("%v\r\n", "start......")
	start := time.Now()
	sp := NewSpider()
	sp.Fetch()
	log.Printf("used time: %v", time.Since(start))
}
