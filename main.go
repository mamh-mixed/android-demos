package main

import (
	"log"
	"net/http"
	"quickpay/handler"

	"github.com/omigo/g"
)

func main() {
	// 日志输出级别
	g.SetLevel(g.LevelTrace)

	addr := ":3000"
	log.Printf("Shoumoney is running at %s", addr)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", handler.Quickpay)

	log.Fatal(http.ListenAndServe(addr, nil))
}
