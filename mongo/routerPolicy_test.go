package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

// func TestPaginationFindRouterPolicy(t *testing.T) {
// 	merId := ""
// 	page, size := 1, 10
//
// 	result, total, err := RouterPolicyColl.PaginationFind(merId, size, page)
// 	if err != nil {
// 		t.Errorf("fail %s", err)
// 	}
//
// 	t.Logf("total is %d; collections are %#v", total, result)
//
// 	t.Logf("count is %d", len(result))
// }

func TestInsertRouterPolicy(t *testing.T) {

	if debug {
		rp := &model.RouterPolicy{
			MerId:     "testMerId",
			CardBrand: "CUP",
			ChanCode:  "testChanCode2",
			ChanMerId: "testChanMerId2",
		}

		if err := RouterPolicyColl.Insert(rp); err != nil {
			t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
		}
	}
}

func TestFindRouterPolicy(t *testing.T) {
	rp := RouterPolicyColl.Find(merId, cardBrand)
	if rp == nil {
		t.Error("Excepted one but get 'nil'")
	}
	t.Logf("RouterPolicy is: %+v", rp)
}

func TestFindAllOfOneMerchant(t *testing.T) {
	r, err := RouterPolicyColl.FindAllOfOneMerchant("001405")

	if err != nil {
		t.Errorf("Error:%s", err)
	}
	t.Logf("Result is %+v", r)
}
