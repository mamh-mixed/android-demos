package core

import (
	"quickpay/model"
	"quickpay/mongo"

	"github.com/omigo/g"
)

// FindRouter 根据源商户号和卡品牌在数据库中查找路由策略
func FindRouter(origMerId, cardBrand string) (r *model.RouterPolicy, err error) {

	r = mongo.RouterPolicyColl.Find(origMerId, cardBrand)
	if err != nil {
		g.Debug("Error message is ", err.Error())
		return nil, err
	}
	g.Debug("Router Policy is %+v", r)

	return r, nil
}
