package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type transCollectionSettle struct {
	name string
}

var SpTransSettleColl = transCollectionSettle{"transSett.sp"}

func (col *transCollectionSettle) FindAll() ([]model.TransSett, error) {
	var result []model.TransSett
	err := database.C(col.name).Find(nil).All(&result)
	return result, err
}

func (col *transCollectionSettle) FindBySettleTime(startTime string, endTime string) ([]model.TransSett, error) {
	var result []model.TransSett

	find := bson.M{
		"trans.payTime": bson.M{"$gte": startTime, "$lte": endTime},
	}
	err := database.C(col.name).Find(find).All(&result)
	return result, err
}

func (col *transCollectionSettle) UpdateSettle(t *model.TransSett) error {

	condition := bson.M{"_id": t.Trans.Id}
	err := database.C(col.name).Update(condition, t)
	if err != nil {
		log.Errorf("UpdateSettle Settle(%+v) fail:%s ", t, err)
	}
	return err
}

func (col *transCollectionSettle) AddSettle(t *model.TransSett) error {
	err := database.C(col.name).Insert(t)
	if err != nil {
		log.Errorf("AddSettle Settle(%+v) fail: %s", t, err)
	}
	return err
}

func (col *transCollectionSettle) Add(t *model.TransSett) error {
	t.Trans.Id = bson.NewObjectId()
	t.Trans.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	err := database.C(col.name).Insert(t)
	if err != nil {
		log.Errorf("add TransSett(%+v) fail: %s", t, err)
	}
	return err
}
