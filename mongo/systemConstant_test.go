package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestUpsertSystemConstant(t *testing.T) {
	cond := &model.SystemConstant{
		RateFloatingUpper: 1.2,
		RateFloatingLower: 0.9,
	}

	err := SysConstColl.Upsert(cond)
	if err != nil {
		t.Logf("Fail is %s", err)
		t.Fail()
	}

}
