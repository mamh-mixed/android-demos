package data

import (
	"testing"
)

func TestReadCsv(t *testing.T) {

	AddSysCodeFromCsv("quickpay.csv")
	AddChanCodeFromScv("cfca", "cfca.csv")
	AddChanCodeFromScv("cil", "cil.csv")
}
