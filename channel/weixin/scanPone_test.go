package weixin

import (
	"fmt"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

const (
	md5Key     = "12sdffjjguddddd2widousldadi9o0i1"
	mch_id     = "1236593202"
	appid      = "wx25ac886b6dac7dd2"
	acqfee     = "0.02"
	merfee     = "0.03"
	fee        = "0.01"
	sub_mch_id = "1247075201"
	url        = "https://api.mch.weixin.qq.com/pay/micropay"
)

func TestProcessBarcodePay(t *testing.T) {
	sp := &model.ScanPay{
		Txndir:      "",
		Busicd:      "",
		Inscd:       "",
		Chcd:        "",
		Sign:        "",
		SysOrderNum: "",

		Mchntid:      mch_id,
		Txamt:        "1",
		GoodsInfo:    "iphone 7s",
		OrderNum:     "111222333444555",
		OrigOrderNum: "111111111111111",
		ScanCodeId:   "130549657517409996", // AuthCode
		NotifyUrl:    url,
		Key:          md5Key,
		Subject:      "iphone", // Body

	}
	fmt.Println(DefaultClient.ProcessBarcodePay(sp))

}
