package main

import (
	"net/http"
	"runtime"

	"github.com/CardInfoLink/quickpay/data"
	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/master"
	"github.com/CardInfoLink/quickpay/pay"
	"github.com/CardInfoLink/quickpay/settle"
	"github.com/omigo/log"

	_ "net/http/pprof"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetOutputLevel(goconf.Config.App.LogLevel)
	log.SetFlags(log.Ldate | log.Ltime | log.Llevel | log.Llongfile)

	// 系统初始化
	master.Initialize()
	settle.Initialize()
	pay.Initialize()

	http.Handle("/", http.FileServer(http.Dir("static/app")))
	// http.HandleFunc("/quickSettle/", settle.QuickSettle)
	http.HandleFunc("/quickpay/", entrance.Quickpay)
	http.HandleFunc("/scanpay/", entrance.Scanpay)
	http.HandleFunc("/quickMaster/", master.QuickMaster)
	http.HandleFunc("/master/", master.MasterRoute)
	http.HandleFunc("/service/", master.Service)
	http.HandleFunc("/qp/back/", entrance.AsyncNotify)
	http.HandleFunc("/import", data.ImportOldTrans)

	log.Infof("Quickpay HTTP is listening, addr=%s", goconf.Config.App.HTTPAddr)
	log.Error(http.ListenAndServe(goconf.Config.App.HTTPAddr, nil))
}
