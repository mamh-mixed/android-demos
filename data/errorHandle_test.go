package data

import (
	"testing"
)

func TestAddSettRoleToTrans(t *testing.T) {
	merIds, err := ReadMerIdsByTxt("/Users/zhiruichen/Desktop/mer.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("read %d merId", len(merIds))

	err = AddSettRoleToTrans(merIds)
	if err != nil {
		t.Error(err)
	}
}
