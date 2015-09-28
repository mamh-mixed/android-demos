package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	// "github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

var SpChanLogsCol = logsCollection{"logs.chan.sp"}
var SpMerLogsCol = logsCollection{"logs.mer.sp"}

type logsCollection struct {
	name string
}

// Add 增加一条日志
func (lc *logsCollection) Add(l *model.SpTransLogs) error {
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
