package mongo

import (
	"errors"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type transSettCollection struct {
	name string
}

// transSettLogCollection
// 记录某台机器某个时间点执行的方法
type transSettLogCollection struct {
	name string
}

var TransSettColl = transSettCollection{"transSett"}
var TransSettLogColl = transSettLogCollection{"transSettLog"}

// AtomUpsert mongodb-findAndModify
func (col *transSettLogCollection) AtomUpsert(l *model.TransSettLog) (int, error) {
	q := bson.M{
		"method":    l.Method,
		"date":      l.Date,
		"$isolated": 1,
	}
	c := mgo.Change{}
	// 开始时status==0
	// 只更新method,date
	if l.Status == 0 {
		c.Update = bson.M{
			"$set": bson.M{
				"method": l.Method,
				"date":   l.Date,
			},
		}
	} else {
		// 结束时
		c.Update = l
	}
	c.Upsert = true

	result := new(model.TransSettLog)
	change, err := database.C(col.name).Find(q).Apply(c, result)

	return change.Updated, err

}

// Summary 清算信息汇总
func (col *transSettCollection) Summary(merId, transDate string) ([]model.SummarySettData, error) {

	//根据商户号、清算时间查找成功清算交易的汇总信息
	var s []model.SummarySettData
	//使用pipe统计
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": bson.M{
			"merId": merId,
			"createTime": bson.M{"$gte": transDate,
				"$lt": tools.NextDay(transDate),
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
	if t == nil {
		return errors.New("transSett is nil")
	}
	t.SettDate = time.Now().Format("2006-01-02 15:04:05")
	return database.C(col.name).Insert(t)
}

// Find 根据商户Id,清分时间查找交易明细
// 按照商户订单号降排序
func (col *transSettCollection) Find(merId, transDate, nextOrderNum string) ([]model.TransSettInfo, error) {

	var transSettInfo []model.TransSettInfo

	p := []bson.M{
		//查找
		{"$match": bson.M{"merId": merId, "settFlag": 1,
			"createTime": bson.M{"$gte": transDate, "$lt": tools.NextDay(transDate)}}},
		//排序
		{"$sort": bson.M{"orderNum": 1}},
	}
	//商户实际拉取为10
	limit := bson.M{"$limit": 11}
	if nextOrderNum != "" {
		p = append(p, bson.M{"$match": bson.M{"orderNum": bson.M{"$gte": nextOrderNum}}}, limit)
	} else {
		p = append(p, limit)
	}
	err := database.C(col.name).Pipe(p).All(&transSettInfo)
	return transSettInfo, err
}

// FindByOrderNum 根据渠道订单号查找
func (col *transSettCollection) FindByOrderNum(sysOrderNum string) (t *model.TransSett, err error) {
	// 订单是uuid 全局唯一
	t = new(model.TransSett)
	q := bson.M{
		"sysOrderNum": sysOrderNum,
	}
	err = database.C(col.name).Find(q).One(t)

	return
}

// Update 更新
func (col *transSettCollection) Update(t *model.TransSett) error {
	if t == nil {
		return errors.New("transSett is nil")
	}
	t.SettDate = time.Now().Format("2006-01-02 15:04:05")
	return database.C(col.name).Update(bson.M{"_id": t.Tran.Id}, t)
}
