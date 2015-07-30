package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/util"
)

func TestPay(t *testing.T) {
	d := &PayReq{
		CommonParams: testCommonParams,

		DeviceInfo:     "",                   // 设备号
		Body:           "product desc",       // 商品描述
		Detail:         "",                   // 商品详情
		Attach:         "",                   // 附加数据
		OutTradeNo:     util.SerialNumber(),  // 商户订单号
		TotalFee:       "3",                  // 总金额
		FeeType:        "",                   // 货币类型
		SpbillCreateIP: util.LocalIP,         // 终端IP
		GoodsGag:       "",                   // 商品标记
		AuthCode:       "130413885648248844", // 授权码
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
