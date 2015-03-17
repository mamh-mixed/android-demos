package mongo

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
	"time"
)

// Add 添加一笔交易
func AddTrans(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.CreateTime = time.Now().Unix()
	t.TransStatus = 0
	return db.trans.Insert(t)
}

// Modify 通过Add时生成的Id来修改
func ModifyTrans(t *model.Trans) error {
	t.UpdateTime = time.Now().Unix()
	return db.trans.Update(bson.M{"_id": t.Id}, t)
}

// CountTrans 通过订单号、商户号查找交易数量
func CountTrans(t *model.Trans) (int, error) {

	if t.MerId == "" {
		return 0, errors.New("商户Id为空。")
	}

	if t.OrderNum == "" {
		return 0, errors.New("订单号为空。")
	}

	q := bson.M{
		"merOrderNum": t.OrderNum,
		"merId":       t.MerId,
	}
	return db.trans.Find(q).Count()
}

// FindTrans 通过订单号、商户号查找
func FindTrans(t *model.Trans) error {
	q := bson.M{
		"merOrderNum": t.OrderNum,
		"merId":       t.MerId,
	}
	return db.trans.Find(q).One(t)
}
