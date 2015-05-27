package data

import (
	// "fmt"
	"testing"
)

func xTestAddRespCodeFromCsv(t *testing.T) {
	// step 1
	err := AddSysCodeFromCsv("respCode_quickpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// step 2
	err = AddChanCodeFromScv("cfca", "respCode_cfca.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// step 3
	err = AddChanCodeFromScv("cil", "respCode_cil.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
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
