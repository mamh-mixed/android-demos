package master

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// settRole
const (
	SR_CHANNEL = "CHANNEL"
	SR_CIL     = "CIL"
	SR_COMPANY = "COMPANY"
	SR_AGENT   = "AGENT"
	SR_GROUP   = "GROUP"
)

var (
	halfwhite         = []byte{0xc2, 0xa0} // ASCII：32被UTF-8编码之后成为ASCII：194 和 160的组合
	noSessionFound    = resultBody("no session found, please retry.", 1)
	maxFee            = 0.03
	settFlagArray     = []string{SR_GROUP, SR_CHANNEL, SR_CIL, SR_AGENT, SR_COMPANY}
	replaceWhitespace = strings.NewReplacer(" ", "", "\r", "", "\t", "", "\n", "", string(halfwhite), "")
	c                 = regexp.MustCompile(`^[0-9a-zA-Z]{15}$`)
)

// importMerchant 接受excel格式文件，导入商户
func importMerchant(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	curSession, err := Session.Get(r)
	if err != nil {
		w.Write(noSessionFound)
		return
	}

	// 获取语言
	locale := GetLocale(curSession.Locale)
	im := &locale.ImportMessage

	// 文件错误
	var fileErr = resultBody(im.FileErr, 1)

	// 调用七牛api获取刚上传的图片
	key := r.FormValue("key")
	resp, err := http.Get(qiniu.MakePrivateUrl(key))
	if err != nil {
		log.Errorf("get file from qiniu err: %s", err)
		w.Write(fileErr)
		return
	}

	defer resp.Body.Close()

	ebytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read body err: %s", err)
		w.Write(fileErr)
		return
	}

	// 判断内容类型
	contentType := resp.Header.Get("content-type")
	if contentType == "application/json" {
		log.Errorf("get file from qiniu err: %s", string(ebytes))
		w.Write(fileErr)
		return
	}

	// 包装成zipReader
	reader := bytes.NewReader(ebytes)
	zipReader, err := zip.NewReader(reader, int64(len(ebytes)))
	if err != nil {
		log.Errorf("zip read bytes err: %s", err)
		w.Write(fileErr)
		return
	}

	// 转换成excel
	file, err := xlsx.ReadZipReader(zipReader)
	if err != nil {
		log.Errorf("zip read excel err: %s", err)
		w.Write(fileErr)
		return
	}

	ip := importer{Sheets: file.Sheets, fileName: key, msg: im}
	info, err := ip.DoImport()
	if err != nil {
		w.Write(resultBody(err.Error(), 2))
		return
	}

	w.Write(resultBody(info, 0))
}

type importer struct {
	A        *operation
	U        *operation
	Sheets   []*xlsx.Sheet
	rowData  []*rowData
	rowMap   map[string]*rowData
	cache    *cache
	fileName string
	IsDebug  bool // 是否调试模式，如果是，会打印结果，不会入库
	msg      *ImportMessage
}

type cache struct {
	ChanMerCache map[string]*model.ChanMer
	CompanyCache map[string]*model.SubAgent
	AgentCache   map[string]*model.Agent
	GroupCache   map[string]*model.Group
	RouterCache  map[string]*model.RouterPolicy
}

func (c *cache) Init() {
	c.ChanMerCache = make(map[string]*model.ChanMer)
	c.CompanyCache = make(map[string]*model.SubAgent)
	c.AgentCache = make(map[string]*model.Agent)
	c.GroupCache = make(map[string]*model.Group)
	c.RouterCache = make(map[string]*model.RouterPolicy)
}

type operation struct {
	Mers                  []model.Merchant
	ChanMers              []model.ChanMer
	RouterPolicys         []model.RouterPolicy
	AppAccts              []model.AppUser
	IsSaveMersSuccess     bool
	IsSaveChanMersSuccess bool
	IsSaveRouterSuccess   bool
	IsSaveAppAcctSuccess  bool
}

// DoImport 执行导入操作
func (i *importer) DoImport() (string, error) {
	before := time.Now()

	// 如果空，返回
	if len(i.Sheets) == 0 {
		return "", errors.New(i.msg.EmptyErr)
	}

	// 初始化map
	i.rowMap = make(map[string]*rowData)

	// 读取数据
	if err := i.read(); err != nil {
		return "", err
	}

	// 成功读取，初始化
	i.cache = new(cache)
	i.cache.Init()
	i.A, i.U = new(operation), new(operation)

	// 数据处理，验证等
	if err := i.dataHandle(); err != nil {
		return "", err
	}

	// 数据入库
	if err := i.persist(); err != nil {
		i.rollback()
		log.Errorf("persist error: %s, rollback ...", err)
		return "", errors.New(i.msg.SysErr)
	}

	return fmt.Sprintf(i.msg.ImportSuccess, len(i.rowData), time.Since(before)), nil
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
			return fmt.Errorf(i.msg.CellMapErr, index+1, err)
		}
	}
	if len(i.rowData) == 0 {
		return errors.New(i.msg.EmptyErr)
	}
	return nil
}

func (i *importer) dataHandle() error {

	m := i.msg.DataHandleErr

	// 数据合法性验证
	for index, r := range i.rowData {

		// 先看是否有填门店号
		if r.MerId == "" {
			return fmt.Errorf(m.NoMerId, index+3)
		}

		if !c.MatchString(r.MerId) {
			return fmt.Errorf(m.MerIdFormatErr, index+3, r.MerId)
		}

		// 字段内容合法验证
		mer, err := mongo.MerchantColl.Find(r.MerId)
		switch r.Operator {
		case "A":
			// 新增找到用户，报错
			if err == nil {
				return fmt.Errorf(m.MerIdExist, r.MerId)
			}
			// 插入验证
			if err = insertValidate(r, i.msg); err != nil {
				return err
			}
		case "U":
			// 修改不存在用户，报错
			if err != nil {
				return fmt.Errorf(m.MerIdNotExist, r.MerId)
			}
			if err = updateValidate(r, i.msg); err != nil {
				return err
			}
			r.Mer = mer
		default:
			// D 先不做删除
			return fmt.Errorf(m.NotSupportOperation, index+3, r.Operator)
		}

		// 处理费率
		if r.AlpAcqFee != "" && r.AlpMerFee != "" {
			var errStr string
			r.AlpMerFeeF, r.AlpAcqFeeF, errStr = feeParse(r.AlpMerFee, r.AlpAcqFee, i.msg)
			if errStr != "" {
				return fmt.Errorf(m.ALPMerchantErr, r.AlpMerId, errStr)
			}
		}
		if r.WxpAcqFee != "" && r.WxpMerFee != "" {
			var errStr string
			r.WxpMerFeeF, r.WxpAcqFeeF, errStr = feeParse(r.WxpMerFee, r.WxpAcqFee, i.msg)
			if errStr != "" {
				return fmt.Errorf(m.WXPMerchantErr, r.WxpMerId, errStr)
			}
		}

		// 处理代理、公司、集团
		if err = handleDegree(r, i.cache, i.msg); err != nil {
			return err
		}

		// 支付宝
		if err = handleAlpMer(r, i.cache, i.msg); err != nil {
			return err
		}

		// 微信
		if err = handleWxpMer(r, i.cache, i.msg); err != nil {
			return err
		}

		// app
		if err = handleAppAcct(r, i.msg); err != nil {
			return err
		}

	}

	// 包装成入库的结构体
	i.doDataWrap()

	return nil
}

// updateValidate 更新验证
func updateValidate(r *rowData, im *ImportMessage) error {
	m, yes, no := im.ValidateErr, im.Yes, im.No
	if r.IsNeedSignStr != "" {
		if r.IsNeedSignStr != yes && r.IsNeedSignStr != no {
			return fmt.Errorf(m.OpenSignValueErr, r.IsNeedSignStr)
		}
		if r.IsNeedSignStr == yes {
			if r.SignKey == "" {
				return fmt.Errorf(m.NoSignKey, r.MerId)
			}
			r.IsNeedSign = true
		}
	}

	if r.IsAddAcctStr != "" {
		if r.IsAddAcctStr != yes && r.IsAddAcctStr != no {
			return fmt.Errorf(m.AddAcctValueErr, r.IsAddAcctStr)
		}
		if r.IsAddAcctStr == yes {
			r.IsAddAcct = true
		}
		// 修改时，是或者否都要填
		if r.AppUsername == "" || r.AppPassword == "" {
			return fmt.Errorf(m.UNOrPWDEmptyErr, r.MerId)
		}
	}

	if r.IsAgentStr != "" {
		if r.IsAgentStr != yes && r.IsAgentStr != no {
			return fmt.Errorf(m.IsAgentStrErr, r.IsAgentStr)
		}
		if r.IsAgentStr == yes {
			r.IsAgent = true
		}

	}
	if r.WxpSettFlag != "" {
		if !util.StringInSlice(r.WxpSettFlag, settFlagArray) {
			return fmt.Errorf(m.WXPSettFlagErr, r.WxpSettFlag)
		}
	}

	if r.AlpSettFlag != "" {
		if !util.StringInSlice(r.AlpSettFlag, settFlagArray) {
			return fmt.Errorf(m.ALPSettFlagErr, r.AlpSettFlag)
		}
	}

	if r.SignKey != "" {
		if len(r.SignKey) != 32 {
			return fmt.Errorf(m.SignLengthErr, r.MerId, r.SignKey)
		}
	}

	if r.AlpMerId != "" {
		// 默认是境内商户
		if r.IsDomesticStr == no {
			// 海外支付宝商户必填字段
			if r.AlpMerName == "" || r.AlpMerNo == "" {
				return fmt.Errorf(m.NoOverseasChanMer, r.MerId)
			}
			if r.AlpSchemeType == "" {
				return fmt.Errorf(m.NoSchemeType, r.MerId)
			}
		} else {
			// 其他情况默认都是国内的
			r.IsDomestic = true
		}
	}

	return nil
}

// insertValidate 插入验证
func insertValidate(r *rowData, im *ImportMessage) error {

	m, yes, no := im.ValidateErr, im.Yes, im.No
	if r.MerName == "" {
		return fmt.Errorf(m.NoMerName, r.MerId)
	}

	if r.AgentCode == "" {
		return fmt.Errorf(m.NoAgentCode, r.MerId)
	}

	if r.IsNeedSignStr != yes && r.IsNeedSignStr != no {
		return fmt.Errorf(m.OpenSignValueErr, r.IsNeedSignStr)
	}

	if r.IsAddAcctStr != "" {
		if r.IsAddAcctStr != yes && r.IsAddAcctStr != no {
			return fmt.Errorf(m.AddAcctValueErr, r.IsAddAcctStr)
		}
		if r.IsAddAcctStr == yes {
			if r.AppUsername == "" || r.AppPassword == "" {
				return fmt.Errorf(m.UNOrPWDEmptyErr, r.MerId)
			}
			r.IsAddAcct = true
		}
	}

	if r.IsNeedSignStr == yes {
		// if r.SignKey == "" {
		// 	return fmt.Errorf("商户：%s 开启验签需要填写签名密钥", r.MerId)
		// }
		if r.SignKey != "" {
			if len(r.SignKey) != 32 {
				return fmt.Errorf(m.SignLengthErr, r.MerId, r.SignKey)
			}
		}
		r.IsNeedSign = true
	}

	if r.CommodityName == "" {
		return fmt.Errorf(m.NoCommodityName, r.MerId)
	}

	if r.WxpSubMerId != "" {
		if r.IsAgentStr != yes && r.IsAgentStr != no {
			return fmt.Errorf(m.IsAgentStrErr, r.IsAgentStr)
		}
		if r.IsAgentStr == yes {
			if r.WxpMerId == "" {
				return fmt.Errorf(m.NoWXPMer, r.MerId)
			}
			r.IsAgent = true
		}
		if !util.StringInSlice(r.WxpSettFlag, settFlagArray) {
			return fmt.Errorf(m.WXPSettFlagErr, r.WxpSettFlag)
		}
	}

	if r.AlpMerId != "" {
		if !util.StringInSlice(r.AlpSettFlag, settFlagArray) {
			return fmt.Errorf(m.ALPSettFlagErr, r.AlpSettFlag)
		}
		if r.IsDomesticStr != "" {
			if r.IsDomesticStr == no {
				// 海外支付宝商户必填字段
				if r.AlpMerName == "" || r.AlpMerNo == "" {
					return fmt.Errorf(m.NoOverseasChanMer, r.MerId)
				}
				if r.AlpSchemeType == "" {
					return fmt.Errorf(m.NoSchemeType, r.MerId)
				}
			} else {
				// 其他情况默认都是国内的
				r.IsDomestic = true
			}
		}
	}

	// 空则说明需要所有权限
	if r.PermissionStr == "" {
		r.Permission = []string{model.Paut, model.Purc, model.Canc, model.Inqy, model.Jszf, model.Qyzf, model.Refd, model.Void}
	} else {
		// TODO:权限格式验证
	}

	// TODO:清算信息格式验证
	return nil
}

// handleAppAcct 处理app账户信息
func handleAppAcct(r *rowData, im *ImportMessage) error {
	m := im.DataHandleErr
	// 填了是或者否再处理
	if r.IsAddAcctStr != "" {
		count, err := mongo.AppUserCol.Count(r.AppUsername)
		if err != nil {
			return errors.New(im.SysErr)
		}
		switch r.Operator {
		case "A":
			if r.IsAddAcct {
				if count > 0 {
					return fmt.Errorf(m.UsernameExist, r.MerId, r.AppUsername)
				}
			}
		case "U":
			if r.IsAddAcct {
				if count > 0 {
					return fmt.Errorf(m.UsernameExist, r.MerId, r.AppUsername)
				}
			} else {
				// 修改密码
				if count == 0 {
					return fmt.Errorf(m.UsernameNotExist, r.MerId, r.AppUsername)
				}
			}
		}
	}
	return nil
}

func handleDegree(r *rowData, c *cache, im *ImportMessage) error {

	m := im.DataHandleErr

	// 验证代理
	if r.AgentCode != "" {
		if _, ok := c.AgentCache[r.AgentCode]; !ok {
			a, err := mongo.AgentColl.Find(r.AgentCode)
			if err != nil {
				return fmt.Errorf(m.AgentNotExist, r.MerId, r.AgentCode)
			}
			// 放入缓存
			c.AgentCache[r.AgentCode] = a
			r.AgentName = a.AgentName
		} else {
			r.AgentName = c.AgentCache[r.AgentCode].AgentName
		}
	}

	// 验证公司
	if r.SubAgentCode != "" {
		if _, ok := c.CompanyCache[r.SubAgentCode]; !ok {
			s, err := mongo.SubAgentColl.Find(r.SubAgentCode)
			if err != nil {
				return fmt.Errorf(m.CompanyNotExist, r.MerId, r.SubAgentCode)
			}
			switch r.Operator {
			case "A":
				if s.AgentCode != r.AgentCode {
					return fmt.Errorf(m.CompanyBelongsErr, r.MerId, r.SubAgentCode)
				}
			case "U":
				if s.AgentCode != r.Mer.AgentCode {
					return fmt.Errorf(m.CompanyBelongsErr+"(%s)", r.MerId, r.SubAgentCode, r.Mer.AgentCode)
				}
			}
			// 放入缓存
			c.CompanyCache[r.SubAgentCode] = s
			r.SubAgentName = s.SubAgentName
		}
	}

	// 验证集团,非空时验证
	if r.GroupCode != "" {
		if _, ok := c.GroupCache[r.GroupCode]; !ok {
			g, err := mongo.GroupColl.Find(r.GroupCode)
			if err != nil {
				return fmt.Errorf(m.GroupNotExist, r.MerId, r.GroupCode)
			}

			// TODO: 先验证是否是属于代理级别的后面加上是否是属于公司级别的
			switch r.Operator {
			case "A":
				if g.AgentCode != r.AgentCode {
					return fmt.Errorf(m.GroupBelongsErr, r.MerId, r.GroupCode)
				}
			case "U":
				if r.Mer.AgentCode != g.AgentCode {
					return fmt.Errorf(m.GroupBelongsErr+"(%s)", r.MerId, r.GroupCode, r.Mer.AgentCode)
				}
			}
			c.GroupCache[r.GroupCode] = g
			r.GroupName = g.GroupName
		}
	}

	return nil
}

func handleAlpMer(r *rowData, c *cache, im *ImportMessage) error {
	m := im.DataHandleErr
	// 支付宝渠道商户
	if r.AlpMerId != "" {
		if _, ok := c.ChanMerCache[r.AlpMerId]; !ok {
			alpMer, err := mongo.ChanMerColl.Find("ALP", r.AlpMerId)
			if err == nil {
				c.ChanMerCache[r.AlpMerId] = alpMer
			} else {
				// 没找到，那么认为此次操作为新增渠道商户
				// 验证必填的信息
				if r.AlpMd5 == "" {
					return fmt.Errorf(m.NoALPKey, r.AlpMerId)
				}
			}
		}
	} else {
		// 可能需要修改路由策略信息
		if r.Operator == "U" {
			// 有填清算标识，但是没有对应的渠道商户号，那么需要处理
			if r.AlpSettFlag != "" && r.AlpMerId == "" {
				rp := mongo.RouterPolicyColl.Find(r.MerId, "ALP")
				if rp == nil {
					return fmt.Errorf(m.NoALPRouteToUdpSf, r.MerId)
				}
				settFlagHandle(r.AlpSettFlag, rp, r.Mer)
				// TODO:处理手续费变更
				c.RouterCache[r.MerId+"ALP"] = rp
			}
		}
	}
	return nil
}

func handleWxpMer(r *rowData, c *cache, im *ImportMessage) error {
	m := im.DataHandleErr
	// 微信渠道商户
	if r.WxpSubMerId != "" {
		if _, ok := c.ChanMerCache[r.WxpSubMerId]; !ok {
			wxpMer, err := mongo.ChanMerColl.Find("WXP", r.WxpSubMerId)
			if err == nil {
				// 系统中存在这个渠道商户，校验信息是否对称
				if r.IsAgent {
					if !wxpMer.IsAgentMode {
						return fmt.Errorf(m.WXPNotAgentMode, r.WxpSubMerId)
					} else {
						if wxpMer.AgentMer == nil {
							log.Errorf("%s:use agentMode but not supply agentMer,please check.", wxpMer.ChanMerId)
							return fmt.Errorf("%s", m.SysConfigErr)
						}
						if wxpMer.AgentMer.ChanMerId != r.WxpMerId {
							return fmt.Errorf(m.AgentMerInfoErr, r.WxpSubMerId, wxpMer.AgentMer.ChanMerId, r.WxpMerId)
						}
					}
				} else {
					if wxpMer.IsAgentMode {
						return fmt.Errorf(m.AgentModeNotMatch, r.WxpSubMerId)
					}
				}
				c.ChanMerCache[r.WxpSubMerId] = wxpMer
			} else {
				// 系统中不存在渠道商户，那么校验必填的信息
				if r.IsAgent {
					if _, ok := c.ChanMerCache[r.WxpMerId]; !ok {
						agent, err := mongo.ChanMerColl.Find("WXP", r.WxpMerId)
						if err != nil {
							return fmt.Errorf(m.NoSuchAgentMer, r.WxpSubMerId, r.WxpMerId)
						}
						c.ChanMerCache[agent.ChanMerId] = agent
					}
				}
				// 不是受理商模式，那么密钥必须要
				if !r.IsAgent {
					if r.WxpMd5 == "" {
						return fmt.Errorf(m.NoWXPKey, r.WxpSubMerId)
					}
				}
			}
		}
	} else {
		// 可能需要修改路由策略信息
		if r.Operator == "U" {
			// 有填清算标识，但是没有对应的渠道商户号，那么需要处理
			if r.WxpSettFlag != "" && r.WxpMerId == "" {
				rp := mongo.RouterPolicyColl.Find(r.MerId, "WXP")
				if rp == nil {
					return fmt.Errorf(m.NoWXPRouteToUdpSf, r.MerId)
				}
				settFlagHandle(r.WxpSettFlag, rp, r.Mer)
				// TODO:处理手续费
				c.RouterCache[r.MerId+"WXP"] = rp
			}
		}
	}
	return nil
}

// feeParse 费率转换
func feeParse(merFee, acqFee string, im *ImportMessage) (mf, af float64, errStr string) {

	m := im.DataHandleErr
	// acqFee
	af64, err := strconv.ParseFloat(acqFee, 10)
	if err != nil {
		errStr = fmt.Sprintf(m.CILFeeErr, acqFee)
		return
	}
	if af64 > maxFee {
		errStr = fmt.Sprintf(m.CILFeeOverMax, acqFee)
		return
	}

	// merFee
	mf64, err := strconv.ParseFloat(merFee, 10)
	if err != nil {
		errStr = fmt.Sprintf(m.MerFeeErr, merFee)
		return
	}
	if mf64 > maxFee {
		errStr = fmt.Sprintf(m.MerFeeOverMax, merFee)
	}

	return mf64, af64, errStr
}

func (i *importer) doDataWrap() {
	for _, r := range i.rowData {
		var mer *model.Merchant
		switch r.Operator {
		case "A":
			// 集团商户
			mer = &model.Merchant{}
			mer.MerId = r.MerId
			mer.UniqueId = util.Confuse(mer.MerId)
			mer.Detail.MerName = r.MerName
			mer.Detail.CommodityName = r.CommodityName
			mer.Detail.ShopID = r.ShopId
			mer.Detail.GoodsTag = r.GoodsTag
			mer.Detail.AcctNum = r.AcctNum
			mer.Detail.AcctName = r.AcctName
			mer.AgentCode = r.AgentCode
			mer.AgentName = r.AgentName
			mer.SubAgentCode = r.SubAgentCode
			mer.SubAgentName = r.SubAgentName
			mer.GroupCode = r.GroupCode
			mer.GroupName = r.GroupName
			mer.SignKey = r.SignKey
			mer.IsNeedSign = r.IsNeedSign
			mer.Permission = r.Permission
			mer.Detail.BankId = r.BankId
			mer.Detail.City = r.City
			mer.Detail.OpenBankName = r.BankName
			mer.Detail.TitleOne = r.TitleOne
			mer.Detail.TitleTwo = r.TitleTwo
			mer.Remark = "add-upload-" + i.fileName
			mer.MerStatus = "Normal"
			// 随机生成密钥
			if mer.IsNeedSign && mer.SignKey == "" {
				mer.SignKey = util.SignKey()
			}
			// 生成账单和支付地址
			if r.TitleOne != "" || r.TitleTwo != "" {
				mer.Detail.BillUrl = fmt.Sprintf("%s/trade.html?merchantCode=%s", webAppUrl, mer.UniqueId)
				mer.Detail.PayUrl = fmt.Sprintf("%s/index.html?merchantCode=%s", webAppUrl, b64Encoding.EncodeToString([]byte(mer.MerId)))
			}

			// 补充境外渠道参数
			if !r.IsDomestic {
				o := model.OverseasParams{}
				o.Bn = r.AlpBusNo
				o.Mcc = r.AlpMcc
				o.MerName = r.AlpMerName
				o.MerNo = r.AlpMerNo
				o.RegionCode = r.AlpRegCode
				o.TerId = r.AlpTermNo
				mer.Options = &o
			}

			i.A.Mers = append(i.A.Mers, *mer)

			// app账户
			if r.IsAddAcct {
				user := model.AppUser{}
				user.UserName = r.AppUsername
				pb := md5.Sum([]byte(r.AppPassword))
				user.Password = fmt.Sprintf("%x", pb[:])
				user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
				user.UpdateTime = user.CreateTime
				user.Activate = "true"
				user.Limit = "false"
				user.MerId = mer.MerId
				user.RegisterFrom = model.PreRegister
				i.A.AppAccts = append(i.A.AppAccts, user)
			}

		case "U":
			mer = r.Mer
			if mer.UniqueId == "" {
				mer.UniqueId = util.Confuse(mer.MerId)
			}
			if r.MerName != "" {
				mer.Detail.MerName = r.MerName
			}
			if r.CommodityName != "" {
				mer.Detail.CommodityName = r.CommodityName
			}
			if r.ShopId != "" {
				mer.Detail.ShopID = r.ShopId
			}
			if r.GoodsTag != "" {
				mer.Detail.GoodsTag = r.GoodsTag
			}
			if r.AcctNum != "" {
				mer.Detail.AcctNum = r.AcctNum
			}
			if r.AcctName != "" {
				mer.Detail.AcctName = r.AcctName
			}
			if r.AgentCode != "" {
				mer.AgentCode = r.AgentCode
			}
			if r.AgentName != "" {
				mer.AgentName = r.AgentName
			}
			if r.SubAgentCode != "" {
				mer.SubAgentCode = r.SubAgentCode
			}
			if r.SubAgentName != "" {
				mer.SubAgentName = r.SubAgentName
			}
			if r.GroupCode != "" {
				mer.GroupCode = r.GroupCode
			}
			if r.GroupName != "" {
				mer.GroupName = r.GroupName
			}
			if r.SignKey != "" {
				mer.SignKey = r.SignKey
			}
			if r.IsNeedSignStr != "" {
				mer.IsNeedSign = r.IsNeedSign
			}
			if r.BankId != "" {
				mer.Detail.BankId = r.BankId
			}
			if r.City != "" {
				mer.Detail.City = r.City
			}
			if r.BankName != "" {
				mer.Detail.OpenBankName = r.BankName
			}
			if r.TitleOne != "" {
				mer.Detail.TitleOne = r.TitleOne
			}
			if r.TitleTwo != "" {
				mer.Detail.TitleTwo = r.TitleTwo
			}
			// 生成账单和支付地址
			if mer.Detail.BillUrl == "" && mer.Detail.PayUrl == "" {
				if r.TitleOne != "" || r.TitleTwo != "" {
					mer.Detail.BillUrl = fmt.Sprintf("%s/trade.html?merchantCode=%s", webAppUrl, mer.UniqueId)
					mer.Detail.PayUrl = fmt.Sprintf("%s/index.html?merchantCode=%s", webAppUrl, b64Encoding.EncodeToString([]byte(mer.MerId)))
				}
			}

			// 修改境外渠道参数
			if !r.IsDomestic {
				if mer.Options != nil {
					if r.AlpBusNo != "" {
						mer.Options.Bn = r.AlpBusNo
					}
					if r.AlpMerNo != "" {
						mer.Options.MerNo = r.AlpMerNo
					}
					if r.AlpMerName != "" {
						mer.Options.MerName = r.AlpMerName
					}
					if r.AlpRegCode != "" {
						mer.Options.RegionCode = r.AlpRegCode
					}
					if r.AlpMcc != "" {
						mer.Options.Mcc = r.AlpMcc
					}
					if r.AlpTermNo != "" {
						mer.Options.TerId = r.AlpTermNo
					}
				} else {
					o := model.OverseasParams{}
					o.Bn = r.AlpBusNo
					o.Mcc = r.AlpMcc
					o.MerName = r.AlpMerName
					o.MerNo = r.AlpMerNo
					o.RegionCode = r.AlpRegCode
					o.TerId = r.AlpTermNo
					mer.Options = &o
				}
			}

			mer.Remark = "update-upload-" + i.fileName
			i.U.Mers = append(i.U.Mers, *mer)

			// app账户
			if r.IsAddAcctStr != "" {
				user := model.AppUser{}
				user.UserName = r.AppUsername
				pb := md5.Sum([]byte(r.AppPassword))
				user.Password = fmt.Sprintf("%x", pb[:])
				user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
				user.UpdateTime = user.CreateTime
				user.Activate = "true"
				user.Limit = "false"
				user.MerId = mer.MerId
				user.RegisterFrom = model.PreRegister
				if r.IsAddAcct {
					i.A.AppAccts = append(i.A.AppAccts, user)
				} else {
					// 修改
					i.U.AppAccts = append(i.U.AppAccts, user)
				}
			}
		}

		// 渠道商户
		// 从缓存中查找支付宝渠道商户，如果没找到，那么增加一个渠道商户，并放到缓存里
		// 如果找到，那么不做任何事
		if r.AlpMerId != "" {
			if _, ok := i.cache.ChanMerCache[r.AlpMerId]; !ok {
				alpChanMer := model.ChanMer{}
				alpChanMer.ChanMerId = r.AlpMerId
				alpChanMer.ChanCode = "ALP"
				alpChanMer.SignKey = r.AlpMd5
				alpChanMer.AgentCode = r.AlpAgentCode
				if !r.IsDomestic {
					alpChanMer.AreaType = channel.Oversea
					alpChanMer.SchemeType, _ = strconv.Atoi(r.AlpSchemeType)
				}

				// TODO:DELETE
				// alpChanMer.AcqFee = r.AlpAcqFeeF
				// alpChanMer.MerFee = r.AlpMerFeeF

				switch r.Operator {
				case "A":
					i.A.ChanMers = append(i.A.ChanMers, alpChanMer)
				case "U":
					i.U.ChanMers = append(i.U.ChanMers, alpChanMer)
				}
				i.cache.ChanMerCache[r.AlpMerId] = &alpChanMer
			}
			// 路由
			alpRoute := model.RouterPolicy{}
			alpRoute.CardBrand = "ALP"
			alpRoute.ChanCode = alpRoute.CardBrand
			alpRoute.MerId = r.MerId
			alpRoute.ChanMerId = r.AlpMerId

			// ADDBY:RUI,DATE:20151012
			// ------------
			alpRoute.MerFee, alpRoute.AcqFee = r.AlpMerFeeF, r.AlpAcqFeeF
			settFlagHandle(r.AlpSettFlag, &alpRoute, mer)
			// ------------

			switch r.Operator {
			case "A":
				i.A.RouterPolicys = append(i.A.RouterPolicys, alpRoute)
			case "U":
				i.U.RouterPolicys = append(i.U.RouterPolicys, alpRoute)
			}
		}

		// 从缓存中查找微信渠道商户，如果没找到，那么增加一个渠道商户，并放到缓存里
		// 如果找到，那么不做任何事
		if r.WxpSubMerId != "" {
			if _, ok := i.cache.ChanMerCache[r.WxpSubMerId]; !ok {
				wxpChanMer := model.ChanMer{}
				wxpChanMer.ChanMerId = r.WxpSubMerId
				wxpChanMer.ChanCode = "WXP"
				wxpChanMer.SignKey = r.WxpMd5
				if r.IsAgent {
					wxpChanMer.IsAgentMode = true
					wxpChanMer.AgentMer = i.cache.ChanMerCache[r.WxpMerId]
				}

				//TODO:DELETE
				// wxpChanMer.AcqFee = r.WxpAcqFeeF
				// wxpChanMer.MerFee = r.WxpMerFeeF

				switch r.Operator {
				case "A":
					i.A.ChanMers = append(i.A.ChanMers, wxpChanMer)
				case "U":
					i.U.ChanMers = append(i.U.ChanMers, wxpChanMer)
				}
				i.cache.ChanMerCache[r.WxpSubMerId] = &wxpChanMer
			}
			// 路由
			wxpRoute := model.RouterPolicy{}
			wxpRoute.CardBrand = "WXP"
			wxpRoute.ChanCode = wxpRoute.CardBrand
			wxpRoute.MerId = r.MerId
			wxpRoute.ChanMerId = r.WxpSubMerId

			// ADDBY:RUI,DATE:20151012
			// --------
			wxpRoute.MerFee, wxpRoute.AcqFee = r.WxpMerFeeF, r.WxpAcqFeeF
			settFlagHandle(r.WxpSettFlag, &wxpRoute, mer)
			// --------

			switch r.Operator {
			case "A":
				i.A.RouterPolicys = append(i.A.RouterPolicys, wxpRoute)
			case "U":
				i.U.RouterPolicys = append(i.U.RouterPolicys, wxpRoute)
			}
		}
	}

	// 更新缓存里的路由策略，如果有的话
	for _, r := range i.cache.RouterCache {
		i.U.RouterPolicys = append(i.U.RouterPolicys, *r)
	}
}

func settFlagHandle(settFlag string, rp *model.RouterPolicy, mer *model.Merchant) {
	rp.SettFlag = settFlag
	switch settFlag {
	case SR_CIL:
		rp.SettRole = SR_CIL
	case SR_CHANNEL:
		rp.SettRole = rp.ChanCode
	case SR_AGENT:
		rp.SettRole = mer.AgentCode
	case SR_COMPANY:
		rp.SettRole = mer.SubAgentCode
	case SR_GROUP:
		rp.SettRole = mer.GroupCode
	}
}

func (o *operation) print() {
	// for _, m := range o.Mers {
	log.Infof("Mers: %d length", len(o.Mers))
	// }
	// for _, c := range o.ChanMers {
	log.Infof("ChanMers: %d length", len(o.ChanMers))
	// }
	// for _, r := range o.RouterPolicys {
	log.Infof("RouterPolicys: %d length", len(o.RouterPolicys))
	// }
	// for _, u := range o.AppAccts {
	log.Infof("AppAccts: %d length", len(o.AppAccts))
	// }
}

func (i *importer) persist() error {

	if i.IsDebug {
		i.A.print()
		i.U.print()
		return nil
	}
	var err error
	// ===============ADD==============
	if len(i.A.Mers) > 0 {
		err = mongo.MerchantColl.BatchAdd(i.A.Mers)
		// mongo insert 操作是没有原子性保证的
		// 所以不管成功或失败，都认为保存成功，后续回退
		// 下同
		i.A.IsSaveMersSuccess = true
		if err != nil {
			return err
		}
	}

	if len(i.A.ChanMers) > 0 {
		err = mongo.ChanMerColl.BatchAdd(i.A.ChanMers)
		i.A.IsSaveChanMersSuccess = true
		if err != nil {
			return err
		}
	}

	if len(i.A.RouterPolicys) > 0 {
		err = mongo.RouterPolicyColl.BatchAdd(i.A.RouterPolicys)
		i.A.IsSaveRouterSuccess = true
		if err != nil {
			return err
		}
	}

	if len(i.A.AppAccts) > 0 {
		err = mongo.AppUserCol.BatchAdd(i.A.AppAccts)
		i.A.IsSaveRouterSuccess = true
		if err != nil {
			return err
		}
	}

	// ===============UPD==============
	for _, m := range i.U.Mers {
		err = mongo.MerchantColl.Update(&m)
		if err != nil {
			return err
		}
	}
	i.U.IsSaveMersSuccess = true

	for _, c := range i.U.ChanMers {
		err = mongo.ChanMerColl.Upsert(&c)
		if err != nil {
			return err
		}
	}
	i.U.IsSaveChanMersSuccess = true

	for _, r := range i.U.RouterPolicys {
		err = mongo.RouterPolicyColl.Insert(&r)
		if err != nil {
			return err
		}
	}
	i.U.IsSaveRouterSuccess = true

	for _, u := range i.U.AppAccts {
		err = mongo.AppUserCol.Update(&u)
		if err != nil {
			return err
		}
	}
	i.U.IsSaveAppAcctSuccess = true

	return nil
}

func (i *importer) rollback() {
	// ===============ADD==============
	if i.A.IsSaveMersSuccess {
		var merIds []string
		for _, m := range i.A.Mers {
			merIds = append(merIds, m.MerId)
		}
		err := mongo.MerchantColl.BatchRemove(merIds)
		if err != nil {
			log.Errorf("rollback merchant error:%s", err)
		}
	}
	if i.A.IsSaveChanMersSuccess {
		err := mongo.ChanMerColl.BatchRemove(i.A.ChanMers)
		if err != nil {
			log.Errorf("rollback chanMer error:%s", err)
		}
	}
	if i.A.IsSaveRouterSuccess {
		err := mongo.RouterPolicyColl.BatchRemove(i.A.RouterPolicys)
		if err != nil {
			log.Errorf("rollback routerPolicy error:%s", err)
		}
	}
	// ===============UPD==============
	// TODO: update的操作如何回滚
}

func (i *importer) cellMapping(cells []*xlsx.Cell) error {

	var col = len(cells)
	if col == 0 {
		return nil
	}

	correctCol := 50
	// 返回某列完整错误信息
	if col != correctCol {
		var order = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		var errStr string
		for k, v := range cells {
			t := k / 26
			var offset string
			if t == 0 {
				offset = string(order[k])
			} else {
				offset = string(order[t-1]) + string(order[k-26])
			}
			errStr += fmt.Sprintf("( %s=%s ), ", offset, v)
		}
		return fmt.Errorf(i.msg.ColNumErr, correctCol, col, errStr)
	}

	r := &rowData{}
	var cell *xlsx.Cell
	if cell = cells[0]; cell != nil {
		r.Operator = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[1]; cell != nil {
		r.AgentCode = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[2]; cell != nil {
		r.AgentName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[5]; cell != nil {
		r.SubAgentCode = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[6]; cell != nil {
		r.SubAgentName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[7]; cell != nil {
		r.GroupCode = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[8]; cell != nil {
		r.GroupName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[9]; cell != nil {
		r.MerId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[10]; cell != nil {
		r.MerName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[11]; cell != nil {
		r.PermissionStr = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[12]; cell != nil {
		r.IsNeedSignStr = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[13]; cell != nil {
		r.SignKey = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[14]; cell != nil {
		r.CommodityName = replaceWhitespace.Replace(cell.Value)
	}

	// -----支付宝商户begin
	if cell = cells[15]; cell != nil {
		r.AlpMerId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[16]; cell != nil {
		r.AlpMd5 = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[17]; cell != nil {
		r.AlpAgentCode = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[18]; cell != nil {
		r.AlpAcqFee = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[19]; cell != nil {
		r.AlpMerFee = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[20]; cell != nil {
		r.AlpSettFlag = replaceWhitespace.Replace(cell.Value)
	}

	// ------海外字段begin
	if cell = cells[21]; cell != nil {
		r.AlpSchemeType = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[22]; cell != nil {
		r.IsDomesticStr = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[23]; cell != nil {
		r.AlpMerName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[24]; cell != nil {
		r.AlpMerNo = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[25]; cell != nil {
		r.AlpBusNo = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[26]; cell != nil {
		r.AlpTermNo = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[27]; cell != nil {
		r.AlpMcc = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[28]; cell != nil {
		r.AlpRegCode = replaceWhitespace.Replace(cell.Value)
	}
	// ------海外字段end
	// -----支付宝商户end

	// ------微信字段begin
	if cell = cells[29]; cell != nil {
		r.WxpMerId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[30]; cell != nil {
		r.WxpSubMerId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[31]; cell != nil {
		r.IsAgentStr = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[32]; cell != nil {
		r.WxpAppId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[33]; cell != nil {
		r.WxpSubAppId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[34]; cell != nil {
		r.WxpMd5 = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[35]; cell != nil {
		r.WxpAcqFee = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[36]; cell != nil {
		r.WxpMerFee = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[37]; cell != nil {
		r.WxpSettFlag = replaceWhitespace.Replace(cell.Value)
	}
	// ------微信字段end

	// ------营销信息begin
	if cell = cells[38]; cell != nil {
		r.ShopId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[39]; cell != nil {
		r.GoodsTag = replaceWhitespace.Replace(cell.Value)
	}
	// ------营销信息end

	// ------清算信息begin
	if cell = cells[40]; cell != nil {
		r.AcctNum = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[41]; cell != nil {
		r.AcctName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[42]; cell != nil {
		r.BankId = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[43]; cell != nil {
		r.BankName = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[44]; cell != nil {
		r.City = replaceWhitespace.Replace(cell.Value)
	}
	// ------清算信息end

	// ------账户信息begin
	if cell = cells[45]; cell != nil {
		r.IsAddAcctStr = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[46]; cell != nil {
		r.AppUsername = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[47]; cell != nil {
		r.AppPassword = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[48]; cell != nil {
		r.TitleOne = replaceWhitespace.Replace(cell.Value)
	}
	if cell = cells[49]; cell != nil {
		r.TitleTwo = replaceWhitespace.Replace(cell.Value)
	}
	// ------账户信息end

	if _, ok := i.rowMap[r.MerId]; ok {
		return fmt.Errorf(i.msg.MerIdRepeat, r.MerId)
	}

	i.rowMap[r.MerId] = r
	i.rowData = append(i.rowData, r)

	return nil
}

type rowData struct {
	Operator  string // A/U/D
	AgentCode string // 机构/代理编号
	AgentName string // 机构/代理名称
	// 机构/代理支付宝成本
	// 机构/代理微信成本
	SubAgentCode  string // 公司编号
	SubAgentName  string // 公司名称
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
	AlpAgentCode  string // 支付宝代理代码
	AlpAcqFee     string // 讯联跟支付宝费率
	AlpMerFee     string // 商户跟讯联费率
	AlpSettFlag   string // 清算标识
	AlpSchemeType string // 计费方式
	IsDomesticStr string // 是否境内渠道
	// ---支付宝海外接口参数
	AlpMerName string
	AlpMerNo   string
	AlpBusNo   string
	AlpTermNo  string
	AlpMcc     string
	AlpRegCode string
	// ---
	WxpAppId     string // 商户appId
	WxpMd5       string // 微信密钥
	WxpMerId     string // 微信商户号
	WxpSubMerId  string // 微信子商户号
	IsAgentStr   string // 是否代理商模式
	WxpSubAppId  string // 子商户AppId
	WxpAcqFee    string // 讯联跟微信费率
	WxpMerFee    string // 商户跟讯联费率(微信)
	WxpSettFlag  string // 是否讯联清算
	ShopId       string // 门店标识
	GoodsTag     string // 商品标识
	AcctNum      string // 开户账户
	AcctName     string // 开户名称
	BankId       string // 行号
	BankName     string // 开户银行名称
	City         string // 城市
	IsAddAcctStr string // 是否新增app账户信息
	IsAddAcct    bool
	AppUsername  string // 用户名
	AppPassword  string // 密码
	TitleOne     string // 标题一
	TitleTwo     string // 标题二
	// ...
	IsAgent    bool
	IsNeedSign bool
	IsDomestic bool
	Permission []string
	AlpAcqFeeF float64
	AlpMerFeeF float64
	WxpAcqFeeF float64
	WxpMerFeeF float64
	Mer        *model.Merchant
}

func resultBody(msg string, status int) []byte {
	result := model.ResultBody{Status: status, Message: msg}
	bs, _ := json.Marshal(result)
	return bs
}
