package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertBindingRelation(t *testing.T) {
	ci := model.BindingCreate{
		BindingId: "12345678901",
		AcctName:  "张三",
		AcctNum:   "6210948000000219",
		IdentType: "0",
		IdentNum:  "36050219880401",
		PhoneNum:  "15600009909",
		AcctType:  "20",
		ValidDate: "1903",
		Cvv2:      "232",
	}
	rp := model.RouterPolicy{
		OrigMerId: "001405",
		CardBrand: "CUP",
		ChanCode:  "CFCA",
		ChanMerId: "20000000002",
	}
	br := &model.BindingRelation{ci, rp, ""}

	if err := InsertBindingRelation(br); err != nil {
		t.Errorf("InsertBindingRelation error,except 'nil',but get '%s'", err.Error())
	}
}

func TestFindBindingRelationByMerCodeAndBindingId(t *testing.T) {
	br, err := FindBindingRelation("001405", "1000000000007")
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
		model.BindingCreate{
			BindingId: "1000000000007",
			AcctName:  "张三",
			AcctNum:   "6210948000000219",
			IdentType: "0",
			IdentNum:  "36050219880401",
			PhoneNum:  "15600009909",
			AcctType:  "20",
			ValidDate: "1903",
			Cvv2:      "232",
		},
		model.RouterPolicy{
			OrigMerId: "001405",
			CardBrand: "CUP",
			ChanCode:  "CFCA",
			ChanMerId: "20000000002",
		},
		"12345678901",
	}

	err := UpdateBindingRelation(br)
	if err != nil {
		t.Errorf("Excepted 'nil',but get  error: '%s'", err.Error())
	}
}

func TestDeleteBindingRelation(t *testing.T) {
	br := &model.BindingRelation{
		model.BindingCreate{
			BindingId: "1000000000007",
			AcctName:  "张三",
			AcctNum:   "6210948000000219",
			IdentType: "0",
			IdentNum:  "36050219880401",
			PhoneNum:  "15600009909",
			AcctType:  "20",
			ValidDate: "1903",
			Cvv2:      "232",
		},
		model.RouterPolicy{
			OrigMerId: "001405",
			CardBrand: "CUP",
			ChanCode:  "CFCA",
			ChanMerId: "20000000002",
		},
		"",
	}

	err := DeleteBindingRelation(br)

	if err != nil {
		t.Errorf("Excepted 'nil',but get one error: '%s'", err.Error())
	}
}
