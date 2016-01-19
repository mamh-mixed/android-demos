package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

type routerPolicy struct{}

var RouterPolicy routerPolicy

// Find 查找路由列表，参数是 merId。
func (i *routerPolicy) Find(merId string, cardBrand string, chanCode string, chanMerId string, pay string, size, page int) (result *model.ResultBody) {
	log.Debugf("merId=%s; cardBrand=%s; chanCode=%s;chanMerId=%s", merId, cardBrand, chanCode, chanMerId)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size <= 0 {
		size = 10
	}

	routers, total, err := mongo.RouterPolicyColl.PaginationFind(merId, cardBrand, chanCode, chanMerId, pay, size, page)
	if err != nil {
		log.Errorf("查询商户(%s)的所有路由失败: %s", merId, err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(routers),
		Data:  routers,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// FindOne 查找路由列表，参数是 merId 和 cardBrand。
func (i *routerPolicy) FindOne(merId, cardBrand string) (result *model.ResultBody) {
	router := mongo.RouterPolicyColl.Find(merId, cardBrand)

	if router == nil {
		// log.Errorf("未找到商户(%s)的一个路由(%s)", merId, cardBrand)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    router,
	}

	return result
}

// Save 新增一个路由策略
func (i *routerPolicy) Save(data []byte) (result *model.ResultBody) {
	r := new(model.RouterPolicy)
	err := json.Unmarshal(data, r)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if r.MerId == "" {
		log.Error("MerId")
		return model.NewResultBody(3, "缺失必要元素 merId")
	}

	if r.ChanCode == "" {
		log.Error("没有 ChanCode")
		return model.NewResultBody(3, "缺失必要元素 chanCode")
	}

	if r.ChanMerId == "" {
		log.Error("没有 ChanMerId")
		return model.NewResultBody(3, "缺失必要元素 chanMerId")
	}

	if r.CardBrand == "" {
		log.Error("没有 CardBrand")
		return model.NewResultBody(3, "缺失必要元素 cardBrand")
	}

	merchant, err := mongo.MerchantColl.FindNotInCache(r.MerId)
	if err != nil {
		if err.Error() == "not found" {
			return model.NewResultBody(4, "merId不存在")
		} else {
			return model.NewResultBody(1, "查询数据库失败")
		}
	}
	// 对清算标识与清算角色做校验
	if r.SettFlag == model.SR_AGENT {
		if r.SettRole != merchant.AgentCode {
			return model.NewResultBody(5, "agentCode错误")
		}
	} else if r.SettFlag == model.SR_COMPANY {
		if r.SettRole != merchant.SubAgentCode {
			return model.NewResultBody(5, "subAgentCode错误")
		}
	} else if r.SettFlag == model.SR_GROUP {
		if r.SettRole != merchant.GroupCode {
			return model.NewResultBody(5, "groupCode错误")
		}
	} else if r.SettFlag == model.SR_CIL {
		if r.SettRole != "CIL" {
			return model.NewResultBody(5, "清算标识与清算角色不匹配")
		}
	} else if r.SettFlag == model.SR_CHANNEL {
		if r.SettFlag == "ALP" && r.SettRole != "ALP" {
			return model.NewResultBody(5, "清算标识与清算角色不匹配")
		} else if r.SettFlag == "WXP" && r.SettRole != "WXP" {
			return model.NewResultBody(5, "清算标识与清算角色不匹配")
		}
	}

	err = mongo.RouterPolicyColl.Insert(r)
	if err != nil {
		log.Errorf("保存路由信息失败:%s", err)
		return model.NewResultBody(1, "保存路由信息失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "保存成功",
		Data:    r,
	}

	return
}

// Delete 删除路由，参数是 merId，chanCode，cardBrand。
func (i *routerPolicy) Delete(merId, chanCode, cardBrand string) (result *model.ResultBody) {

	err := mongo.RouterPolicyColl.Remove(merId, chanCode, cardBrand)

	if err != nil {
		log.Errorf("删除路由失败: %s", err)
		return model.NewResultBody(1, "删除路由失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return result
}

func (i *routerPolicy) Update(data []byte) (result *model.ResultBody) {
	r := new(model.RouterPolicy)
	err := json.Unmarshal(data, r)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if r.MerId == "" {
		log.Error("MerId")
		return model.NewResultBody(3, "缺失必要元素 merId")
	}

	if r.ChanCode == "" {
		log.Error("没有 ChanCode")
		return model.NewResultBody(3, "缺失必要元素 chanCode")
	}

	if r.ChanMerId == "" {
		log.Error("没有 ChanMerId")
		return model.NewResultBody(3, "缺失必要元素 chanMerId")
	}

	if r.CardBrand == "" {
		log.Error("没有 CardBrand")
		return model.NewResultBody(3, "缺失必要元素 cardBrand")
	}

	merchant, err := mongo.MerchantColl.FindNotInCache(r.MerId)
	if err != nil {
		if err.Error() == "not found" {
			return model.NewResultBody(4, "merId不存在")
		} else {
			return model.NewResultBody(1, "查询数据库失败")
		}
	}
	// 对清算标识与清算角色做校验
	if r.SettFlag == model.SR_AGENT {
		if r.SettRole != merchant.AgentCode {
			return model.NewResultBody(5, "agentCode错误")
		}
	} else if r.SettFlag == model.SR_COMPANY {
		if r.SettRole != merchant.SubAgentCode {
			return model.NewResultBody(5, "subAgentCode错误")
		}
	} else if r.SettFlag == model.SR_GROUP {
		if r.SettRole != merchant.GroupCode {
			return model.NewResultBody(5, "groupCode错误")
		}
	} else if r.SettFlag == model.SR_CIL {
		if r.SettRole != "CIL" {
			return model.NewResultBody(5, "清算标识与清算角色不匹配")
		}
	} else if r.SettFlag == model.SR_CHANNEL {
		if r.SettFlag == "ALP" && r.SettRole != "ALP" {
			return model.NewResultBody(5, "清算标识与清算角色不匹配")
		} else if r.SettFlag == "WXP" && r.SettRole != "WXP" {
			return model.NewResultBody(5, "清算标识与清算角色不匹配")
		}
	}

	err = mongo.RouterPolicyColl.Update(r)
	if err != nil {
		log.Errorf("保存路由信息失败:%s", err)
		return model.NewResultBody(1, "保存路由信息失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "保存成功",
		Data:    r,
	}

	return
}
