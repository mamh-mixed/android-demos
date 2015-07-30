package scanpay

import "github.com/CardInfoLink/quickpay/util"

var testCommonParams = CommonParams{
	Appid:    "wx25ac886b6dac7dd2", // 公众账号ID
	MchID:    "1236593202",         // 商户号
	SubMchId: "1247075201",         // 子商户号（文档没有该字段）
	NonceStr: util.Nonce(32),       // 随机字符串
	Sign:     "",                   // 签名

	WeixinMD5Key: "12sdffjjguddddd2widousldadi9o0i1",
}
