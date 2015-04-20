package alp

import (
	"github.com/omigo/log"
	"testing"
)

var barCodePayReq = &AlpRequest{
	Service:       "alipay.acquire.createandpay",
	Partner:       "",
	Charset:       "UTF-8",
	NotifyUrl:     "",
	OutTradeNo:    "440583111100000",
	Subject:       "",
	GoodsDetail:   "",
	ProductCode:   "BARCODE_PAY_OFFLINE",
	TotalFee:      "0.01",
	SellerId:      "string",
	Currency:      "156",
	ExtendParams:  "",
	ItBPay:        "2m",
	DynamicIdType: "bar_code",
	DynamicId:     "",
	Key:           "",
}

func TestProcessBarcodePay(t *testing.T) {
	resp := Obj.ProcessBarcodePay(barCodePayReq)
	log.Debugf("%+v", resp)
}
