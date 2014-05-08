package main

import (
	"code.google.com/p/mahonia"
	"encoding/json"
	"fmt"
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

var numFile int

// 把文件从gb2312编码转换成utf8编码
func walkFunc(path string, info os.FileInfo, err error) error {
	file, err := os.Open(path)
	checkError(err)
	defer file.Close()

	decoder := mahonia.NewDecoder("gb2312")
	data, err := ioutil.ReadAll(decoder.NewReader(file))
	ioutil.WriteFile(path, data, os.ModePerm)

	numFile++
	fmt.Println("numFile:", numFile)

	return nil
}

func main() {
	numFile = 0
	start := time.Now()
	conf := NewConfig()
	filepath.Walk(conf.ExtractWebpagePath, walkFunc)
	fmt.Println("Used time:", time.Since(start))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
