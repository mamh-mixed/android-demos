package master

import (
	"testing"
)

func TestTradeSettleReportQuery(t *testing.T) {
	role, date := "", ""
	size, page := 10, 1
	result := tradeSettleReportQuery(role, date, size, page)

	t.Logf("result is %+v", result.Data)
}
