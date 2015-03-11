package mongo

import (
	// "errors"
	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

type RouterPolicy struct {
	OrigMerCode    string `json:"origMerCode" bson:"origMerCode,omitempty"`       // 源商户号
	CardBrand      string `json:"cardBrand" bson:"cardBrand,omitempty"`           // 卡品牌
	ChannelCode    string `json:"channelCode" bson:"channelCode,omitempty"`       // 渠道代码
	ChannelMerCode string `json:"channelMerCode" bson:"channelMerCode,omitempty"` // 渠道商户号
}

// 根据源商户号和卡品牌在数据库中查找路由策略
func FindRouter(origMerCode, cardBrand string) (rp *RouterPolicy, err error) {
	rp = new(RouterPolicy)
	err = db.routerPolicy.Find(bson.M{"origMerCode": origMerCode, "cardBrand": cardBrand}).One(rp)
	if err != nil {
		g.Debug("Error message is ", err.Error())
		return nil, err
	}
	g.Debug("Router Policy is %+v", rp)
	return rp, nil
}

// 插入一个路由策略到数据库中
func InsertOneRouterPolicy(rp *RouterPolicy) error {
	if err := db.routerPolicy.Insert(rp); err != nil {
		return err
	}
	return nil
}
