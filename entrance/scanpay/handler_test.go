package scanpay

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
)

var scanPayBarcodePay = &model.ScanPay{
	GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
	OrderNum:   tools.Millisecond(),
	ScanCodeId: "289434710505996982",
	Inscd:      "CIL00002",
	Txamt:      "5",
	Busicd:     "purc",
	Mchntid:    "CIL0001",
}

var scanPayEnquiry = &model.ScanPay{
	Busicd:       "inqy",
	Mchntid:      "CIL0001",
	OrigOrderNum: "1431414042270",
}

func TestScanBarcodePay(t *testing.T) {
	// mongo.Connect()
	reqBytes, _ := json.Marshal(scanPayBarcodePay)
	respBytes := Router(reqBytes)
	fmt.Println(string(respBytes))
}

func TestScanPayEnquiry(t *testing.T) {
	reqBytes, _ := json.Marshal(scanPayEnquiry)
	respBytes := Router(reqBytes)
	fmt.Println(string(respBytes))
}
