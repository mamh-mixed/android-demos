package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/util"
)

func TestRefund(t *testing.T) {
	d := &RefundReq{
		CommonParams: testCommonParams,

		DeviceInfo:    "",                  // 设备号
		TransactionId: "",                  // 微信的订单号，优先使用
		OutTradeNo:    "1438137518988",     // 商户系统内部的订单号，当没提供transaction_id时需要传这个
		OutRefundNo:   util.SerialNumber(), // 商户退款单号
		TotalFee:      "1",                 // 总金额
		RefundFee:     "1",                 // 退款金额
		RefundFeeType: "",                  // 货币种类
		OpUserId:      "migo",              // 操作员
	}

	r := &RefundResp{}
	err := base(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
	}
}
