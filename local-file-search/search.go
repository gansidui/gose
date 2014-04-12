package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	db, err := sql.Open("sqlite3", "./localfile.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := fmt.Sprintf("%%%s%%", scanner.Text()) // 子串查询

		start := time.Now()
		rows, err := db.Query("select * from info where name like ?", s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("query used time:", time.Since(start))

		num := 0
		for rows.Next() {
			var name string
			var path string
			rows.Scan(&name, &path)
			fmt.Println(name)
			fmt.Println(path)
			fmt.Println("--------------------------------------------------")
			num++
			if num > 5 {
				break
			}
		}
		// 必须释放查询结果，不然内存暴涨
		rows.Close()

		// start = time.Now()
		// rows, err = db.Query("select count(*) from info where name like ?", s)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// rows.Next()
		// rows.Scan(&num)
		// rows.Close()
		// fmt.Println("total:", num, "used time:", time.Since(start))
	}
}
