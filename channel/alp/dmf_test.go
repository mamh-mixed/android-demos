package alp

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"testing"
)

var qrCodePay = &model.QrCodePay{
	GoodsInfo:    "鞋子,1000,2;衣服,1500,3",
	ChanOrderNum: "awdajwdadn",
	Key:          "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ScanCodeId:   "23131242413",
	Txamt:        "0.01",
}

func TestProcessBarcodePay(t *testing.T) {
	resp := DefaultClient.ProcessBarcodePay(qrCodePay)
	log.Debugf("%+v", resp)
}
