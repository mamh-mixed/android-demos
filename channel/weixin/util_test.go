package weixin

var microPay = &MicropayRequest{
	AppId:          "sdfsd",
	MchId:          "werwer1231",
	NonceStr:       "xxxxoooo",
	TotalFee:       12,
	OutTradeNo:     "fdsfsfdsf",
	FeeType:        "CNY",
	SpbillCreateIp: "10.10.10.1",
	Body:           "sdfdfsdds",
	AuthCode:       "sfdsfafd",
}

/*
func TestToMap(t *testing.T) {

	m := toMapWithValueNotNil(microPay)
	fmt.Println("m: ", m)

	if len(m) != 9 {
		fmt.Println("lenth of map: ", len(m))
		t.FailNow()
	}
}

func TestCalculateSign(t *testing.T) {
	setSign(microPay, md5Key)
	if microPay.Sign == "" {
		t.FailNow()
	}
}
*/
