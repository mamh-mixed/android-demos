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

// Find
func (lc *logsCollection) Find(merId, orderNum string) ([]model.SpTransLogs, error) {
	query := bson.M{
		"merId":    merId,
		"orderNum": orderNum,
	}
	var result []model.SpTransLogs
	err := database.C(lc.name).Find(query).All(&result)
	return result, err
}
