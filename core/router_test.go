package core

import "testing"

func TestFindRouterPolicyByMerCodeAndCardBrand(t *testing.T) {
	rp, err := FindRouter("001405", "CUP")
	if err != nil {
		t.Errorf("Excepted an error,but get nil")
	} else {
		t.Logf("error message is %+v", err)
		t.Logf("router policy is %+v", rp)
	}
}
