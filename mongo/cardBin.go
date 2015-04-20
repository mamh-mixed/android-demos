package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

// CardBinColl 卡Bin Collection
var CardBinColl = cardBinCollection{"cardBin"}

type cardBinCollection struct {
	name string
}

// Find 根据卡长度查找卡BIN列表
func (c *cardBinCollection) Find(cardBin string, length int) (cb *model.CardBin, err error) {
	cb = new(model.CardBin)
	// q := bson.M{
	// 	"cardLen":  len(cardNum),
	// 	"bin":      bson.M{"$lte": cardNum},
	// 	"overflow": bson.M{"$gt": cardNum},
	// }
	// err = database.C(c.name).Find(q).Sort("-bin", "overflow").Limit(1).One(&cb)
	// if err != nil {
	// 	log.Errorf("Find CardBin ERROR! error message is: %s; condition is: %+v", err.Error(), q)
	// 	return nil, err
	// }
	// return cb, err
	q := bson.M{
		"bin":     cardBin,
		"cardLen": length,
	}
	err = database.C(c.name).Find(q).One(cb)

	return
}

// LoadAll 加载所有卡bin
func (c *cardBinCollection) LoadAll() ([]*model.CardBin, error) {
	var cardBins []*model.CardBin
	err := database.C(c.name).Find(nil).All(&cardBins)
	return cardBins, err
}
