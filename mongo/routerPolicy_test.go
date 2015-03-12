package mongo

import (
	"testing"
)

func TestInsertOneRouterPolicy(t *testing.T) {
	rp := &RouterPolicy{
		OrigMerId: "001405",
		CardBrand: "VIS",
		ChanCode:  "CFCA",
		ChanMerId: "20000000002",
	}

	if err := InsertOneRouterPolicy(rp); err != nil {
		t.Errorf("Excepted no erro,but get one,error message is %s", err.Error())
	}
}

func TestFindRouterPolicyByMerCodeAndCardBrand(t *testing.T) {
	rp, err := FindRouter("001405", "CUP")
	if err != nil {
		t.Errorf("Excepted an error,but get nil")
	} else {
		t.Logf("error message is %+v", err)
		t.Logf("router policy is %+v", rp)
	}
}
