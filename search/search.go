package search

import (
	"database/sql"
	"encoding/json"
	"github.com/gansidui/gose/indexing/participleutil"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// 配置文件
type Config struct {
	DownUrlDataPath    string
	ExtractWebpagePath string
	ExtractUrlDataPath string
	DicPath            string
	TfIdfPath          string
	DocMarkPath        string
}

// 关键词索引信息
type WordIndexInfo struct {
	id    int     // 文档编号
	tfidf float32 // TF-IDF值
}

// 网页信息
type UrlInfo struct {
	url  string // 原始url
	path string // 本地保存路径
}

// 搜索结果信息
type SearchResultInfo struct {
	id     int     // 文档编号
	tfidfs float32 // 该文档针对搜索串中的每个关键词的TF_IDF值之和
}

// 文章信息
type ArticleInfo struct {
	Title   template.HTML // 标题
	Summary template.HTML // 摘要
	Url     string        // 原始url
	Path    string        // 本地路径
}

// word --> []*WordIndexInfo 按 id 从小到大排序，便于二分查找
type ById []*WordIndexInfo

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].id < a[j].id }

// 按 tfidfs 从大到小排序
type ByTfIdfs []*SearchResultInfo

func (a ByTfIdfs) Len() int           { return len(a) }
func (a ByTfIdfs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTfIdfs) Less(i, j int) bool { return a[i].tfidfs > a[j].tfidfs }

// 按字符串长度从大到小排序
type ByLength []string

func (a ByLength) Len() int           { return len(a) }
func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLength) Less(i, j int) bool { return len(a[i]) > len(a[j]) }

var conf *Config                                 // 配置
var wordMapIndexInfo map[string][]*WordIndexInfo // word --> []*WordIndexInfo
var idMapMd5 map[int]string                      // id --> md5
var md5MapUrlInfo map[string]*UrlInfo            // md5 --> *UrlInfo

// 读取配置文件
func ReadConfig(confpath string) {
	file, err := ioutil.ReadFile(confpath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
}

// 初始化
func InitSearch() {
	participleutil.LoadDic(conf.DicPath)
	wordMapIndexInfo = make(map[string][]*WordIndexInfo)
	idMapMd5 = make(map[int]string)
	md5MapUrlInfo = make(map[string]*UrlInfo)
	readWordMapIndexInfo()
	readIdMapMd5()
	readMd5MapUrlInfo()
}

// 读取 word --> []*WordIndexInfo, 并排序
func readWordMapIndexInfo() {
	db, err := sql.Open("sqlite3", conf.TfIdfPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var (
		word  string
		id    int
		tfidf float32
	)

	for rows.Next() {
		rows.Scan(&word, &id, &tfidf)
		wordMapIndexInfo[word] = append(wordMapIndexInfo[word], &WordIndexInfo{id: id, tfidf: tfidf})
	}

	// 排序
	for _, v := range wordMapIndexInfo {
		sort.Sort(ById(v))
	}
}

// 读取 id --> md5
func readIdMapMd5() {
	db, err := sql.Open("sqlite3", conf.DocMarkPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var (
		md5 string
		id  int
	)

	for rows.Next() {
		rows.Scan(&md5, &id)
		idMapMd5[id] = md5
	}
}

// 读取 md5 --> *UrlInfo
func readMd5MapUrlInfo() {
	// 先读取已经分析过了的网页数据
	extracted := make(map[string]bool)
	exdb, err := sql.Open("sqlite3", conf.ExtractUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer exdb.Close()

	exrows, err := exdb.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer exrows.Close()

	var md5 string

	for exrows.Next() {
		exrows.Scan(&md5)
		extracted[md5] = true
	}

	// 读取 已经分析过了的网页数据 的 UrlInfo
	db, err := sql.Open("sqlite3", conf.DownUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var (
		url  string
		path string
	)

	for rows.Next() {
		rows.Scan(&md5, &url, &path)
		if extracted[md5] {
			md5MapUrlInfo[md5] = &UrlInfo{url: url, path: path}
		}
	}
}

// 根据id得到 *UrlInfo
func getUrlInfoById(id int) *UrlInfo {
	md5 := idMapMd5[id]
	return md5MapUrlInfo[md5]
}

// 根据word和id得到 tfidf, 二分查找
func getTfIdfByWordId(word string, id int) (result float32) {
	left, right := 0, len(wordMapIndexInfo[word])-1
	var mid, tmpid int

	for left <= right {
		mid = (left + right) >> 1
		tmpid = wordMapIndexInfo[word][mid].id
		if tmpid < id {
			left = mid + 1
		} else if tmpid == id {
			result = wordMapIndexInfo[word][mid].tfidf
			return result
		} else {
			right = mid - 1
		}
	}
	return -1
}

// 求并集
func union(sets [][]int) (result []int) {
	m := make(map[int]bool)
	for _, v := range sets {
		for _, vv := range v {
			m[vv] = true
		}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// 求交集
func inter(sets [][]int) (result []int) {
	m := make(map[int]int)
	for _, v := range sets {
		for _, vv := range v {
			m[vv]++
		}
	}
	setNum := len(sets)
	for k, v := range m {
		if v >= setNum {
			result = append(result, k)
		}
	}
	return result
}

// 对字符串数组去重
func clearRepeat(ss []string) (result []string) {
	m := make(map[string]bool)
	for _, v := range ss {
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}

// 搜索，返回 [start, end] 之间的结果(文档id集合)以及搜索到的文档总数
func search(searchString string, start, end int) (result []int, total int) {
	// 先按空格将搜索串分成多个句子，并过滤掉空句子
	var sentences []string
	texts := strings.Split(searchString, " ")
	for _, sen := range texts {
		if sen != "" {
			sentences = append(sentences, sen)
		}
	}

	var (
		tempResult       []int               // 临时结果
		words            []string            // 搜索串的关键词集合
		searchResultInfo []*SearchResultInfo // 用来根据tfidfs排序
	)

	// 对每个句子进行分词，句子内对每个词的id集合求并集，句子间对id集合求交集
	var outidsets [][]int
	for _, sen := range sentences {
		ws := participleutil.Participle(sen)
		var inidsets [][]int

		for _, w := range ws {
			var ids []int
			for _, v := range wordMapIndexInfo[w] {
				ids = append(ids, v.id)
			}
			inidsets = append(inidsets, ids)
			words = append(words, w)
		}

		outidsets = append(outidsets, union(inidsets)) // 对句内的集合求并集
	}
	tempResult = inter(outidsets) // 对句间的集合求交集

	// 对tempResult进行排序
	words = clearRepeat(words) // 去重
	for _, id := range tempResult {
		var tfidfs float32 = 0.0
		for _, w := range words {
			tfidfs += getTfIdfByWordId(w, id)
		}
		searchResultInfo = append(searchResultInfo, &SearchResultInfo{id: id, tfidfs: tfidfs})
	}
	sort.Sort(ByTfIdfs(searchResultInfo))

	// 选取 [start, end] 之间的id作为结果
	if start < 0 {
		start = 0
	}
	if end >= len(searchResultInfo) {
		end = len(searchResultInfo) - 1
	}
	for i := start; i <= end; i++ {
		result = append(result, searchResultInfo[i].id)
	}
	total = len(searchResultInfo)

	return result, total
}

// 将关键词标红
func markRedKeywords(content string, keywords []string) (result string) {
	patterns := []string{}
	sort.Sort(ByLength(keywords))
	for _, oldstr := range keywords {
		patterns = append(patterns, oldstr)
		newstr := "<font color=\"red\">" + oldstr + "</font>"
		patterns = append(patterns, newstr)
	}
	replacer := strings.NewReplacer(patterns...)
	result = replacer.Replace(content)
	return result
}

// 向调用者返回搜索结果
// searchString为搜索串，返回第[start, end]篇文档的信息，result保存文档信息, total为搜索到的文档总数
func GetSearchResult(searchString string, start, end int) (result []ArticleInfo, total int) {
	// 定义一个提取摘要的函数, 提取前100个rune
	getSummary := func(content string) string {
		num := 0
		var res string
		for _, v := range content {
			s := string(v)
			if s != " " && s != "\t" && s != "\r" && s != "\n" {
				num++
				res = res + s
				if num > 100 {
					break
				}
			}
		}
		return res
	}

	// 得到分词用来结果标红
	keywords := participleutil.Participle(searchString)

	res, tot := search(searchString, start, end)
	total = tot
	var articleInfo ArticleInfo
	for _, id := range res {
		md5 := idMapMd5[id]
		urlInfo := md5MapUrlInfo[md5]

		articleInfo.Url = urlInfo.url
		articleInfo.Path = urlInfo.path

		content, err := ioutil.ReadFile(conf.ExtractWebpagePath + md5 + "_body.txt")
		if err != nil {
			log.Printf("%v\r\n", err)
		} else {
			articleInfo.Summary = template.HTML(markRedKeywords(getSummary(string(content)), keywords))
		}

		title, err := ioutil.ReadFile(conf.ExtractWebpagePath + md5 + "_title.txt")
		if err != nil {
			log.Printf("%v\r\n", err)
		} else {
			articleInfo.Title = template.HTML(markRedKeywords(string(title), keywords))
		}

		result = append(result, articleInfo)
	}

	return result, total
}
