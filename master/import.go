package master

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

const (
	domain = "7xl02q.com1.z0.glb.clouddn.com"
)

var client = kodo.Client{}

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
		return errors.New("系统错误，请重新上传。")
	}

	return nil
}

func (i *importer) read() error {
	for _, s := range i.Sheets {
		for index, _ := range s.Rows {
			// 跳过标题
			if index < 2 {
				continue
			}
		}
	}
	return nil
}

func (i *importer) dataHandle() error {

	return nil
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

type rowData struct {
	Operator  string // A/U/D
	AgentCode string // 机构/代理编号
	AgentName string // 机构/代理名称
	// 机构/代理支付宝成本
	// 机构/代理微信成本
	GroupCode     string // 集团商户编号
	GroupName     string // 集团商户名称
	MerId         string // 商户编号
	MerName       string // 商户名称
	Permission    string // 权限（空即默认全部开放）
	IsNeedSign    string // 是否开启验签
	SignKey       string // 签名密钥
	CommodityName string // 商户商品名称
	AlpMerId      string // 支付宝商户号（PID）
	AlpMd5        string // 支付宝密钥
	AlpAcqFee     string // 讯联跟支付宝费率
	AlpMerFee     string // 商户跟讯联费率
	WxpMerId      string // 微信商户号
	WxpSubMerId   string // 微信子商户号
	IsAgent       string // 是否代理商模式
	WxpSubAppId   string // 子商户AppId
	WxpAcqFee     string // 讯联跟微信费率
	WxpMerFee     string // 商户跟讯联费率(微信)
	ShopId        string // 门店标识
	// 商品标识
	AcctNum  string // 开户账户
	AcctName string // 开户名称
	BankId   string // 行号
	BankName string // 开户银行名称
	City     string // 城市
}
