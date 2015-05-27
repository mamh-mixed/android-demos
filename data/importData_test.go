package data

import (
	// "fmt"
	"testing"

	// "github.com/omigo/log"
)

// xTestAddRespCodeFromCsv 新应答码插入的方法，去掉x，然后执行go test
func xTestAddRespCodeFromCsv(t *testing.T) {
	// 插入quickpay的应答码
	err := AddSysCodeFromCsv("respCode_quickpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("quickpay的应答码插入完成")

	// 插入中金和quickpay的应答码转换数据
	err = AddChanCodeFromScv("cfca", "respCode_cfca.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("中金和quickpay的应答码转换数据插入完成")

	// 插入线下网关和quickpay的应答码转换数据
	err = AddChanCodeFromScv("cil", "respCode_cil.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("线下网关和quickpay的应答码转换数据插入完成")

}

func TestAddSettSchemeCd(t *testing.T) {

	// import settSchemeCd
	// err := AddSettSchemeCdFromCsv("settSchemeCd.csv")
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
}

func TestAddCardBinFromCsv(t *testing.T) {

	// import cardBin !!!
	// err := AddCardBinFromCsv("cardBin.csv", false)
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
}
