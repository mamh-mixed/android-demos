package query

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestCouponTransQuery(t *testing.T) {

	q := &model.QueryCondition{
		Page: 1,
		Size: 10,
	}
	result := CouponTransQuery(q)

	t.Logf("result is %#v", result)

	t.Logf("pagination is %#v", result.Data)

	t.Logf("data is %#v", result.Data.(*model.Pagination).Data)
}
