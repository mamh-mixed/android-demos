package core

import (
	"quickpay/model"
	"testing"
)

func TestProcessBindingCreate(t *testing.T) {
	bc := &model.BindingCreate{
		MerId:     "001405",
		BindingId: "1000000000006",
		AcctName:  "张三",
		AcctNum:   "6222020302062061908",
		IdentType: "0",
		IdentNum:  "350583199009153732",
		PhoneNum:  "18205960039",
		AcctType:  "10",
		ValidDate: "",
		Cvv2:      "",
		SendSmsId: "1000000000009",
		SmsCode:   "12353",
		BankId:    "102",
	}

	ret := ProcessBindingCreate(bc)

	t.Logf("%+v", ret)

	if ret.RespCode != "000000" {
		t.Errorf("Excpeted success,but get failure: %+v", ret.RespMsg)
	}

}
