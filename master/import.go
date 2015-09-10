package master

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sysErr = errors.New("系统错误，请重新上传。")
var emptyErr = errors.New("上传表格为空，请检查。")
var maxFee = 0.03

// importMerchant 接受excel格式文件，导入商户
func importMerchant(w http.ResponseWriter, r *http.Request) {

	// 调用七牛api获取刚上传的图片
	key := r.FormValue("key")
	resp, err := http.Get(makePrivateUrl(key))
	if err != nil {
		log.Error(err)
		w.Write(resultBody("无法获取文件，请重新上传。", 1))
		return
	}

	ebytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		w.Write(resultBody("无法获取文件，请重新上传。", 1))
		return
	}

	// 判断内容类型
	contentType := resp.Header.Get("content-type")
	if contentType == "application/json" {
		log.Error(string(ebytes))
		w.Write(resultBody("无法获取文件，请重新上传。", 1))
		return
	}

	// 包装成zipReader
	reader := bytes.NewReader(ebytes)
	zipReader, err := zip.NewReader(reader, int64(len(ebytes)))
	if err != nil {
		log.Error(err)
		w.Write(resultBody("无法获取文件，请重新上传。", 1))
		return
	}

	// 转换成excel
	file, err := xlsx.ReadZipReader(zipReader)
	if err != nil {
		log.Error(err)
		w.Write(resultBody("无法获取文件，请重新上传。", 1))
		return
	}

	ip := importer{Sheets: file.Sheets, IsDebug: false, fileName: key}
	err = ip.DoImport()
	if err != nil {
		w.Write(resultBody(err.Error(), 2))
		return
	}

	w.Write(resultBody("处理成功。", 0))
}

type importer struct {
	Mers                  []model.Merchant
	ChanMers              []model.ChanMer
	RouterPolicys         []model.RouterPolicy
	Sheets                []*xlsx.Sheet
	rowData               []*rowData
	chanMerCache          map[string]*model.ChanMer
	agentCache            map[string]*model.Agent
	groupCache            map[string]*model.Group
	fileName              string
	IsSaveMersSuccess     bool
	IsSaveChanMersSuccess bool
	IsSaveRouterSuccess   bool
	IsDebug               bool // 是否调试模式，如果是，会打印结果，不会入库
}

// DoImport 执行导入操作
func (i *importer) DoImport() error {
	before := time.Now()
	if len(i.Sheets) == 0 {
		return emptyErr
	}
	// 初始化
	i.chanMerCache = make(map[string]*model.ChanMer)
	i.agentCache = make(map[string]*model.Agent)
	i.groupCache = make(map[string]*model.Group)

	if err := i.read(); err != nil {
		return err
	}
	log.Debugf("read over, len(row)=%d", len(i.rowData))

	// 数据处理，验证等
	if err := i.dataHandle(); err != nil {
		return err
	}
	log.Debug("data handle over")

	// 数据入库
	if err := i.persist(); err != nil {
		i.rollback()
		return sysErr
	}
	after := time.Now()
	log.Debugf("import spent time %s", after.Sub(before))
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
	if len(i.rowData) == 0 {
		return emptyErr
	}
	return nil
}

func (i *importer) dataHandle() error {

	// 数据合法性验证
	for _, r := range i.rowData {

		// 先看是否有填商户号
		if r.MerId == "" {
			return fmt.Errorf("%s", "商户代码为空")
		}
		// 字段内容合法验证
		mer, err := mongo.MerchantColl.Find(r.MerId)
		switch r.Operator {
		case "A":
			// 新增找到用户，报错
			if err == nil {
				return fmt.Errorf("商户：%s 已存在", r.MerId)
			}
			// 插入验证
			if err = insertValidate(r); err != nil {
				return err
			}
		case "U":
			// 修改不存在用户，报错
			if err != nil {
				return fmt.Errorf("商户：%s 不存在", r.MerId)
			}
			log.Debug(mer)
		default:
			// D 先不做删除
			return fmt.Errorf("暂不支持 %s 操作。", r.Operator)
		}

		// 处理代理、集团
		if err = handleAgentAndGroup(r, i.agentCache, i.groupCache); err != nil {
			return err
		}

		// 支付宝
		if err = handleAlpMer(r, i.chanMerCache); err != nil {
			return err
		}

		// 微信
		if err = handleWxpMer(r, i.chanMerCache); err != nil {
			return err
		}
	}

	// 包装成入库的结构体
	i.doDataWrap()

	return nil
}

// updateValidate 更新验证
func updateValidate(r *rowData) error {
	return nil
}

// insertValidate 插入验证
func insertValidate(r *rowData) error {

	if r.MerName == "" {
		return fmt.Errorf("商户：%s 商户名称为空", r.MerId)
	}

	if r.AgentCode == "" {
		return fmt.Errorf("商户：%s 代理代码为空", r.MerId)
	}

	if r.IsNeedSignStr != "是" && r.IsNeedSignStr != "否" {
		return fmt.Errorf("是否开启验签：%s 取值错误，应为【是】或【否】", r.IsNeedSignStr)
	}

	if r.IsNeedSignStr == "是" {
		if r.SignKey == "" {
			return fmt.Errorf("商户：%s 开启验签需要填写签名密钥", r.MerId)
		}
		r.IsNeedSign = true
	}

	if r.CommodityName == "" {
		return fmt.Errorf("商户：%s 商品名称为空", r.MerId)
	}

	if r.WxpSubMerId != "" {
		if r.IsAgentStr != "是" && r.IsAgentStr != "否" {
			return fmt.Errorf("是否代理商模式：%s 取值错误，应为【是】或【否】", r.IsAgentStr)
		}
		if r.IsAgentStr == "是" {
			if r.WxpMerId == "" {
				return fmt.Errorf("商户：%s 代理商模式需要填写微信商户号", r.MerId)
			}
			if r.WxpSubMerId == "" {
				return fmt.Errorf("商户：%s 代理商模式需要填写微信子商户号", r.MerId)
			}
			r.IsAgent = true
		}
		if r.IsWxpCilSett != "是" && r.IsWxpCilSett != "否" {
			return fmt.Errorf("微信商户是否讯联清算：%s 取值错误，应为【是】或【否】", r.IsWxpCilSett)
		}
	}

	if r.AlpMerId != "" {
		if r.IsAlpCilSett != "是" && r.IsAlpCilSett != "否" {
			return fmt.Errorf("支付宝商户是否讯联清算：%s 取值错误，应为【是】或【否】", r.IsAlpCilSett)
		}
	}

	// 空则说明需要所有权限
	if r.PermissionStr == "" {
		r.Permission = []string{model.Paut, model.Purc, model.Canc, model.Inqy, model.Jszf, model.Qyfk, model.Refd, model.Void}
	} else {
		// TODO:权限格式验证
	}

	// TODO:清算信息格式验证
	return nil
}

func handleAgentAndGroup(r *rowData, agentCache map[string]*model.Agent, groupCache map[string]*model.Group) error {
	// 验证代理
	if r.AgentCode != "" {
		if _, ok := agentCache[r.AgentCode]; !ok {
			a, err := mongo.AgentColl.Find(r.AgentCode)
			if err != nil {
				return fmt.Errorf("商户：%s 代理代码(%s)不存在", r.MerId, r.AgentCode)
			}
			// 放入缓存
			agentCache[r.AgentCode] = a
			r.AgentName = a.AgentName
		}
	}

	// 验证集团,非空时验证
	if r.GroupCode != "" {
		if _, ok := groupCache[r.GroupCode]; !ok {
			g, err := mongo.GroupColl.Find(r.GroupCode)
			if err != nil {
				return fmt.Errorf("商户：%s 集团代码(%s)不存在", r.MerId, r.GroupCode)
			}
			if g.AgentCode != r.AgentCode {
				return fmt.Errorf("商户：%s 集团代码不属于该代理", r.MerId)
			}
			groupCache[r.GroupCode] = g
			r.GroupName = g.GroupName
		}
	}
	return nil
}

func handleAlpMer(r *rowData, chanMerCache map[string]*model.ChanMer) error {
	// 支付宝渠道商户
	if r.AlpMerId != "" {
		if _, ok := chanMerCache[r.AlpMerId]; !ok {
			alpMer, err := mongo.ChanMerColl.Find("ALP", r.AlpMerId)
			if err == nil {
				chanMerCache[r.AlpMerId] = alpMer
			} else {
				// 没找到，那么认为此次操作为新增渠道商户
				// 验证必填的信息
				if r.AlpMd5 == "" {
					return fmt.Errorf("支付宝商户：%s 密钥为空", r.AlpMerId)
				}
				// 费率转换
				f64, err := strconv.ParseFloat(r.AlpAcqFee, 10)
				if err != nil {
					return fmt.Errorf("支付宝商户：%s 讯联跟支付宝费率格式错误(%s)", r.AlpMerId, r.AlpAcqFee)
				}
				if f64 > maxFee {
					return fmt.Errorf("支付宝商户：%s 讯联跟支付宝费率超过最大值 %0.2f (%s)", r.AlpMerId, maxFee, r.AlpAcqFee)
				}
				r.AlpAcqFeeF = float32(f64)
				f64, err = strconv.ParseFloat(r.AlpMerFee, 10)
				if err != nil {
					return fmt.Errorf("支付宝商户：%s 商户跟讯联费率格式错误(%s)", r.AlpMerId, r.AlpMerFee)
				}
				if f64 > maxFee {
					return fmt.Errorf("支付宝商户：%s 商户跟讯联费率超过最大值 %0.2f (%s)", r.AlpMerId, maxFee, r.AlpMerFee)
				}
				r.AlpMerFeeF = float32(f64)
			}
		}
	}
	return nil
}

func handleWxpMer(r *rowData, chanMerCache map[string]*model.ChanMer) error {
	// 微信渠道商户
	if r.WxpSubMerId != "" {
		if _, ok := chanMerCache[r.WxpSubMerId]; !ok {
			wxpMer, err := mongo.ChanMerColl.Find("WXP", r.WxpSubMerId)
			if err == nil {
				// 系统中存在这个渠道商户，校验信息是否对称
				if r.IsAgent {
					if !wxpMer.IsAgentMode {
						return fmt.Errorf("微信商户：%s 并不是受理商模式", r.WxpSubMerId)
					} else {
						if wxpMer.AgentMer == nil {
							log.Errorf("%s:use agentMode but not supply agentMer,please check.", wxpMer.ChanMerId)
							return fmt.Errorf("%s", "系统错误配置，请联系管理员。")
						}
						if wxpMer.AgentMer.ChanMerId != r.WxpMerId {
							return fmt.Errorf("微信商户：%s 代理商商户号填写错误，应为 %s，实际为 %s", r.WxpSubMerId, wxpMer.AgentMer.ChanMerId, r.WxpMerId)
						}
					}
				} else {
					if wxpMer.IsAgentMode {
						return fmt.Errorf("微信商户：%s 为受理商模式", r.WxpSubMerId)
					}
				}
				chanMerCache[r.WxpSubMerId] = wxpMer
			} else {
				// 系统中不存在渠道商户，那么校验必填的信息
				if r.IsAgent {
					agent, err := mongo.ChanMerColl.Find("WXP", r.WxpMerId)
					if err != nil {
						return fmt.Errorf("微信商户：%s 系统中没有代码为 %s 的代理商商户", r.WxpSubMerId, r.WxpMerId)
					}
					chanMerCache[agent.ChanMerId] = agent
				}
				// 不是受理商模式，那么密钥必须要
				if !r.IsAgent {
					if r.WxpMd5 == "" {
						return fmt.Errorf("微信商户：%s 密钥为空", r.WxpSubMerId)
					}
				}

				// 费率转换
				f64, err := strconv.ParseFloat(r.WxpAcqFee, 10)
				if err != nil {
					return fmt.Errorf("微信商户：%s 讯联跟微信费率格式错误(%s)", r.WxpSubMerId, r.WxpAcqFee)
				}
				if f64 > maxFee {
					return fmt.Errorf("微信商户：%s 讯联跟微信费率超过最大值 3% (%s)", r.WxpSubMerId, r.WxpAcqFee)
				}
				r.WxpAcqFeeF = float32(f64)
				f64, err = strconv.ParseFloat(r.WxpMerFee, 10)
				if err != nil {
					return fmt.Errorf("微信商户：%s 商户跟讯联费率格式错误(%s)", r.WxpSubMerId, r.WxpMerFee)
				}
				if f64 > maxFee {
					return fmt.Errorf("微信商户：%s 商户跟讯联费率超过最大值 3% (%s)", r.WxpSubMerId, r.WxpMerFee)
				}
				r.WxpMerFeeF = float32(f64)
			}
		}
	}
	return nil
}

func (i *importer) doDataWrap() {
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
		mer.Detail.OpenBankName = r.BankName
		mer.Remark = "upload-" + i.fileName
		mer.MerStatus = "Normal"
		i.Mers = append(i.Mers, mer)

		// 渠道商户
		// 从缓存中查找支付宝渠道商户，如果没找到，那么增加一个渠道商户，并放到缓存里
		// 如果找到，那么不做任何事
		if r.AlpMerId != "" {
			if _, ok := i.chanMerCache[r.AlpMerId]; !ok {
				alpChanMer := model.ChanMer{}
				alpChanMer.ChanMerId = r.AlpMerId
				alpChanMer.ChanCode = "ALP"
				alpChanMer.SignKey = r.AlpMd5
				alpChanMer.AcqFee = r.AlpAcqFeeF
				alpChanMer.MerFee = r.AlpMerFeeF
				if r.IsAlpCilSett == "是" {
					alpChanMer.SettFlag = "1"
					alpChanMer.SettRole = "99911888" // TODO:check
				}
				i.ChanMers = append(i.ChanMers, alpChanMer)
				i.chanMerCache[r.AlpMerId] = &alpChanMer
			}
			// 路由
			alpRoute := model.RouterPolicy{}
			alpRoute.CardBrand = "ALP"
			alpRoute.ChanCode = alpRoute.CardBrand
			alpRoute.MerId = r.MerId
			alpRoute.ChanMerId = r.AlpMerId
			i.RouterPolicys = append(i.RouterPolicys, alpRoute)
		}

		// 从缓存中查找微信渠道商户，如果没找到，那么增加一个渠道商户，并放到缓存里
		// 如果找到，那么不做任何事
		if r.WxpSubMerId != "" {
			if _, ok := i.chanMerCache[r.WxpSubMerId]; !ok {
				wxpChanMer := model.ChanMer{}
				wxpChanMer.ChanMerId = r.WxpSubMerId
				wxpChanMer.ChanCode = "WXP"
				wxpChanMer.SignKey = r.WxpMd5
				if r.IsAgent {
					wxpChanMer.IsAgentMode = true
					wxpChanMer.AgentMer = i.chanMerCache[r.WxpMerId]
				}
				wxpChanMer.AcqFee = r.WxpAcqFeeF
				wxpChanMer.MerFee = r.WxpMerFeeF
				if r.IsWxpCilSett == "是" {
					wxpChanMer.SettFlag = "1"
					wxpChanMer.SettRole = "99911888" // TODO:check
				}
				i.ChanMers = append(i.ChanMers, wxpChanMer)
				i.chanMerCache[r.WxpSubMerId] = &wxpChanMer
			}
			// 路由
			wxpRoute := model.RouterPolicy{}
			wxpRoute.CardBrand = "WXP"
			wxpRoute.ChanCode = wxpRoute.CardBrand
			wxpRoute.MerId = r.MerId
			wxpRoute.ChanMerId = r.WxpSubMerId
			i.RouterPolicys = append(i.RouterPolicys, wxpRoute)
		}
	}
}

func (i *importer) persist() error {

	if i.IsDebug {
		for _, m := range i.Mers {
			log.Debugf("%+v", m)
		}
		for _, c := range i.ChanMers {
			log.Debugf("%+v", c)
		}
		for _, r := range i.RouterPolicys {
			log.Debugf("%+v", r)
		}
		return nil
	}

	// save mers
	err := mongo.MerchantColl.BatchAdd(i.Mers)
	if err != nil {
		return err
	}
	i.IsSaveMersSuccess = true

	// save chanMers
	// 数组长度可能为空
	if len(i.ChanMers) > 0 {
		err = mongo.ChanMerColl.BatchAdd(i.ChanMers)
		if err != nil {
			return err
		}
		i.IsSaveChanMersSuccess = true
	}

	// save routers
	// 数组长度可能为空
	if len(i.RouterPolicys) > 0 {
		err = mongo.RouterPolicyColl.BatchAdd(i.RouterPolicys)
		if err != nil {
			return err
		}
		i.IsSaveRouterSuccess = true
	}

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
	r := &rowData{}
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
		r.IsAlpCilSett = strings.TrimSpace(cell.Value)
	}
	if cell = cells[18]; cell != nil {
		r.WxpMerId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[19]; cell != nil {
		r.WxpSubMerId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[20]; cell != nil {
		r.IsAgentStr = strings.Trim(cell.Value, " ")
	}
	if cell = cells[21]; cell != nil {
		r.WxpAppId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[22]; cell != nil {
		r.WxpSubAppId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[23]; cell != nil {
		r.WxpMd5 = strings.Trim(cell.Value, " ")
	}
	if cell = cells[24]; cell != nil {
		r.WxpAcqFee = strings.Trim(cell.Value, " ")
	}
	if cell = cells[25]; cell != nil {
		r.WxpMerFee = strings.Trim(cell.Value, " ")
	}
	if cell = cells[26]; cell != nil {
		r.IsWxpCilSett = strings.TrimSpace(cell.Value)
	}
	if cell = cells[27]; cell != nil {
		r.ShopId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[28]; cell != nil {
		r.GoodsTag = strings.Trim(cell.Value, " ")
	}
	if cell = cells[29]; cell != nil {
		r.AcctNum = strings.Trim(cell.Value, " ")
	}
	if cell = cells[30]; cell != nil {
		r.AcctName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[31]; cell != nil {
		r.BankId = strings.Trim(cell.Value, " ")
	}
	if cell = cells[32]; cell != nil {
		r.BankName = strings.Trim(cell.Value, " ")
	}
	if cell = cells[33]; cell != nil {
		r.City = strings.Trim(cell.Value, " ")
	}
	i.rowData = append(i.rowData, r)
	return nil
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
	PermissionStr string // 权限（空即默认全部开放）
	IsNeedSignStr string // 是否开启验签
	SignKey       string // 签名密钥
	CommodityName string // 商户商品名称
	AlpMerId      string // 支付宝商户号（PID）
	AlpMd5        string // 支付宝密钥
	AlpAcqFee     string // 讯联跟支付宝费率
	AlpMerFee     string // 商户跟讯联费率
	IsAlpCilSett  string // 是否讯联清算
	WxpAppId      string // 商户appId
	WxpMd5        string // 微信密钥
	WxpMerId      string // 微信商户号
	WxpSubMerId   string // 微信子商户号
	IsAgentStr    string // 是否代理商模式
	WxpSubAppId   string // 子商户AppId
	WxpAcqFee     string // 讯联跟微信费率
	WxpMerFee     string // 商户跟讯联费率(微信)
	IsWxpCilSett  string // 是否讯联清算
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
	AlpAcqFeeF float32
	AlpMerFeeF float32
	WxpAcqFeeF float32
	WxpMerFeeF float32
}

func resultBody(msg string, status int) []byte {
	result := model.ResultBody{Status: status, Message: msg}
	bs, _ := json.Marshal(result)
	return bs
}
