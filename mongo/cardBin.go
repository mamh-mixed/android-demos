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
func FindByCardLen(cardLen int) *[]CardBin {
	result := &[]CardBin{}
	db.cardBin.Find(bson.M{"cardLen": cardLen}).All(result)
	return result
}
