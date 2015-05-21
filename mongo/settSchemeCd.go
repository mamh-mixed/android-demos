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
