package data

import (
	// "fmt"
	"testing"
)

func TestAddRespCodeFromCsv(t *testing.T) {
	// step 1
	// err := AddSysCodeFromCsv("quickpay.csv")
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }

	// step 2
	// err := AddChanCodeFromScv("cfca", "cfca.csv")
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }

	// step 3
	err := AddChanCodeFromScv("cil", "cil.csv")
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
