package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestInsertFindUpdateMerchant(t *testing.T) {
	merId := "CIL0001"
	m := &model.Merchant{
		MerId:      merId,
		MerStatus:  "Normal",
		TransCurr:  "156",
		SignKey:    "",
		EncryptKey: "",
	}
	t.Log("Insert------")
	err := MerchantColl.Insert(m)

	if err != nil {
		t.Errorf("Insert merchant error: %s", err)
	}

	t.Log("Find------")
	m, err = MerchantColl.Find(merId)

	if err != nil {
		t.Errorf("Find merchant error: %s", err)
	}

	t.Log("Update------")
	m.MerStatus = "Deleted"
	err = MerchantColl.Update(m)

	if err != nil {
		t.Errorf("Update merchant error: %s", err)
	}

}
