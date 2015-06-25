package scanpay

import (
	"encoding/xml"
	"testing"

	"github.com/omigo/log"
	"github.com/omigo/validator"
)

func TestScanPayGenSign(t *testing.T) {
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
