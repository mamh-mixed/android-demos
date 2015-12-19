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

// FindRate 查找一对货币的汇率
func (c *exchangeRateCollection) FindRate(localCurrency, targetCurrency string) (rate *model.ExchangeRate, err error) {
	rate, currencyPair := new(model.ExchangeRate), localCurrency+"<=>"+targetCurrency

	err = database.C(c.name).Find(bson.M{"currencyPair": currencyPair}).One(rate)
	return rate, err
}
