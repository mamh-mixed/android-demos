package mongo

import (
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
	"testing"
	"time"
)

func TestTransAdd(t *testing.T) {
	trans := &model.Trans{
		CreateTime:  time.Now().Unix(),
		TransStatus: 0,
	}
	err := AddTrans(trans)
	if err != nil {
		t.Errorf("add trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("add trans success %s", trans)
}

func TestTransModify(t *testing.T) {
	objectId := bson.ObjectIdHex("55004d5d6a3dd74ef8000001")
	trans := &model.Trans{
		CreateTime:  time.Now().Unix(),
		Id:          objectId,
		TransStatus: 0,
	}
	err := ModifyTrans(trans)
	if err != nil {
		t.Errorf("modify trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("modify trans success %s", trans)

}

func TestCountTrans(t *testing.T) {

	trans := &model.Trans{
		MerId:    "4353332424",
		OrderNum: "22323232323233",
	}
	c, err := CountTrans(trans)
	if err != nil {
		t.Errorf("count trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("count trans success %s", c)
}

func TestFindTrans(t *testing.T) {
	trans := &model.Trans{
		MerId:    "4353332424",
		OrderNum: "22323232323233",
	}
	err := FindTrans(trans)
	if err != nil {
		t.Errorf("find trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("find trans success %s", trans)
}
