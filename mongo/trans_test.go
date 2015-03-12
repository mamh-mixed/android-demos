package mongo

import (
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"

	"github.com/omigo/g"
)

func TestTransAdd(t *testing.T) {
	trans := Trans{
		Time: time.Now().Unix(),
		Flag: 0,
	}
	err := trans.Add()
	if err != nil {
		t.Errorf("add trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("add trans success %s", trans)
}

func TestTransModify(t *testing.T) {
	objectId := bson.ObjectIdHex("55004d5d6a3dd74ef8000001")
	trans := Trans{
		Time: time.Now().Unix(),
		Id:   objectId,
		Flag: 0,
	}
	err := trans.Modify()
	if err != nil {
		t.Errorf("modify trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("modify trans success %s", trans)

}
