package mongo

import (
	// "errors"
	"quickpay/model"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

// FindRouter 根据源商户号和卡品牌在数据库中查找路由策略 [moved]
func FindRouter(origMerId, cardBrand string) (rp *model.RouterPolicy, err error) {
	rp = &model.RouterPolicy{}
	q := bson.M{"origMerId": origMerId, "cardBrand": cardBrand}
	g.Debug("'FindRouter' condition: %+v", q)

	err = db.routerPolicy.Find(q).One(rp)
	if err != nil {
		g.Debug("Error message is ", err.Error())
		return nil, err
	}
	g.Debug("Router Policy is %+v", rp)

	return rp, nil
}

// InsertRouterPolicy 插入一个路由策略到数据库中
func InsertRouterPolicy(rp *model.RouterPolicy) error {
	if err := db.routerPolicy.Insert(rp); err != nil {
		return err
	}
	return nil
}

// FindRouterPolicy 根据源商户Id 和 卡品牌查找路由
func FindRouterPolicy(origMerId, cardBrand string) (r *model.RouterPolicy) {
	r = &model.RouterPolicy{}
	q := bson.M{"origMerId": origMerId, "cardBrand": cardBrand}
	db.routerPolicy.Find(q).One(r)

	g.Debug("'FindRouter' condition: %+v, result %#v", q, r)
	return r
}
