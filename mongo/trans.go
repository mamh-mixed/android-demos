package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
	"time"
)

type Trans struct {
	Id      bson.ObjectId        `bson:"_id"`
	Chan    ChanMer              //渠道信息
	Payment model.BindingPayment //支付信息
	Time    int64                //时间
	Flag    int8                 //交易状态
}

// Add 添加一笔交易
func (t *Trans) Add() error {
	t.Id = bson.NewObjectId()
	t.Time = time.Now().Unix()
	t.Flag = 0
	return db.trans.Insert(t)
}

// Modify 通过Add时生成的Id来修改
func (t *Trans) Modify() error {
	return db.trans.Update(bson.M{"_id": t.Id}, t)
}
