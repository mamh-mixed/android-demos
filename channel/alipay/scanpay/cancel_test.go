package scanpay

import "testing"

func TestCancel(t *testing.T) {
	req := &CancelReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo: "42608e0f54f940624a86d4696da83f7d",
	}

	resp := &CancelResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	if resp.Code != "10000" {
		t.Errorf("cancel failed")
		t.FailNow()
	}
	t.Logf("%+v", resp)
}
