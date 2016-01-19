package scanpay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"io"
	// "net"
	"strconv"
	"testing"
	// "time"
)

func TestListenTcp(t *testing.T) {
	// 扫码 TCP 接口，UTF-8 编码传输，UTF-8 签名
	port := goconf.Config.App.TCPAddr
	ListenScanPay(port)

	// 扫码 TCP 接口，GBK 编码传输，UTF-8 签名
	port = goconf.Config.App.TCPGBKAddr
	ListenScanPay(port, true)
}

func TestDailTcp(t *testing.T) {
	var err error
	// addr := "overseas.show.money:6000"
	addr := "52.192.213.82:6000"

	config := &tls.Config{}
	config.InsecureSkipVerify = true

	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Errorf("can't connect to cil-online tcp://%s: %s", addr, err)
		return
	}
	defer conn.Close()
	req := new(model.ScanPayRequest)
	req.Busicd = "PURC"
	encoded, _ := json.Marshal(req)
	head := fmt.Sprintf("%04d", len(encoded))
	io.WriteString(conn, head+string(encoded))

	for {
		dl := make([]byte, 4)
		conn.Read(dl)
		if string(dl) != "" {
			l, _ := strconv.Atoi(string(dl))
			bs := make([]byte, l)
			conn.Read(bs)
			t.Logf("%s", string(bs))
			return
		}
		return
	}

}
