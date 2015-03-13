package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertBindingRelation(t *testing.T) {
	br := &model.BindingRelation{
		BindingId:     "2000000000001",
		MerId:         "001405",
		AcctName:      "张三",
		AcctNum:       "6222020302062061908",
		IdentType:     "0",
		IdentNum:      "350583199009153732",
		PhoneNum:      "18205960039",
		AcctType:      "20",
		ValidDate:     "1903",
		Cvv2:          "232",
		BankId:        "102",
		CardBrand:     "CUP",
		ChanCode:      "CFCA",
		ChanMerId:     "001405",
		SysBindingId:  "12345678901",
		BindingStatus: "",
	}

	if err := InsertBindingRelation(br); err != nil {
		t.Errorf("InsertBindingRelation error,except 'nil',but get '%s'", err.Error())
	}
}

func TestFindBindingRelationByMerCodeAndBindingId(t *testing.T) {
	br, err := FindBindingRelation("001405", "2000000000001")
	if err != nil {
		t.Errorf("Excepted 'nil',but get  error: '%s'", err.Error())
	}

	if br == nil {
		t.Errorf("Excepted  BindingRelation,but get nil")
	}

	t.Logf("BindingRelation:%+v", br)
}

func TestUpdateBindingRelation(t *testing.T) {
	br := &model.BindingRelation{
		BindingId:     "2000000000001",
		MerId:         "001405",
		AcctName:      "张三",
		AcctNum:       "6222020302062061908",
		IdentType:     "0",
		IdentNum:      "350583199009153732",
		PhoneNum:      "18205960039",
		AcctType:      "20",
		ValidDate:     "1903",
		Cvv2:          "232",
		BankId:        "102",
		CardBrand:     "CUP",
		ChanCode:      "CFCA",
		ChanMerId:     "001405",
		SysBindingId:  "12345678901",
		BindingStatus: "000000",
	}

	err := UpdateBindingRelation(br)
	if err != nil {
		t.Errorf("Excepted 'nil',but get  error: '%s'", err.Error())
	}
}

func TestDeleteBindingRelation(t *testing.T) {
	br := &model.BindingRelation{
		BindingId:     "2000000000001",
		MerId:         "001405",
		AcctName:      "张三",
		AcctNum:       "6222020302062061908",
		IdentType:     "0",
		IdentNum:      "350583199009153732",
		PhoneNum:      "18205960039",
		AcctType:      "20",
		ValidDate:     "1903",
		Cvv2:          "232",
		BankId:        "102",
		CardBrand:     "CUP",
		ChanCode:      "CFCA",
		ChanMerId:     "001405",
		SysBindingId:  "12345678901",
		BindingStatus: "000000",
	}

	err := DeleteBindingRelation(br)

	if err != nil {
		t.Errorf("Excepted 'nil',but get one error: '%s'", err.Error())
	}
}
