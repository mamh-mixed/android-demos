package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
)

func TestRefundQuery(t *testing.T) {
	d := &RefundQueryReq{
		CommonParams: testCommonParams,

		DeviceInfo:    "xxx",                              // 设备号
		TransactionId: "1002080115201507300508405226",     // 微信订单号
		OutTradeNo:    "c0048d6aff60453a4dd7c66ea74cbcc5", // 商户订单号
		OutRefundNo:   "f08c673a804a441947464733761bba7b", // 商户退款单号
		RefundId:      "2002080115201507300022072313",     // 微信退款单号
	}

	r := &RefundQueryResp{}
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
