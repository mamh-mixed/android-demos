package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertOneBindingRelation(t *testing.T) {
	br := &BindingRelation{
		CardInfo: model.BindingCreate{
			BindingId: "1000000000001",
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
			OrigMerCode:    "M1000000001",
			CardBrand:      "CUP",
			ChannelCode:    "CFCA",
			ChannelMerCode: "20000000002",
		},
		ChannelBindingId: "",
	}

	if err := InsertOneBindingRelation(br); err != nil {
		t.Errorf("InsertOneBindingRelation error,except 'nil',but get '%s'", err.Error())
	}

	// br := new(BindingRelation)
	// br.BindingId
}

func TestFindOneBindingRelationByMerCodeAndBindingId(t *testing.T) {
	br, err := FindOneBindingRelation("M1000000001", "1000000000001")
	if err != nil {
		t.Errorf("Excepted 'nil',but get one error: '%s'", err.Error())
	}

	if br == nil {
		t.Errorf("Excepted one BindingRelation,but get nil")
	}

	t.Logf("BindingRelation:%+v", br)
}

func TestUpdateOneBindingRelation(t *testing.T) {
	br := &BindingRelation{
		CardInfo: model.BindingCreate{
			BindingId: "1000000000001",
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
			OrigMerCode:    "M1000000001",
			CardBrand:      "CUP",
			ChannelCode:    "CFCA",
			ChannelMerCode: "20000000002",
		},
		ChannelBindingId: "CB11110000011",
	}

	err := UpdateOneBindingRelation(br)
	if err != nil {
		t.Errorf("Excepted 'nil',but get one error: '%s'", err.Error())
	}
}
