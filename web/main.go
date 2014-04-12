package main

import (
	"github.com/gansidui/gose/search"
	"github.com/gansidui/gose/web/controllers"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	search.ReadConfig("../search/search.conf")
	search.InitSearch()
}

func main() {
	println("start......")

	http.HandleFunc("/", controllers.HomePage)
	http.HandleFunc("/search", controllers.SearchPage)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
