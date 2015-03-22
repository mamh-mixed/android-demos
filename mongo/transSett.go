package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"gopkg.in/mgo.v2/bson"
)

type transSettCollection struct {
	name string
}

var TransSettColl = transSettCollection{"transSett"}

func (col *transSettCollection) Summary(merId, settDate string) ([]model.SummarySettData, error) {

	//根据商户号、清算时间查找成功清算交易的汇总信息
	var s []model.SummarySettData
	//使用pipe统计
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": bson.M{
			"merId": merId,
			"settDate": bson.M{"$gt": settDate,
				"$lte": tools.NextDay(settDate),
			},
			"settFlag": 1,
		}},
		{"$group": bson.M{
			"_id":           "$transType",
			"totalTransAmt": bson.M{"$sum": "$transAmt"},
			"totalSettAmt":  bson.M{"$sum": "$settAmt"},
			"totalMerFee":   bson.M{"$sum": "$merFee"},
			"totalTransNum": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{"transType": "$_id",
			"totalTransAmt": 1,
			"totalSettAmt":  1,
			"totalMerFee":   1,
			"totalTransNum": 1}},
	}).All(&s)

	return s, err
}

// Add 增加一条清分记录
func (col *transSettCollection) Add(t *model.TransSett) error {
	return database.C(col.name).Insert(t)
}

// Find 根据商户Id,清分时间查找交易明细
// 按照清分时间降排序
// TODO确定返回的struct
func (col *transSettCollection) Find(merId, settDate string) ([]model.TransInfo, error) {

	var transInfo []model.TransInfo
	q := bson.M{
		"merId":    merId,
		"settFlag": 1,
		"settDate": bson.M{"$gt": settDate, "$lte": tools.NextDay(settDate)},
	}
	err := database.C(col.name).Find(q).Sort("-settDate").All(&transInfo)

	return transInfo, err
}
