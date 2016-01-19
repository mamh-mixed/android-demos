package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

func TestFindAndGroupBy(t *testing.T) {
	g, a, _ := SpTransSettColl.FindAndGroupBy(&model.QueryCondition{
		StartTime: "2015-12-01 00:00:00",
		EndTime:   "2015-12-01 23:59:59",
		MerId:     "100000000000017",
	})

	t.Logf("%+v", g)
	t.Logf("%+v", a)
}

func TestAtomUpsert(t *testing.T) {

	l := &model.TransSettLog{
		Method: "doSettWork",
		Date:   "2015-04-27",
	}
	c := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			updated, _ := TransSettLogColl.AtomUpsert(l)
			log.Debugf("%d", updated)
			c <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-c
	}
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }

}

func TestTransSettSummary(t *testing.T) {

	all, err := TransSettColl.Summary("1000000000002", "2015-04-13")
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
			Trans:      tran,
			SettDate:   "2015-03-23 23:59:59",
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
	trans, err := TransSettColl.FindByDate("001405", "2015-03-23", "")
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

func TestTransSettFindByOrderNum(t *testing.T) {
	transSett, err := TransSettColl.FindByOrderNum("86b3a23d495048fa5de2e3643464f116")
	if err != nil {
		t.Errorf("find trans fail : %s", err)
		t.FailNow()
	}
	log.Debug("%+v", transSett)
}
