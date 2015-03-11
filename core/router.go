package core

import (
	"github.com/omigo/g"
	"quickpay/mongo"
	"regexp"
	"strings"
)

// 根据卡号查找卡品牌，首先找出卡号长度匹配的卡BIN数组，然后匹配卡BIN，最后把卡BIN长度最长的返回
func FindCardBrandByCardNum(cardNum string) (cardBrand string) {
	g.Info("card number: %s", cardNum)

	// 根据长度查找出来的卡BIN列表
	tempArrays := mongo.FindByCardLen(len(cardNum))

	if len(*tempArrays) == 0 {
		return ""
	}

	binMatchedArrays := []mongo.CardBin{}
	// 遍历上一步查找出来的卡BIN列表，找出BIN值匹配的所有卡BIN，放到 binMatchedArrays 列表中
	for _, item := range *tempArrays {
		bin := item.Bin
		if matched, _ := regexp.MatchString(`^`+bin, cardNum); matched {
			binMatchedArrays = append(binMatchedArrays, item)
		}
	}

	if len(binMatchedArrays) == 1 {
		return binMatchedArrays[0].CardBrand
	}

	flag := 0
	for _, item := range binMatchedArrays {
		if item.BinLen > flag {
			flag, cardBrand = item.BinLen, item.CardBrand
		}
	}
	g.Info("card brand: %s", cardBrand)

	return cardBrand
}

// 根据卡号判断是否是银联卡
func IsUnionPayCard(cardNum, cardBrand string) bool {
	if cardBrand == "" {
		cardBrand = FindCardBrandByCardNum(cardNum)
	}
	return strings.EqualFold("CUP", cardBrand) || strings.EqualFold("UPI", cardBrand)
}
