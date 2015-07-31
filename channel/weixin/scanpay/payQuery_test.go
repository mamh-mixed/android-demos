package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
)

func TestPayQuery(t *testing.T) {
	d := &PayQueryReq{
		CommonParams: testCommonParams,

		TransactionId: "",                                 // 微信的订单号，优先使用
		OutTradeNo:    "723f7c7a90f14dfb47356d4cfdebb212", // 商户系统内部的订单号，当没提供transaction_id时需要传这个
	}

	r := &PayQueryResp{}
	err := weixin.Execute(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
		t.FailNow()
	}

	if r.ResultCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
		t.FailNow()
	}
}
