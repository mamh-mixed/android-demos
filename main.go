package main

import (
	"net/http"

	"github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/omigo/log"
)

func main() {
	// 日志输出级别
	log.SetOutputLevel(log.Ldebug)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	addr := ":3000"
	log.Debugf("QuickPay is running on %s", addr)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/quickpay/", bindingpay.BindingPay)

	log.Fatal(http.ListenAndServe(addr, nil))
}
