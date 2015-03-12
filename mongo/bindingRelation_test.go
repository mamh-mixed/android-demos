package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertBindingRelation(t *testing.T) {
	br := &BindingRelation{
		CardInfo: model.BindingCreate{
			BindingId: "12345678901",
			AcctName:  "张三",
			AcctNum:   "6210948000000219",
			IdentType: "0",
			IdentNum:  "36050219880401",
			PhoneNum:  "15600009909",
			AcctType:  "20",
			ValidDate: "1903",
			Cvv2:      "232",
		},
		Router: RouterPolicy{
			OrigMerId: "20000000002",
			CardBrand: "CUP",
			ChanCode:  "CFCA",
			ChanMerId: "001405",
		},
		ChanBindingId: "12345678901",
	}

	if err := InsertBindingRelation(br); err != nil {
		t.Errorf("InsertBindingRelation error,except 'nil',but get '%s'", err.Error())
	}
}

func TestFindBindingRelationByMerCodeAndBindingId(t *testing.T) {
	br, err := FindBindingRelation("001405", "1000000000001")
	if err != nil {
		t.Errorf("Excepted 'nil',but get  error: '%s'", err.Error())
	}

	if br == nil {
		t.Errorf("Excepted  BindingRelation,but get nil")
	}

	t.Logf("BindingRelation:%+v", br)
}

func TestUpdateBindingRelation(t *testing.T) {
	br := &BindingRelation{
		CardInfo: model.BindingCreate{
			BindingId: "12345678901",
			AcctName:  "张三",
			AcctNum:   "6210948000000219",
			IdentType: "0",
			IdentNum:  "36050219880401",
			PhoneNum:  "15600009909",
			AcctType:  "20",
			ValidDate: "1903",
			Cvv2:      "232",
		},
		Router: RouterPolicy{
			OrigMerId: "20000000002",
			CardBrand: "CUP",
			ChanCode:  "CFCA",
			ChanMerId: "001405",
		},
		ChanBindingId: "12345678901",
	}

	err := UpdateBindingRelation(br)
	if err != nil {
		t.Errorf("Excepted 'nil',but get  error: '%s'", err.Error())
	}
}
