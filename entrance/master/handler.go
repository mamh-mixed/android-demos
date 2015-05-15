package master

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// AllMerchant 处理查找所有商户的请求
func AllMerchant() (result *model.ResultBody) {
	merchants, err := mongo.MerchantColl.FindAllMerchant()

	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    merchants,
	}

	return
}

// AllRouterOfOneMerchant 处理查找商户的所有路由的请求
func AllRouterOfOneMerchant(merId string) (result *model.ResultBody) {
	routers, err := mongo.RouterPolicyColl.FindAllOfOneMerchant(merId)

	if err != nil {
		log.Errorf("查询商户(%s)的所有路由失败: %s", merId, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    routers,
	}

	return
}
