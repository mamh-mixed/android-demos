package scanpay

import (
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/util"
)

func TestPay(t *testing.T) {

	// 设置失效时间
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)

	d := &PayReq{
		CommonParams: testCommonParams,

		DeviceInfo:     "xxx",                // 设备号
		Body:           "test",               // 商品描述
		Detail:         "xxx",                // 商品详情
		Attach:         "xxx",                // 附加数据
		OutTradeNo:     util.SerialNumber(),  // 商户订单号
		TotalFee:       "300",                // 总金额
		FeeType:        "CNY",                // 货币类型
		SpbillCreateIP: util.LocalIP,         // 终端IP
		GoodsGag:       "xxx",                // 商品标记
		AuthCode:       "130470441880647678", // 授权码

		TimeStart:  startTime.Format("20060102150405"), // 交易起始时间
		TimeExpire: endTime.Format("20060102150405"),   // 交易结束时间
	}

	r := &PayResp{}
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
