package scanpay

import "testing"

func TestRefundQuery(t *testing.T) {
	d := &RefundQueryReq{
		CommonParams: testCommonParams,

		DeviceInfo:    "",                                 // 设备号
		TransactionId: "",                                 // 微信订单号
		OutTradeNo:    "7a5d8c60e1284fe8697af775c60d15d7", // 商户订单号
		OutRefundNo:   "005ffbd5f14e4fda429a745d7987b0be", // 商户退款单号
		RefundId:      "2010070115201506260012210940",     // 微信退款单号
	}

	r := &RefundQueryResp{}
	err := base(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
	}
}
