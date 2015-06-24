package scanpay

import (
	"encoding/xml"
	"testing"

	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

func TestScanPayQueryGenSign(t *testing.T) {
	d := &ScanPayQueryReqData{
		Appid:         "wx25ac886b6dac7dd2", // 公众账号ID
		MchID:         "1236593202",         // 商户号
		SubMchId:      "1247075201",
		TransactionId: "1010070115201506230291458545", // 微信支付订单号
		OutTradeNo:    "",                             // 商户订单号
		NonceStr:      tools.Nonce(32),                // 商品详情
		Sign:          "",
		WeixinMD5Key:  "12sdffjjguddddd2widousldadi9o0i1",
	}

	d.GenSign()

	t.Log(d.Sign)

	xmlBytes, err := xml.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Errorf("struct(%#v) to xml error: %s", d, err)
	}

	t.Log(string(xmlBytes))
}
