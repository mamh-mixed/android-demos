package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type merchant struct{}

var Merchant merchant

// Find 根据条件分页查找商户。
func (m *merchant) FindOne(merId string) (result *model.ResultBody) {
	log.Debugf("merId is %s", merId)

	merchant, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		log.Errorf("查询一个商户(%s)出错: %s", merId, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    merchant,
	}

	return result
}

// Find 根据条件分页查找商户。
func (m *merchant) Find(merId, merStatus, merName, groupCode, groupName, agentCode, agentName string, size, page int) (result *model.ResultBody) {
	log.Debugf("merId is %s; merName is %s;groupCode is %s, groupName is %s, agentCode is %s, agentName is %s",
		merId, merStatus, merName, groupCode, groupName, agentCode, agentName)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	merchants, total, err := mongo.MerchantColl.PaginationFind(merId, merStatus, merName, groupCode, groupName, agentCode, agentName, size, page)
	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(merchants),
		Data:  merchants,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// Save 保存商户信息，能同时用于新增或者修改的时候
func (i *merchant) Save(data []byte) (result *model.ResultBody) {
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

	return result
}

// Delete 删除机构商户
func (i *merchant) Delete(merId string) (result *model.ResultBody) {
	log.Debugf("delete merchant by merId,merId=%s", merId)
	if merId == "" {
		log.Errorf("merId为空")
		return model.NewResultBody(2, "merId不能为空")
	}

	err := mongo.MerchantColl.Remove(merId)

	if err != nil {
		log.Errorf("删除机构商户失败: %s", err)
		return model.NewResultBody(1, "删除机构商户失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return result
}
