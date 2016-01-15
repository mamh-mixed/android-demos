package mongo

import (
	"testing"
)

func TestFindOneByCode(t *testing.T) {
	code := "898d1b79dc8e4f2243f6e831fcd4e487"
	email, err := EmailCol.FindOneByCode(code)
	if err != nil {
		t.Errorf("test fail: %s", err)
	}

	t.Logf("email is %+v", email)
}
