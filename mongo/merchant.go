package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type merchantCollection struct {
	name string
}

var MerchantColl = merchantCollection{"merchant"}

var merCache = cache.New(model.Cache_Merchant)

// Upsert 插入一个商户信息。如果存在则更新，不存在则插入。@WonSikin
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

// FindNoUniqueId 查找商户信息
func (c *merchantCollection) FindNoUniqueId() ([]*model.Merchant, error) {
	var ms []*model.Merchant
	q := bson.M{"uniqueId": bson.M{"$exists": false}}
	err := database.C(c.name).Find(q).All(&ms)
	return ms, err
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

// FindNotCache 根据merId查找商户信息，非缓存
func (c *merchantCollection) FindNotInCache(merId string) (m *model.Merchant, err error) {

	m = new(model.Merchant)
	q := bson.M{"merId": merId}
	err = database.C(c.name).Find(q).One(m)
	if err != nil {
		return nil, err
	}

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
	m.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
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
func (c *merchantCollection) PaginationFind(merchant model.Merchant, pay, createStartTime, createEndTime string, size, page int) (results []*model.Merchant, total int, err error) {
	results = make([]*model.Merchant, 1)

	match := bson.M{}
	if merchant.MerId != "" {
		match["merId"] = bson.RegEx{merchant.MerId, "i"}
	}
	if merchant.MerStatus != "" {
		match["merStatus"] = merchant.MerStatus
	}
	if merchant.Detail.MerName != "" {
		match["merDetail.merName"] = bson.RegEx{merchant.Detail.MerName, "i"}
	}
	if merchant.AgentCode != "" {
		match["agentCode"] = bson.RegEx{merchant.AgentCode, "i"}
	}
	if merchant.AgentName != "" {
		match["agentName"] = bson.RegEx{merchant.AgentName, "i"}
	}
	if merchant.SubAgentCode != "" {
		match["subAgentCode"] = bson.RegEx{merchant.SubAgentCode, "i"}
	}
	if merchant.SubAgentName != "" {
		match["subAgentName"] = bson.RegEx{merchant.SubAgentName, "i"}
	}
	if merchant.GroupCode != "" {
		match["groupCode"] = bson.RegEx{merchant.GroupCode, "i"}
	}
	if merchant.GroupName != "" {
		match["groupName"] = bson.RegEx{merchant.GroupName, "i"}
	}
	if merchant.Detail.CommodityName != "" {
		match["merDetail.commodityName"] = bson.RegEx{merchant.Detail.CommodityName, "i"}
	}
	if merchant.IsNeedSign == true {
		match["isNeedSign"] = merchant.IsNeedSign
	}

	if merchant.Detail.AcctNum != "" {
		match["merDetail.acctNum"] = bson.RegEx{merchant.Detail.AcctNum, "i"}
	}
	if merchant.Detail.GoodsTag != "" {
		match["merDetail.goodsTag"] = bson.RegEx{merchant.Detail.GoodsTag, "i"}
	}

	if pay == "bp" {
		match["encryptKey"] = bson.M{"$exists": true}
	} else {
		match["encryptKey"] = bson.M{"$exists": false}
	}

	if createStartTime != "" && createEndTime != "" {
		match["createTime"] = bson.M{"$gte": createStartTime, "$lte": createEndTime}
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
		m.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		m.UpdateTime = m.CreateTime
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

// Insert 创建一个机构商户
func (c *merchantCollection) Insert(m *model.Merchant) error {

	m.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = m.CreateTime
	err := database.C(c.name).Insert(m)
	if err != nil {
		log.Errorf("'Insert Merchant ERROR!' Merchant is (%+v); error is (%s)", m, err)
		return err
	}
	return nil
}

// findMaxMerId 查询merId最大值
func (c *merchantCollection) FindMaxMerId(prefix string) (merId string, err error) {

	m := new(model.Merchant)
	length := fmt.Sprintf("%d", 15-len(prefix))
	regex := bson.RegEx{"^" + prefix + "\\d{" + length + "}$", ""}
	query := bson.M{"merId": regex}

	err = database.C(c.name).Find(query).Sort("-merId").Select(bson.M{"merId": 1}).One(m)

	return m.MerId, err
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
