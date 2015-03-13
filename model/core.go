package model

import (
	"gopkg.in/mgo.v2/bson"
)

// RouterPolicy 路由策略
type RouterPolicy struct {
	MerId     string `json:"merId" bson:"merId,omitempty"`         // 商户号
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	ChanCode  string `json:"chanCode" bson:"chanCode,omitempty"`   // 渠道代码
	ChanMerId string `json:"chanMerId" bson:"chanMerId,omitempty"` // 渠道商户号
}

// BindingRelation 绑定关系
type BindingRelation struct {
	BindingId     string `json:"bindingId" bson:"bindingId,omitempty"`         // 银行卡绑定ID
	MerId         string `json:"merId" bson:"merId,omitempty"`                 // 商户ID
	AcctName      string `json:"acctName" bson:"acctName,omitempty"`           // 账户名称
	AcctNum       string `json:"acctNum" bson:"acctNum,omitempty"`             // 账户号码
	IdentType     string `json:"identType" bson:"identType,omitempty"`         // 证件类型
	IdentNum      string `json:"identNum" bson:"identNum,omitempty"`           // 证件号码
	PhoneNum      string `json:"phoneNum" bson:"phoneNum,omitempty"`           // 手机号
	AcctType      string `json:"acctType" bson:"acctType,omitempty"`           // 账户类型
	ValidDate     string `json:"validDate" bson:"validDate,omitempty"`         // 信用卡有限期
	Cvv2          string `json:"cvv2" bson:"cvv2,omitempty"`                   // CVV2
	BankId        string `json:"bankId" bson:"bankId,omitempty"`               // 银行ID
	CardBrand     string `json:"cardBrand" bson:"cardBrand,omitempty"`         // 卡品牌
	ChanCode      string `json:"chanCode" bson:"chanCode,omitempty"`           // 渠道代码
	ChanMerId     string `json:"chanMerId" bson:"chanMerId,omitempty"`         // 渠道商户号
	ChanBindingId string `json:"chanBindingId" bson:"chanBindingId,omitempty"` // 渠道绑定ID
	BindingStatus string `json:"bindingStatus" bson:"bindingStatus,omitempty"` // 绑定状态，成功，失败，或者处理中
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
