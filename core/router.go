package core

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/CardInfoLink/log"
)

// FindRouter 根据源商户号和卡品牌在数据库中查找路由策略
func FindRouter(origMerId, cardBrand string) (r *model.RouterPolicy, err error) {

	r = mongo.RouterPolicyColl.Find(origMerId, cardBrand)
	if err != nil {
		log.Debugf("Error message is ", err.Error())
		return nil, err
	}
	log.Debugf("Router Policy is %+v", r)

	return r, nil
}
