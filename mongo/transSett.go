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
			"settDate": bson.M{"$gte": settDate,
				"$lt": tools.NextDay(settDate),
			},
			"settFlag": 1,
		}},
		{"$group": bson.M{
			"_id":           "$transType",
			"totalTransAmt": bson.M{"$sum": "$transAmt"},
			"totalSettAmt":  bson.M{"$sum": "$merSettAmt"},
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
// 按照商户订单号降排序
func (col *transSettCollection) Find(merId, settDate, nextOrderNum string) ([]model.TransSettInfo, error) {

	var transSettInfo []model.TransSettInfo

	p := []bson.M{
		//查找
		{"$match": bson.M{"merId": merId, "settFlag": 1,
			"settDate": bson.M{"$gte": settDate, "$lt": tools.NextDay(settDate)}}},
		//排序
		{"$sort": bson.M{"orderNum": -1}},
	}
	//商户实际拉取为10
	limit := bson.M{"$limit": 11}
	if nextOrderNum != "" {
		p = append(p, bson.M{"$match": bson.M{"orderNum": bson.M{"$lte": nextOrderNum}}}, limit)
	} else {
		p = append(p, limit)
	}
	err := database.C(col.name).Pipe(p).All(&transSettInfo)
	return transSettInfo, err
}
