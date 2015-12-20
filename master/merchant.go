package master

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
)

type merchant struct{}

var Merchant merchant

var b64Encoding = base64.StdEncoding
var webAppUrl = goconf.Config.MobileApp.WebAppUrl

// Find 根据条件分页查找商户。
func (m *merchant) FindOne(merId string) (result *model.ResultBody) {
	log.Debugf("merId=%s", merId)

	merchant, err := mongo.MerchantColl.FindNotInCache(merId)

	if err != nil {
		log.Errorf("查询一个商户(%s)出错: %s", merId, err)
		return model.NewResultBody(1, "查询失败")
	}

	merchant.SignKey = ProcessSensitiveInfo(merchant.SignKey)
	merchant.EncryptKey = ProcessSensitiveInfo(merchant.EncryptKey)

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    merchant,
	}

	return result
}

// Find 根据条件分页查找商户。
func (m *merchant) Find(merchant model.Merchant, pay, createStartTime, createEndTime string, size, page int) (result *model.ResultBody) {

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	merchants, total, err := mongo.MerchantColl.PaginationFind(merchant, pay, createStartTime, createEndTime, size, page)
	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	for _, merchant := range merchants {
		merchant.SignKey = ProcessSensitiveInfo(merchant.SignKey)
		merchant.EncryptKey = ProcessSensitiveInfo(merchant.EncryptKey)
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
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if m.MerId == "" {
		log.Error("no MerId")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	// 签名密钥和加密密钥长度不能小于8
	if m.SignKey != "" && len(m.SignKey) < 8 {
		return model.NewResultBody(3, "SIGN_KEY_LOWEST_LENGTH")
	}
	if m.EncryptKey != "" && len(m.EncryptKey) < 8 {
		return model.NewResultBody(3, "加密密钥长度不能小于8")
	}

	if m.MerStatus == "" {
		m.MerStatus = model.MerStatusNormal
	}
	// 用户是否存在
	num, err := mongo.MerchantColl.FindCountByMerId(m.MerId)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.NewResultBody(4, "SELECT_ERROR")
	}
	if num != 0 {
		return model.NewResultBody(5, "MER_ID_EXIST")
	}

	if m.EncryptKey == "" {
		m.UniqueId = util.Confuse(m.MerId)
		// 有填相关信息才需要生成两个连接地址
		if m.Detail.TitleOne != "" || m.Detail.TitleTwo != "" {
			billUrl := fmt.Sprintf("%s/trade.html?merchantCode=%s", webAppUrl, m.UniqueId)
			payUrl := fmt.Sprintf("%s/index.html?merchantCode=%s", webAppUrl, b64Encoding.EncodeToString([]byte(m.MerId)))
			m.Detail.BillUrl = billUrl
			m.Detail.PayUrl = payUrl
		}
	}

	err = mongo.MerchantColl.Insert(m)
	if err != nil {
		log.Errorf("create merchant error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "CREATE_MERCHANT_SUCCESS",
		Data:    m,
	}

	return result
}

func (i *merchant) Update(data []byte) (result *model.ResultBody) {
	m := new(model.Merchant)
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if m.MerId == "" {
		log.Error("NO MerId")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}
	// 签名密钥和加密密钥长度不能小于8
	if m.SignKey != "" && len(m.SignKey) < 8 {
		return model.NewResultBody(3, "SIGN_KEY_LOWEST_LENGTH")
	}
	if m.EncryptKey != "" && len(m.EncryptKey) < 8 {
		return model.NewResultBody(3, "加密密钥长度不能小于8")
	}

	if m.MerStatus == "" {
		m.MerStatus = model.MerStatusNormal
	}

	merchant, err := mongo.MerchantColl.FindNotInCache(m.MerId)

	if err != nil {
		log.Errorf("select merchant (%s)error: %s", m.MerId, err)
		return model.NewResultBody(1, "SELECT_ERROR")
	}

	log.Debugf("newSignKey:%s,oldSignKey:%s", m.SignKey, merchant.SignKey)
	log.Debugf("newEncryptKey:%s,oldEncryptKey:%s", m.EncryptKey, merchant.EncryptKey)

	// 修改签名密钥和加密密钥
	if strings.Contains(m.SignKey, "*") {
		m.SignKey = merchant.SignKey
	}

	if strings.Contains(m.EncryptKey, "*") {
		m.EncryptKey = merchant.EncryptKey
	}

	// 扫码商户
	if m.EncryptKey == "" {
		if m.Detail.TitleOne != "" || m.Detail.TitleTwo != "" {
			if m.Detail.BillUrl == "" {
				if m.UniqueId == "" {
					m.UniqueId = util.Confuse(m.MerId)
				}
				m.Detail.BillUrl = fmt.Sprintf("%s/trade.html?merchantCode=%s", webAppUrl, m.UniqueId)
			}
			if m.Detail.PayUrl == "" {
				b64 := base64.StdEncoding.EncodeToString([]byte(m.MerId))
				m.Detail.PayUrl = fmt.Sprintf("%s/index.html?merchantCode=%s", webAppUrl, b64)
			}
		}
	}

	err = mongo.MerchantColl.Update(m)
	if err != nil {
		log.Errorf("update merchant error:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "UPDATE_MERCHANT_SUCCESS",
		Data:    m,
	}

	return result
}

// Delete 删除机构商户
func (i *merchant) Delete(merId string) (result *model.ResultBody) {
	log.Debugf("delete merchant by merId,merId=%s", merId)
	if merId == "" {
		log.Errorf("no merId")
		return model.NewResultBody(2, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err := mongo.MerchantColl.Remove(merId)

	if err != nil {
		log.Errorf("delete merchant error: %s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "DELETE_MERCHANT_ERROR",
	}

	return result
}

// ProcessSensitiveInfo 处理敏感信息
func ProcessSensitiveInfo(value string) string {
	valueLen := len(value)
	if value == "" || valueLen < 9 {
		return value
	} else {
		starString := strings.Repeat("*", valueLen-8)
		value = fmt.Sprintf("%s%s%s", value[:4], starString, value[valueLen-4:valueLen])
		return value
	}
}

func (m *merchant) Export(w http.ResponseWriter, merchant model.Merchant, pay, filename, createStartTime, createEndTime string, session *model.Session) {
	size := 10000
	page := 1
	var file = xlsx.NewFile()

	merchants, total, err := mongo.MerchantColl.PaginationFind(merchant, pay, createStartTime, createEndTime, size, page)
	if err != nil {
		log.Errorf("select all merchants error:%s", err)
		return
	}
	log.Debugf("total:%d", total)
	exportMerchant(file, merchants, session.Locale)
	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, filename))
	file.Write(w)
}

func exportMerchant(file *xlsx.File, merchants []*model.Merchant, locale string) {
	merchantLocale := GetLocale(locale).Merchant
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 可能有多个sheet
	// sheet, _ = file.AddSheet("商户表")
	sheet, _ = file.AddSheet(merchantLocale.Title)
	// 生成title
	row = sheet.AddRow()
	headRow := &struct {
		MerId      string
		MerName    string
		IsNeedSign string
		SignKey    string
		BillUrl    string
		PayUrl     string
		// }{"商户号", "商户名称", "是否验签", "签名密钥", "账单链接", "支付链接"}
	}{merchantLocale.MerId, merchantLocale.MerName, merchantLocale.IsNeedSign, merchantLocale.SignKey, merchantLocale.BillUrl, merchantLocale.PayUrl}
	row.WriteStruct(headRow, -1)
	for _, v := range merchants {
		// 商户号 商户名称 是否签名 签名密钥
		row = sheet.AddRow()
		// 商户号
		cell = row.AddCell()
		cell.Value = v.MerId
		// 商户名称
		cell = row.AddCell()
		cell.Value = v.Detail.MerName

		isNeedSign := merchantLocale.Yes
		if !v.IsNeedSign {
			isNeedSign = merchantLocale.No
		}
		//  是否���签
		cell = row.AddCell()
		cell.Value = isNeedSign

		//  签名密钥
		cell = row.AddCell()
		cell.Value = v.SignKey

		//  账单链接
		cell = row.AddCell()
		cell.Value = v.Detail.BillUrl

		//  支付链接
		cell = row.AddCell()
		cell.Value = v.Detail.PayUrl

	}

	// 设置列宽
	sheet.SetColWidth(0, 3, 18)
}
