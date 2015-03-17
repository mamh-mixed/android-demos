package mongo

import (
	"quickpay/model"
	"quickpay/tools"
	"testing"
)

func TestBindingInfo(t *testing.T) {
	merId, bindingId := "MI"+tools.Millisecond(), tools.SerialNumber()
	bi := &model.BindingInfo{
		MerId:     merId,
		BindingId: bindingId,
		CardBrand: "CUP",
		AcctType:  "10",
		AcctName:  "Arthur",
		AcctNum:   "6222020302062061908",
		BankId:    "102",
		IdentType: "0",
		IdentNum:  "350583199009153732",
		PhoneNum:  "18205960039",
		ValidDate: "",
		Cvv2:      "",
	}
	if err := InsertBindingInfo(bi); err != nil {
		t.Errorf("insert bindinginfo error: (%s)", err.Error())
	}

	bi.AcctName = "WonSikin"

	if err := UpdateBindingInfo(bi); err != nil {
		t.Errorf("update bindinginfo error: (%s)", err.Error())
	}

}

func TestBindingMap(t *testing.T) {
	merId, bindingId, chanBindingId := "MI12345", "BI"+tools.Millisecond(), tools.SerialNumber()
	bm := &model.BindingMap{
		MerId:         merId,
		BindingId:     bindingId,
		ChanCode:      "CFCA",
		ChanMerId:     "001045",
		ChanBindingId: chanBindingId,
		BindingStatus: "000009",
	}
	if err := InsertBindingMap(bm); err != nil {
		t.Errorf("insert bindinginfo error: (%s)", err.Error())
	}

	bm.BindingStatus = "000000"

	if err := UpdateBindingMap(bm); err != nil {
		t.Errorf("update bindinginfo error: (%s)", err.Error())
	}

}
