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
var ChanMerColl = chanMerCollection{"chanMer.old"}

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

// CountByKey 验证是否存在
func (col *chanMerCollection) CountByKey(chanCode, chanMerId string) (int, error) {
	bo := bson.M{
		"chanCode":  chanCode,
		"chanMerId": chanMerId,
	}
	return database.C(col.name).Find(bo).Count()
}

// BatchAdd 批量增加渠道商户
func (col *chanMerCollection) BatchAdd(cms []model.ChanMer) error {

	var temp []interface{}
	for _, cm := range cms {
		temp = append(temp, cm)
	}
	err := database.C(col.name).Insert(temp...)
	return err
}

// BatchRemove 批量删除渠道商户
func (col *chanMerCollection) BatchRemove(cms []model.ChanMer) error {

	var keys []bson.M
	for _, cm := range cms {
		keys = append(keys, bson.M{"chanCode": cm.ChanCode, "chanMerId": cm.ChanMerId})
	}

	selector := bson.M{
		"$or": keys,
	}
	change, err := database.C(col.name).RemoveAll(selector)
	if change.Removed != len(cms) {
		log.Warnf("expect remove %d records,but %d removed", len(cms), change.Removed)
	}
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
	err = database.C(col.name).Find(cond).All(&results)
	if err != nil {
		log.Errorf("Find all merchant error: %s", err)
		return nil, err
	}

	return
}

// PaginationFind 分页查找渠道商户的信息
func (col *chanMerCollection) PaginationFind(chanCode, chanMerId, chanMerName string, size, page int) (results []model.ChanMer, total int, err error) {
	results = make([]model.ChanMer, 1)

	match := bson.M{}
	if chanCode != "" {
		match["chanCode"] = chanCode
	}
	if chanMerId != "" {
		match["chanMerId"] = chanMerId
	}
	if chanMerName != "" {
		match["chanMerName"] = chanMerName
	}

	// 计算总数
	total, err = database.C(col.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, skip, limit)

	err = database.C(col.name).Pipe(cond).All(&results)

	return results, total, err
}
