package scanpay

import (
	"encoding/xml"
	"testing"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"github.com/omigo/validator"
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

func TestScanPayGenSign(t *testing.T) {
	d := &PayReq{
		CommonParams: testCommonParams,

		Body:           "micropay test",      // 商品描述
		Detail:         "",                   // 商品详情
		Attach:         "attach data",        // 附加数据
		OutTradeNo:     "1415757673",         // 商户订单号
		TotalFee:       "1",                  // 总金额
		FeeType:        "CNY",                // 货币类型
		SpbillCreateIP: "14.17.22.52",        // 终端IP
		GoodsGag:       "",                   // 商品标记
		AuthCode:       "120269300684844649", // 授权码
	}

	d.GenSign()

	t.Log(d.Sign)

	xmlBytes, err := xml.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Errorf("struct(%#v) to xml error: %s", d, err)
	}

	t.Log(string(xmlBytes))
}

func TestValidateScanPayReqData(t *testing.T) {
	d := &PayReq{
		CommonParams: testCommonParams,

		DeviceInfo:     "1000",               // 设备号
		Body:           "被扫支付测试",             // 商品描述
		Detail:         "",                   // 商品详情
		Attach:         "订单额外描述",             // 附加数据
		OutTradeNo:     "1415757673",         // 商户订单号
		TotalFee:       "1",                  // 总金额
		FeeType:        "CNY",                // 货币类型
		SpbillCreateIP: "14.17.22.52",        // 终端IP
		GoodsGag:       "",                   // 商品标记
		AuthCode:       "120269300684844649", // 授权码
	}

	if err := validator.Validate(d); err != nil {
		t.Error(err)
	}
}
