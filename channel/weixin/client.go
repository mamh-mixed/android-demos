package weixin

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/CardInfoLink/log"
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

	ipPort := wxScanpayIP + ":" + wxScanpayPort

	log.Debugf("%s: dail %s %s\n", time.Now(), network, ipPort)
	conn, err := dial.Dial(network, ipPort)
	if err != nil {
		return conn, err
	}
	log.Debugf("%s: connect done, local %s, remote %s\n", time.Now(),
		conn.LocalAddr().String(), conn.RemoteAddr().String())

	return conn, err
}
