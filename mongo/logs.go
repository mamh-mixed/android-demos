package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
	"time"
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

// Find 查找莫个订单的日志
func (lc *logsCollection) Find(q *model.QueryCondition) ([]model.SpTransLogs, error) {

	query := bson.M{}
	if q.StartTime != "" && q.EndTime != "" {
		query["transTime"] = bson.M{"$gte": q.StartTime, "$lte": q.EndTime}
	}

	if q.MerId != "" {
		query["$or"] = []bson.M{bson.M{"merId": q.MerId}, bson.M{"origOrderNum": q.OrigOrderNum}}
	}

	if q.OrderNum != "" {
		query["orderNum"] = q.OrderNum
	}

	if q.OrigOrderNum != "" {
		query["origOrderNum"] = q.OrigOrderNum
	}

	var result []model.SpTransLogs
	err := database.C(lc.name).Find(query).Sort("-transTime").Skip((q.Page - 1) * q.Size).Limit(q.Size).All(&result)
	return result, err
}
