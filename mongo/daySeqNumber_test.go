package mongo

import (
	"testing"
)

// GetDaySN 返回一个当天唯一的六位数字
func TestGetDaySN(t *testing.T) {
	s := DaySNColl.GetDaySN("M126", "T126")
	t.Log(s)

	if s == "" {
		t.Error("Error")
	}

	s = DaySNColl.GetDaySN("M127", "T127")
	t.Log(s)

	if s == "" {
		t.Error("Error")
	}

	if s != "000000" {
		t.Error("Error")
	}
}
