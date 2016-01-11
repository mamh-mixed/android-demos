package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
)

func TestRefundQuery(t *testing.T) {
	spReq := &model.ScanPayRequest{}

	d := &RefundQueryReq{
		CommonParams: testCommonParams,
		OutTradeNo:   "1452477096056", // 商户订单号
		// OutRefundNo:  "1452479745377", // 商户退款单号
	}
	d.CommonParams.Req = spReq

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
