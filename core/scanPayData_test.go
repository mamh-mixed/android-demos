package core

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestTransStatistics(t *testing.T) {

	q := &model.QueryCondition{
		// MerName:     "讯联数据测试商户",
		// AgentCode:   "123131",
		StartTime:    "2015-08-13 00:00:00",
		EndTime:      "2015-08-13 23:59:59",
		TransStatus:  model.TransSuccess,
		TransType:    model.PayTrans,
		RefundStatus: model.TransRefunded,
		// MerId:       "10000",
		Size: 10,
		Page: 1,
	}
	ret := TransStatistics(q)
	bytes, err := json.Marshal(ret)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(string(bytes))
}
