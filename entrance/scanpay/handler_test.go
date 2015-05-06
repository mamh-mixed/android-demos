package scanpay

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
)

var scanPay = &model.ScanPay{
	GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
	OrderNum:   tools.Millisecond(),
	ScanCodeId: "289434710505996982",
	Inscd:      "CIL00002",
	Txamt:      "0.01",
	Busicd:     "purc",
	Mchntid:    "CIL0001",
}

func TestScanPay(t *testing.T) {
	// mongo.Connect()
	reqBytes, _ := json.Marshal(scanPay)
	respBytes := Router(reqBytes)
	fmt.Println(string(respBytes))
}
