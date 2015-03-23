package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestTransSettSummary(t *testing.T) {

	all, err := TransSettColl.Summary("001405", "2015-03-23")
	if err != nil {
		t.Errorf("get transSett summary fail : (%s)", err)
	}
	g.Debug("summary transSett : (%s)", all)

}

func TestTransSettAdd(t *testing.T) {
	Id := bson.NewObjectId()

	tran := model.Trans{
		Id:        Id,
		MerId:     "001405",
		TransAmt:  11111,
		TransType: 2,
	}
	transSett := &model.TransSett{
		Tran:       tran,
		SettDate:   "2015-03-23 23:59:59",
		SettFlag:   1,
		MerSettAmt: 100,
		MerFee:     100,
	}
	// transSett.MerId = "001405"
	// transSett.TransAmt = 10000
	// transSett.Id = Id
	// transSett.TransType = 1
	err := TransSettColl.Add(transSett)
	if err != nil {
		t.Errorf("add transSett fail : (%s)", err)
	}
}

func TestTransSettFind(t *testing.T) {
	trans, err := TransSettColl.Find("001405", "2015-03-23")
	if err != nil {
		t.Errorf("find trans fail : %s", err)
	}
	g.Debug("%+v", trans)
}
