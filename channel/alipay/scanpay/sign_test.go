package scanpay

import "testing"

func TestVerify(t *testing.T) {
	// resp := `{"code":"40002","msg":"Invalid Arguments","sub_code":"isv.invalid-signature","sub_msg":"无效签名"}`
	// sign := "aUududgICSy+Wr5NppUHX3FedZeQJCZ+L+xYy7q/TwwXLT5PEkF1PGlbG8hohV0kLZO9iFkWt7mhJeozz5GAp9cqVCnvEdiOT1ausq1EdBXIWF4kPb8C2ssoq6OpgBoHgd+cPbpaekbwbyHxWIHvfZopDUCG+G8aTUxJjLrTVsY="
	sign := "vdIUy7rMmbPviBgBJKHvWZDv72cUQwuU2eqhVCEETcX5WWzhPk0BrXKOj5qUMxO7TQuaiKgKKymndkmEXUhraIcOImDm+HX/w2x+SAOjMVPyDSQ6DgUlMWoSqkIF+38OoNgg132thaBBSBpvMkhbbbtZv/w9Aft3qJebi4E/1L4="
	resp := `appId=2015032400038629&biz_content={"out_trade_no":"14141341234"}&charset=utf-8&method=alipay.trade.cancel&sign_type=RSA&timestamp=2015-07-21 15:00:25`

	err := Verify([]byte(resp), sign)
	if err != nil {
		t.Error(err)
	}
}

func TestSign(t *testing.T) {
	// data := `appId=2015032400038629&biz_content={"out_trade_no":"14141341234"}&charset=utf-8&method=alipay.trade.cancel&sign_type=RSA&timestamp=2015-07-21 15:00:25`
	data := `app_id=2015032400038629&biz_content={"out_trade_no":"2015072017250000","total_amount":"0.01","subject":"讯联数据测试"}&charset=utf-8&method=alipay.trade.precreate&sign_type=RSA&timestamp=2015-07-21 16:36:55&version=1.0`
	sign := `rR4jK0oXbJ+HAdhE0siuRK41UA6u1mXG1nc6gbvhz1q4GYyzCBGau/oTTvb6Sc4Y54ZSB1rxpipA4qk45qXo4OJg2ugdbyTK1T0de5noQIzZJBdjyydkGYsYnBv/KZQaiY+D4eQDshBKQ3dH1/uuuFuep83E9CTVIQ3wnTuCY+4=`

	b64Sign, err := Sha1WithRsa([]byte(data), LoadPrivateKey([]byte(privateKeyPem)))
	if err != nil {
		t.Error(err)
	}

	t.Log(b64Sign)

	if b64Sign != sign {
		t.Error("not valid")
	}
}
