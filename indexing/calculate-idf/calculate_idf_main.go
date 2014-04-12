package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gansidui/gose/indexing/participleutil"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

// 配置文件
type Config struct {
	ExtractWebpagePath string
	ExtractUrlDataPath string
}

var conf *Config           // 配置
var numFile int            // 语料库的文档总数
var wordMap map[string]int // 包含该词的文档数

// 初始化
func init() {
	setLogOutput()
	participleutil.LoadDic("../participleutil/mydic.txt")
	conf = NewConfig()
	numFile = 0
	wordMap = make(map[string]int)
}

// 读取配置文件
func NewConfig() *Config {
	file, err := ioutil.ReadFile("./idf.conf")
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
	logfile, err := os.OpenFile("./idf.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	// 注意这里不能关闭logfile
	if err != nil {
		log.Printf("%v\r\n", err)
	}
	log.SetOutput(logfile)
}

// 计算逆文档频率
func calculateIDF() {
	start := time.Now()

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

		// 读取文档
		content, _ := ioutil.ReadFile(conf.ExtractWebpagePath + md5 + "_body.txt")
		// 得到分词结果
		ss := participleutil.Participle(string(content))
		// 去重
		m := make(map[string]bool)
		for _, v := range ss {
			if !m[v] {
				m[v] = true
			}
		}
		// 保存结果
		for k, _ := range m {
			wordMap[k]++
		}
	}

	fmt.Println("calculateIDF used time:", time.Since(start))
}

// 将逆文档频率信息保存到数据库
func writeDatabase() {
	start := time.Now()

	// 打开数据库
	db, err := sql.Open("sqlite3", "./idf.db")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer db.Close()

	// 先删表再建表
	db.Exec("drop table data")
	_, err = db.Exec("create table data(word varchar(30), idf float)")
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	// 启动事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	for k, v := range wordMap {
		idf := math.Log10(float64(numFile) / float64(v+1))
		if idf < 0 {
			idf = 0
		}
		_, err := tx.Exec("insert into data(word, idf) values(?, ?)", k, idf)
		if err != nil {
			log.Fatal(err, "\r\n")
		}
	}
	tx.Commit()

	fmt.Println("write database used time:", time.Since(start))
}

// 将逆文档频率信息保存到文件
func writeFile() {
	start := time.Now()

	file, err := os.OpenFile("./idf.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer file.Close()

	for k, v := range wordMap {
		idf := math.Log10(float64(numFile) / float64(v+1))
		if idf < 0 {
			idf = 0
		}
		file.WriteString(k + "-----" + fmt.Sprintf("%f", idf) + "\r\n")
	}

	fmt.Println("write file used time:", time.Since(start))
}

func main() {
	log.Printf("%v\r\n", "start......")
	start := time.Now()
	calculateIDF()
	writeDatabase()
	writeFile()
	log.Printf("used time: %v", time.Since(start))
}
