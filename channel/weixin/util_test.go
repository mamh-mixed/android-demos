package weixin

import (
	"fmt"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

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

func TestToMap(t *testing.T) {

	m := toMapWithKeySortedAndValueNotNil(microPay)
	fmt.Println("m: ", m)

	if len(m) != 9 {
		fmt.Println("lenth of map: ", len(m))
		t.FailNow()
	}

}

func TestCalculateSign(t *testing.T) {

	microPay.Sign = CalculateSign(microPay, md5Key)
	if microPay.Sign == "" {
		t.FailNow()
	}
}

func TestPerpareRequestStruct(t *testing.T) {
	req1 := PerpareRequestStruct(&model.ScanPay{})

	if req1.Sign == "" {
		t.FailNow()
	}

}
