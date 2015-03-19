package core

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

func TestUnionPayCardValidity(t *testing.T) {
	bc := &model.BindingCreate{
		BindingId: "1000000000001",
		AcctName:  "张三",
		AcctNum:   "6210948000000219",
		IdentType: "0",
		IdentNum:  "36050219880401",
		PhoneNum:  "15600009909",
		AcctType:  "20",
		ValidDate: "1903",
		Cvv2:      "232",
		SendSmsId: "1000000000009",
		SmsCode:   "12353",
	}

	ret := UnionPayCardValidity(bc)
	if ret != nil {
		t.Errorf("Excepted 'nil',but get code is '%s' and msg is '%s' ", ret.RespCode, ret.RespMsg)
	}
	bc.IdentType = "12"
	ret = UnionPayCardValidity(bc)
	if ret == nil {
		t.Error("Excepted code is '200111' and msg is '证件类型不正确',but get 'nil'")
	}
	t.Logf("code is '%s' and msg is '%s' ", ret.RespCode, ret.RespMsg)

	bc.IdentType = "X"
	bc.PhoneNum = "1234567899"
	ret = UnionPayCardValidity(bc)
	if ret == nil {
		t.Errorf("Excepted code is '200114' and msg is '手机号码格式错误',but get 'nil'")
	}
	t.Logf("code is '%s' and msg is '%s' ", ret.RespCode, ret.RespMsg)
}
