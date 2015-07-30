package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/util"
)

func TestPrePay(t *testing.T) {
	d := &PrePayReq{
		CommonParams: testCommonParams,

		DeviceInfo:     "xxx",               // 设备号
		Body:           "test test",         // 商品描述
		Detail:         "xxx",               // 商品详情
		Attach:         "xxx",               // 附加数据
		OutTradeNo:     util.SerialNumber(), // 商户订单号
		TotalFee:       "1",                 // 总金额
		FeeType:        "CNY",               // 货币类型
		SpbillCreateIP: util.LocalIP,        // 终端IP
		TimeStart:      "",                  // 交易起始时间
		TimeExpire:     "",                  // 交易结束时间
		GoodsGag:       "",                  // 商品标记
		NotifyURL:      weixinNotifyURL,     // 通知地址
		TradeType:      "NATIVE",            // 交易类型
		ProductID:      "123123123",         // 商品ID
		Openid:         "",                  // 用户标识
		SubOpenid:      "",                  // 用户标识
	}

	r := &PrePayResp{}
	err := base(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
	}
}
