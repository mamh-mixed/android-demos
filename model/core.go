package model

// RouterPolicy 路由策略
type RouterPolicy struct {
	OrigMerId string `json:"origMerId" bson:"origMerId,omitempty"` // 源商户号
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	ChanCode  string `json:"chanCode" bson:"chanCode,omitempty"`   // 渠道代码
	ChanMerId string `json:"chanMerId" bson:"chanMerId,omitempty"` // 渠道商户号
}
