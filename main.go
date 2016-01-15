package main

import (
	"net/http"
	"runtime"

	"github.com/CardInfoLink/quickpay/app"
	"github.com/CardInfoLink/quickpay/auth"
	"github.com/CardInfoLink/quickpay/bindingpay"
	"github.com/CardInfoLink/quickpay/check"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/crontab"
	_ "github.com/CardInfoLink/quickpay/flags"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/master"
	"github.com/CardInfoLink/quickpay/scanpay"
	"github.com/CardInfoLink/quickpay/settle"
	"github.com/omigo/log"

	// _ "net/http/pprof"
)

// 生成版本文件
//go:generate version

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetOutputLevel(goconf.Config.App.LogLevel)
	// log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llevel | log.Lprojectfile)

	startScanpay()    // 扫码支付
	startBindingpay() // 绑定支付
	startSettle()     // 清分任务
	startMaster()     // 管理平台
	startApp()        // App接口
	startContab()     // 定时任务
	startAlipayAuth() // 支付宝授权

	log.Infof("Quickpay HTTP is listening, addr=%s", goconf.Config.App.HTTPAddr)
	log.Error(http.ListenAndServe(goconf.Config.App.HTTPAddr, nil))
}

func startScanpay() {
	// 扫码支付 HTTP 接口，包括微信、支付宝的异步通知
	http.Handle("/scanpay/", scanpay.Route())

	// 扫码 TCP 接口，UTF-8 编码传输，UTF-8 签名
	port := goconf.Config.App.TCPAddr
	scanpay.ListenScanPay(port)

	// 扫码 TCP 接口，GBK 编码传输，UTF-8 签名
	port = goconf.Config.App.TCPGBKAddr
	scanpay.ListenScanPay(port, true)
}

func startBindingpay() {
	// 初始化卡 Bin 树
	core.BuildTree()
	// 检查数据配置是否有变化
	check.DoCheck()

	http.Handle("/bindingpay/", bindingpay.Route())
}

func startSettle() {
	settle.DoSettWork()
}

func startMaster() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/index/", http.StripPrefix("/index/", http.FileServer(http.Dir("admin"))))
	// http.Handle("/", handlers.LoggingHandler(os.Stdout, http.FileServer(http.Dir("static"))))

	http.Handle("/master/", master.Route())
}

func startApp() {
	http.Handle("/app/", app.Route())
	http.Handle("/app/v3/", app.RouteV3())
	app.StartPush()
}

func startContab() {
	crontab.Start()
}

func startAlipayAuth() {
	http.HandleFunc("/alipay/auth", auth.AuthHandle)
}
