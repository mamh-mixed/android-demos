package mongo

import (
	"quickpay/model"
	"testing"
)

func TestInsertRouterPolicy(t *testing.T) {
	rp := &model.RouterPolicy{
		OrigMerId: "001405",
		CardBrand: "VIS",
		ChanCode:  "CFCA",
		ChanMerId: "20000000002",
	}

	if err := InsertRouterPolicy(rp); err != nil {
		t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
	}
}
