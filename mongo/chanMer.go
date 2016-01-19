package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
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
		// log.Errorf("Find ChanMer condition is: %+v;error is %s", bo, err)
		return nil, err
	}
	return c, nil
}

func (col *chanMerCollection) FindByArea(areaType int) ([]model.ChanMer, error) {
	var result []model.ChanMer
	err := database.C(col.name).Find(bson.M{"areaType": areaType}).All(&result)
	return result, err
}

// Add 增加一个渠道商户
func (col *chanMerCollection) Add(c *model.ChanMer) error {
	c.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	bo := bson.M{
		"chanMerId": c.ChanMerId,
		"chanCode":  c.ChanCode,
	}
	_, err := database.C(col.name).Upsert(bo, c)

	return err
}

// BatchAdd 批量增加渠道商户
func (col *chanMerCollection) BatchAdd(cms []model.ChanMer) error {

	var temp []interface{}
	for _, cm := range cms {
		cm.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		cm.UpdateTime = cm.CreateTime
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
	c.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	bo := bson.M{
		"chanCode":  c.ChanCode,
		"chanMerId": c.ChanMerId,
	}
	return database.C(col.name).Update(bo, c)
}

// Upsert
func (col *chanMerCollection) Upsert(c *model.ChanMer) error {
	bo := bson.M{
		"chanCode":  c.ChanCode,
		"chanMerId": c.ChanMerId,
	}
	_, err := database.C(col.name).Upsert(bo, c)
	return err
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
func (col *chanMerCollection) PaginationFind(chanCode, chanMerId, chanMerName, pay string, size, page int) (results []*model.ChanMer, total int, err error) {
	results = make([]*model.ChanMer, 0)

	match := bson.M{}

	if pay == "bp" {
		match["chanCode"] = bson.M{"$in": []string{"CFCA", "CIL", "Mock"}}
	} else {
		match["chanCode"] = bson.M{"$in": []string{"ALP", "WXP", "ULIVE"}}
	}

	if chanCode != "" {
		match["$and"] = []interface{}{bson.M{"chanCode": chanCode}}
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
	sort := bson.M{"$sort": bson.M{"chanMerId": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, sort, skip, limit)

	err = database.C(col.name).Pipe(cond).All(&results)

	return results, total, err
}

// FuzzyFind 模糊查找
func (col *chanMerCollection) FuzzyFind(chanCode, chanMerId, chanMerName string, maxSize int) (results []*model.ChanMer, err error) {
	results = make([]*model.ChanMer, 0)

	q := bson.M{}

	if chanCode != "" {
		cc := []bson.M{}
		cc = append(cc, bson.M{"chanCode": bson.RegEx{chanCode, "i."}})
		q["$and"] = cc
	}

	if chanMerId != "" {
		cmi := []bson.M{}
		cmi = append(cmi, bson.M{"chanMerId": bson.RegEx{chanMerId, "i."}})
		q["$and"] = cmi
	}

	if chanMerName != "" {
		cmn := []bson.M{}
		cmn = append(cmn, bson.M{"chanMerName": bson.RegEx{chanMerName, "i."}})
		q["$and"] = cmn
	}

	err = database.C(col.name).Find(q).Sort("chanMerId", "chanCode").All(&results)
	if err != nil {
		log.Errorf("fuzzy find channel merchant error: %s", err)
		return nil, err
	}

	if len(results) > maxSize {
		return results[:maxSize], err
	}

	return results, err
}

// Remove 删除渠道商户
func (col *chanMerCollection) Remove(chanCode, chanMerId string) error {

	q := bson.M{}
	if chanCode != "" {
		q["chanCode"] = chanCode
	}
	if chanMerId != "" {
		q["chanMerId"] = chanMerId
	}
	err := database.C(col.name).Remove(q)

	return err
}

func (col *chanMerCollection) Insert(c *model.ChanMer) error {
	c.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	c.UpdateTime = c.CreateTime
	err := database.C(col.name).Insert(c)
	return err
}

// FindWXPAgent 查找微信代理
func (col *chanMerCollection) FindWXPAgent() ([]model.ChanMer, error) {

	var result []model.ChanMer

	find := bson.M{
		"chanCode":    "WXP",
		"isAgentMode": false,
	}
	err := database.C(col.name).Find(find).All(&result)

	return result, err
}
