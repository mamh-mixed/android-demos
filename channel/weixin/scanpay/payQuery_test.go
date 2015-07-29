package scanpay

import (
	"encoding/xml"
	"testing"

	"github.com/omigo/log"
)

func TestPayQuery(t *testing.T) {
	d := &PayQueryReq{
		CommonParams: testCommonParams,

		TransactionId: "",              // 微信的订单号，优先使用
		OutTradeNo:    "1438137518988", // 商户系统内部的订单号，当没提供transaction_id时需要传这个
	}

	r := &PayQueryResp{}
	err := base(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
	}
}

func TestScanPayQueryGenSign(t *testing.T) {
	d := &PayQueryReq{
		CommonParams: testCommonParams,

		TransactionId: "1010070115201506230291458545", // 微信支付订单号
		OutTradeNo:    "",                             // 商户订单号
	}

	d.GenSign()

	t.Log(d.Sign)

	xmlBytes, err := xml.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Errorf("struct(%#v) to xml error: %s", d, err)
	}

	t.Log(string(xmlBytes))
}
