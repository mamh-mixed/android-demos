package scanpay2

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
)

func TestPrecreate(t *testing.T) {
	num := util.SerialNumber()
	req := &PrecreateReq{
		CommonParams: CommonParams{
			AppID:      "2014122500021754",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
			Req:        &model.ScanPayRequest{},
		},
		OutTradeNo:  num,
		Subject:     "2024-" + num,
		TotalAmount: "0.01",
	}

	resp := &PrecreateResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	if resp.Code != "10000" {
		t.Errorf("precreate failed")
		t.FailNow()
	}

	t.Logf("%+v", resp)
}
