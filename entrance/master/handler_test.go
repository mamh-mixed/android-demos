package master

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

// BindingCreateHandle 建立绑定关系
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
