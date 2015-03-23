package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestTransAdd(t *testing.T) {
	trans := &model.Trans{
		TransStatus: 0,
		MerId:       "testTransMerId",
		OrderNum:    "testTransOrderNum",
		TransType:   1,
	}
	g.Debug("%+v", TransColl)

	err := TransColl.Add(trans)
	if err != nil {
		t.Errorf("add trans unsunccessful: %s", err)
		t.FailNow()
	}
	g.Debug("add trans success %s", trans)
}

func TestTransUpdate(t *testing.T) {
	objectId := bson.ObjectIdHex("550ee5e36a3dd74f93000001")
	trans := &model.Trans{
		// CreateTime:  time.Now().Unix(),
		Id:          objectId,
		MerId:       "111111110000000",
		OrderNum:    "222222220000000",
		TransType:   1,
		TransStatus: 0,
	}
	err := TransColl.Update(trans)
	if err != nil {
		t.Errorf("modify trans unsunccessful: %s", err)
		t.FailNow()
	}
	g.Debug("modify trans success %s", trans)

}

func TestCountTrans(t *testing.T) {

	c, err := TransColl.Count("111111110000000", "222222220000000")
	if err != nil {
		t.Errorf("count trans unsunccessful: %s", err)
		t.FailNow()
	}
	g.Debug("count trans success %d", c)
}

func TestFindTrans(t *testing.T) {
	trans, err := TransColl.Find("111111110000000", "222222220000000")
	if err != nil {
		t.Errorf("find trans unsunccessful: %s", err)
		t.FailNow()
	}
	g.Debug("find trans success %s", trans)
}

func TestFindByTime(t *testing.T) {
	trans, err := TransColl.FindByTime("2015-03-21")
	if err != nil {
		t.Errorf("find trans unsunccessful: %s", err)
		t.FailNow()
	}
	g.Debug("find trans success %s", trans)
}
