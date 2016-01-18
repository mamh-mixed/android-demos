package mongo

import (
	"testing"
)

func TestFindPromoteLimit(t *testing.T) {
	users, err := AppUserCol.FindPromoteLimit("2015-12-31")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", users[0].MerId)
}
