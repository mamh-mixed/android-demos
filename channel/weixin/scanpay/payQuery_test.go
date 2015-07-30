package scanpay

import "testing"

func TestPayQuery(t *testing.T) {
	d := &PayQueryReq{
		CommonParams: testCommonParams,

		TransactionId: "",              // 微信的订单号，优先使用
		OutTradeNo:    "1438140832742", // 商户系统内部的订单号，当没提供transaction_id时需要传这个
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
