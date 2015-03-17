package core

import (
	"quickpay/model"
	"testing"
)

func TestProcessBindingCreate(t *testing.T) {
	bc := &model.BindingCreate{
		MerId:     "001405",
		BindingId: "1000000000011",
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

func TestProcessBindingEnquiry(t *testing.T) {
	be := &model.BindingEnquiry{
		MerId:     "001405",
		BindingId: "1000000000011",
	}

	// br, err := mongo.FindBindingRelation(be.MerId, be.BindingId)
	//
	// if err != nil {
	// 	t.Errorf("'FindBindingRelation' error: %s", err.Error())
	// }
	//
	// br.BindingStatus = "000009"
	// if err = mongo.UpdateBindingRelation(br); err != nil {
	// 	t.Errorf("'UpdateBindingRelation' error: %s", err.Error())
	// }

	ret := ProcessBindingEnquiry(be)

	t.Logf("%+v", ret)
}

func TestProcessBindingRemove(t *testing.T) {
	be := &model.BindingRemove{
		MerId:     "1426562901844",
		BindingId: "1426562901897",
	}

	ret := ProcessBindingReomve(be)

	t.Logf("%+v", ret)
}
