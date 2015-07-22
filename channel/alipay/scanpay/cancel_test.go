package scanpay

import "testing"

func TestCancel(t *testing.T) {
	req := &CancelReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo: "14141341234",
	}

	resp := &CancelResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	t.Logf("%+v", resp)
}
