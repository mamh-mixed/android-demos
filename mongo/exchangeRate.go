package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type exchangeRateCollection struct {
	name string
}

var ExchangeRateColl = exchangeRateCollection{"exchangeRate"}

// Upsert 新增或者更新
func (c *exchangeRateCollection) Upsert(cond *model.ExchangeRate) (err error) {
	_, err = database.C(c.name).Upsert(bson.M{"currencyPair": cond.CurrencyPair}, cond)

	return err
}

// FindOne 查询一对货币中的币种
func (c *exchangeRateCollection) FindOne(currencyPair string) (rate *model.ExchangeRate, err error) {
	rate = new(model.ExchangeRate)

	err = database.C(c.name).Find(bson.M{"currencyPair": currencyPair}).One(rate)

	return rate, err
}
