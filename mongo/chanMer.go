package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type chanMerCollection struct {
	name string
}

// ChanMerColl 渠道商户 Collection
var ChanMerColl = chanMerCollection{"chanMer"}

// Find 根据渠道代码、商户号查找
func (col *chanMerCollection) Find(chanCode, chanMerId string) (c *model.ChanMer, err error) {

	bo := bson.M{
		"chanCode":  chanCode,
		"chanMerId": chanMerId,
	}
	c = new(model.ChanMer)
	err = database.C(col.name).Find(bo).One(c)
	if err != nil {
		log.Errorf("Find ChanMer condition is: %+v;error is %s", bo, err)
	}
	return
}

// Add 增加一个渠道商户
func (col *chanMerCollection) Add(c *model.ChanMer) error {
	return database.C(col.name).Insert(c)
}

// Modify 更新渠道商户信息
func (col *chanMerCollection) Update(c *model.ChanMer) error {
	bo := bson.M{
		"chanCode":  c.ChanCode,
		"chanMerId": c.ChanMerId,
	}
	return database.C(col.name).Update(bo, c)
}
