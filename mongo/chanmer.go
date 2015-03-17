package mongo

import (
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

type chanMerCollection struct {
	name string
}

var ChanMerColl = chanMerCollection{"chanMer"}

// Find 根据渠道代码、商户号查找
func (col *chanMerCollection) Find(c *model.ChanMer) error {

	bo := bson.M{
		"chancode":  c.ChanCode,
		"chanmerid": c.ChanMerId,
	}
	return database.C(col.name).Find(bo).One(c)
}

// Add 增加一个渠道商户
func (col *chanMerCollection) Add(c *model.ChanMer) error {
	return database.C(col.name).Insert(c)
}

// Modify 更新渠道商户信息
func (col *chanMerCollection) Update(c *model.ChanMer) error {
	bo := bson.M{
		"chancode":  c.ChanCode,
		"chanmerid": c.ChanMerId,
	}
	return database.C(col.name).Update(bo, c)
}
