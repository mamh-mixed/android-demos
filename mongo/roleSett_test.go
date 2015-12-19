package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestFindOneRoleSettCol(t *testing.T) {
	role, date := "CIL", "2015-12-10"

	err := RoleSettCol.Upsert(&model.RoleSett{SettDate: date, SettRole: role, ReportName: "sett/report/20151210/IC202_CIL_20151210.xlsx"})
	if err != nil {
		t.Errorf("test FAIL: %s", err)
	}

	// t.Logf("result is %#v", rs)
}

func TestPaginationFindRoleSettCol(t *testing.T) {
	role, date, reportType := "ALP", "", 1

	size, page := 10, 1

	result, total, err := RoleSettCol.PaginationFind(role, date, 0, size, page)

	if err != nil {
		t.Errorf("test FAIL: %s", err)
	}

	t.Logf("resul length is %d", len(result))

	t.Logf("total is %d", total)

	for i, item := range result {
		t.Logf("index %d, item %#v", i, item)
	}
}
