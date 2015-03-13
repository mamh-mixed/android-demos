package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertRouterPolicy(t *testing.T) {
	rp := &model.RouterPolicy{
		MerId:     "001405",
		CardBrand: "VIS",
		ChanCode:  "CFCA",
		ChanMerId: "10000000007",
	}

	if err := InsertRouterPolicy(rp); err != nil {
		t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
	}
}

func TestFindRouterPolicy(t *testing.T) {
	rp := FindRouterPolicy("001405", "VIS")
	if rp == nil {
		t.Error("Excepted one but get 'nil'")
	}
	t.Logf("RouterPolicy is: %+v", rp)
}
