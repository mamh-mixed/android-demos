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
func (c *merchantCollection) Insert(m *model.Merchant) error {
	q := bson.M{"merId": m.MerId}

	_, err := database.C(c.name).Upsert(q, m)
	if err != nil {
		log.Errorf("'Upsert Merchant ERROR!' Merchant is (%+v); error is (%s)", m, err)
	}
	return err
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
		log.Errorf("'Find Merchant ERROR!' Condition is (%+v);error is(%s)", q, err)
		return nil, err
	}
	// save
	merCache.Set(merId, m, cache.DefaultExpiration)

	return m, nil
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
	}

	var mers []*model.Merchant
	err := database.C(c.name).Find(q).Select(bson.M{"merId": 1, "merDetail.merName": 1, "agentName": 1}).All(&mers)
	return mers, err
}

// PaginationFind 分页查找机构商户
func (c *merchantCollection) PaginationFind(merId, merStatus string, size, page int) (results []model.Merchant, total int, err error) {
	results = make([]model.Merchant, 1)

	match := bson.M{}
	if merId != "" {
		match["merId"] = merId
	}
	if merStatus != "" {
		match["merStatus"] = merStatus
	}

	// 计算总数
	total, err = database.C(c.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, skip, limit)

	err = database.C(c.name).Pipe(cond).All(&results)

	return results, total, err
}
