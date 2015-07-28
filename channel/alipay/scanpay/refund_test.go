package scanpay

import "testing"

func TestRefund(t *testing.T) {
	req := &RefundReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		TradeNo:      "2015072821001004000056602110",
		RefundAmount: "0.01",
	}

	resp := &RefundResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	if resp.Code != "10000" {
		t.Errorf("refund failed")
		t.FailNow()
	}

	t.Logf("%+v", resp)
}
