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

type Trans struct {
	Id      bson.ObjectId  `bson:"_id"`
	Chan    ChanMer        //渠道信息
	Payment BindingPayment //支付信息
	Time    int64          //时间
	Flag    int8           //交易状态
}

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
