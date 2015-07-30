package scanpay

import "testing"

func TestReverse(t *testing.T) {
	d := &ReverseReq{
		CommonParams: testCommonParams,

		TransactionId: "1002080115201507300508612145",     // 微信的订单号，优先使用
		OutTradeNo:    "ffef191aca3e4cf15e7672a5eb2113da", // 商户系统内部的订单号，当没提供transaction_id时需要传这个
	}

	r := &ReverseResp{}
	err := base(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
	}
}
