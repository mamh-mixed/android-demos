package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/util"
)

func TestPay(t *testing.T) {
	d := &PayReq{
		CommonParams: testCommonParams,

		DeviceInfo:     "xxx",                // 设备号
		Body:           "test",               // 商品描述
		Detail:         "xxx",                // 商品详情
		Attach:         "xxx",                // 附加数据
		OutTradeNo:     util.SerialNumber(),  // 商户订单号
		TotalFee:       "2",                  // 总金额
		FeeType:        "CNY",                // 货币类型
		SpbillCreateIP: util.LocalIP,         // 终端IP
		GoodsGag:       "xxx",                // 商品标记
		AuthCode:       "130755126399220600", // 授权码
	}

	r := &PayResp{}
	err := base(d, r)
	if err != nil {
		t.Errorf("weixin scan pay error: %s", err)
	}
	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin scanpay return: %#v", r)
	}
}
