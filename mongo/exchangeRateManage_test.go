package mongo

import (
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
)

func TestExchangeRatePaginationFind(t *testing.T) {
	cond := &model.ExchangeRateManage{
		LocalCurrency:  "JPY",
		TargetCurrency: "CNY",
		IsEnforced:     true,
	}

	size, page := 20, 1

	results, total, err := ExchangeRateManageColl.PaginationFind(cond, size, page)
	if err != nil {
		t.Logf("FAIL %s", err)
	}

	t.Logf("total is %d, result is %#v", total, results)
}

func TestAddExchangeRate(t *testing.T) {
	rate := &model.ExchangeRateManage{
		EId:            util.SerialNumber(),
		LocalCurrency:  "TWD",
		TargetCurrency: "USD",
		// Rate:                0.0306,
		PlanEnforcementTime: "2015-12-15 00:00:00",
		CreateTime:          time.Now().Format("2006-01-02 15:04:05"),
		CreateUser:          "admin",
		UpdateTime:          time.Now().Format("2006-01-02 15:04:05"),
		UpdateUser:          "admin",
		// ActualEnforcementTime: "",
	}

	err := ExchangeRateManageColl.Add(rate)
	t.Logf("error is %s", err)
}
