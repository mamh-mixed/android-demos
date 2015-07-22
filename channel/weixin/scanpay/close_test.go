package scanpay

import (
	"encoding/xml"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/validator"
	"testing"
)

func TestClose(t *testing.T) {
	// TODO 需要补充单元测试
	d := &CloseReq{
		Appid:         "wx25ac886b6dac7dd2",
		MchID:         "1236593202",
		SubMchId:      "1247075201",
		NonceStr:      util.Nonce(32),
		Sign:          "C380BEC2BFD727A4B6845133519F3AD6",
		WeixinMD5Key:  "12sdffjjguddddd2widousldadi9o0i1",
		TransactionId: "",
		OutTradeNo:    util.Millisecond(),
	}

	r := &PayResp{}

	err := base(d, r)
	if err != nil {
		t.Errorf("weixin close error: %s", err)
	}

	if r.ReturnCode != "SUCCESS" {
		t.Logf("weixin close return: %#v", r)
	}
}

func TestCloseGenSign(t *testing.T) {
	d := &CloseReq{
		Appid:         "wx2421b1c4370ec43b",
		MchID:         "10000100",
		SubMchId:      "1247075201",
		NonceStr:      util.Nonce(32),
		Sign:          "C380BEC2BFD727A4B6845133519F3AD6",
		WeixinMD5Key:  "0123435657",
		TransactionId: "",
		OutTradeNo:    util.Millisecond(),
	}

	d.GenSign()

	t.Log(d.Sign)

	xmlBytes, err := xml.MarshalIndent(d, "", "\t")
	if err != nil {
		t.Logf("struct(%#v) to xml error: %s", d, err)
	}

	t.Log(string(xmlBytes))
}

func TestValidateCloseReqData(t *testing.T) {
	d := &CloseReq{
		Appid:         "wx2421b1c4370ec43b",
		MchID:         "10000100",
		SubMchId:      "1247075201",
		NonceStr:      util.Nonce(32),
		Sign:          "C380BEC2BFD727A4B6845133519F3AD6",
		WeixinMD5Key:  "0123435657",
		TransactionId: "",
		OutTradeNo:    util.Millisecond(),
	}

	if err := validator.Validate(d); err != nil {
		t.Error(err)
	}

}
