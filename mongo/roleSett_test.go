package mongo

import (
	"testing"
)

func TestFindOneRoleSettCol(t *testing.T) {
	role, date := "CIL", "2015-10-13"

	rs, err := RoleSettCol.FindOne(role, date)
	if err != nil {
		t.Errorf("test FAIL: %s", err)
	}

	t.Logf("result is %#v", rs)
}

func TestPaginationFindRoleSettCol(t *testing.T) {
	role, date := "ALP", ""

	size, page := 10, 1

	result, total, err := RoleSettCol.PaginationFind(role, date, size, page)
	if err != nil {
		t.Errorf("test FAIL: %s", err)
	}

	t.Logf("resul length is %d", len(result))

	t.Logf("total is %d", total)

	for i, item := range result {
		t.Logf("index %d, item %#v", i, item)
	}
}
