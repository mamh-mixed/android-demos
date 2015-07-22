package scanpay

import "testing"

func TestPrecreate(t *testing.T) {
	req := &PrecreateReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo:  "14141341234",
		Subject:     "2024-14141341234",
		TotalAmount: "0.01",
	}

	resp := &PrecreateResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	t.Logf("%+v", resp)
}
