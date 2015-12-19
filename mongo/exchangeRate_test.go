package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
	"time"
)

func TestUpdateExchangeRate(t *testing.T) {
	cond := &model.ExchangeRate{
		CurrencyPair:    "JPY<=>CNY",
		Rate:            1.33,
		EnforcementTime: time.Now().Format("2006-01-02 15:04:05"),
		EnforceUser:     "admin",
	}

	err := ExchangeRateColl.Upsert(cond)
	if err != nil {
		t.Logf("error is %s", err)
		t.FailNow()
	}
}

func TestFindExchangeRate(t *testing.T) {
	currPair := "TWD<=>CNY"

	result, err := ExchangeRateColl.FindOne(currPair)
	if err != nil {
		t.Logf("error is %s", err)
		t.FailNow()
	}

	t.Logf("result is %#v", result)
}

func TestTimeLocal(t *testing.T) {
	t.Logf("local is %s", time.Local)

	tm, err := time.ParseInLocation("2006-01-02 15:04:05", "2015-12-18 00:00:00", time.Local)

	t.Logf("time is %s, error is %s", tm, err)
}
