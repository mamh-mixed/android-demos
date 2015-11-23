package currency

import (
	"fmt"
)

var CurMap map[string]int

func init() {
	initAvailableCurrency()
}

func initAvailableCurrency() {
	// 币种、精度
	CurMap = make(map[string]int)
	CurMap["GBP"] = 2
	CurMap["HKD"] = 2
	CurMap["USD"] = 2
	CurMap["CHF"] = 2
	CurMap["SGD"] = 2
	CurMap["SEK"] = 2
	CurMap["DKK"] = 2
	CurMap["NOK"] = 2
	CurMap["JPY"] = 0
	CurMap["CAD"] = 2
	CurMap["AUD"] = 2
	CurMap["EUR"] = 2
	CurMap["KRW"] = 0
}

func TransferAmt(transAmt int64, cur string) float64 {
	return float64(transAmt)
}
