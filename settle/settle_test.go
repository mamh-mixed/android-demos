package settle

import (
	"github.com/CardInfoLink/quickpay/model"
	// "github.com/CardInfoLink/quickpay/mongo"
	"testing"

	"github.com/CardInfoLink/quickpay/tools"
	"gopkg.in/mgo.v2/bson"
)

func TestAddTrans(t *testing.T) {

	tran := &model.Trans{
		SysOrderNum: tools.SerialNumber(),
		Id:          bson.NewObjectId(),
		OrderNum:    tools.Millisecond(),
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
