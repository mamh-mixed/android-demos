package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type merchant struct{}

var Merchant merchant

// Find 根据条件查找商户。
func (m *merchant) Find(merId, merStatus string) (result *model.ResultBody) {
	log.Debugf("merId is %s; merStatus is %s", merId, merStatus)

	cond := new(model.Merchant)

	if merId != "" {
		cond.MerId = merId
	}

	if merStatus != "" {
		cond.MerStatus = merStatus
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
