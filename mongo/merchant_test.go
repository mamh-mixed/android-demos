package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

func TestMerchantFind(t *testing.T) {
	_, err := MerchantColl.Find("012345678901234")
	_, err = MerchantColl.Find("012345678901234")
	if err != nil {
		t.Error(err)
	}
}

func TestInsertFindUpdateMerchant(t *testing.T) {
	m := &model.Merchant{
		MerId:      "CIL0001",
		MerStatus:  merStatus,
		TransCurr:  transCurr,
		IsNeedSign: false,
		SignKey:    "signKey",
		EncryptKey: "encryptKey",
		Permission: []string{"purc", "paut", "inqy", "refd", "void"},
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

func TestFindAllMerchant(t *testing.T) {
	cond := &model.Merchant{
		MerStatus: "Test",
	}
	result, err := MerchantColl.FindAllMerchant(cond)

	if err != nil {
		t.Errorf("失败了%s", err)
	} else {
		t.Logf("查找结果是： %+v", result)
		t.Logf("长度： %d", len(result))
	}

}

func TestPaginationFind(t *testing.T) {
	merId, merStatus := "", ""
	page, size := 1, 10

	results, total, err := MerchantColl.PaginationFind(merId, merStatus, size, page)
	if err != nil {
		t.Errorf("fail %s", err)
	}

	t.Logf("total is %d; collections are %#v", total, results)

	t.Logf("collections' length is %d", len(results))
}
