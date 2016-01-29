package mongo

import (
	// "errors"
	"time"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type routerPolicyCollection struct {
	name string
}

var RouterPolicyColl = routerPolicyCollection{"routerPolicy"}

// Insert 插入一个路由策略到数据库中，如果路由中已经存在一模一样的，就更新
func (c *routerPolicyCollection) Insert(rp *model.RouterPolicy) error {
	rp.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	cond := bson.M{
		"merId":     rp.MerId,
		"cardBrand": rp.CardBrand,
		"chanCode":  rp.ChanCode,
	}
	if _, err := database.C(c.name).Upsert(cond, rp); err != nil {
		return err
	}
	return nil
}

// BatchAdd 批量增加路由策略
func (c *routerPolicyCollection) BatchAdd(routers []model.RouterPolicy) error {
	var temp []interface{}
	for _, r := range routers {
		r.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		r.UpdateTime = r.CreateTime
		temp = append(temp, r)
	}
	return database.C(c.name).Insert(temp...)
}

// BatchRemove 批量删除路由策略
func (c *routerPolicyCollection) BatchRemove(routers []model.RouterPolicy) error {
	var rs []bson.M
	for _, r := range routers {
		rs = append(rs, bson.M{"merId": r.MerId, "cardBrand": r.CardBrand})
	}
	selector := bson.M{
		"$in": rs,
	}
	change, err := database.C(c.name).RemoveAll(selector)
	if change.Removed != len(routers) {
		log.Warnf("expect remove %d records,but %d removed", len(routers), change.Removed)
	}
	return err
}

// Find 根据源商户Id 和 卡品牌查找路由
func (c *routerPolicyCollection) Find(merId, cardBrand string) (r *model.RouterPolicy) {
	r = &model.RouterPolicy{}
	q := bson.M{"merId": merId, "cardBrand": cardBrand}
	err := database.C(c.name).Find(q).One(r)
	if err != nil {
		// log.Errorf("FindRouter Error message is: %s", err)
		return nil
	}
	return r
}

// PaginationFind 分页查找
func (c *routerPolicyCollection) PaginationFind(merId string, cardBrand string, chanCode string, chanMerId string, pay string, size, page int) (results []model.RouterPolicy, total int, err error) {
	results = make([]model.RouterPolicy, 0)

	match := bson.M{}
	if pay == "bp" {
		match["chanCode"] = bson.M{"$in": []string{"CFCA", "CIL", "Mock"}}
	} else {
		match["chanCode"] = bson.M{"$in": []string{"ALP", "WXP", "ULIVE"}}
	}

	if chanCode != "" {
		match["$and"] = []interface{}{bson.M{"chanCode": chanCode}}
	}

	if merId != "" {
		match["merId"] = merId
	}
	if cardBrand != "" {
		match["cardBrand"] = cardBrand
	}

	if chanMerId != "" {
		match["chanMerId"] = chanMerId
	}

	// 计算总数
	total, err = database.C(c.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}
	sort := bson.M{"$sort": bson.M{"merId": 1}}

	cond := []bson.M{
		{"$match": match},
	}

	skip := bson.M{
		"$skip": (page - 1) * size,
	}

	limit := bson.M{
		"$limit": size,
	}

	cond = append(cond, sort, skip, limit)

	err = database.C(c.name).Pipe(cond).All(&results)

	return results, total, err
}

// FindAllOfOneMerchant 根据源商户Id查找该商户下的所有路由信息
func (c *routerPolicyCollection) FindAllOfOneMerchant(merId string) (r []model.RouterPolicy, err error) {
	r = make([]model.RouterPolicy, 0)
	q := bson.M{"merId": merId}

	if merId == "" {
		q = nil
	}

	err = database.C(c.name).Find(q).All(&r)
	if err != nil {
		log.Errorf("FindAllOfOneMerchant Error message is: %s\n", err)
		return nil, err
	}
	return r, nil
}

// Remove 删除路由策略
func (c *routerPolicyCollection) Remove(merId, chanCode, cardBrand string) (err error) {
	q := bson.M{}
	if merId != "" {
		q["merId"] = merId
	}
	if cardBrand != "" {
		q["cardBrand"] = cardBrand
	}
	if chanCode != "" {
		q["chanCode"] = chanCode
	}

	err = database.C(c.name).Remove(q)

	return err
}

func (c *routerPolicyCollection) Update(rp *model.RouterPolicy) error {
	rp.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	cond := bson.M{
		"merId":     rp.MerId,
		"cardBrand": rp.CardBrand,
		"chanCode":  rp.ChanCode,
	}
	return database.C(c.name).Update(cond, rp)
}

// Find 根据渠道商户Id和渠道chanCode获取路由信息
func (c *routerPolicyCollection) FindByChanMerId(chanMerId, chanCode string) (r *model.RouterPolicy) {
	r = &model.RouterPolicy{}
	q := bson.M{"chanMerId": chanMerId, "chanCode": chanCode}
	err := database.C(c.name).Find(q).One(r)
	if err != nil {
		// log.Errorf("FindRouter Error message is: %s", err)
		return nil
	}
	return r
}
