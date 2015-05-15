package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

func TestInsertRouterPolicy(t *testing.T) {

	if debug {
		rp := &model.RouterPolicy{
			MerId:     merId,
			CardBrand: cardBrand,
			ChanCode:  chanCode,
			ChanMerId: chanMerId,
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
