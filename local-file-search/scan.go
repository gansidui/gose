package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	name string
	path string
}

var isSkipDir map[string]bool // 需要过滤的目录
var infoSlice []*FileInfo     // 保存所有文件信息

// 读取需要过滤掉的目录
func ReadSkipDir() {
	file, err := os.Open("./skipdir.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	isSkipDir = make(map[string]bool, 0)
	re := bufio.NewReader(file)
	for {
		line, _, err := re.ReadLine()
		if err != nil {
			break
		}
		isSkipDir[string(line)] = true
	}
}

func WalkFunc(path string, info os.FileInfo, err error) error {
	// 有错误就跳过这个目录
	if err != nil {
		return filepath.SkipDir
	}
	// 需要过滤的path
	if isSkipDir[path] {
		return filepath.SkipDir
	}

	infoSlice = append(infoSlice, &FileInfo{name: info.Name(), path: path})
	return nil
}

func main() {
	// 为log添加短文件名，方便查看行数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// 初始化
	ReadSkipDir()
	infoSlice = make([]*FileInfo, 0)

	// 遍历所有磁盘，获取文件信息
	start := time.Now()
	for i := 0; i < 26; i++ {
		root := fmt.Sprintf("%c:", 'A'+i)
		filepath.Walk(root, WalkFunc)
	}
	fmt.Println("aquire file info used time:", time.Since(start))

	// 打开数据库
	db, err := sql.Open("sqlite3", "./localfile.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 先删表，再建表, 如果info不存在，drop时返回的err应该忽略
	db.Exec("drop table info")
	_, err = db.Exec("create table info(name nvarchar(256), path nvarchar(256))")
	if err != nil {
		log.Fatal(err)
	}

	// 将数据写入到info表中
	start = time.Now()
	// 开始一个事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range infoSlice {
		_, err := tx.Exec("insert into info(name, path) values(?, ?)", v.name, v.path)
		if err != nil {
			break
		}
	}
	tx.Commit()
	fmt.Println("write data to database used time:", time.Since(start))
}
