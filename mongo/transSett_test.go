package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestSummary(t *testing.T) {

	all, err := TranSettColl.Summary("001405", "2015-03-21")
	if err != nil {
		t.Errorf("get transSett summary fail : (%s)", err)
	}
	g.Debug("summary transSett : (%s)", all)

}

func TestAdd(t *testing.T) {
	Id := bson.NewObjectId()
	transSett := &model.TransSett{
		SettDate: "2015-03-21 15:00:01",
		SettFlag: 1,
		SettAmt:  70,
		MerFee:   70,
	}
	transSett.MerId = "001405"
	transSett.TransAmt = 1000
	transSett.Id = Id
	transSett.TransType = 2
	err := TranSettColl.Add(transSett)
	if err != nil {
		t.Errorf("add transSett fail : (%s)", err)
	}
}
