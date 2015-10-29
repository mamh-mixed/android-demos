package settle

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"

	"github.com/CardInfoLink/quickpay/util"
	"gopkg.in/mgo.v2/bson"
)

func TestAddTrans(t *testing.T) {

	tran := &model.Trans{
		SysOrderNum: util.SerialNumber(),
		Id:          bson.NewObjectId(),
		OrderNum:    util.Millisecond(),
		ChanCode:    "CFCA",
		MerId:       "001405",
		TransAmt:    1700,
		TransType:   model.PayTrans,
		TransStatus: model.TransSuccess,
	}

	// 测试计算费率
	addTransSett(tran, model.SettSuccess)
	// mongo.TransColl.Add(tran)
}

func TestDoSettWork(t *testing.T) {
	yesterday = "2015-05-22"
	doTransSett()
}

func TestDoScanpaySettReport(t *testing.T) {
	yesterday = "2015-10-28"
	err := doScanpaySettReport(yesterday)
	if err != nil {
		t.Error(err)
	}
}
