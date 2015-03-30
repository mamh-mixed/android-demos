package mongo

import (
	"errors"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type merchantCollection struct {
	name string
}

var MerchantColl = merchantCollection{"merchant"}

func (c *merchantCollection) Insert(m *model.Merchant) error {
	m1 := new(model.Merchant)
	q := bson.M{"merId": m.MerId}
	err := database.C(c.name).Find(q).One(m1)
	if err == nil {
		return errors.New("MerId is existed!")
	}

	err = database.C(c.name).Insert(m)
	if err != nil {
		log.Errorf("'Insert Merchant ERROR!' Merchant is (%+v);error is (%s)", m, err)
	}
	return err
}

func (c *merchantCollection) Find(merId string) (m *model.Merchant, err error) {
	m = new(model.Merchant)
	q := bson.M{"merId": merId}
	err = database.C(c.name).Find(q).One(m)
	if err != nil {
		log.Errorf("'Find Merchant ERROR!' Condition is (%+v);error is(%s)", q, err)
		return nil, err
	}
	return m, nil
}

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
