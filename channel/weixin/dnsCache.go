package weixin

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/log"
)

var wxScanpayIP, wxScanpayPort string

func init() {
	wxScanpayURL := goconf.Config.WeixinScanPay.URL
	u, err := url.Parse(wxScanpayURL)
	if err != nil {
		fmt.Printf("parse weixin scanpay url %s error: %s\n", wxScanpayURL, err)
		os.Exit(1)
	}

	re, _ := regexp.Compile(`:\d+$`)
	if re.MatchString(u.Host) {
		wxScanpayIP, wxScanpayPort, err = net.SplitHostPort(u.Host) // 初始值是个域名不是IP，这样避免空值错误
		if err != nil {
			fmt.Printf("splitHostPort %s error: %s\n", u.Host, err)
			os.Exit(2)
		}
	} else {
		wxScanpayIP = u.Host // 初始值是个域名不是IP，这样避免空值错误
		if u.Scheme == "https" {
			wxScanpayPort = "443"
		} else {
			wxScanpayPort = "80"
		}
	}

	go refreshIP(wxScanpayIP)
}

func refreshIP(host string) {
	var dnsCacheRefreshTime = time.Duration(goconf.Config.WeixinScanPay.DNSCacheRefreshTime)

	tick := time.Tick(dnsCacheRefreshTime)
	for {
		wxScanpayIP = lookupIP(host)

		<-tick
	}
}

func lookupIP(host string) string {
	start := time.Now()
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Errorf("lookupIP for server name %s error: %s", host, err)
		return host
	}
	end := time.Now()
	log.Infof("=== %s === host: %s => %s", end.Sub(start), host, ips)

	if len(ips) == 0 {
		return host
	}

	idx := rand.Intn(len(ips)) // 随机取一个ip
	ip := ips[idx].String()
	log.Infof("lookupIP for server name %s => %s", host, ip)

	return ip
}
