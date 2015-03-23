package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"

	"gopkg.in/mgo.v2/bson"
)

type cardBinCollection struct {
	name string
}

// CardBinColl 卡Bin Collection
var CardBinColl = cardBinCollection{"cardBin"}

// Find 根据卡长度查找卡BIN列表
func (c *cardBinCollection) Find(cardNum string) (cb *model.CardBin, err error) {
	cb = new(model.CardBin)
	q := bson.M{
		"cardLen":  len(cardNum),
		"bin":      bson.M{"$lte": cardNum},
		"overflow": bson.M{"$gt": cardNum},
	}
	err = database.C(c.name).Find(q).Sort("-bin", "overflow").Limit(1).One(&cb)
	if err != nil {
		g.Error("Find CardBin ERROR! error message is: %s; condition is: %+v", err.Error(), q)
		return nil, err
	}

	return cb, nil
}
