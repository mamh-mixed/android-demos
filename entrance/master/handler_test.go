package master

import (
	"testing"
)

// BindingCreateHandle 建立绑定关系
func TestAllMerchant(t *testing.T) {
	result := AllMerchant()

	if result == nil {
		t.Error("出错啦")
	}
	t.Logf("所有的商户:%+v", result)
}
