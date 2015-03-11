package mongo

import (
	"testing"
)

func TestFindRouterPolicyByMerCodeAndCardBrand(t *testing.T) {
	rp, err := FindRouterPolicyByMerCodeAndCardBrand("1000000001", "CUP")
	if err != nil {
		t.Errorf("Excepted an error,but get nil")
	} else {
		t.Logf("error message is %+v", err)
		t.Logf("router policy is %+v", rp)
	}
}

func TestInsertOneRouterPolicy(t *testing.T) {
	rp := &RouterPolicy{
		OrigMerCode:    "1000000001",
		CardBrand:      "CUP",
		ChannelCode:    "CFCA",
		ChannelMerCode: "20000000002",
	}

	if err := InsertOneRouterPolicy(rp); err != nil {
		t.Errorf("Excepted no erro,but get one,error message is %s", err.Error())
	}
}
