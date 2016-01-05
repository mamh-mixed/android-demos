package settle

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"

	"github.com/CardInfoLink/quickpay/util"
	"gopkg.in/mgo.v2/bson"
)

func TestSpSettle(t *testing.T) {
	s := &scanpayDomestic{}

	s.Reconciliation("2015-12-23")
}

func TestGenReport(t *testing.T) {
	SpSettReport("2015-12-08")
}

func TestSpReconciliatReport(t *testing.T) {
	SpReconciliatReport("2015-12-10")
}

func TestSpSettReport(t *testing.T) {
	SpSettReport("2015-12-27")
}

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
