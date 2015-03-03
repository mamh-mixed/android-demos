package main

import (
	"log"
	"net/http"
	"quickpay/handler"
)

func main() {
	// 为log添加短文件名,方便查看行数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	addr := ":3000"
	log.Printf("Shoumoney is running at %s", addr)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", handler.Quickpay)

	log.Fatal(http.ListenAndServe(addr, nil))
}
