package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestInsertRouterPolicy(t *testing.T) {
	rp := &model.RouterPolicy{
		MerId:     "211111001405",
		CardBrand: "CUP",
		ChanCode:  "CFCA",
		ChanMerId: "001405",
	}

	if err := RouterPolicyColl.Insert(rp); err != nil {
		t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
	}
}

func TestFindRouterPolicy(t *testing.T) {
	rp := RouterPolicyColl.Find("111111001405", "CUP")
	if rp == nil {
		t.Error("Excepted one but get 'nil'")
	}
	t.Logf("RouterPolicy is: %+v", rp)
}
