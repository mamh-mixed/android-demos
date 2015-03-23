package mongo

import (
	"errors"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"gopkg.in/mgo.v2/bson"
)

type transCollection struct {
	name string
}

// TransColl 交易 Collection
var TransColl = transCollection{"trans"}

// Add 添加一笔交易
func (col *transCollection) Add(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	t.TransStatus = "00"
	return database.C(col.name).Insert(t)
}

// Update 通过Add时生成的Id来修改
func (col *transCollection) Update(t *model.Trans) error {
	t.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	return database.C(col.name).Update(bson.M{"_id": t.Id}, t)
}

// Count 通过订单号、商户号查找交易数量
func (col *transCollection) Count(merId, orderNum string) (count int, err error) {

	if merId == "" {
		return 0, errors.New("商户Id为空。")
	}

	if orderNum == "" {
		return 0, errors.New("订单号为空。")
	}
	// 不需要交易类型，对同一个商户下的订单号不能重复
	// if transType != 1 && transType != 2 {
	// 	return 0, errors.New("交易类型错误。")
	// }

	q := bson.M{
		"orderNum": orderNum,
		"merId":    merId,
		// "transType": transType,
	}
	count, err = database.C(col.name).Find(q).Count()
	return
}

// FindPay 通过订单号、商户号查找一条交易记录
func (col *transCollection) Find(merId, orderNum string) (t *model.Trans, err error) {

	q := bson.M{
		"orderNum": orderNum,
		"merId":    merId,
		// "transType": transType,
	}
	t = new(model.Trans)
	err = database.C(col.name).Find(q).One(t)
	return
}

// FindByTime 查找某天的交易记录
func (col *transCollection) FindByTime(time string) ([]model.Trans, error) {

	var ts []model.Trans
	q := bson.M{
		"createTime": bson.M{
			"$gt":  time,
			"$lte": tools.NextDay(time),
		},
	}
	err := database.C(col.name).Find(q).All(&ts)
	return ts, err
}
