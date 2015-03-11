package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type ChanMer struct {
	ChanCode       string //渠道代码
	ChanMerId      string //商户号
	ChanMerName    string //商户名称
	SettlementFlag string //清算标识
	SettlementRole string //清算角色
	SignCert       string //签名证书
	CheckSignCert  string //验签证书
	//...
}

// Init 根据渠道代码、商户号初始化
func (c *ChanMer) Init() error {

	bo := bson.M{
		"chancode":  c.ChanCode,
		"chanmerid": c.ChanMerId,
	}
	return db.chanMer.Find(bo).One(c)

}

// Add 增加一个渠道商户
func (c *ChanMer) Add() error {
	return db.chanMer.Insert(c)
}
