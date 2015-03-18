package main

import (
	"log"
	"net/http"
	"github.com/CardInfoLink/quickpay/entrance/bindingpay"

	"github.com/omigo/g"
)

func main() {
	// 日志输出级别
	g.SetLevel(g.LevelDebug)

	addr := ":3000"
	log.Printf("QuickPay is running on %s", addr)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", bindingpay.BindingPay)

	log.Fatal(http.ListenAndServe(addr, nil))
}
