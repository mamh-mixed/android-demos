package mongo

import (
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

// 根据卡长度查找卡BIN列表
func FindCardBin(cardNum string) *model.CardBin {
	// var card = "6222801932062061908";
	// db.cardBin.find({"cardLen": card.length, "bin": {"$lte": card}, "overflow": {"$gt": card}}).sort({"bin": -1, "overflow": 1}).limit(1)

	var b model.CardBin
	db.cardBin.Find(bson.M{
		"cardLen":  len(cardNum),
		"bin":      bson.M{"$lte": cardNum},
		"overflow": bson.M{"$gt": cardNum},
	}).Sort("-bin", "overflow").Limit(1).One(&b)

	return &b
}
