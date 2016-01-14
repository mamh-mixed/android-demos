package scanpay2

import (
	"testing"

	"github.com/CardInfoLink/quickpay/util"
)

func TestPay(t *testing.T) {
	num := util.SerialNumber()
	req := &PayReq{
		CommonParams: CommonParams{
			AppID:      "2014122500021754",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo:  num,
		Scene:       "bar_code",
		AuthCode:    "286507303116520037",
		Subject:     "2024-" + num,
		TotalAmount: "0.01",
	}

	resp := &PayResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	// if resp.Code != "10000" {
	// 	t.Errorf("pay failed")
	// 	t.FailNow()
	// }

	t.Logf("%+v", resp)
}
