package main

import (
	"log"
	"net/http"
	"quickpay/enter/bindingpay"

	"github.com/omigo/g"
)

// test

func main() {
	// 日志输出级别
	g.SetLevel(g.LevelTrace)

	addr := ":3000"
	log.Printf("Shoumoney is running at %s", addr)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", bindingpay.BindingPay)

	log.Fatal(http.ListenAndServe(addr, nil))
}
