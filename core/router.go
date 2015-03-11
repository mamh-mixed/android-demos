package core

import (
	"quickpay/mongo"
	"strings"
)

// db.cardBin.find().forEach(
//     function(item) {
//             db.cardBin.update(
//                 {"_id": item._id},
//                 {"$set": {"overflow":  String(NumberLong(item.bin) + 1)}},
//                 false, true
//             )
//     }
// );
//
// db.cardBin.find().sort({"bin": -1}).skip(100);
//
// var card = "6222801932062061908";
// db.cardBin.find({"cardLen": card.length, "bin": {"$lte": card}, "overflow": {"$gt": card}}).sort({"bin": -1, "overflow": 1}).limit(1)
//

// 根据卡号判断是否是银联卡
func IsUnionPayCard(cardNum string) bool {
	cardBin := mongo.FindCardBin(cardNum)

	return strings.EqualFold("CUP", cardBin.CardBrand) || strings.EqualFold("UPI", cardBin.CardBrand)
}
