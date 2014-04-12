package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// 配置文件
type Config struct {
	ExtractUrlDataPath string
}

var conf *Config // 配置

// 初始化
func init() {
	setLogOutput()
	conf = NewConfig()
}

// 读取配置文件
func NewConfig() *Config {
	file, err := ioutil.ReadFile("./docmark.conf")
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
	logfile, err := os.OpenFile("./docmark.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	// 注意这里不能关闭logfile
	if err != nil {
		log.Printf("%v\r\n", err)
	}
	log.SetOutput(logfile)
}

func main() {
	log.Printf("%v\r\n", "start......")
	start := time.Now()

	var curMaxId int = 0      // 目前的最大标号
	docs := make([]string, 0) // 待标号的文档的md5值

	// 读取待标号文档
	docdb, err := sql.Open("sqlite3", conf.ExtractUrlDataPath)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer docdb.Close()

	rows, err := docdb.Query("select * from data")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer rows.Close()

	var md5 string
	for rows.Next() {
		rows.Scan(&md5)
		docs = append(docs, md5)
	}

	// 开始标号
	markdb, err := sql.Open("sqlite3", "./docmark.db")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	defer markdb.Close()

	markdb.Exec("create table data(md5 varchar(32) not null primary key, id integer not null)")
	// 获取最大标号
	rows, err = markdb.Query("select max(id) from data")
	if err != nil {
		curMaxId = 0
	}
	rows.Next()
	rows.Scan(&curMaxId)
	rows.Close()

	tx, err := markdb.Begin()
	for _, v := range docs {
		curMaxId++
		tx.Exec("insert into data(md5, id) values(?, ?)", v, curMaxId)
	}
	tx.Commit()

	log.Printf("used time: %v", time.Since(start))
}
