package mongo

import (
	// "errors"
	"quickpay/model"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

type routerPolicyCollection struct {
	name string
}

var RouterPolicyColl = routerPolicyCollection{"routerPolicy"}

// Insert 插入一个路由策略到数据库中
func (c *routerPolicyCollection) Insert(rp *model.RouterPolicy) error {
	if err := database.C(c.name).Insert(rp); err != nil {
		return err
	}
	return nil
}

// Find 根据源商户Id 和 卡品牌查找路由
func (c *routerPolicyCollection) Find(merId, cardBrand string) (r *model.RouterPolicy) {
	r = &model.RouterPolicy{}
	q := bson.M{"merId": merId, "cardBrand": cardBrand}
	err := database.C(c.name).Find(q).One(r)
	if err != nil {
		g.Debug("'FindRouter' condition: %+v\n", q)
		g.Debug("Error message is: \n", err.Error())
		return nil
	}
	g.Debug("'FindRouter' condition: %+v, result %#v", q, r)
	return r
}
