package mongo

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
	"strings"
	"time"
)

type transCollection struct {
	trans *mgo.Collection
}

var TransColl = transCollection{database.C("trans")}

// Add 添加一笔交易
func (col *transCollection) Add(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.CreateTime = time.Now().Unix()
	t.TransStatus = 0
	return col.trans.Insert(t)
}

// Update 通过Add时生成的Id来修改
func (col *transCollection) Update(t *model.Trans) error {
	t.UpdateTime = time.Now().Unix()
	return col.trans.Update(bson.M{"_id": t.Id}, t)
}

// Count 通过订单号、商户号、交易类型查找交易数量
func (col *transCollection) Count(t *model.Trans) (int, error) {

	if t.MerId == "" {
		return 0, errors.New("商户Id为空。")
	}

	if t.OrderNum == "" {
		return 0, errors.New("订单号为空。")
	}

	if t.TransType != 1 || t.TransTyoe != 2 {
		return 0, errors.New("交易类型错误。")
	}

	q := bson.M{
		"merOrderNum": t.OrderNum,
		"merId":       t.MerId,
		"transType":   t.TransType,
	}
	return col.trans.Find(q).Count()
}

// FindPay 通过订单号、商户号、交易类型查找一条交易记录
func (col *transCollection) Find(t *model.Trans) error {

	q := bson.M{
		"merOrderNum": t.OrderNum,
		"merId":       t.MerId,
		"transType":   t.TransType,
	}
	return col.trans.Find(q).One(t)
}
