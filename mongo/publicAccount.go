package mongo

import (
	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

var PulicAccountCol = &publicAccountCollecton{"publicAccount"}

var paCache = cache.New("public_account")

type publicAccountCollecton struct {
	name string
}

// Add 添加一条公众号信息
func (p *publicAccountCollecton) Add(a *model.PublicAccount) error {
	return database.C(p.name).Insert(a)
}

// Get
func (p *publicAccountCollecton) Get(chanMerId string) (a *model.PublicAccount, err error) {

	// get from cache
	o, found := paCache.Get(chanMerId)
	if found {
		a = o.(*model.PublicAccount)
		return a, nil
	}

	a = new(model.PublicAccount)
	err = database.C(p.name).Find(bson.M{"chanMerId": chanMerId}).One(a)

	if err != nil {
		return nil, err
	}

	paCache.Set(chanMerId, a, cache.DefaultExpiration)

	return a, err
}
