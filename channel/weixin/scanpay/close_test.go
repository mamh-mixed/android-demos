package scanpay

import (
	"encoding/xml"

	"testing"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/validator"
)

func TestClose(t *testing.T) {
	// TODO 需要补充单元测试
	d := &CloseReq{
		CommonParams: testCommonParams,

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
		CommonParams: testCommonParams,

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
		CommonParams: testCommonParams,

		TransactionId: "",
		OutTradeNo:    util.Millisecond(),
	}

	if err := validator.Validate(d); err != nil {
		t.Error(err)
	}

}
