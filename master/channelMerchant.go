package master

import (
	"encoding/json"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type chanMer struct{}

var ChanMer chanMer

// Find 根据条件查找商户。
func (c *chanMer) Find(chanCode, chanMerId, chanMerName, pay string, size, page int) (result *model.ResultBody) {
	log.Debugf("chanCode=%s; chanMerId=%s; chanMerName=%s", chanCode, chanMerId, chanMerName)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	chanMers, total, err := mongo.ChanMerColl.PaginationFind(chanCode, chanMerId, chanMerName, pay, size, page)
	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	for _, chanMer := range chanMers {
		chanMer.SignKey = ProcessSensitiveInfo(chanMer.SignKey)
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(chanMers),
		Data:  chanMers,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// FindByMerIdAndCardBrand 通过机构商户的id和卡品牌查找渠道商户的信息
func (i *chanMer) FindByMerIdAndCardBrand(merId, cardBrand string) (result *model.ResultBody) {
	router := mongo.RouterPolicyColl.Find(merId, cardBrand)

	if router == nil {
		log.Errorf("未找到商户(%s)的一个路由(%s)", merId, cardBrand)
		return model.NewResultBody(1, "查询失败")
	}

	mer, err := mongo.ChanMerColl.Find(router.ChanCode, router.ChanMerId)
	if err != nil {
		log.Errorf("未找到渠道商户(%s)失败：(%s)", router.ChanCode, err)
		return model.NewResultBody(1, "查询失败")
	}

	mer.SignKey = ProcessSensitiveInfo(mer.SignKey)

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    mer,
	}

	return result
}

// Save 保存商户信息，能同时用于新增或者修改的时候
func (i *chanMer) Save(data []byte) (result *model.ResultBody) {
	c := new(model.ChanMer)
	err := json.Unmarshal(data, c)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if c.ChanCode == "" {
		log.Error("没有chanCode")
		return model.NewResultBody(3, "缺失必要元素chanCode")
	}

	if c.ChanMerId == "" {
		log.Error("没有chanMerId")
		return model.NewResultBody(3, "缺失必要元素chanMerId")
	}

	if c.SignKey != "" && len(c.SignKey) < 8 {
		log.Debugf("签名密钥长度不能小于8，signKey=%s", c.SignKey)
		return model.NewResultBody(3, "签名密钥长度不能小于")
	}
	if c.SignKey != "" {
		isCreate := false
		channel, err := mongo.ChanMerColl.Find(c.ChanCode, c.ChanMerId)
		if err != nil {
			if err.Error() == "not found" {
				isCreate = true
			} else {
				log.Errorf("find database err,%s", err)
				return model.NewResultBody(1, "查找数据库失败")
			}
		}
		if !isCreate {
			log.Debugf("newSignCert:%s,oldSignCert:%s", c.SignKey, channel.SignKey)
		}

		if !isCreate && strings.Contains(c.SignKey, "*") {
			c.SignKey = channel.SignKey
		}
	}

	// 将微信大商户的签名密钥带*号的改为不带*号的
	if c.ChanCode == "WXP" && c.AgentMer != nil && c.AgentMer.SignKey != "" {
		bigChannel, err := mongo.ChanMerColl.Find(c.AgentMer.ChanCode, c.AgentMer.ChanMerId)
		if err != nil {
			log.Errorf("find database err,%s", err)
			return model.NewResultBody(1, "查找数据库失败")
		}
		c.AgentMer.SignKey = bigChannel.SignKey
		log.Debugf("bigChannel signCert:%s", c.AgentMer.SignKey)
	}

	err = mongo.ChanMerColl.Add(c)
	if err != nil {
		log.Errorf("新增渠道商户失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    c,
	}

	return result
}

// Match 模糊查找渠道商户
func (i *chanMer) Match(chanCode, chanMerId, chanMerName string, maxSize int) (result *model.ResultBody) {
	if maxSize <= 0 {
		maxSize = 10
	}

	chanMers, err := mongo.ChanMerColl.FuzzyFind(chanCode, chanMerId, chanMerName, maxSize)
	if err != nil {
		log.Errorf("未找到渠道商户(chanCode: %s; chanMerId: %s; chanMerName: %s)失败：(%s)", chanCode, chanMerId, chanMerName, err)
		return model.NewResultBody(1, "查询失败")
	}
	for _, chanMer := range chanMers {
		chanMer.SignKey = ProcessSensitiveInfo(chanMer.SignKey)
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    chanMers,
	}

	return result
}

// Delete 删除渠道商户
func (i *chanMer) Delete(chanCode, chanMerId string) (result *model.ResultBody) {
	if chanCode == "" || chanMerId == "" {
		return model.NewResultBody(2, "chanCode和chanMerId不能为空")
	}
	err := mongo.ChanMerColl.Remove(chanCode, chanMerId)

	if err != nil {
		log.Errorf("删除渠道商户失败: %s", err)
		return model.NewResultBody(1, "删除渠道商户失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return result
}
