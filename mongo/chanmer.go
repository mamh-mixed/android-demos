package mongo

import (
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

// Find 根据渠道代码、商户号查找
func FindChanMer(c *model.ChanMer) error {

	bo := bson.M{
		"chancode":  c.ChanCode,
		"chanmerid": c.ChanMerId,
	}
	return db.chanMer.Find(bo).One(c)
}

// Add 增加一个渠道商户
func AddChanMer(c *model.ChanMer) error {
	return db.chanMer.Insert(c)
}

// Modify 更新渠道商户信息
func ModifyChanMer(c *model.ChanMer) error {
	bo := bson.M{
		"chancode":  c.ChanCode,
		"chanmerid": c.ChanMerId,
	}
	return db.chanMer.Update(bo, c)
}
