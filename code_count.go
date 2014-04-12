package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var numCodeRows int

func getRowsFromFile(path string) int {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	num := 0
	re := bufio.NewReader(file)
	for {
		_, _, err = re.ReadLine()
		if err != nil {
			break
		}
		num++
	}
	fmt.Println(path, "------------", num)
	return num
}

func WalkFunc(path string, info os.FileInfo, err error) error {
	if path == "calculate_code_rows.go" {
		return nil
	}
	if filepath.Ext(path) == ".go" {
		numCodeRows += getRowsFromFile(path)
	}
	return nil
}

func main() {
	numCodeRows = 0
	filepath.Walk("./", WalkFunc)
	fmt.Println("total:", numCodeRows)
}
