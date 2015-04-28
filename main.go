package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/master"
	"github.com/CardInfoLink/quickpay/pay"
	"github.com/CardInfoLink/quickpay/settle"
	"github.com/omigo/log"
)

func main() {
	// 日志输出级别
	log.SetOutputLevel(log.Ldebug)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var (
		argMaster, argPay, argSettle bool
		port                         int
	)

	flag.BoolVar(&argMaster, "master", false, "Startup QuickMaster")
	flag.BoolVar(&argPay, "pay", false, "Startup Quickpay")
	flag.BoolVar(&argSettle, "settle", false, "Startup QuickSettle")
	flag.IntVar(&port, "port", 3800, "server listen port, default QuickMaster 3700, Quickpay 3800, QuickSettle 3900")

	flag.Parse()

	if (!argMaster && !argPay && !argSettle) || port == 0 {
		flag.Usage()
		return
	}
	if (argMaster && argPay) || (argMaster && argSettle) || (argPay && argSettle) {
		flag.Usage()
		fmt.Println("`master` `pay` `settle` must only one to be set")
		return
	}

	if argMaster {
		if port == 3800 {
			fmt.Println("`QuickMaster` must not on port 3800, use 3700")
			port = 3700
		}
		quickMaster(port)
		return
	}

	if argPay {
		quickpay(port)
		return
	}

	if argSettle {
		if port == 3800 {
			fmt.Println("`QuickSettle` must not on port 3800, use 3900")
			port = 3900
		}
		quickSettle(port)
		return
	}
}

func quickMaster(port int) {
	// 系统初始化
	master.Initialize()

	http.Handle("/", http.FileServer(http.Dir("static")))

	// http.HandleFunc("/quickMaster/", master.Quickpay)

	addr := fmt.Sprintf(":%d", port)
	log.Debugf("QuickMaster is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func quickpay(port int) {
	// 系统初始化
	pay.Initialize()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Quickpay is running"))
	})

	http.HandleFunc("/quickpay/", entrance.Quickpay)

	addr := fmt.Sprintf(":%d", port)
	log.Debugf("Quickpay is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func quickSettle(port int) {
	// 系统初始化
	settle.Initialize()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("QuickSettle is running"))
	})

	// http.HandleFunc("/quickSettle/", settle.QuickSettle)

	addr := fmt.Sprintf(":%d", port)
	log.Debugf("QuickSettle is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
