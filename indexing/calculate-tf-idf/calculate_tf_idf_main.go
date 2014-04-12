package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gansidui/gose/indexing/participleutil"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// 配置文件
type Config struct {
	ExtractWebpagePath string
	ExtractUrlDataPath string
}

var conf *Config               // 配置
var numFile int                // 文档数量
var numWord int                // 关键词数量
var wordIDF map[string]float32 // word -- IDF值
var docID map[string]int       // doc -- id 文档对应的标号
var words []string             // 关键词
var docids []int               // 文档标号
var tfidfs []float32           // TF-IDF值

// 初始化
func init() {
	setLogOutput()
	participleutil.LoadDic("../participleutil/mydic.txt")
	conf = NewConfig()
	numFile = 0
	numWord = 0
	wordIDF = make(map[string]float32)
	docID = make(map[string]int)
	words = make([]string, 0)
	docids = make([]int, 0)
	tfidfs = make([]float32, 0)
	readIDF()
	readDocId()
	initDatabase()
}

// 读取配置文件
func NewConfig() *Config {
	file, err := ioutil.ReadFile("./tfidf.conf")
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
	logfile, err := os.OpenFile("./tfidf.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	// 注意这里不能关闭logfile
	if err != nil {
		log.Printf("%v\r\n", err)
	}
	log.SetOutput(logfile)
}

// 读取 IDF 值
func readIDF() {
	db, err := sql.Open("sqlite3", "../calculate-idf/idf.db")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var word string
	var idf float32
	for rows.Next() {
		rows.Scan(&word, &idf)
		wordIDF[word] = idf
	}
}

// 读取文档编号
func readDocId() {
	db, err := sql.Open("sqlite3", "../doc-mark/docmark.db")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var md5 string
	var id int
	for rows.Next() {
		rows.Scan(&md5, &id)
		docID[md5] = id
	}
}

// 计算TF-IDF
func calculateTFIDF() {
	start := time.Now()
	// 读取文档数据
	db, err := sql.Open("sqlite3", conf.ExtractUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	rows, err := db.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var md5 string
	for rows.Next() {
		rows.Scan(&md5)
		// 计数
		numFile++
		fmt.Println("numFile:", numFile)

		// 读取正文文档
		content, _ := ioutil.ReadFile(conf.ExtractWebpagePath + md5 + "_body.txt")
		// 得到分词结果
		ss := participleutil.Participle(string(content))
		totalWord := len(ss) // 文档的总词数
		// 统计每个词在这篇文档中出现的次数
		m := make(map[string]int)
		for _, v := range ss {
			m[v]++
		}

		// 读取标题文档
		content, _ = ioutil.ReadFile(conf.ExtractWebpagePath + md5 + "_title.txt")
		ss = participleutil.Participle(string(content))
		for _, v := range ss {
			m[v] += 5
		}

		docid := docID[md5] // 文档ID

		for k, v := range m {
			tf := float32(float32(v) / float32(totalWord)) // 词频
			idf := wordIDF[k]                              // 逆文档频率
			words = append(words, k)                       // 关键词
			docids = append(docids, docid)                 // 文档标号
			tfidfs = append(tfidfs, tf*idf)                // TF-IDF值
			numWord++
			if numWord%2000000 == 0 {
				writeDatabase()
				words = []string{}
				docids = []int{}
				tfidfs = []float32{}
			}
		}
	}

	writeDatabase()
	fmt.Println("calculateTFIDF used time:", time.Since(start))
}

// 初始化数据库
func initDatabase() {
	// 打开数据库
	db, err := sql.Open("sqlite3", "./tfidf.db")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 先删表再建表
	db.Exec("drop table data")
	_, err = db.Exec("create table data(word varchar(30), docid integer, tfidf float)")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
}

// 将TF-IDF信息保存到数据库
func writeDatabase() {
	// 打开数据库
	db, err := sql.Open("sqlite3", "./tfidf.db")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 启动事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	for k, _ := range words {
		tx.Exec("insert into data(word, docid, tfidf) values(?, ?, ?)", words[k], docids[k], tfidfs[k])
	}

	tx.Commit()
}

func main() {
	log.Printf("%v\r\n", "start......")
	start := time.Now()

	calculateTFIDF()

	log.Printf("记录总数: %d\r\n", numWord)
	log.Printf("平均每篇文档有%d个不同的词\r\n", numWord/numFile)
	log.Printf("used time: %v", time.Since(start))
}
