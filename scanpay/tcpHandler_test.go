package scanpay

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
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
	addr := ":3000"
	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
		log.Errorf("can't connect to cil-online tcp://%s: %s", addr, err)
		return
	}
	req := new(model.ScanPay)
	req.Busicd = "purc"
	encoded, _ := json.Marshal(req)
	head := fmt.Sprintf("%04d", len(encoded))
	io.WriteString(conn, head+string(encoded))
}
