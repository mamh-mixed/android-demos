package scanpay

import "testing"

func TestClose(t *testing.T) {
	d := &CloseReq{
		CommonParams: testCommonParams,

		TransactionId: "",
		OutTradeNo:    "21212",
	}

	r := &PayResp{}

	err := base(d, r)
	if err != nil {
		t.Errorf("weixin close error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin close return: %#v", r)
	}
}
