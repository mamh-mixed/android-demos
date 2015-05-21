package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type settSchemeCdColletion struct {
	name string
}

var SettSchemeCdCol = &settSchemeCdColletion{"settSchemeCd"}

// Upsert
func (col *settSchemeCdColletion) Upsert(update *model.SettSchemeCd) error {

	selector := bson.M{
		"schemeCd": update.SchemeCd,
	}
	_, err := database.C(col.name).Upsert(selector, update)
	return err
}

// Find 计费方案
func (col *settSchemeCdColletion) Find(cd string) (*model.SettSchemeCd, error) {

	result := new(model.SettSchemeCd)
	q := bson.M{
		"schemeCd": cd,
	}
	err := database.C(col.name).Find(q).One(result)
	return result, err
}
