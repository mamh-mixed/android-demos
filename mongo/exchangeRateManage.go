package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"

	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

// 录入、复核汇率
type exchangeRateManageCollection struct {
	name string
}

var ExchangeRateManageColl = exchangeRateManageCollection{"exchangeRateManage"}

// PaginationFind 分页查找
func (c *exchangeRateManageCollection) PaginationFind(cond *model.ExchangeRateManage, size, page int) (results []model.ExchangeRateManage, total int, err error) {
	results = make([]model.ExchangeRateManage, 0)

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

	if cond.CheckedUser != "" {
		match["checkedUser"] = cond.CheckedUser
	}

	if cond.Status != "" {
		match["status"] = cond.Status
	}

	if cond.CreateTime != "" {
		match["createTime"] = bson.M{
			"$gt":  cond.CreateTime,
			"$lte": util.NextDay(cond.CreateTime),
		}
	}

	if cond.PlanEnforcementTime != "" {
		match["planEnforcementTime"] = bson.M{
			"$gt":  cond.PlanEnforcementTime,
			"$lte": util.NextDay(cond.PlanEnforcementTime),
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
func (c *exchangeRateManageCollection) FindOne(eId string) (rate *model.ExchangeRateManage, err error) {
	cond := bson.M{
		"eId": eId,
	}

	rate = new(model.ExchangeRateManage)
	err = database.C(c.name).Find(cond).One(rate)
	return rate, err
}

// Add 新增一条费率记录
func (c *exchangeRateManageCollection) Add(rate *model.ExchangeRateManage) (err error) {
	err = database.C(c.name).Insert(rate)
	return err
}

// Update 更新
func (c *exchangeRateManageCollection) Update(rate *model.ExchangeRateManage) (err error) {
	b := bson.M{
		"eId": rate.EId,
	}

	err = database.C(c.name).Update(b, rate)
	return err
}
