package scanpay

import "testing"

func TestPreCreate(t *testing.T) {
	req := NewPreCreateReq("2015051100069108", LoadPrivateKey([]byte(privateKeyPem)))
	req.OutTradeNo = "2015072017250000"
	req.Subject = "讯联数据测试"
	req.TotalAmount = "0.01"

	body := &PreCreateBody{}
	resp := &PreCreateResp{}
	err := base(req, body, resp)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
		t.FailNow()
	}

	t.Logf("%+v", body)
	t.Logf("%+v", resp)
}
