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
	bo := bson.M{
		"chanMerId": c.ChanMerId,
		"chanCode":  c.ChanCode,
	}
	_, err := database.C(col.name).Upsert(bo, c)

	return err
}

// Modify 更新渠道商户信息
func (col *chanMerCollection) Update(c *model.ChanMer) error {
	bo := bson.M{
		"chanCode":  c.ChanCode,
		"chanMerId": c.ChanMerId,
	}
	return database.C(col.name).Update(bo, c)
}

// FindByCode 得到某个渠道所有商户
func (col *chanMerCollection) FindByCode(chanCode string) ([]*model.ChanMer, error) {
	var cs []*model.ChanMer
	err := database.C(col.name).Find(bson.M{"chanCode": chanCode}).All(&cs)
	return cs, err
}

// FindByCondition 根据渠道商户的条件查找渠道商户
func (col *chanMerCollection) FindByCondition(cond *model.ChanMer) (results []model.ChanMer, err error) {
	results = make([]model.ChanMer, 1)
	err = database.C(col.name).Find(nil).All(&results)
	if err != nil {
		log.Errorf("Find all merchant error: %s", err)
		return nil, err
	}

	return
}
