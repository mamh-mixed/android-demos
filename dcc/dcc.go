package dcc

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"math"
)

// defaultDcc 默认台币对美元
var defaultDcc = Dcc{localCurr: "TWD", targetCurr: "USD"}

type Dcc struct {
	localCurr  string
	targetCurr string
}

// New 指定转换币种
func New(localCurr, targetCurr string) Dcc {
	return Dcc{localCurr: localCurr, targetCurr: targetCurr}
}

// NewUSD 对美元
func NewUSD(localCurr string) Dcc {
	return Dcc{localCurr: localCurr, targetCurr: "USD"}
}

// Do 汇率转换
func Do(txamt int64) (amt int64, rate float64, err error) {
	return defaultDcc.Do(txamt)
}

// Do 汇率转换
func (d Dcc) Do(txamt int64) (amt int64, rate float64, err error) {

	if d.localCurr == "" || d.targetCurr == "" {
		return amt, rate, fmt.Errorf("%s", "params not found, need localCurr and targetCurr")
	}

	// 同一币种
	if d.localCurr == d.targetCurr {
		return txamt, 1.00, nil
	}
	var retry int
	var er *model.ExchangeRate
	for {
		// 防止偶尔的查询失败影响交易
		er, err = mongo.ExchangeRateColl.FindOne(d.localCurr, d.targetCurr)
		if err != nil {
			retry++
			if retry == 3 {
				log.Errorf("find rate error: %s, dcc: %s ==> %s", err, d.localCurr, d.targetCurr)
				return amt, rate, err
			}
			continue
		}
		break
	}
	rate = er.Rate
	amt = int64(math.Floor(float64(txamt)/rate + 0.5))
	return amt, rate, nil
}
