package mongo

import (
	"errors"

	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type merchantCollection struct {
	name string
}

var MerchantColl = merchantCollection{"merchant"}

var merCache = cache.New(model.Cache_Merchant)

// Insert 插入一个商户信息。如果存在则更新，不存在则插入。@WonSikin
func (c *merchantCollection) Upsert(m *model.Merchant) error {
	q := bson.M{"merId": m.MerId}

	_, err := database.C(c.name).Upsert(q, m)
	if err != nil {
		log.Errorf("'Upsert Merchant ERROR!' Merchant is (%+v); error is (%s)", m, err)
	}
	return err
}

// FindByUniqueId 查找商户信息
func (c *merchantCollection) FindByUniqueId(uniqueId string) (m *model.Merchant, err error) {
	m = new(model.Merchant)
	q := bson.M{"uniqueId": uniqueId}
	err = database.C(c.name).Find(q).One(m)
	return
}

// Find 查找商户信息
// 先从缓存里取，没有再访问数据库
func (c *merchantCollection) Find(merId string) (m *model.Merchant, err error) {

	// get from cache
	o, found := merCache.Get(merId)
	if found {
		m = o.(*model.Merchant)
		return m, nil
	}
	m = new(model.Merchant)
	q := bson.M{"merId": merId}
	err = database.C(c.name).Find(q).One(m)
	if err != nil {
		return nil, err
	}
	// save
	merCache.Set(merId, m, cache.DefaultExpiration)

	return m, nil
}

// CountById 检查商户是否存在
func (c *merchantCollection) CountById(merId string) (int, error) {
	q := bson.M{"merId": merId}
	return database.C(c.name).Find(q).Count()
}

// Update 更新一个商户信息。
func (c *merchantCollection) Update(m *model.Merchant) error {
	if m.MerId == "" {
		return errors.New("MerId is required!")
	}
	q := bson.M{"merId": m.MerId}
	err := database.C(c.name).Update(q, m)
	if err != nil {
		log.Errorf("'Update Merchant ERROR!' condition is (%+v);error is (%s)", q, err)
	}
	return err
}

// FindAllMerchant 查找所有的商户信息。
func (c *merchantCollection) FindAllMerchant(cond *model.Merchant) (results []model.Merchant, err error) {
	results = make([]model.Merchant, 1)
	err = database.C(c.name).Find(cond).All(&results)
	if err != nil {
		log.Errorf("Find all merchant error: %s", err)
		return nil, err
	}

	return
}

// FuzzyFind 模糊查询拿到merId
func (c *merchantCollection) FuzzyFind(cond *model.QueryCondition) ([]*model.Merchant, error) {

	q := bson.M{}
	if cond.AgentCode != "" {
		q["agentCode"] = cond.AgentCode
	}

	if cond.MerName != "" {
		or := []bson.M{}
		or = append(or, bson.M{"merDetail.merName": bson.RegEx{cond.MerName, "."}})
		or = append(or, bson.M{"merDetail.shortName": bson.RegEx{cond.MerName, "."}})
		q["$or"] = or
	}
	if cond.MerId != "" {
		and := []bson.M{}
		and = append(and, bson.M{"merId": bson.RegEx{cond.MerId, "."}})
		q["$and"] = and
	} else {
		if len(cond.MerIds) != 0 {
			q["merId"] = bson.M{"$in": cond.MerIds}
		}
	}

	var mers []*model.Merchant
	err := database.C(c.name).Find(q).All(&mers)
	return mers, err
}

// PaginationFind 分页查找机构商户
func (c *merchantCollection) PaginationFind(merId, merStatus, merName, groupCode, groupName, agentCode, agentName, pay string, size, page int) (results []*model.Merchant, total int, err error) {
	results = make([]*model.Merchant, 1)

	match := bson.M{}
	if merId != "" {
		match["merId"] = merId
	}
	if merStatus != "" {
		match["merStatus"] = merStatus
	}
	if merName != "" {
		match["merDetail.merName"] = merName
	}
	if groupCode != "" {
		match["groupCode"] = groupCode
	}
	if groupName != "" {
		match["groupName"] = groupName
	}
	if agentCode != "" {
		match["agentCode"] = agentCode
	}
	if agentName != "" {
		match["agentName"] = agentName
	}
	if pay == "bp" {
		match["encryptKey"] = bson.M{"$exists": true}
	} else {
		match["encryptKey"] = bson.M{"$exists": false}
	}
	// 计算总数
	total, err = database.C(c.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	sort := bson.M{"$sort": bson.M{"merId": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, sort, skip, limit)

	err = database.C(c.name).Pipe(cond).All(&results)

	return results, total, err
}

// BatchAdd 批量增加商户
func (c *merchantCollection) BatchAdd(mers []model.Merchant) error {

	var temps []interface{}
	for _, m := range mers {
		temps = append(temps, m)
	}
	err := database.C(c.name).Insert(temps...)
	return err
}

// BatchRemove 批量删除
func (c *merchantCollection) BatchRemove(merIds []string) error {

	selector := bson.M{
		"$in": merIds,
	}
	change, err := database.C(c.name).RemoveAll(selector)
	if change.Removed != len(merIds) {
		log.Warnf("expect remove %d records,but %d removed", len(merIds), change.Removed)
	}
	return err
}

// Remove 删除机构商户
func (col *merchantCollection) Remove(merId string) (err error) {
	bo := bson.M{}
	if merId != "" {
		bo["merId"] = merId
	}
	err = database.C(col.name).Remove(bo)
	return err
}

// Insert2 创建一个机构商户
func (c *merchantCollection) Insert2(m *model.Merchant) error {

	err := database.C(c.name).Insert(m)
	if err != nil {
		log.Errorf("'Insert Merchant ERROR!' Merchant is (%+v); error is (%s)", m, err)
		return err
	}
	return nil
}

// findMaxMerId 查询merId最大值
func (c *merchantCollection) FindMaxMerId(prefix string) (merId string, err error) {

	// match := bson.M{}
	// match["merId"] = bson.RegEx{prefix + ".", "\\d+"}
	cond := []bson.M{
		{"$match": bson.M{"merId": bson.M{"$regex": prefix + "\\d+"}}},
	}
	sort := bson.M{"$sort": bson.M{"merId": -1}}
	limit := bson.M{"$limit": 1}

	cond = append(cond, sort, limit)

	m := new(model.Merchant)
	err = database.C(c.name).Pipe(cond).One(m)
	if err != nil {
		log.Errorf("select maxMerId err,%s", err)
		return "", err
	}
	return m.MerId, nil
}

func (col *merchantCollection) FindCountByMerId(merId string) (num int, err error) {
	bo := bson.M{
		"merId": merId,
	}
	num, err = database.C(col.name).Find(bo).Count()
	if err != nil {
		log.Errorf("find count by merId err,merId=%s", err)
		return 0, err
	}
	return num, nil
}
