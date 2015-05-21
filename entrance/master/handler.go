package master

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

const (
	NormalMerStatus = "Normal"
)

// AllMerchant 处理查找所有商户的请求
func AllMerchant(data []byte) (result *model.ResultBody) {
	cond := new(model.Merchant)
	err := json.Unmarshal(data, cond)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	merchants, err := mongo.MerchantColl.FindAllMerchant(cond)

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

// AddMerchant 处理新增一个商户的请求
func AddMerchant(data []byte) (result *model.ResultBody) {
	m := new(model.Merchant)
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if m.MerId == "" {
		log.Error("没有MerId")
		return model.NewResultBody(3, "缺失必要元素merId")
	}

	if m.MerStatus == "" {
		m.MerStatus = NormalMerStatus
	}

	err = mongo.MerchantColl.Insert(m)
	if err != nil {
		log.Errorf("新增商户失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    m,
	}

	return
}

// AddChannelMerchant 处理新增一个渠道商户的请求
func AddChannelMerchant(data []byte) (result *model.ResultBody) {
	m := new(model.ChanMer)
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if m.ChanCode == "" {
		log.Error("没有ChanCode")
		return model.NewResultBody(3, "缺失必要元素chanCode")
	}

	err = mongo.ChanMerColl.Add(m)
	if err != nil {
		log.Errorf("新增渠道商户失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    m,
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
