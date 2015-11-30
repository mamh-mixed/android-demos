package mongo

import (
	"fmt"
	"testing"
)

func TestFindBySettleTime(t *testing.T) {
	startTime := "2015-10-01 00:00:00"
	endTime := "2015-11-09 23:59:59"
	coll, _ := SpTransSettleColl.FindBySettleTime(startTime, endTime)
	fmt.Printf("the values is :", coll)
}
