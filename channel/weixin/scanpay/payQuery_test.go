package scanpay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestPayQuery(t *testing.T) {
	d := &PayQueryReq{
		CommonParams: testCommonParams,

		TransactionId: "",                                 // 微信的订单号，优先使用
		OutTradeNo:    "df1b5161b785431942031e0c93ebe7ba", // 商户系统内部的订单号，当没提供transaction_id时需要传这个
		// Req:           &model.ScanPayRequest{},
	}
	d.CommonParams.Req = &model.ScanPayRequest{}

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
