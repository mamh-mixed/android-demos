package master

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/omigo/validator"
	"github.com/tealeg/xlsx"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

const (
	domain = "7xl02q.com1.z0.glb.clouddn.com"
)

var client = kodo.Client{}
var sysErr = errors.New("系统错误，请重新上传。")

func init() {
	conf.ACCESS_KEY = "-OOrgfZJbxz29kiW6HQsJ_OQJcjX6gaPRDf6xOcc"
	conf.SECRET_KEY = "rgBxbGeGJluv8ApEjY1RL2vq9IIfXcQAQqH4ttGo"
}

// importMerchant 接受excel格式文件，导入商户
func importMerchant(w http.ResponseWriter, r *http.Request) {

	// 调用七牛api获取刚上传的图片
	key := r.FormValue("key")
	baseUrl := kodo.MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl, nil)

	resp, err := http.Get(privateUrl)
	if err != nil {
		log.Error(err)
		return
	}

	ebytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	// 判断内容类型
	contentType := resp.Header.Get("content-type")
	if contentType == "application/json" {
		log.Error(string(ebytes))
		return
	}

	// 包装成zipReader
	reader := bytes.NewReader(ebytes)
	zipReader, err := zip.NewReader(reader, int64(len(ebytes)))
	if err != nil {
		log.Error(err)
		return
	}

	// 转换成excel
	file, err := xlsx.ReadZipReader(zipReader)
	if err != nil {
		log.Error(err)
		return
	}

	ip := importer{Sheets: file.Sheets}
	err = ip.DoImport()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("上传成功"))
}

type importer struct {
	Mers                  []model.Merchant
	ChanMers              []model.ChanMer
	RouterPolicys         []model.RouterPolicy
	Sheets                []*xlsx.Sheet
	rowData               []rowData
	IsSaveMersSuccess     bool
	IsSaveChanMersSuccess bool
	IsSaveRouterSuccess   bool
}

// DoImport 执行导入操作
func (i *importer) DoImport() error {
	if len(i.Sheets) == 0 {
		return errors.New("上传表格为空，请检查。")
	}

	if err := i.read(); err != nil {
		return err
	}
	// 数据处理，验证等
	if err := i.dataHandle(); err != nil {
		return err
	}
	// 数据入库
	if err := i.persist(); err != nil {
		i.rollback()
		return sysErr
	}

	return nil
}

func (i *importer) read() error {
	s := i.Sheets[0]
	for index, r := range s.Rows {
		// 跳过标题
		if index < 2 {
			continue
		}
		err := i.cellMapping(r.Cells)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *importer) dataHandle() error {

	// 数据合法性验证
	for _, r := range i.rowData {
		if err := validator.Validate(r); err != nil {
			// TODO: 自定义错误消息
			return err
		}

		if r.IsNeedSignStr != "是" && r.IsNeedSignStr != "否" {
			return errors.New("是否开启验签：" + r.IsNeedSignStr + " 取值错误，应为【是】或【否】")
		}

		if r.IsAgentStr != "是" && r.IsAgentStr != "否" {
			return errors.New("是否代理商模式：" + r.IsAgentStr + " 取值错误，应为【是】或【否】")
		}

		// 商户号去重
		count, err := mongo.MerchantColl.CountById(r.MerId)
		if err != nil {
			return sysErr
		}
		if count > 0 {
			return errors.New("商户ID：" + r.MerId + " 已存在。")
		}

		// 验证代理、集团信息
		g, err := mongo.GroupColl.Find(r.GroupCode)
		if err != nil {
			return errors.New("集团编号错误：" + r.GroupCode)
		}
		if g.AgentCode != r.AgentCode {
			return errors.New("代理编号错误：" + r.AgentCode)
		}

		// 验证支付宝渠道商户是否已存在
		count, err = mongo.ChanMerColl.CountByKey("ALP", r.AlpMerId)
		if err != nil {
			return sysErr
		}
		if count > 0 {
			return errors.New("支付宝商户号（PID）：" + r.AlpMerId + " 已存在。")
		}

		wxpChanMerId := ""
		if r.IsAgentStr == "是" {
			if r.WxpMerId != "1236593202" {
				return errors.New("微信商户号：" + r.WxpMerId + " 错误。")
			}
			r.IsAgent = true
			wxpChanMerId = r.WxpSubMerId
		} else {
			wxpChanMerId = r.WxpMerId
		}

		// 验证微信渠道商户是否已存在
		count, err = mongo.ChanMerColl.CountByKey("ALP", wxpChanMerId)
		if err != nil {
			return sysErr
		}
		if count > 0 {
			return errors.New("微信子商户号错误：" + r.WxpSubMerId + " 已存在。")
		}

		// 数据处理
		if r.IsNeedSignStr == "是" {
			r.IsNeedSign = true
		}

		// 空则说明需要所有权限
		if r.PermissionStr == "" {
			r.Permission = []string{model.Paut, model.Purc, model.Canc, model.Inqy, model.Jszf, model.Qyfk, model.Refd, model.Void}
		} else {
			// TODO: 确认格式
		}
	}

	i.doDataMapping()

	return nil
}

func (i *importer) doDataMapping() {
	for _, r := range i.rowData {
		// 集团商户
		mer := model.Merchant{}
		mer.MerId = r.MerId
		mer.Detail.MerName = r.MerName
		mer.Detail.CommodityName = r.CommodityName
		mer.Detail.ShopID = r.ShopId
		mer.Detail.GoodsTag = r.GoodsTag
		mer.Detail.AcctNum = r.AcctNum
		mer.Detail.AcctName = r.AcctName
		mer.AgentCode = r.AgentCode
		mer.AgentName = r.AgentName
		mer.GroupCode = r.GroupCode
		mer.GroupName = r.GroupName
		mer.SignKey = r.SignKey
		mer.IsNeedSign = r.IsNeedSign
		mer.Permission = r.Permission
		mer.Detail.BankId = r.BankId
		mer.Detail.City = r.City
		mer.Detail.BankName = r.BankName
		i.Mers = append(i.Mers, mer)

		// 渠道商户
		alpChanMer := model.ChanMer{}
		alpChanMer.ChanMerId = r.AlpMerId
		alpChanMer.ChanCode = "ALP"
		alpChanMer.SignKey = r.AlpMd5
		// TODO 费率

		wxpChanMer := model.ChanMer{}
		if r.IsAgent {
			wxpChanMer.ChanMerId = r.WxpSubMerId
			wxpChanMer.WxpAppId = r.WxpSubAppId
		} else {
			wxpChanMer.ChanMerId = r.WxpMerId
			wxpChanMer.WxpAppId = r.WxpAppId
		}
		wxpChanMer.SignKey = r.WxpMd5
		wxpChanMer.ChanCode = "WXP"
		// TODO 费率
		i.ChanMers = append(i.ChanMers, alpChanMer, wxpChanMer)

		// 路由
		alpRoute := model.RouterPolicy{}
		alpRoute.CardBrand = "ALP"
		alpRoute.ChanCode = alpRoute.CardBrand
		alpRoute.MerId = r.MerId
		alpRoute.ChanMerId = r.AlpMerId

		wxpRoute := model.RouterPolicy{}
		wxpRoute.CardBrand = "WXP"
		wxpRoute.ChanCode = wxpRoute.CardBrand
		wxpRoute.MerId = r.MerId
		wxpRoute.ChanMerId = r.WxpMerId
		// wxpRoute.IsAgent = r.IsAgent
		// wxpRoute.SubMerId = r.WxpSubMerId
		i.RouterPolicys = append(i.RouterPolicys, alpRoute, wxpRoute)
	}
}

func (i *importer) persist() error {

	// save mers
	err := mongo.MerchantColl.BatchAdd(i.Mers)
	if err != nil {
		return err
	}
	i.IsSaveMersSuccess = true

	// save chanMers
	err = mongo.ChanMerColl.BatchAdd(i.ChanMers)
	if err != nil {
		return err
	}
	i.IsSaveChanMersSuccess = true

	// save routers
	err = mongo.RouterPolicyColl.BatchAdd(i.RouterPolicys)
	if err != nil {
		return err
	}
	i.IsSaveRouterSuccess = true
	return nil
}

func (i *importer) rollback() {
	if i.IsSaveMersSuccess {
		var merIds []string
		for _, m := range i.Mers {
			merIds = append(merIds, m.MerId)
		}
		err := mongo.MerchantColl.BatchRemove(merIds)
		if err != nil {
			log.Errorf("rollback merchant error:%s", err)
		}
	}
	if i.IsSaveChanMersSuccess {
		err := mongo.ChanMerColl.BatchRemove(i.ChanMers)
		if err != nil {
			log.Errorf("rollback chanMer error:%s", err)
		}
	}
	if i.IsSaveRouterSuccess {
		err := mongo.RouterPolicyColl.BatchRemove(i.RouterPolicys)
		if err != nil {
			log.Errorf("rollback routerPolicy error:%s", err)
		}
	}
}

func (i *importer) cellMapping(cells []*xlsx.Cell) error {

	if len(cells) == 0 {
		return nil
	}
	r := rowData{}
	var cell *xlsx.Cell
	if cell = cells[0]; cell != nil {
		r.Operator = strings.Trim(cell.Value, " ")
	}
	if cell = cells[1]; cell != nil {
		r.AgentCode = strings.Trim(cell.Value, " ")
	}
	if cell = cells[2]; cell != nil {
		r.AgentName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[5]; cell != nil {
		r.GroupCode = strings.Trim(cell.Value, " ")
	}
	if cell = cells[6]; cell != nil {
		r.GroupName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[7]; cell != nil {
		r.MerId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[8]; cell != nil {
		r.MerName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[9]; cell != nil {
		r.PermissionStr = strings.Trim(cell.Value, " ")
	}
	if cell = cells[10]; cell != nil {
		r.IsNeedSignStr = strings.Trim(cell.Value, " ")
	}
	if cell = cells[11]; cell != nil {
		r.SignKey = strings.Trim(cell.Value, " ")
	}
	if cell = cells[12]; cell != nil {
		r.CommodityName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[13]; cell != nil {
		r.AlpMerId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[14]; cell != nil {
		r.AlpMd5 = strings.Trim(cell.Value, " ")
	}
	if cell = cells[15]; cell != nil {
		r.AlpAcqFee = strings.Trim(cell.Value, " ")
	}
	if cell = cells[16]; cell != nil {
		r.AlpMerFee = strings.Trim(cell.Value, " ")
	}
	if cell = cells[17]; cell != nil {
		r.WxpMerId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[18]; cell != nil {
		r.WxpSubMerId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[19]; cell != nil {
		r.IsAgentStr = strings.Trim(cell.Value, " ")
	}
	if cell = cells[20]; cell != nil {
		r.WxpSubAppId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[21]; cell != nil {
		r.WxpAcqFee = strings.Trim(cell.Value, " ")
	}
	if cell = cells[22]; cell != nil {
		r.WxpMerFee = strings.Trim(cell.Value, " ")
	}
	if cell = cells[23]; cell != nil {
		r.ShopId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[24]; cell != nil {
		r.GoodsTag = strings.Trim(cell.Value, " ")
	}
	if cell = cells[25]; cell != nil {
		r.AcctNum = strings.Trim(cell.Value, " ")
	}
	if cell = cells[26]; cell != nil {
		r.AcctName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[27]; cell != nil {
		r.BankId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[28]; cell != nil {
		r.BankName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[29]; cell != nil {
		r.City = strings.Trim(cell.Value, " ")
	}
	i.rowData = append(i.rowData, r)
	return nil
}

type rowData struct {
	Operator  string `validate:"regexp=^[A|U|D]$"` // A/U/D
	AgentCode string // 机构/代理编号
	AgentName string // 机构/代理名称
	// 机构/代理支付宝成本
	// 机构/代理微信成本
	GroupCode     string `validate:"nonzero"` // 集团商户编号
	GroupName     string `validate:"nonzero"` // 集团商户名称
	MerId         string `validate:"nonzero"` // 商户编号
	MerName       string `validate:"nonzero"` // 商户名称
	PermissionStr string // 权限（空即默认全部开放）
	IsNeedSignStr string `validate:"nonzero"` // 是否开启验签
	SignKey       string // 签名密钥
	CommodityName string `validate:"nonzero"` // 商户商品名称
	AlpMerId      string `validate:"nonzero"` // 支付宝商户号（PID）
	AlpMd5        string `validate:"nonzero"` // 支付宝密钥
	AlpAcqFee     string `validate:"nonzero"` // 讯联跟支付宝费率
	AlpMerFee     string `validate:"nonzero"` // 商户跟讯联费率
	WxpAppId      string // 商户appId
	WxpMd5        string `validate:"nonzero"` // 微信密钥
	WxpMerId      string `validate:"nonzero"` // 微信商户号
	WxpSubMerId   string // 微信子商户号
	IsAgentStr    string `validate:"nonzero"` // 是否代理商模式
	WxpSubAppId   string // 子商户AppId
	WxpAcqFee     string // 讯联跟微信费率
	WxpMerFee     string // 商户跟讯联费率(微信)
	ShopId        string // 门店标识
	GoodsTag      string // 商品标识
	AcctNum       string // 开户账户
	AcctName      string // 开户名称
	BankId        string // 行号
	BankName      string // 开户银行名称
	City          string // 城市
	// ...
	IsAgent    bool
	IsNeedSign bool
	Permission []string
}
