package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
)

type Trans struct {
	Id                   string `bson:"_id"`
	ChanMer                     //TODO
	model.BindingPayment        //TODO
	Time                 int64  //时间
	Flag                 int8   //交易状态
}

// Add 添加一笔交易
func (t *Trans) Add() error {

	return mgodb.another.Insert(t)
}

func (t *Trans) Modify() error {
	return mgodb.another.Update(bson.M{"_id": t.Id}, t)
}
