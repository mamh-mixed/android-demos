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

		// httpCert和httpKey敏感信息处理
		processSensitiveKey(chanMer)
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
		// log.Errorf("未找到商户(%s)的一个路由(%s)", merId, cardBrand)
		return model.NewResultBody(1, "查询失败")
	}

	mer, err := mongo.ChanMerColl.Find(router.ChanCode, router.ChanMerId)
	if err != nil {
		// log.Errorf("未找到渠道商户(%s)失败：(%s)", router.ChanCode, err)
		return model.NewResultBody(1, "查询失败")
	}

	mer.SignKey = ProcessSensitiveInfo(mer.SignKey)

	// httpCert和httpKey敏感信息处理
	processSensitiveKey(mer)

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
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if c.ChanCode == "" {
		log.Error("no chanCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if c.ChanMerId == "" {
		log.Error("no chanMerId")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if c.SignKey != "" && len(c.SignKey) < 8 {
		log.Debugf("签名密钥长度不能小于8，signKey=%s", c.SignKey)
		return model.NewResultBody(3, "SIGN_KEY_LOWEST_LENGTH")
	}
	// if c.SignKey != "" {
	isCreate := false
	channel, err := mongo.ChanMerColl.Find(c.ChanCode, c.ChanMerId)
	if err != nil {
		if err.Error() == "not found" {
			isCreate = true
		} else {
			log.Errorf("find database err,%s", err)
			return model.NewResultBody(1, "ERROR")
		}
	}
	if isCreate {
		if strings.Contains(c.SignKey, "*") {
			return model.NewResultBody(4, "CHAN_MERCHANT.SIGN_KEY_CANNOT_CONTAIN_STAR")
		}
		if strings.Contains(c.HttpCert, "*") {
			return model.NewResultBody(4, "CHAN_MERCHANT.HTTP_CERT_CANNOT_CONTAIN_STAR")
		}
		if strings.Contains(c.HttpKey, "*") {
			return model.NewResultBody(4, "CHAN_MERCHANT.HTTP_KEY_CANNOT_CONTAIN_STAR")
		}
		if strings.Contains(c.PrivateKey, "*") {
			return model.NewResultBody(4, "商户私钥不能包含*")
		}
	}
	if !isCreate {
		log.Debugf("newSignCert:%s,oldSignCert:%s", c.SignKey, channel.SignKey)
	}
	if !isCreate {
		if strings.Contains(c.SignKey, "*") {
			c.SignKey = channel.SignKey
		}
		if strings.Contains(c.HttpCert, "*") {
			c.HttpCert = channel.HttpCert
		}
		if strings.Contains(c.HttpKey, "*") {
			c.HttpKey = channel.HttpKey
		}
		if strings.Contains(c.PrivateKey, "*") {
			c.PrivateKey = channel.PrivateKey
		}
	}

	// }

	// 将微信大商户的签名密钥带*号的改为不带*号的
	if c.ChanCode == "WXP" && c.AgentMer != nil && c.AgentMer.SignKey != "" {
		bigChannel, err := mongo.ChanMerColl.Find(c.AgentMer.ChanCode, c.AgentMer.ChanMerId)
		if err != nil {
			log.Errorf("find database err,%s", err)
			return model.NewResultBody(1, "SELECT_ERROR")
		}
		c.AgentMer.SignKey = bigChannel.SignKey
		c.AgentMer.HttpCert = bigChannel.HttpCert
		c.AgentMer.HttpKey = bigChannel.HttpKey
		log.Debugf("bigChannel signCert:%s", c.AgentMer.SignKey)
	}

	err = mongo.ChanMerColl.Add(c)
	if err != nil {
		log.Errorf("create chan merchant error:%s", err)
		return model.NewResultBody(1, "ERRORDefaultLocale")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "CREATE_CHAN_MERCHANT_SUCCESS",
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

		// httpCert和httpKey敏感信息处理
		processSensitiveKey(chanMer)
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
		return model.NewResultBody(2, "REQUIRED_FILED_NOT_BE_EMPTY")
	}
	err := mongo.ChanMerColl.Remove(chanCode, chanMerId)

	if err != nil {
		log.Errorf("delete chan merchant error: %s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "DELETE_CHAN_MERCHANT_SUCCESS",
	}

	return result
}

func processSensitiveKey(chanMer *model.ChanMer) {
	// httpCert,httpKey，privateKey敏感信息处理
	if chanMer.HttpCert != "" {
		httpCertStrs := strings.Split(chanMer.HttpCert, "\n")
		newHttpCert := ""
		for _, httpCert := range httpCertStrs {
			if strings.Contains(httpCert, "BEGIN CERTIFICATE") || strings.Contains(httpCert, "END CERTIFICATE") {
				newHttpCert += httpCert + "\n"
			} else {
				newHttpCert += ProcessSensitiveInfo(httpCert) + "\n"
			}

		}
		chanMer.HttpCert = newHttpCert
	}
	if chanMer.HttpKey != "" {
		httpKeyStrs := strings.Split(chanMer.HttpKey, "\n")
		newHttpKey := ""
		for _, httpKey := range httpKeyStrs {
			if strings.Contains(httpKey, "BEGIN RSA PRIVATE KEY") || strings.Contains(httpKey, "END RSA PRIVATE KEY") {
				newHttpKey += httpKey + "\n"
			} else {
				newHttpKey += ProcessSensitiveInfo(httpKey) + "\n"
			}

		}
		chanMer.HttpKey = newHttpKey
	}
	if chanMer.PrivateKey != "" {
		privateKeyStrs := strings.Split(chanMer.PrivateKey, "\n")
		newPrivateKey := ""
		for _, privateKey := range privateKeyStrs {
			if strings.Contains(privateKey, "BEGIN RSA PRIVATE KEY") || strings.Contains(privateKey, "END RSA PRIVATE KEY") {
				newPrivateKey += privateKey + "\n"
			} else {
				newPrivateKey += ProcessSensitiveInfo(privateKey) + "\n"
			}

		}
		chanMer.PrivateKey = newPrivateKey
	}
}

func (i *chanMer) Update(data []byte) (result *model.ResultBody) {
	c := new(model.ChanMer)
	err := json.Unmarshal(data, c)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if c.ChanCode == "" {
		log.Error("没有chanCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if c.ChanMerId == "" {
		log.Error("没有chanMerId")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if c.SignKey != "" && len(c.SignKey) < 8 {
		log.Debugf("签名密钥长度不能小于8，signKey=%s", c.SignKey)
		return model.NewResultBody(3, "CHAN_MERCHANT.SIGN_KEY_LOWEST_LENGTH")
	}
	channel, err := mongo.ChanMerColl.Find(c.ChanCode, c.ChanMerId)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.NewResultBody(1, "SELECT_ERROR")
	}

	log.Debugf("newSignCert:%s,oldSignCert:%s", c.SignKey, channel.SignKey)

	if strings.Contains(c.SignKey, "*") {
		c.SignKey = channel.SignKey
	}
	if strings.Contains(c.HttpCert, "*") {
		c.HttpCert = channel.HttpCert
	}
	if strings.Contains(c.HttpKey, "*") {
		c.HttpKey = channel.HttpKey
	}
	if strings.Contains(c.PrivateKey, "*") {
		c.PrivateKey = channel.PrivateKey
	}

	// 将微信大商户的签名密钥带*号的改为不带*号的
	if c.ChanCode == "WXP" && c.AgentMer != nil && c.AgentMer.SignKey != "" {
		bigChannel, err := mongo.ChanMerColl.Find(c.AgentMer.ChanCode, c.AgentMer.ChanMerId)
		if err != nil {
			log.Errorf("find database err,%s", err)
			return model.NewResultBody(1, "SELECT_ERROR")
		}
		c.AgentMer.SignKey = bigChannel.SignKey
		c.AgentMer.HttpCert = bigChannel.HttpCert
		c.AgentMer.HttpKey = bigChannel.HttpKey
		log.Debugf("bigChannel signCert:%s", c.AgentMer.SignKey)
	}

	err = mongo.ChanMerColl.Update(c)
	if err != nil {
		log.Errorf("update chan merchant error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "UPDATE_CHAN_MERCHANT_SUCCESS",
		Data:    c,
	}

	return result
}
