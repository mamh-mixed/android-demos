package model

import (
	"gopkg.in/mgo.v2/bson"
)

// RouterPolicy 路由策略
type RouterPolicy struct {
	OrigMerId string `json:"origMerId" bson:"origMerId,omitempty"` // 源商户号
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	ChanCode  string `json:"chanCode" bson:"chanCode,omitempty"`   // 渠道代码
	ChanMerId string `json:"chanMerId" bson:"chanMerId,omitempty"` // 渠道商户号
}

// BindingRelation 绑定关系
type BindingRelation struct {
	BindingCreate `json:"cardInfo" bson:"cardInfo,omitempty,inline"` //卡片信息
	RouterPolicy  `json:"router" bson:"router,omitempty,inline"`     //路由信息
	ChanBindingId string                                             `json:"chanBindingId" bson:"chanBindingId,omitempty"` //渠道绑定ID
}

// 卡BIN
type CardBin struct {
	Bin       string `json:"bin" bson:"bin,omitempty"`             // 卡BIN
	BinLen    int    `json:"binLen" bson:"binLen,omitempty"`       // 卡BIN长度
	CardLen   int    `json:"cardLen" bson:"cardLen,omitempty"`     // 卡号长度
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
}

// 渠道商户
type ChanMer struct {
	ChanCode      string //渠道代码
	ChanMerId     string //商户号
	ChanMerName   string //商户名称
	SettFlag      string //清算标识
	SettRole      string //清算角色
	SignCert      string //签名证书
	CheckSignCert string //验签证书
	//...
}
type Trans struct {
	Id      bson.ObjectId  `bson:"_id"`
	Chan    ChanMer        //渠道信息
	Payment BindingPayment //支付信息
	Time    int64          //时间
	Flag    int8           //交易状态
}
