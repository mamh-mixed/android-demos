package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

func TestFindTransQuery(t *testing.T) {

	q := &model.QueryCondition{
		StartTime: "2015-07-01 00:00:00",
		EndTime:   "2015-07-30 00:00:00",
		Mchntid:   "100000000000203",
		Page:      1,
		Size:      10,
	}

	transInfo, total, err := SpTransColl.Find(q)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("total : %d", total)
	for _, v := range transInfo {
		t.Logf("%s,%s", v.OrderNum, v.TransTime)
	}

	// t.Logf("%d", len(transInfo))
}

func TestTransAdd(t *testing.T) {
	if debug {
		trans := &model.Trans{
			TransStatus: transStatus,
			MerId:       merId,
			OrderNum:    orderNum,
			TransType:   int8(transType),
		}
		log.Debugf("%+v", TransColl)

		err := TransColl.Add(trans)
		if err != nil {
			t.Errorf("add trans unsunccessful: %s", err)
			t.FailNow()
		}
		log.Debugf("add trans success %s", trans)
	}
}

func TestTransUpdate(t *testing.T) {
	objectId := bson.ObjectIdHex(hexId)
	trans := &model.Trans{
		// CreateTime:  time.Now().Unix(),
		Id:          objectId,
		MerId:       merId,
		OrderNum:    orderNum,
		TransType:   int8(transType),
		TransStatus: transStatus,
	}
	err := TransColl.Update(trans)
	if err != nil {
		t.Errorf("modify trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("modify trans success %s", trans)

}

func TestCountTrans(t *testing.T) {

	c, err := TransColl.Count(merId, orderNum)
	if err != nil {
		t.Errorf("count trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("count trans success %d", c)
}

func TestFindTrans(t *testing.T) {
	trans, err := TransColl.FindOne(merId, orderNum)
	if err != nil {
		t.Errorf("find trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("find trans success %s", trans)
}

func TestFindByTime(t *testing.T) {
	trans, err := TransColl.FindByTime(createTime)
	if err != nil {
		t.Errorf("find trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("find trans success %s", trans)
}

func TestFindTransRefundAmt(t *testing.T) {
	total, err := TransColl.FindTransRefundAmt("1000000000002", "DqfTuPvvvTWDfD0Ke9DGOqbT")
	if err != nil {
		t.FailNow()
	}
	log.Debug(total)
}
