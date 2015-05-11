package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "github.com/CardInfoLink/quickpay/config"
	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/master"
	"github.com/CardInfoLink/quickpay/pay"
	"github.com/CardInfoLink/quickpay/settle"
	"github.com/omigo/log"

	_ "net/http/pprof"
)

func main() {
	// 日志输出级别
	log.SetOutputLevel(log.Linfo)
	// log.SetFlags(log.Ldate | log.Ltime)

	var (
		argAll, argMaster, argPay, argSettle bool
		port                                 int
	)

	flag.BoolVar(&argAll, "all", false, "Startup QuickAll")
	flag.BoolVar(&argMaster, "master", false, "Startup QuickMaster")
	flag.BoolVar(&argPay, "pay", false, "Startup Quickpay")
	flag.BoolVar(&argSettle, "settle", false, "Startup QuickSettle")
	flag.IntVar(&port, "port", 6800, "server listen port, default QuickMaster 6700, Quickpay 6800, QuickSettle 6900")

	flag.Parse()

	if (!argAll && !argMaster && !argPay && !argSettle) || port == 0 {
		flag.Usage()
		return
	}

	if argAll {
		quickAll(port)
		return
	}

	if (argMaster && argPay) || (argMaster && argSettle) || (argPay && argSettle) {
		flag.Usage()
		fmt.Println("`master` `pay` `settle` must only one to be set")
		return
	}

	if argMaster {
		if port == 6800 {
			fmt.Println("`QuickMaster` must not on port 6800, use 6700")
			port = 6700
		}
		quickMaster(port)
		return
	}

	if argPay {
		quickpay(port)
		return
	}

	if argSettle {
		if port == 6800 {
			fmt.Println("`QuickSettle` must not on port 6800, use 6900")
			port = 6900
		}
		quickSettle(port)
		return
	}
}

func quickAll(port int) {
	// 系统初始化
	master.Initialize()
	settle.Initialize()
	pay.Initialize()

	http.Handle("/", http.FileServer(http.Dir("static")))
	// http.HandleFunc("/quickSettle/", settle.QuickSettle)
	http.HandleFunc("/quickpay/", entrance.Quickpay)

	addr := fmt.Sprintf(":%d", port)
	log.Infof("Quickpay is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func quickMaster(port int) {
	// 系统初始化
	master.Initialize()

	http.Handle("/", http.FileServer(http.Dir("static")))

	// http.HandleFunc("/quickMaster/", master.Quickpay)

	addr := fmt.Sprintf(":%d", port)
	log.Infof("QuickMaster is running on %s", addr)
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
	log.Infof("Quickpay is running on %s", addr)
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
	log.Infof("QuickSettle is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
