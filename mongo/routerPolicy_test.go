package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertRouterPolicy(t *testing.T) {
	rp := &model.RouterPolicy{
		MerId:     "001405",
		CardBrand: "UPI",
		ChanCode:  "CFCA",
		ChanMerId: "001405",
	}

	if err := InsertRouterPolicy(rp); err != nil {
		t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
	}
}

func TestFindRouterPolicy(t *testing.T) {
	rp := FindRouterPolicy("001405", "CUP")
	if rp == nil {
		t.Error("Excepted one but get 'nil'")
	}
	t.Logf("RouterPolicy is: %+v", rp)
}
