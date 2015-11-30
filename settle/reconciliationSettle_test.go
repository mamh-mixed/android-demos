package settle

import (
	"testing"
)

func TestReconciliat(t *testing.T) {
	startTime := "2015-10-01 00:00:00"
	endTime := "2015-11-09 23:59:59"

	Reconciliat(startTime, endTime)
}
