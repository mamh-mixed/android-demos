package currency

import (
	"fmt"
	"math"
)

var CurMap map[string]Cur

// Cur 币种、精度
type Cur struct {
	Cur       string
	Precision int
}

func init() {
	initAvailableCurrency()
}

func initAvailableCurrency() {
	CurMap = make(map[string]Cur)
	CurMap["GBP"] = Cur{"GBP", 2}
	CurMap["HKD"] = Cur{"HKD", 2}
	CurMap["USD"] = Cur{"USD", 2}
	CurMap["CHF"] = Cur{"CHF", 2}
	CurMap["SGD"] = Cur{"SGD", 2}
	CurMap["SEK"] = Cur{"SEK", 2}
	CurMap["DKK"] = Cur{"DKK", 2}
	CurMap["NOK"] = Cur{"NOK", 2}
	CurMap["JPY"] = Cur{"JPY", 0}
	CurMap["CAD"] = Cur{"CAD", 2}
	CurMap["AUD"] = Cur{"AUD", 2}
	CurMap["EUR"] = Cur{"EUR", 2}
	CurMap["KRW"] = Cur{"KRW", 0}
	CurMap["CNY"] = Cur{"CNY", 2}
}

// Get 得到一种币种
func Get(currency string) Cur {
	if cur, ok := CurMap[currency]; ok {
		return cur
	}
	return CurMap["CNY"]
}

// F64 根据币种转换成实际的金额单位
// 比如 CNY-transAmt/100
// 比如 JPY-transAmt/1
func F64(currency string, amt int64) float64 {
	return Get(currency).F64(amt)
}

// Str CNY-1->0.01,JPY-1->1
func Str(currency string, amt int64) string {
	return Get(currency).Str(amt)
}

// F64 根据币种转换成实际的金额单位
// 比如 CNY-transAmt/100
// 比如 JPY-transAmt/1
func (c Cur) F64(amt int64) float64 {
	return float64(amt) / math.Pow(10, float64(c.Precision))
}

// Str CNY-1->0.01,JPY-1->1
func (c Cur) Str(amt int64) string {
	if c.Precision == 0 {
		return fmt.Sprintf("%d", amt)
	}
	return fmt.Sprintf("%0."+fmt.Sprintf("%d", c.Precision)+"f", c.F64(amt))
}
