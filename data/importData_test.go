// 导入数据时注意数据库，若要更新生产环境的数据，请先在develop，testing环境测试。
package data

import (
	"testing"

	// "github.com/omigo/log"
)

func TestAddScanPayRespFromCSV(t *testing.T) {

	data, err := readQuickpayCSV("respCode_scanpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(len(data))
}

// 导入系统应答码，存在时跳过，不存在插入
func xTestAddRespCodeFromCSV(t *testing.T) {
	// 插入quickpay的应答码
	err := AddSysCodeFromCSV("respCode_quickpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("系统应答码导入成功")
}

// 导入渠道应答码
// 导入前先确定系统应答码表是否已存在
// 存在则更新，但不会删除
func xTestAddChanCodeFromCSV(t *testing.T) {
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

// TestAddCardBinFromCSV 导入卡 bin 表，执行单元测试时去掉 x
func TestAddCardBinFromCSV(t *testing.T) {
	// false 时会更新数据，但不会删除
	// true 时会丢掉集合，重新建立
	err := AddCardBinFromCSV("cardBin.csv", true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("卡 BIN 数据插入完成")
}
