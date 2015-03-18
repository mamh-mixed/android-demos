package mongo

import (
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
	"time"
)

func TestTransAdd(t *testing.T) {
	trans := &model.Trans{
		CreateTime:  time.Now().Unix(),
		TransStatus: 0,
		MerId:       "11031012",
		OrderNum:    "20150317",
		TransType:   1,
	}
	g.Debug("%+v", TransColl)
	err := TransColl.Add(trans)
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
	err := TransColl.Update(trans)
	if err != nil {
		t.Errorf("modify trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("modify trans success %s", trans)

}

func TestCountTrans(t *testing.T) {

	// trans := &model.Trans{
	// 	MerId:     "4353332424",
	// 	OrderNum:  "22323232323233",
	// 	TransType: 1,
	// }
	c, err := TransColl.Count("4353332424", "22323232323233", 1)
	if err != nil {
		t.Errorf("count trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("count trans success %s", c)
}

func TestFindTrans(t *testing.T) {
	// trans := &model.Trans{
	// 	MerId:     "4353332424",
	// 	OrderNum:  "22323232323233",
	// 	TransType: 1,
	// }
	trans, err := TransColl.Find("11031012", "20150317", 1)
	if err != nil {
		t.Errorf("find trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("find trans success %s", trans)
}
