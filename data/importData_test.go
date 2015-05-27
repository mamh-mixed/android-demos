package data

import (
	"testing"
)

// 导入系统应答码，存在时跳过，不存在插入
func xTestAddRespCodeFromCsv(t *testing.T) {
	// respCode
	err := AddSysCodeFromCsv("respCode_quickpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 导入渠道应答码
// 导入前先确定系统应答码表是否已存在
// 存在则更新，但不会删除
func xTestAddChanCodeFromScv(t *testing.T) {
	// cfca
	err := AddChanCodeFromScv("cfca", "respCode_cfca.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// cil
	err = AddChanCodeFromScv("cil", "respCode_cil.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 导入卡bin表，执行单元测试时去掉x
// false时会更新数据，但不会删除
// true时会丢掉集合，重新建立
func xTestAddCardBinFromCsv(t *testing.T) {

	// import cardBin !!!
	err := AddCardBinFromCsv("cardBin.csv", false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
