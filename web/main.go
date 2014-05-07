package main

import (
	"flag"
	"fmt"
	"github.com/gansidui/gose/search"
	"github.com/gansidui/gose/web/controllers"
	"log"
	"net"
	"net/http"
)

var (
	ip   string
	port string
)

func init() {
	defaultIP := "127.0.0.1"
	defaultPort := "9090"

	// 获取本机的IP(A global unicast address)
	addr, _ := net.InterfaceAddrs()
	for _, v := range addr {
		IP := net.ParseIP(v.String())
		if IP.IsGlobalUnicast() {
			defaultIP = v.String()
			break
		}
	}

	flag.StringVar(&ip, "ip", defaultIP, "ip address")
	flag.StringVar(&port, "port", defaultPort, "port number")
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	search.ReadConfig("../search/search.conf")
	search.InitSearch()
}

func main() {
	http.HandleFunc("/", controllers.HomePage)
	http.HandleFunc("/search", controllers.SearchPage)

	addr := ip + ":" + port
	fmt.Println("Listenning:", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
