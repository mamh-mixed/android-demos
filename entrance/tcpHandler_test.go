package entrance

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"io"
	"net"
	"testing"
)

func TestListenTcp(t *testing.T) {
	Listen()
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
	req := new(model.QrCodePay)
	req.Busicd = "purc"
	encoded, _ := json.Marshal(req)
	head := fmt.Sprintf("%04d", len(encoded))
	io.WriteString(conn, head+string(encoded))
}
