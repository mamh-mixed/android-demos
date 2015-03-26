package main

import (
	"github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"log"
	"net/http"

	"github.com/omigo/g"
)

func main() {
	// 日志输出级别
	g.SetLevel(g.LevelDebug)

	// 测试端口
	addr := ":3009"
	// 生产端口
	// addr := ":3000"
	log.Printf("QuickPay is running on %s", addr)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", bindingpay.BindingPay)

	log.Fatal(http.ListenAndServe(addr, nil))
}
