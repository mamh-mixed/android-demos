package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

func TestTransSettSummary(t *testing.T) {

	all, err := TransSettColl.Summary("001405", "2015-03-23")
	if err != nil {
		t.Errorf("get transSett summary fail : (%s)", err)
	}
	log.Debugf("summary transSett : (%s)", all)

}

func TestTransSettAdd(t *testing.T) {

	debug := false
	if debug {
		Id := bson.NewObjectId()
		orderNum := "1000000003"
		tran := model.Trans{
			Id:         Id,
			MerId:      "001405",
			TransAmt:   11111,
			TransType:  1,
			OrderNum:   orderNum,
			CreateTime: "2015-03-23 23:59:59",
		}
		transSett := &model.TransSett{
			Tran:       tran,
			SettDate:   "2015-03-23 23:59:59",
			SettFlag:   1,
			MerSettAmt: 100,
			MerFee:     100,
		}
		err := TransSettColl.Add(transSett)
		if err != nil {
			t.Errorf("add transSett fail : (%s)", err)
		}
	}
}

func TestTransSettFind(t *testing.T) {
	trans, err := TransSettColl.Find("001405", "2015-03-23", "")
	if len(trans) == 11 {
		data := trans[:len(trans)-1]
		log.Debugf("%+v", data)
		lastOrderNum := trans[len(trans)-1].OrderNum
		log.Debugf("%s", lastOrderNum)
	}
	if err != nil {
		t.Errorf("find trans fail : %s", err)
	}
	log.Debugf("%+v", trans)
}
