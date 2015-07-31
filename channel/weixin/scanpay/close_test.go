package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
)

func TestClose(t *testing.T) {
	d := &CloseReq{
		CommonParams: testCommonParams,

		TransactionId: "",
		OutTradeNo:    "723f7c7a90f14dfb47356d4cfdebb212",
	}

	r := &CloseResp{}

	err := weixin.Execute(d, r)
	if err != nil {
		t.Errorf("weixin close error: %s", err)
		t.FailNow()
	}

	if r.ResultCode != "SUCCESS" {
		t.Logf("weixin close return: %#v", r)
		t.FailNow()
	}
}
