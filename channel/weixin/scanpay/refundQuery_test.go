package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
)

func TestRefundQuery(t *testing.T) {
	d := &RefundQueryReq{
		CommonParams: testCommonParams,

		DeviceInfo:    "xxx",           // 设备号
		TransactionId: "",              // 微信订单号
		OutTradeNo:    "1441615439898", // 商户订单号
		OutRefundNo:   "",              // 商户退款单号
		RefundId:      "",              // 微信退款单号
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
