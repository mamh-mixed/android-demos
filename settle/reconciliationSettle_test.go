package settle

import (
	"testing"
)

func TestReconciliat(t *testing.T) {
	startTime := "2015-12-01 00:00:00"
	endTime := "2015-12-01 23:59:59"

	Reconciliat(startTime, endTime)
}
