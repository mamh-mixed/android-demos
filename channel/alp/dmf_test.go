package alp

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"testing"
)

var qrCodePay = &model.QrCodePay{
	GoodsInfo: "鞋子,1000,2;衣服,1500,3",
}

func TestProcessBarcodePay(t *testing.T) {
	resp := DefaultClient.ProcessBarcodePay(qrCodePay)
	log.Debugf("%+v", resp)
}
