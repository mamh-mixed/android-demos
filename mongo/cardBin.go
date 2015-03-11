package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type CardBin struct {
	Bin       string `json:"bin" bson:"bin,omitempty"`             // 卡BIN
	BinLen    int    `json:"binLen" bson:"binLen,omitempty"`       // 卡BIN长度
	CardLen   int    `json:"cardLen" bson:"cardLen,omitempty"`     // 卡号长度
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
}

// 根据卡长度查找卡BIN列表
func FindCardBin(cardNum string) *CardBin {
	// var card = "6222801932062061908";
	// db.cardBin.find({"cardLen": card.length, "bin": {"$lte": card}, "overflow": {"$gt": card}}).sort({"bin": -1, "overflow": 1}).limit(1)

	var b CardBin
	db.cardBin.Find(bson.M{
		"cardLen":  len(cardNum),
		"bin":      bson.M{"$lte": cardNum},
		"overflow": bson.M{"$gt": cardNum},
	}).Sort("-bin", "overflow").Limit(1).One(&b)

	return &b
}
