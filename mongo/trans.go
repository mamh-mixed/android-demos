package mongo

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
	"time"
)

type transCollection struct {
	name string
}

var TransColl = transCollection{"trans"}

// Add 添加一笔交易
func (col *transCollection) Add(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.CreateTime = time.Now().Unix()
	t.TransStatus = 0
	return database.C(col.name).Insert(t)
}

// Update 通过Add时生成的Id来修改
func (col *transCollection) Update(t *model.Trans) error {
	t.UpdateTime = time.Now().Unix()
	return database.C(col.name).Update(bson.M{"_id": t.Id}, t)
}

// Count 通过订单号、商户号、交易类型查找交易数量
func (col *transCollection) Count(merId, orderNum string, transType int8) (count int, err error) {

	if merId == "" {
		return 0, errors.New("商户Id为空。")
	}

	if orderNum == "" {
		return 0, errors.New("订单号为空。")
	}

	if transType != 1 && transType != 2 {
		return 0, errors.New("交易类型错误。")
	}

	q := bson.M{
		"orderNum":  orderNum,
		"merId":     merId,
		"transType": transType,
	}
	count, err = database.C(col.name).Find(q).Count()
	return
}

// FindPay 通过订单号、商户号、交易类型查找一条交易记录
func (col *transCollection) Find(merId, orderNum string, transType int8) (t *model.Trans, err error) {

	q := bson.M{
		"orderNum":  orderNum,
		"merId":     merId,
		"transType": transType,
	}
	t = &model.Trans{}
	err = database.C(col.name).Find(q).One(t)
	return
}
