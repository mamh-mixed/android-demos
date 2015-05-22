package core

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

func TestUnionPayCardValidity(t *testing.T) {
	bc := &model.BindingCreate{
		BindingId:       "1000000000001",
		AcctName:        "张三",
		AcctNum:         "6210948000000219",
		IdentType:       "0",
		IdentNum:        "440583199111031012",
		IdentNumDecrypt: "44058319911103101X",
		PhoneNum:        "15600009909",
		PhoneNumDecrypt: "15600009909",
		AcctType:        "20",
		ValidDate:       "1903",
		Cvv2:            "232",
		SendSmsId:       "1000000000009",
		SmsCode:         "12353",
	}

	ret := UnionPayCardCommonValidity(bc.IdentType, bc.IdentNumDecrypt, bc.PhoneNumDecrypt)
	if ret != nil {
		t.Errorf("Excepted 'nil',but get code is '%s' and msg is '%s' ", ret.RespCode, ret.RespMsg)
	}
}
