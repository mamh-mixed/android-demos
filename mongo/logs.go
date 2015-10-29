package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

var SpChanLogsCol = logsCollection{"logs.chan.sp"}
var SpMerLogsCol = logsCollection{"logs.mer.sp"}

type logsCollection struct {
	name string
}

// Add 增加一条日志
func (lc *logsCollection) Add(l *model.SpTransLogs) error {
	l.TransTime = time.Now().Format("2006-01-02 15:04:05")
	return database.C(lc.name).Insert(l)
}

// FindOne 查找一条记录
func (lc *logsCollection) FindOne(q *model.QueryCondition) (*model.SpTransLogs, error) {
	var result = new(model.SpTransLogs)
	err := database.C(lc.name).Find(lc.query(q)).One(result)
	return result, err
}

// Find 查找莫个订单的日志
func (lc *logsCollection) Find(q *model.QueryCondition) ([]model.SpTransLogs, error) {

	var result []model.SpTransLogs
	var err error
	if len(q.ReqIds) > 0 {
		err = database.C(lc.name).Find(lc.query(q)).Sort("transTime").All(&result)
	} else {
		err = database.C(lc.name).Find(lc.query(q)).Sort("transTime").Skip((q.Page - 1) * q.Size).Limit(q.Size).All(&result)
	}

	return result, err
}

// Count 总数
func (lc *logsCollection) Count(q *model.QueryCondition) (int, error) {
	// 计算total
	var total = struct {
		Value int `bson:"total"`
	}{}
	err := database.C(lc.name).Pipe([]bson.M{
		{"$match": lc.query(q)},
		{"$group": bson.M{"_id": "$reqId"}},
		{"$group": bson.M{"_id": "null", "total": bson.M{"$sum": 1}}},
	}).One(&total)
	return total.Value, err
}

func (lc *logsCollection) query(q *model.QueryCondition) bson.M {
	query := bson.M{}
	if len(q.ReqIds) > 0 {
		query["reqId"] = bson.M{"$in": q.ReqIds}
	}

	if q.Direction != "" {
		query["direction"] = q.Direction
	}

	if q.MerId != "" {
		query["merId"] = q.MerId
	}

	if q.Busicd != "" {
		query["transType"] = q.Busicd
	}

	if q.OrderNum != "" {
		query["$or"] = []bson.M{bson.M{"orderNum": q.OrderNum}, bson.M{"origOrderNum": q.OrderNum, "transType": model.Inqy}}
	}

	if q.OrigOrderNum != "" {
		query["origOrderNum"] = q.OrigOrderNum
	}
	return query
}
