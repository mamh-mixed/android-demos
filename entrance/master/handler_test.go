package master

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestAllMerchant(t *testing.T) {
	cond := model.Merchant{
		MerStatus: "Test",
	}

	data, err := json.Marshal(cond)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	result := AllMerchant(data)

	if result == nil {
		t.Error("出错啦")
	}
	t.Logf("所有的商户:%+v", result)
}

func TestAddMerchant(t *testing.T) {
	cond := model.Merchant{
		MerId:      "TEST001",
		TransCurr:  "156",
		MerStatus:  "Deleted",
		SignKey:    "TEST001",
		EncryptKey: "TEST001",
		Remark:     "TEST001",
	}
	t.Log("============")
	data, err := json.Marshal(cond)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	result := AddMerchant(data)

	if result == nil {
		t.Error("出错啦")
	}
	t.Logf("结果:%+v", result)
}

func TestAddChannelMerchant(t *testing.T) {
	cond := &model.ChanMer{
		ChanCode:      "TEST2",
		ChanMerId:     "TEST",
		ChanMerName:   "TEST",
		SettFlag:      "TEST",
		SettRole:      "TEST",
		SignCert:      "TEST",
		CheckSignCert: "TEST",
	}

	t.Log("============")
	data, err := json.Marshal(cond)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	result := AddChannelMerchant(data)

	if result == nil {
		t.Error("出错啦")
	}
	t.Logf("结果:%+v", result)
}

func TestAllChannelMerchant(t *testing.T) {
	cond := model.ChanMer{}

	t.Log("============")
	data, err := json.Marshal(cond)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	result := AllChannelMerchant(data)

	if result == nil {
		t.Error("出错啦")
	}
	t.Logf("结果:%+v", result)

}

func TestAddRouter(t *testing.T) {
	cond := model.RouterPolicy{
		MerId:     "Test1",
		CardBrand: "Test1",
		ChanCode:  "Test2",
		ChanMerId: "Test2",
	}

	t.Log("============")
	data, err := json.Marshal(cond)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	result := AddRouter(data)

	if result == nil {
		t.Error("出错啦")
	}
	t.Logf("结果:%+v", result)

}
