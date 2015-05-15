package mongo

import (
	// "errors"
	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type routerPolicyCollection struct {
	name string
}

var RouterPolicyColl = routerPolicyCollection{"routerPolicy"}

// Insert 插入一个路由策略到数据库中，如果路由中已经存在一模一样的，就更新
func (c *routerPolicyCollection) Insert(rp *model.RouterPolicy) error {
	if _, err := database.C(c.name).Upsert(rp, rp); err != nil {
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
		log.Errorf("FindRouter Error message is: %s", err)
		return nil
	}
	return r
}

// FindAllOfOneMerchant 根据源商户Id查找该商户下的所有路由信息
func (c *routerPolicyCollection) FindAllOfOneMerchant(merId string) (r []model.RouterPolicy, err error) {
	r = make([]model.RouterPolicy, 0)
	q := bson.M{"merId": merId}

	err = database.C(c.name).Find(q).All(&r)
	if err != nil {
		log.Errorf("FindAllOfOneMerchant Error message is: %s\n", err)
		return nil, err
	}
	return r, nil
}
