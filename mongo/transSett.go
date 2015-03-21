package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type transSettCollection struct {
	name string
}

var TranSettColl = transSettCollection{"transSett"}

func (col *transSettCollection) Summary(merId, settDate string) ([]model.SummarySettData, error) {

	//根据商户号、清算时间查找成功清算交易的汇总信息
	var s []model.SummarySettData
	//使用pipe统计
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": bson.M{"merId": merId}},
		// {"$project": bson.M{"transType": 1, "totalTransAmt": 1}},
		{"$group": bson.M{
			"_id":           "$transType",
			"totalTransAmt": bson.M{"$sum": "$transAmt"},
			"totalSettAmt":  bson.M{"$sum": "$settAmt"},
			"totalMerFee":   bson.M{"$sum": "$merFee"},
			"totalTransNum": bson.M{"$sum": 1},
			"transType":     bson.M{"$avg": "$transType"},
		}},
	}).All(&s)

	return s, err
}

// Add 增加一条清分记录
func (col *transSettCollection) Add(t *model.TransSett) error {
	return database.C(col.name).Insert(t)
}
