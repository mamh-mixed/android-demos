package scanpay

import "testing"

func TestQuery(t *testing.T) {
	req := &QueryReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo: "5b1fbd21a9334e68431337f4884bc061",
	}

	resp := &QueryResp{}
	err := Execute(req, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	if resp.Code != "10000" {
		t.Errorf("query failed")
		t.FailNow()
	}

	t.Logf("%+v", resp)
}
