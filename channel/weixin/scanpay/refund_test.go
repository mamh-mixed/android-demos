package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/util"
)

func TestRefundPay(t *testing.T) {
	d := &RefundReq{
		CommonParams: testCommonParams,

		DeviceInfo:    "xxx",               // 设备号
		TransactionId: "",                  // 微信的订单号，优先使用
		OutTradeNo:    "1441615439898",     // 商户系统内部的订单号，当没提供transaction_id时需要传这个
		OutRefundNo:   util.SerialNumber(), // 商户退款单号
		TotalFee:      "300",               // 总金额
		RefundFee:     "300",               // 退款金额
		RefundFeeType: "CNY",               // 货币种类
		OpUserId:      "migo",              // 操作员
	}

	r := &RefundResp{}
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
