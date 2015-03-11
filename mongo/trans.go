package mongo

import (
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

type Trans struct {
	Id                   string `bson:"_id,omitempty"`
	ChanCode             string
	ChanMer                    //TODO
	model.BindingPayment       //TODO
	Time                 int64 //时间
	Flag                 int8  //交易状态
}

// Add 添加一笔交易
func (t *Trans) Add() error {

	return db.trans.Insert(t)
}

func (t *Trans) Modify() error {
	return db.trans.Update(bson.M{"_id": t.Id}, t)
}
