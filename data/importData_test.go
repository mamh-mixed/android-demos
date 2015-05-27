// 导入数据时注意数据库，若要更新生产环境的数据，请先在develop，testing环境测试。
package data

import (
	"testing"

	// "github.com/omigo/log"
)

// 导入系统应答码，存在时跳过，不存在插入
func xTestAddRespCodeFromCsv(t *testing.T) {
	// 插入quickpay的应答码
	err := AddSysCodeFromCsv("respCode_quickpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("系统应答码导入成功")
}

// 导入渠道应答码
// 导入前先确定系统应答码表是否已存在
// 存在则更新，但不会删除
func xTestAddChanCodeFromScv(t *testing.T) {

	// 插入中金和quickpay的应答码转换数据
	err := AddChanCodeFromScv("cfca", "respCode_cfca.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("中金应答码转换数据插入完成")

	// 插入线下网关和quickpay的应答码转换数据
	err = AddChanCodeFromScv("cil", "respCode_cil.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("线下网关应答码转换数据插入完成")

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
