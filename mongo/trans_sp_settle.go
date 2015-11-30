package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
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
