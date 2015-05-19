package data

import (
	// "fmt"
	"testing"
)

func TestReadCsv(t *testing.T) {
	// AddSysCodeFromCsv("quickpay.csv")
	// AddChanCodeFromScv("cfca", "cfca.csv")
	// AddChanCodeFromScv("cil", "cil.csv")
	// InitTestMer(1, 100, "CUP")
	// err := AddCardBinFromCsv("cardBin.csv", true)
	// fmt.Println(err)
	cbs, err := ReadCardBinCsv("cardBin.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%d", len(cbs))
}
