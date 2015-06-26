package scanpay

import (
	"encoding/xml"
	"testing"

	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"github.com/omigo/validator"
)

func TestPay(t *testing.T) {
	d := &PayReq{
		// 公共字段
		Appid:    "wx25ac886b6dac7dd2", // 公众账号ID
		MchID:    "1236593202",         // 商户号
		SubMchId: "1247075201",         // 子商户号（文档没有该字段）
		NonceStr: tools.Nonce(32),      // 随机字符串
		Sign:     "",                   // 签名

		WeixinMD5Key: "12sdffjjguddddd2widousldadi9o0i1",

		DeviceInfo:     "",                   // 设备号
		Body:           "product desc",       // 商品描述
		Detail:         "",                   // 商品详情
		Attach:         "",                   // 附加数据
		OutTradeNo:     tools.SerialNumber(), // 商户订单号
		TotalFee:       "3",                  // 总金额
		FeeType:        "",                   // 货币类型
		SpbillCreateIP: tools.LocalIP,        // 终端IP
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
		Appid:          "wx2421b1c4370ec43b",               // 公众账号ID
		MchID:          "10000100",                         // 商户号
		SubMchId:       "1247075201",                       // 子商户
		DeviceInfo:     "1000",                             // 设备号
		NonceStr:       "8aaee146b1dee7cec9100add9b96cbe2", // 随机字符串
		Sign:           "C380BEC2BFD727A4B6845133519F3AD6", // 签名
		Body:           "micropay test",                    // 商品描述
		Detail:         "",                                 // 商品详情
		Attach:         "attach data",                      // 附加数据
		OutTradeNo:     "1415757673",                       // 商户订单号
		TotalFee:       "1",                                // 总金额
		FeeType:        "CNY",                              // 货币类型
		SpbillCreateIP: "14.17.22.52",                      // 终端IP
		GoodsGag:       "",                                 // 商品标记
		AuthCode:       "120269300684844649",               // 授权码
		WeixinMD5Key:   "0123435657",
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
		Appid:          "wx2421b1c4370ec43b", // 公众账号ID
		MchID:          "10000100",           // 商户号
		SubMchId:       "1247075201",
		DeviceInfo:     "1000",                             // 设备号
		NonceStr:       "8aaee146b1dee7cec9100add9b96cbe2", // 随机字符串
		Sign:           "C380BEC2BFD727A4B6845133519F3AD6", // 签名
		Body:           "被扫支付测试",                           // 商品描述
		Detail:         "",                                 // 商品详情
		Attach:         "订单额外描述",                           // 附加数据
		OutTradeNo:     "1415757673",                       // 商户订单号
		TotalFee:       "1",                                // 总金额
		FeeType:        "CNY",                              // 货币类型
		SpbillCreateIP: "14.17.22.52",                      // 终端IP
		GoodsGag:       "",                                 // 商品标记
		AuthCode:       "120269300684844649",               // 授权码
		WeixinMD5Key:   "0123435657",
	}

	if err := validator.Validate(d); err != nil {
		t.Error(err)
	}
}
