package main

import (
	"net/http"

	"github.com/CardInfoLink/quickpay/conf"
	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/omigo/log"

	_ "net/http/pprof"
)

func main() {
	// 日志输出级别
	log.SetOutputLevel(log.Ldebug)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 系统初始化
	conf.Initialize()

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", entrance.Quickpay)

	addr := ":3009"
	log.Debugf("QuickPay is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
