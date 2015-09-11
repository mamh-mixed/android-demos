package weixin

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/omigo/log"
)

var defaultWeixinClient *http.Client

func getDefaultWeixinClient() *http.Client {
	if defaultWeixinClient == nil {
		defaultWeixinClient = &http.Client{
			Transport: &http.Transport{
				Dial: tcpDail,
			},
		}
	}
	return defaultWeixinClient
}
func getPrivateWeixinClient(cliCrt *tls.Certificate) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: tcpDail,
			TLSClientConfig: &tls.Config{
				// InsecureSkipVerify: true, // only for testing
				Certificates: []tls.Certificate{*cliCrt},
			},
		},
	}
}

func tcpDail(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	ipPort := fetchDns(addr)
	fmt.Printf("%s: dail %s %s\n", time.Now(), network, ipPort)
	conn, err := dial.Dial(network, ipPort)
	if err != nil {
		return conn, err
	}
	fmt.Printf("%s: connect done, local %s, remote %s\n", time.Now(),
		conn.LocalAddr().String(), conn.RemoteAddr().String())
	return conn, err
}

var (
	dnsLock        = &sync.RWMutex{}
	dnsCache       = map[string]string{}
	lastLookupTime time.Time
)

func fetchDns(hostport string) string {

	host, port, _ := net.SplitHostPort(hostport)

	dnsLock.RLock()
	ip := dnsCache[host]
	dnsLock.RUnlock()
	if ip != "" {
		if time.Now().Sub(lastLookupTime) > 4*time.Hour {
			go refreshDNS(host, port)
		}
		return ip + ":" + port
	}

	return refreshDNS(host, port)
}

func refreshDNS(host, port string) string {
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Errorf("lookupIP for server name %s error: %s", host, err)
		return ""
	}

	if len(ips) == 0 {
		return ""
	}
	idx := rand.Intn(len(ips)) // 随机取一个ip
	ip := ips[idx].String()

	dnsLock.Lock()
	dnsCache[host] = ip
	dnsLock.Unlock()

	log.Infof("lookupIP for server name %s => %s:%s, %s", host, ip, port, ips)
	return ip + ":" + port
}
