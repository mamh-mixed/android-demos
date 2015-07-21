package scanpay

import "testing"

func TestCancel(t *testing.T) {
	req := NewCancelReq("2015051100069108", LoadPrivateKey([]byte(privateKeyPem)))
	req.OutTradeNo = "14141341234"
	req.Subject = "2024-14141341234"

	body := &CancelBody{}
	resp := &CancelResp{}
	err := base(req, body, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	t.Logf("%+v", body)
	t.Logf("%+v", resp)
}
