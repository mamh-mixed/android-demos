package mongo

import (
	"testing"
)

func TestNotifyGetAll(t *testing.T) {

	cans, err := NotifyColl.GetAll()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v", cans)
}
