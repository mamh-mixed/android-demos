package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"

	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type exchangeRateCollection struct {
	name string
}

var ExchangeRateColl = exchangeRateCollection{"exchangeRate"}

// PaginationFind 分页查找
func (c *exchangeRateCollection) PaginationFind(cond *model.ExchangeRate, size, page int) (results []model.ExchangeRate, total int, err error) {
	results = make([]model.ExchangeRate, 0)

	match := bson.M{}

	if cond.LocalCurrency != "" {
		match["localCurrency"] = cond.LocalCurrency
	}

	if cond.TargetCurrency != "" {
		match["targetCurrency"] = cond.TargetCurrency
	}

	if cond.Rate != 0.0 {
		match["rate"] = cond.Rate
	}

	if cond.EnforceUser != "" {
		match["enforceUser"] = cond.EnforceUser
	}

	if cond.IsEnforced {
		match["isEnforced"] = cond.IsEnforced
	}

	if cond.CreateTime != "" {
		match["createTime"] = bson.M{
			"$gt":  cond.CreateTime,
			"$lte": util.NextDay(cond.CreateTime),
		}
	}

	if cond.ActualEnforcementTime != "" {
		match["actualEnforcementTime"] = bson.M{
			"$gt":  cond.ActualEnforcementTime,
			"$lte": util.NextDay(cond.ActualEnforcementTime),
		}
	}

	log.Debugf("match is %#v", match)

	total, err = database.C(c.name).Find(match).Count()
	log.Debugf("total is %d, error is %s", total, err)
	if err != nil {
		return nil, 0, err
	}

	q := []bson.M{
		{"$match": match},
	}

	sort := bson.M{"$sort": bson.M{"localCurrency": 1, "targetCurrency": 1, "updateTime": -1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	q = append(q, sort, skip, limit)

	err = database.C(c.name).Pipe(q).All(&results)

	return results, total, err
}

// FindOne 查找一个
func (c *exchangeRateCollection) FindOne(eId string) (rate *model.ExchangeRate, err error) {
	cond := bson.M{
		"eId": eId,
	}

	rate = new(model.ExchangeRate)
	err = database.C(c.name).Find(cond).One(rate)
	return rate, err
}

// Add 新增一条费率记录
func (c *exchangeRateCollection) Add(rate *model.ExchangeRate) (err error) {
	err = database.C(c.name).Insert(rate)
	return err
}

// Update 更新
func (c *exchangeRateCollection) Update(rate *model.ExchangeRate) (err error) {
	b := bson.M{
		"eId": rate.EId,
	}

	err = database.C(c.name).Update(b, rate)
	return err
}
