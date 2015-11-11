package scanpay2

import (
	"testing"

	"github.com/dawei888/quickpay/util"
)

func TestPay(t *testing.T) {
	num := util.SerialNumber()
	req := &PayReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo:  num,
		Scene:       "bar_code",
		AuthCode:    "282528673824835509",
		Subject:     "2024-" + num,
		TotalAmount: "0.01",
	}

	resp := &PayResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	if resp.Code != "10000" {
		t.Errorf("pay failed")
		t.FailNow()
	}

	t.Logf("%+v", resp)
}
