package scanpay

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"testing"
)

var scanPayBarcodePay = &model.ScanPay{
	GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
	OrderNum:   tools.Millisecond(),
	ScanCodeId: "284339495400486450",
	Inscd:      "CIL00002",
	Txamt:      "0.01",
	Busicd:     "purc",
	Mchntid:    "CIL0001",
}

var scanPayQrCodeOfflinePay = &model.ScanPay{
	GoodsInfo: "鞋子,1000,2;衣服,1500,3",
	OrderNum:  tools.Millisecond(),
	Inscd:     "CIL00002",
	Txamt:     "5",
	Busicd:    "paut",
	Mchntid:   "CIL0001",
	Chcd:      "ALP",
}

var scanPayEnquiry = &model.ScanPay{
	Busicd:       "inqy",
	Mchntid:      "CIL0001",
	OrigOrderNum: "1431414042270",
}

var scanPayRefund = &model.ScanPay{
	Busicd:       "refd",
	Mchntid:      "CIL0001",
	OrderNum:     tools.Millisecond(),
	OrigOrderNum: "1435117750095",
	Inscd:        "CIL00002",
	Txamt:        "0.01",
}

func TestScanPay(t *testing.T) {
	log.SetOutputLevel(log.Ldebug)
	reqBytes, _ := json.Marshal(scanPayRefund)
	respBytes := Router(reqBytes)
	fmt.Println(string(respBytes))
}
