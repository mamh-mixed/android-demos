package scanpay

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

var scanPayBarcodePay = &model.ScanPay{
	GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
	OrderNum:   tools.Millisecond(),
	ScanCodeId: "283112524283364763",
	Inscd:      "CIL00002",
	Txamt:      "000000000002",
	Busicd:     "purc",
	Mchntid:    "CIL0001",
}

var scanPayQrCodeOfflinePay = &model.ScanPay{
	GoodsInfo: "鞋子,1000,2;衣服,1500,3",
	OrderNum:  tools.Millisecond(),
	Inscd:     "CIL00002",
	Txamt:     "000000000001",
	Busicd:    "paut",
	Mchntid:   "CIL0001",
	Chcd:      "ALP",
}

var scanPayEnquiry = &model.ScanPay{
	Busicd:       "inqy",
	Mchntid:      "CIL0001",
	Inscd:        "CIL00002",
	OrigOrderNum: "1435198837472",
}

var scanPayRefund = &model.ScanPay{
	Busicd:       "refd",
	Mchntid:      "CIL0001",
	OrderNum:     tools.Millisecond(),
	OrigOrderNum: "1435199869254",
	Inscd:        "CIL00002",
	Txamt:        "000000000001",
}

func TestScanPay(t *testing.T) {
	log.SetOutputLevel(log.Ldebug)
	reqBytes, _ := json.Marshal(scanPayRefund)
	respBytes := Router(reqBytes)
	fmt.Println(string(respBytes))
}
