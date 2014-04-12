package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gansidui/gose/analyze/extractutil"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 配置文件
type Config struct {
	DownUrlDataPath    string
	ExtractWebpagePath string
	ExtractUrlDataPath string
}

// 分析
type Analyze struct {
	conf     *Config         // 配置信息
	doneUrls map[string]bool // 已经分析过的网页，用md5标记
}

// 返回一个初始化了的Analyze实例
func NewAnalyze() *Analyze {
	var an Analyze
	an.conf = NewConfig()
	an.doneUrls = make(map[string]bool)
	an.readExtractUrlData()
	return &an
}

// 读取已经分析过的网页数据
func (this *Analyze) readExtractUrlData() {
	// 创建父目录
	err := os.MkdirAll(filepath.Dir(this.conf.ExtractUrlDataPath), os.ModePerm)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", this.conf.ExtractUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 先建表，若表已经存在则建表失败
	db.Exec("create table data(md5 varchar(32))")

	// 读取数据
	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var md5 string
	for rows.Next() {
		rows.Scan(&md5)
		this.doneUrls[md5] = true
	}
}

// 开始网页分析
func (this *Analyze) Do() {
	// 创建父目录
	err := os.MkdirAll(filepath.Dir(this.conf.ExtractWebpagePath), os.ModePerm)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	// 打开保存下载网页信息数据的数据库
	db, err := sql.Open("sqlite3", this.conf.DownUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 打开保存分析记录的数据库
	exdb, err := sql.Open("sqlite3", this.conf.ExtractUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer exdb.Close()

	// 读取数据
	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	tx, err := exdb.Begin() // 启动事务
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	var md5, url, path string
	for rows.Next() {
		rows.Scan(&md5, &url, &path)
		if !this.doneUrls[md5] {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err, "\r\n")
			}

			// 抽取标题和正文，并写入文件
			title := extractutil.ExtractTitle(string(content))
			ioutil.WriteFile(this.conf.ExtractWebpagePath+md5+"_title.txt", []byte(title), os.ModePerm)

			body := extractutil.ExtractBody(string(content))
			ioutil.WriteFile(this.conf.ExtractWebpagePath+md5+"_body.txt", []byte(body), os.ModePerm)

			// 标记已经分析过了并写入数据库
			this.doneUrls[md5] = true
			_, err = tx.Exec("insert into data(md5) values(?)", md5)
			if err != nil {
				log.Fatal(err, "\r\n")
			}
		}
	}
	tx.Commit()
}

// 初始化
func init() {
	setLogOutput()
}

// 读取配置文件
func NewConfig() *Config {
	file, err := ioutil.ReadFile("./analyze.conf")
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
	logfile, err := os.OpenFile("./analyze.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	// 注意这里不能关闭logfile
	if err != nil {
		log.Printf("%v\r\n", err)
	}
	log.SetOutput(logfile)
}

func main() {
	log.Printf("%v\r\n", "start......")
	start := time.Now()
	an := NewAnalyze()
	an.Do()
	log.Printf("used time: %v", time.Since(start))
}
