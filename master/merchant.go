package master

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type merchant struct{}

var Merchant merchant

// Find 根据条件分页查找商户。
func (m *merchant) FindOne(merId string) (result *model.ResultBody) {
	log.Debugf("merId=%s", merId)

	merchant, err := mongo.MerchantColl.Find(merId)

	if err != nil {
		log.Errorf("查询一个商户(%s)出错: %s", merId, err)
		return model.NewResultBody(1, "查询失败")
	}

	merchant.SignKey = processSensitiveInfo(merchant.SignKey)
	merchant.EncryptKey = processSensitiveInfo(merchant.EncryptKey)

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    merchant,
	}

	return result
}

// Find 根据条件分页查找商户。
func (m *merchant) Find(merId, merStatus, merName, groupCode, groupName, agentCode, agentName, pay string, size, page int) (result *model.ResultBody) {
	log.Debugf("merId=%s,merStatus=%s, merName=%s, groupCode=%s, groupName=%s, agentCode=%s, agentName=%s",
		merId, merStatus, merName, groupCode, groupName, agentCode, agentName)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	merchants, total, err := mongo.MerchantColl.PaginationFind(merId, merStatus, merName, groupCode, groupName, agentCode, agentName, pay, size, page)
	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	for _, merchant := range merchants {
		merchant.SignKey = processSensitiveInfo(merchant.SignKey)
		merchant.EncryptKey = processSensitiveInfo(merchant.EncryptKey)
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
		m.MerStatus = model.MerStatusNormal
	}
	// 用户是否存在
	num, err := mongo.MerchantColl.FindCountByMerId(m.MerId)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.NewResultBody(4, "系统错误")
	}
	if num != 0 {
		return model.NewResultBody(5, "merId已存在")
	}

	if m.EncryptKey == "" {
		uniqueId := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31())
		b64 := base64.StdEncoding.EncodeToString([]byte(m.MerId))
		billUrl := fmt.Sprintf("http://qrcode.cardinfolink.net/payment/trade.html?merchantCode=%s", b64)
		userInfoUrl := fmt.Sprintf("http://qrcode.cardinfolink.net/payment/index.html?merchantCode=%s", uniqueId)
		m.UniqueId = uniqueId
		m.Detail.BillUrl = billUrl
		m.Detail.UserInfoUrl = userInfoUrl
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

func (i *merchant) Update(data []byte) (result *model.ResultBody) {
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
		m.MerStatus = model.MerStatusNormal
	}

	merchant, err := mongo.MerchantColl.Find(m.MerId)

	if err != nil {
		log.Errorf("查询一个商户(%s)出错: %s", m.MerId, err)
		return model.NewResultBody(1, "查询失败")
	}

	signKeyPro := processSensitiveInfo(merchant.SignKey)
	encryptKeyPro := processSensitiveInfo(merchant.EncryptKey)

	if m.SignKey == signKeyPro {
		m.SignKey = merchant.SignKey
	}

	if m.EncryptKey == encryptKeyPro {
		m.EncryptKey = merchant.EncryptKey
	}

	err = mongo.MerchantColl.Insert(m)
	if err != nil {
		log.Errorf("更新商户失败:%s", err)
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

// processSensitiveInfo 处理敏感信息
func processSensitiveInfo(value string) string {
	valueLen := len(value)
	if value == "" || valueLen < 9 {
		return value
	} else {
		value = fmt.Sprintf("%s************************%s", value[:4], value[valueLen-4:valueLen])
		return value
	}
}
