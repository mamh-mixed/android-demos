package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/util"
)

func TestRefund(t *testing.T) {
	d := &RefundReq{
		// 公共字段
		Appid:    "wx25ac886b6dac7dd2", // 公众账号ID
		MchID:    "1236593202",         // 商户号
		SubMchId: "1247075201",         // 子商户号（文档没有该字段）
		NonceStr: util.Nonce(32),       // 随机字符串
		Sign:     "",                   // 签名

		WeixinMD5Key: "12sdffjjguddddd2widousldadi9o0i1",

		DeviceInfo:    "",                                 // 设备号
		TransactionId: "",                                 // 微信的订单号，优先使用
		OutTradeNo:    "7a5d8c60e1284fe8697af775c60d15d7", // 商户系统内部的订单号，当没提供transaction_id时需要传这个
		OutRefundNo:   util.SerialNumber(),                // 商户退款单号
		TotalFee:      "3",                                // 总金额
		RefundFee:     "1",                                // 退款金额
		RefundFeeType: "",                                 // 货币种类
		OpUserId:      "migo",                             // 操作员
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
