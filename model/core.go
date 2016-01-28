package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// status
const (
	// refundStatus
	TransRefunded       = 1 // 已退款
	TransPartRefunded   = 2 // 部分退款
	TransMerClosed      = 3 // 商户发起关闭
	TransOverTimeClosed = 4 // 交易超时系统关闭

	// refundType
	CurrentDayRefund = 1
	OtherDayRefund   = 2
	NoLimitRefund    = 0

	// transType
	PayTrans    = 1 // 支付交易
	RefundTrans = 2 // 退款交易
	// ... 预授权
	CancelTrans     = 4 // 撤销交易
	CloseTrans      = 5 // 关单
	EnterpriseTrans = 6 // 企业付款
	SettTrans       = 7 // 结算交易
	PurchaseCoupons = 8 // 卡券核销

	// settStatus
	SettOK         = 0 //对账标记
	SettSuccess    = 1 // 勾兑成功
	SettSysRemain  = 2 // 系统多出的
	SettChanRemain = 3 // 渠道多出的

	// transStatus
	TransHandling = "10" // 交易处理中
	TransFail     = "20" // 交易失败
	TransSuccess  = "30" // 交易成功
	TransClosed   = "40" // 交易已关闭
	TransNotPay   = "50" // 交易待支付

	// bindingStatus
	BindingHandling = "10" // 绑定处理中
	BindingFail     = "20" // 绑定失败
	BindingSuccess  = "30" // 绑定成功
	BindingRemoved  = "40" // 已解绑（绑定成功过，后续解绑也成功）

	// transMode
	MarketMode = 2 // 市场模式
	MerMode    = 1 // 商户模式

	// 清算标识
	SR_CHANNEL = "CHANNEL"
	SR_CIL     = "CIL"
	SR_COMPANY = "COMPANY"
	SR_AGENT   = "AGENT"
	SR_GROUP   = "GROUP"
)

// cache name
const (
	Cache_RespCode            = "respCode"
	Cache_ScanPayResp         = "scanPayResp"
	Cache_CardBin             = "cardBin"
	Cache_Merchant            = "merchant"
	Cache_ChanMer             = "chanMer"
	Cache_CfcaBankMap         = "cfcaBankMap"
	Cache_CfcaMerRSAPrivKey   = "cfcaMerRSAPrivKey"
	Cache_AlipayMerRSAPrivKey = "alipayMerRSAPrivKey"
)

// RouterPolicy 路由策略
type RouterPolicy struct {
	MerId      string  `json:"merId" bson:"merId,omitempty"`         // 商户号
	CardBrand  string  `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	CardType   string  `json:"cardType" bson:"cardType,omitempty"`   // 卡类型
	TransType  string  `json:"transType" bson:"transType,omitempty"` // 交易类型
	BinGroup   string  `json:"binGroup" bson:"binGroup,omitempty"`   // 卡Bin组
	InputWay   string  `json:"inputWay" bson:"inputWay,omitempty"`   // 输入方式
	MinAmount  string  `json:"minAmount" bson:"minAmount,omitempty"` // 起始金额
	MaxAmount  string  `json:"maxAmount" bson:"maxAmount,omitempty"` // 最大金额（与起始金额配套使用，该金额范围）
	ChanCode   string  `json:"chanCode" bson:"chanCode,omitempty"`   // 渠道代码
	ChanMerId  string  `json:"chanMerId" bson:"chanMerId,omitempty"` // 渠道商户号
	SettFlag   string  `json:"settFlag" bson:"settFlag,omitempty"`
	SettRole   string  `json:"settRole" bson:"settRole,omitempty"`
	MerFee     float64 `json:"merFee" bson:"merFee,omitempty"`
	AcqFee     float64 `json:"acqFee" bson:"acqFee,omitempty"`
	CreateTime string  `bson:"createTime,omitempty" json:"createTime,omitempty"` // 创建时间
	UpdateTime string  `bson:"updateTime,omitempty" json:"updateTime,omitempty"` // 更新时间
}

// BindingInfo 商家绑定信息
type BindingInfo struct {
	MerId     string `json:"merId" bson:"merId,omitempty"`         // 商户ID
	BindingId string `json:"bindingId" bson:"bindingId,omitempty"` // 银行卡绑定ID
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	AcctType  string `json:"acctType" bson:"acctType,omitempty"`   // 账户类型
	AcctName  string `json:"acctName" bson:"acctName,omitempty"`   // 账户名称
	AcctNum   string `json:"acctNum" bson:"acctNum,omitempty"`     // 账户号码
	BankId    string `json:"bankId" bson:"bankId,omitempty"`       // 银行ID
	IdentType string `json:"identType" bson:"identType,omitempty"` // 证件类型
	IdentNum  string `json:"identNum" bson:"identNum,omitempty"`   // 证件号码
	PhoneNum  string `json:"phoneNum" bson:"phoneNum,omitempty"`   // 手机号
	ValidDate string `json:"validDate" bson:"validDate,omitempty"` // 信用卡有效期
	Cvv2      string `json:"cvv2" bson:"cvv2,omitempty"`           // CVV2
}

// BindingMap 绑定关系映射
type BindingMap struct {
	BindingId     string `json:"bindingId" bson:"bindingId,omitempty"`         // 银行卡绑定ID
	MerId         string `json:"merId" bson:"merId,omitempty"`                 // 商户ID
	ChanCode      string `json:"chanCode" bson:"chanCode,omitempty"`           // 渠道代码
	ChanMerId     string `json:"chanMerId" bson:"chanMerId,omitempty"`         // 渠道商户号
	ChanBindingId string `json:"chanBindingId" bson:"chanBindingId,omitempty"` // 目标渠道绑定ID，系统生成的
	BindingStatus string `json:"bindingStatus" bson:"bindingStatus,omitempty"` // 绑定状态，成功，失败，或者处理中
}

// CardBin 卡BIN
type CardBin struct {
	Bin       string `json:"bin" bson:"bin,omitempty"`             // 卡BIN
	BinLen    int    `json:"binLen" bson:"binLen,omitempty"`       // 卡BIN长度
	CardLen   int    `json:"cardLen" bson:"cardLen,omitempty"`     // 卡号长度
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	InsCode   string `json:"insCode" bson:"insCode,omitempty"`     // 发卡行代码
	InsName   string `json:"insName" bson:"insName,omitempty"`     // 发卡行名称
	CardName  string `json:"cardName" bson:"cardName,omitempty"`   // 卡名
	AcctType  string `json:"acctType" bson:"acctType,omitempty"`   // 账户类型
}

const MerStatusNormal = "Normal"

// Merchant 商户基本信息
type Merchant struct {
	MerId        string          `bson:"merId,omitempty" json:"merId,omitempty"`               // 商户号
	UniqueId     string          `bson:"uniqueId,omitempty" json:"uniqueId,omitempty"`         // 唯一标识
	AgentCode    string          `bson:"agentCode,omitempty" json:"agentCode,omitempty"`       // 公司代码
	SubAgentCode string          `bson:"subAgentCode,omitempty" json:"subAgentCode,omitempty"` // 代理/机构代码
	GroupCode    string          `bson:"groupCode,omitempty" json:"groupCode,omitempty"`       // 集团商户代码
	MerStatus    string          `bson:"merStatus,omitempty" json:"merStatus,omitempty"`       // 商户状态（Normal，Deleted，Test）
	AgentName    string          `bson:"agentName,omitempty" json:"agentName,omitempty"`       // 代理/机构名称
	SubAgentName string          `bson:"subAgentName,omitempty" json:"subAgentName,omitempty"` // 公司名称
	GroupName    string          `bson:"groupName,omitempty" json:"groupName,omitempty"`       // 集团/机构名称
	TransCurr    string          `bson:"transCurr,omitempty" json:"transCurr,omitempty"`       // 商户交易币种
	SignKey      string          `bson:"signKey,omitempty" json:"signKey,omitempty"`           // 商户签名密钥
	IsNeedSign   bool            `bson:"isNeedSign" json:"isNeedSign"`                         // 是否开启验签
	EncryptKey   string          `bson:"encryptKey,omitempty" json:"encryptKey,omitempty"`     // 商户加密密钥
	Remark       string          `bson:"remark,omitempty" json:"remark,omitempty"`             // 备注信息
	Permission   []string        `bson:"permission,omitempty" json:"permission,omitempty"`     // 接口权限
	RefundType   int             `bson:"refundType" json:"refundType"`                         // 0-无限制 1-只能当日退 2-只能隔日退
	Detail       MerDetail       `bson:"merDetail,omitempty" json:"detail,omitempty"`          // 商户详细信息
	CreateTime   string          `bson:"createTime,omitempty" json:"createTime,omitempty"`     // 创建时间
	UpdateTime   string          `bson:"updateTime,omitempty" json:"updateTime,omitempty"`     // 更新时间
	JsPayVersion string          `bson:"jsPayVersion,omitempty" json:"jsPayVersion,omitempty"`
	Options      *OverseasParams `bson:"options,omitempty"`
}

//  TODO :OverseasParams
type OverseasParams struct {
	MerName    string `json:"merchant_name,omitempty" bson:"merName"`
	MerNo      string `json:"merchant_no,omitempty" bson:"merNo"`
	Bn         string `json:"business_no,omitempty" bson:"busNo"`
	TerId      string `json:"terminal_id,omitempty" bson:"termNo"`
	Mcc        string `json:"mcc,omitempty" bson:"mcc"`
	RegionCode string `json:"region_code,omitempty" bson:"regCode"`
}

// MerDetail 商户详细信息
type MerDetail struct {
	// MerId         string `bson:"merId,omitempty"`         // 商户号
	MerName       string   `bson:"merName,omitempty" json:"merName,omitempty"`             // 商户名称
	GoodsTag      string   `bson:"goodsTag,omitempty" json:"goodsTag,omitempty"`           // 商品标识
	CommodityName string   `bson:"commodityName,omitempty" json:"commodityName,omitempty"` // 商品名称
	ShortName     string   `bson:"shortName,omitempty" json:"shortName,omitempty"`         // 商户简称
	Area          string   `bson:"area,omitempty" json:"area,omitempty"`
	TitleOne      string   `bson:"titleOne,omitempty" json:"titleOne,omitempty"`           // 微信扫固定码支付页面的标题1
	TitleTwo      string   `bson:"titleTwo,omitempty" json:"titleTwo,omitempty"`           // 微信扫固定码支付页面的标题2
	SuccBtnTxt    string   `bson:"succBtnTxt,omitempty" json:"succBtnTxt,omitempty"`       // 微信扫固定码支付成功后的按钮text
	SuccBtnLink   string   `bson:"succBtnLink,omitempty" json:"succBtnLink,omitempty"`     // 微信扫固定码支付成功后的按钮连接
	IsPostAmount  bool     `bson:"isPostAmount,omitempty" json:"isPostAmount,omitempty"`   // 微信扫固定码支付成功后的按钮连接是否传输金额
	Province      string   `bson:"province,omitempty" json:"province,omitempty"`           // 商户省份
	City          string   `bson:"city,omitempty" json:"city,omitempty"`                   // 商户城市
	Nation        string   `bson:"nation,omitempty" json:"nation,omitempty"`               // 商户国家
	MerType       string   `bson:"merType,omitempty" json:"merType,omitempty"`             // 商户类型
	BillingScheme string   `bson:"billingScheme,omitempty" json:"billingScheme,omitempty"` // 商户计费方案代码
	SettCurr      string   `bson:"SettCurr,omitempty" json:"SettCurr,omitempty"`           // 商户清算币种
	AcctName      string   `bson:"acctName,omitempty" json:"acctName,omitempty"`           // 商户账户名称
	AcctNum       string   `bson:"acctNum,omitempty" json:"acctNum,omitempty"`             // 商户账户
	Corp          string   `bson:"corp,omitempty" json:"corp,omitempty"`                   // 法人代表
	Master        string   `bson:"master,omitempty" json:"master,omitempty"`               // 商户负责人
	Contact       string   `bson:"contact,omitempty" json:"contact,omitempty"`             // 商户联系人
	ContactTel    string   `bson:"contactTel,omitempty" json:"contactTel,omitempty"`       // 商户联系电话
	Fax           string   `bson:"fax,omitempty" json:"fax,omitempty"`                     // 商户传真
	Email         string   `bson:"email,omitempty" json:"email,omitempty"`                 // 商户邮箱
	Addr          string   `bson:"addr,omitempty" json:"addr,omitempty"`                   // 商户地址
	Postcode      string   `bson:"postcode,omitempty" json:"postcode,omitempty"`           // 商户邮编
	Password      string   `bson:"password,omitempty" json:"password,omitempty"`           // 商户密码
	ShopID        string   `bson:"shopID,omitempty" json:"shopID,omitempty"`               // 门店id
	ShopType      string   `bson:"shopType,omitempty" json:"shopType,omitempty"`           // 门店类型
	BrandNum      string   `bson:"brandNum,omitempty" json:"brandNum,omitempty"`           // 品牌编号
	BankId        string   `bson:"bankId,omitempty" json:"bankId,omitempty"`               // 行号
	OpenBankName  string   `bson:"openBankName,omitempty" json:"openBankName,omitempty"`   // 开户银行名称
	BankName      string   `bson:"bankName,omitempty" json:"bankName,omitempty"`           // 银行名称
	BillUrl       string   `bson:"billUrl,omitempty" json:"billUrl,omitempty"`             // 扫固定码URL,获取账单
	PayUrl        string   `bson:"payUrl,omitempty" json:"payUrl,omitempty"`               // 扫固定码URL,支付地址
	Images        []string `bson:"images,omitempty" json:"images,omitempty"`               // 有关商户的一些图片路径
}

// ChanMer 渠道商户
type ChanMer struct {
	ChanCode    string   `bson:"chanCode,omitempty" json:"chanCode,omitempty"`       // 渠道代码
	ChanMerId   string   `bson:"chanMerId,omitempty" json:"chanMerId,omitempty"`     // 商户号
	ChanMerName string   `bson:"chanMerName,omitempty" json:"chanMerName,omitempty"` // 商户名称
	SettFlag    string   `bson:"settFlag,omitempty" json:"settFlag,omitempty"`       // 清算标识
	SettRole    string   `bson:"settRole,omitempty" json:"settRole,omitempty"`       // 清算角色
	SignKey     string   `bson:"signCert,omitempty" json:"signCert,omitempty"`       // 签名密钥 !!!!数据库存的是signCert
	PrivateKey  string   `bson:"privateKey,omitempty" json:"privateKey,omitempty"`   // 渠道商户私钥
	PublicKey   string   `bson:"publicKey,omitempty" json:"publicKey,omitempty"`     // 渠道商户公钥
	WxpAppId    string   `bson:"wxpAppId,omitempty" json:"wxpAppId,omitempty"`       // 微信支付App Id
	InsCode     string   `bson:"insCode,omitempty" json:"insCode,omitempty"`         // 机构号，Apple Pay支付需要把该字段对应到线下网关的chcd域
	TerminalId  string   `bson:"terminalId,omitempty" json:"terminalId,omitempty"`   // 终端号，Apple Pay支付需要把该字段对应到线下网关的terminalid域
	AcqFee      float64  `bson:"acqFee,omitempty" json:"acqFee,omitempty"`           // 讯联跟渠道费率
	HttpCert    string   `bson:"httpCert,omitempty" json:"httpCert,omitempty"`       // http cert证书
	HttpKey     string   `bson:"httpKey,omitempty" json:"httpKey,omitempty"`         // http key 证书
	AgentCode   string   `bson:"agentCode,omitempty" json:"agentCode,omitempty"`     // 支付宝代理代码
	IsAgentMode bool     `bson:"isAgentMode" json:"isAgentMode"`                     // 是否受理商模式
	AgentMer    *ChanMer `bson:"agentMer,omitempty" json:"agentMer,omitempty"`       // 受理商商户
	TransMode   int      `bson:"transMode,omitempty" json:"transMode,omitempty"`     // 交易模式 1-商户模式 2-市场模式
	AreaType    int      `bson:"areaType,omitempty" json:"areaType,omitempty"`       // 境内外区分字段0-境内 1-境外
	CreateTime  string   `bson:"createTime,omitempty" json:"createTime,omitempty"`   // 创建时间
	UpdateTime  string   `bson:"updateTime,omitempty" json:"updateTime,omitempty"`   // 更新时间

	// 0. 渠道退手续费，手续费原路返还，支付宝→机构→商户，统计报表及清算报表中的交易金额 =  负的原交易金额；
	// 1. 渠道不退手续费，机构承担手续费，统计报表及清算报表中的交易金额 =  负的原交易金额；
	// 2. 渠道不退手续费，商户承担手续费，统计报表及清算报表中的交易金额 =  负的（原交易金额 – 手续费）；
	// 3. 渠道不退手续费（预留），机构商户按比例承担手续费，这个模式目前不会有，先不统计在报表里。
	SchemeType int          `bson:"schemeType" json:"schemeType"` // 计费方案
	Sftp       *SftpAccount `bson:"sftp,omitempty" json:"-"`
	// ...
}

// SftpAccount 登录sftp的帐号
type SftpAccount struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email,omitempty"`
}

type Agent struct {
	AgentCode  string  `bson:"agentCode,omitempty" json:"agentCode,omitempty"`   // 代理代码
	AgentName  string  `bson:"agentName,omitempty" json:"agentName,omitempty"`   // 代理名称
	WxpCost    float64 `bson:"wxpCost,omitempty" json:"wxpCost,omitempty"`       // 微信成本
	AlpCost    float64 `bson:"alpCost,omitempty" json:"alpCost,omitempty"`       // 支付宝成本
	CreateTime string  `bson:"createTime,omitempty" json:"createTime,omitempty"` // 创建时间
	UpdateTime string  `bson:"updateTime,omitempty" json:"updateTime,omitempty"` // 更新时间
}

type SubAgent struct {
	AgentCode    string `bson:"agentCode,omitempty" json:"agentCode,omitempty"`       // 代理代码
	AgentName    string `bson:"agentName,omitempty" json:"agentName,omitempty"`       // 代理名称
	SubAgentCode string `bson:"subAgentCode,omitempty" json:"subAgentCode,omitempty"` // 二级代理代码
	SubAgentName string `bson:"subAgentName,omitempty" json:"subAgentName,omitempty"` // 二级代理名称
	CreateTime   string `bson:"createTime,omitempty" json:"createTime,omitempty"`     // 创建时间
	UpdateTime   string `bson:"updateTime,omitempty" json:"updateTime,omitempty"`     // 更新时间
}

type Group struct {
	GroupCode    string `bson:"groupCode,omitempty" json:"groupCode,omitempty"`       // 集团代码
	GroupName    string `bson:"groupName,omitempty" json:"groupName,omitempty"`       // 集团名称
	AgentCode    string `bson:"agentCode,omitempty" json:"agentCode,omitempty"`       // 代理代码
	AgentName    string `bson:"agentName,omitempty" json:"agentName,omitempty"`       // 代理名称
	SubAgentCode string `bson:"subAgentCode,omitempty" json:"subAgentCode,omitempty"` // 二级代理代码
	SubAgentName string `bson:"subAgentName,omitempty" json:"subAgentName,omitempty"` // 二级代理名称
	CreateTime   string `bson:"createTime,omitempty" json:"createTime,omitempty"`     // 创建时间
	UpdateTime   string `bson:"updateTime,omitempty" json:"updateTime,omitempty"`     // 更新时间
}

// SettSchemeCd 计费方案代码
type SettSchemeCd struct {
	SchemeCd  string `bson:"schemeCd"`
	FitBitMap string `bson:"fitBitMap"`
	Nm        string `bson:"nm"`
	Descs     string `bson:"descs"`
	OperIn    int    `bson:"operIn"`
	EventId   int    `bson:"eventId"`
	RecId     int    `bson:"recId"`
	RecUpdTs  string `bson:"recUpdTs"`
	RecCrtTs  string `bson:"recCrtTs"`
}

// Trans 支付、退款交易记录
type Trans struct {
	// 基本字段
	Id           bson.ObjectId `bson:"_id" json:"-"`
	OrderNum     string        `bson:"orderNum,omitempty" json:"orderNum"`                   // 商户订单流水号、退款流水号
	SysOrderNum  string        `bson:"sysOrderNum,omitempty" json:"-"`                       // 系统订单流水号、退款流水号
	ChanOrderNum string        `bson:"chanOrderNum,omitempty" json:"chanOrderNum,omitempty"` // 渠道返回订单流水号
	OrigOrderNum string        `bson:"origOrderNum,omitempty" json:"origOrderNum,omitempty"` // 源订单号 当交易类型为退款/撤销/关单时
	RespCode     string        `bson:"respCode,omitempty" json:"respcd,omitempty"`           // 网关应答码
	MerId        string        `bson:"merId,omitempty" json:"merId"`                         // 商户号
	TransAmt     int64         `bson:"transAmt" json:"transAmt"`                             // 交易金额 没有即为0
	TransStatus  string        `bson:"transStatus,omitempty" json:"transStatus"`             // 交易状态 10-处理中 20-失败 30-成功 40-已关闭
	TransType    int8          `bson:"transType,omitempty" json:"transType"`                 // 交易类型 1-支付 2-退款 3-预授权 4-撤销 5-关单
	ChanMerId    string        `bson:"chanMerId,omitempty" json:"chanMerId,omitempty"`       // 渠道商户号
	ChanCode     string        `bson:"chanCode,omitempty" json:"chanCode"`                   // 渠道代码
	ChanRespCode string        `bson:"chanRespCode,omitempty" json:"-"`                      // 渠道应答码
	CreateTime   string        `bson:"createTime,omitempty" json:"transTime,omitempty"`      // 交易创建时间 yyyy-mm-dd hh:mm:ss
	UpdateTime   string        `bson:"updateTime,omitempty" json:"updateTime,omitempty"`     // 交易更新时间 yyyy-mm-dd hh:mm:ss
	RefundStatus int8          `bson:"refundStatus,omitempty" json:"refundStatus"`           // 退款状态 当交易类型为支付时 0-正常 1-已退款/已撤销 2-部分退款
	RefundAmt    int64         `bson:"refundAmt,omitempty" json:"-"`                         // 已退款金额
	Remark       string        `bson:"remark,omitempty" json:"-"`                            // 备注
	Fee          int64         `bson:"fee" json:"-"`                                         // 手续费
	NetFee       int64         `bson:"netFee" json:"-"`                                      // 净手续费 方便计算费率
	TradeFrom    string        `bson:"tradeFrom,omitempty" json:"-"`                         // 交易来源
	LockFlag     int           `bson:"lockFlag" json:"-"`                                    // 是否加锁 1-锁住 0-无锁
	SettRole     string        `bson:"settRole,omitempty" json:"settRole,omitempty"`         // 清算角色
	PayTime      string        `bson:"payTime,omitempty" json:"payTime,omitempty"`           // 支付时间
	Currency     string        `bson:"currency,omitempty" json:"currency"`
	ExchangeRate string        `bson:"exchangeRate,omitempty" json:"-"`

	// 快捷支付
	AcctNum       string `bson:"acctNum,omitempty" json:"-"`                     // 交易账户
	SendSmsId     string `bson:"sendSmsId,omitempty" json:"-"`                   // 短信流水号
	SmsCode       string `bson:"smsCode,omitempty" json:"-"`                     // 短信验证码
	SubMerId      string `bson:"subMerId,omitempty" json:"-"`                    // 子商户id
	BindingId     string `bson:"bindingId,omitempty" json:"bindingId,omitempty"` // 商户绑定ID
	ChanBindingId string `bson:"chanBindingId,omitempty" json:"-"`               // 渠道绑定ID
	TransCurr     string `bson:"transCurr,omitempty" json:"transCurr,omitempty"` // 交易币种
	SettOrderNum  string `bson:"settOrderNum,omitempty" json:"-"`                // 结算订单号
	AcctName      string `bson:"acctName,omitempty" json:"-"`
	Province      string `bson:"province,omitempty" json:"-"`
	City          string `bson:"city,omitempty" json:"-"`
	BranchName    string `bson:"branchName,omitempty" json:"-"`

	// 扫码交易字段
	ChanDiscount    string `bson:"chanDiscount,omitempty" json:"chanDiscount,omitempty"`       // 渠道折扣 支付宝、微信
	MerDiscount     string `bson:"merDiscount,omitempty" json:"merDiscount,omitempty"`         // 商户折扣 支付宝、微信
	ConsumerAccount string `bson:"consumerAccount,omitempty" json:"consumerAccount,omitempty"` // 消费帐号 支付宝、微信
	ConsumerId      string `bson:"consumerId,omitempty" json:"consumerId,omitempty"`           // 消费id  支付宝、微信
	Busicd          string `bson:"busicd,omitempty" json:"busicd"`                             // 业务id
	AgentCode       string `bson:"agentCode,omitempty" json:"agentCode,omitempty"`             // 代理/机构号
	QrCode          string `bson:"qrCode,omitempty" json:"-"`                                  // 预下单时的二维码
	PrePayId        string `bson:"prePayId,omitempty" json:"-"`                                // 预支付凭证
	Terminalid      string `bson:"terminalid,omitempty" json:"terminalid,omitempty"`           // 终端号
	ErrorDetail     string `bson:"errorDetail,omitempty" json:"errorDetail"`                   // 错误信息
	GatheringId     string `bson:"gatheringId,omitempty" json:"-"`                             // 收款号
	GatheringName   string `bson:"gatheringName,omitempty" json:"-"`                           // 收款人
	NotifyUrl       string `bson:"notifyUrl,omitempty" json:"-"`                               // 异步通知地址
	VeriCode        string `bson:"veriCode,omitempty" json:"-"`                                // 交易凭证
	GoodsInfo       string `bson:"goodsInfo,omitempty" json:"goodsInfo,omitempty"`             // 商品详情
	NickName        string `bson:"nickName,omitempty" json:"-"`
	HeadImgUrl      string `bson:"headImgUrl,omitempty" json:"-"`
	Attach          string `bson:"attach,omitempty" json:"-"`

	// APP
	TicketNum string `bson:"ticketNum,omitempty" json:"ticketNum,omitempty"` // 关联的小票号

	// 可用于关联查询字段
	MerName      string `bson:"merName,omitempty" json:"merName,omitempty"` // 商户名称
	AgentName    string `bson:"agentName,omitempty" json:"agentName,omitempty"`
	GroupCode    string `bson:"groupCode,omitempty" json:"groupCode,omitempty"`
	GroupName    string `bson:"groupName,omitempty" json:"groupName,omitempty"`
	ShortName    string `bson:"shortName,omitempty" json:"shortName,omitempty"`
	SubAgentCode string `bson:"subAgentCode,omitempty" json:"subAgentCode,omitempty"`
	SubAgentName string `bson:"subAgentName,omitempty" json:"subAgentName,omitempty"`

	// 批导辅助字段
	MerFee float64 `bson:"-" json:"-"` // 商户费率，方便计算

	//卡券字段
	CouponsNo       string `bson:"couponsNo,omitempty" json:"couponsNo,omitempty"`              // 卡券号
	Prodname        string `bson:"prodname,omitempty" json:"prodname,omitempty"`                // 卡券名称
	WriteoffStatus  string `bson:"writeoffStatus,omitempty" json:"writeoffStatus,omitempty"`    // 核销状态
	VeriTime        string `json:"veriTime,omitempty" bson:"veriTime,omitempty"`                // 核销次数
	CardInfo        string `json:"cardInfo,omitempty"  bson:"cardInfo,omitempty"`               // 卡券详情
	AvailCount      string `json:"availCount,omitempty"  bson:"availCount,omitempty"`           // 卡券剩余可用次数
	ExpDate         string `json:"expDate,omitempty"  bson:"expDate,omitempty"`                 // 卡券有效期
	Authcode        int    `json:"authcode,omitempty"  bson:"authcode,omitempty"`               // 卡券有效期
	VoucherType     string `json:"voucherType,omitempty"  bson:"voucherType,omitempty"`         // 券类型
	SaleMinAmount   string `json:"saleMinAmount,omitempty" bson:"saleMinAmount,omitempty"`      // 满足优惠条件的最小金额
	SaleDiscount    string `json:"saleDiscount,omitempty"  bson:"saleDiscount,omitempty"`       // 抵扣值
	Cardbin         string `json:"cardbin,omitempty" bson:"cardbin,omitempty"`                  // 银行卡cardbin或者用户标识等
	TransAmount     string `json:"transAmount,omitempty"  bson:"transAmount,omitempty"`         // 交易原始金额
	PayType         string `json:"payType,omitempty"  bson:"payType,omitempty"`                 // 支付方式
	ActualPayAmount string `json:"actualPayAmount,omitempty" bson:"actualPayAmount,omitempty"`  // 实际支付金额
	ChannelTime     string `json:"channelTime,omitempty"  bson:"channelTime,omitempty"`         // 渠道处理时间
	OrigRespCode    string `bson:"origRespCode,omitempty" json:"origRespcd,omitempty"`          // 网关应答码
	OrigErrorDetail string `json:"origErrorDetail,omitempty"  bson:"origErrorDetail,omitempty"` // 原错误信息   C

}

// SummarySettData 交易汇总
type SummarySettData struct {
	TransType     int8  `bson:"transType" json:"transType"`         // 交易类型
	TotalTransNum int8  `bson:"totalTransNum" json:"totalTransNum"` // 总交易数量
	TotalTransAmt int64 `bson:"totalTransAmt" json:"totalTransAmt"` // 总交易金额
	TotalSettAmt  int64 `bson:"totalSettAmt" json:"totalSettAmt"`   // 总清算金额
	TotalMerFee   int64 `bson:"totalMerFee" json:"totalMerFee"`     // 总手续费
}

// TransInfo 交易明细 对商户
type TransInfo struct {
	TransType    int8   `json:"transType,omitempty" bson:"transType,omitempty"`
	TransAmt     int64  `json:"transAmt,omitempty" bson:"transAmt,omitempty"`
	RefundStatus int8   `json:"refundStatus,omitempty"`
	RefundAmt    int64  `json:"refundAmt,omitempty"`
	PayOrderNum  string `json:"payOrderNum,omitempty"`
}

// CfcaBankMap 中金支持银行映射表，为了和卡BIN表的银行匹配
type CfcaBankMap struct {
	InsCode  string `bson:"insCode,omitempty"`  // 机构号
	BankId   string `bson:"bankId,omitempty"`   // 银行ID
	BankName string `bson:"bankName,omitempty"` // 银行名称
}

// NewTransInfo TransInfo 构造方法
func NewTransInfo(t Trans) (info *TransInfo) {
	info = new(TransInfo)
	info.TransType = t.TransType
	switch info.TransType {
	case PayTrans:
		info.TransAmt = t.TransAmt
		info.RefundStatus = t.RefundStatus
		// 退款金额暂默认等于支付金额
		if info.RefundStatus == TransRefunded {
			info.RefundAmt = t.TransAmt
		}
	case RefundTrans:
		info.TransAmt = t.TransAmt
		info.PayOrderNum = t.OrigOrderNum
	}
	return
}

// TransSett 清算信息
type TransSett struct {
	Trans      Trans  `bson:"trans"`              // 清算的交易
	SettRole   string `bson:"settRole,omitempty"` // 清算角色
	SettDate   string `bson:"settDate,omitempty"` // 清算日期
	SettTime   string `bson:"settTime,omitempty"` // 清算具体时间
	MerFee     int64  `bson:"merFee"`             // 商户手续费
	MerSettAmt int64  `bson:"merSettAmt"`         // 商户清算金额
	AcqFee     int64  `bson:"acqFee"`             // 讯联成本
	AcqSettAmt int64  `bson:"acqSettAmt"`         // 讯联应收
	InsFee     int64  `bson:"InsFee"`             // 机构、代理手续费
	InsSettAmt int64  `bson:"InsSettAmt"`         // 机构、代理应收金额
	BlendType  int    `bson:"blendType"`          // 勾兑标识
}

// TransSettInfo 清分信息明细
type TransSettInfo struct {
	OrderNum   string `bson:"orderNum" json:"orderNum"`     // 订单号
	TransType  int8   `bson:"transType" json:"transType"`   // 交易类型
	CreateTime string `bson:"createTime" json:"transTime"`  // 交易时间
	TransAmt   int64  `bson:"transAmt" json:"transAmt"`     // 交易金额
	MerFee     int64  `bson:"merFee" json:"merFee"`         // 商户手续费
	MerSettAmt int64  `bson:"merSettAmt" json:"merSettAmt"` // 商户清算金额
	// TODO check 交易日期
}

// SN 每个终端对应的当日唯一的6位序列号
type SN struct {
	Type   string `bson:"type"`   // 类型
	MerId  string `bson:"merId"`  // 商户号
	TermId string `bson:"termId"` // 终端号
	Sn     int64  `bson:"sn"`     // 序列号
}

// Version 版本信息
type Version struct {
	Vn     string `bson:"vn"`     // 版本号 yyyymmddhhmmss
	LastVn string `bson:"lastVn"` // 上一个版本号
	VnType string `bson:"vnType"` // 版本类型
}

// ChanCSV 渠道文件csv
type ChanCSV struct {
	Code     string `bson:"code" json:"code"`
	Msg      string `bson:"msg" json:"msg"`
	RespCode string `bson:",omitempty" json:"respCode"`
	RespMsg  string `bson:",omitempty" json:"respMsg"`
}

// QuickpayCSV 系统应答码
type QuickpayCSV struct {
	RespCode string     `bson:"respCode" json:"respCode"`
	RespMsg  string     `bson:"respMsg" json:"respMsg"`
	Cfca     []*ChanCSV `bson:"cfca,omitempty" json:"cfca"`
	Cil      []*ChanCSV `bson:"cil,omitempty" json:"cil"`
	// ...
}

// TransSettLog 清算日志
type TransSettLog struct {
	Status     int8   `bson:"status,omitempty"`     // 状态值 1-成功 2-失败
	Addr       string `bson:"addr,omitempty"`       // 地址
	Date       string `bson:"date,omitempty"`       // 日期 yyyy-mm-dd
	CreateTime string `bson:"createTime,omitempty"` // 开始时间 yyyy-mm-dd hh:mm:ss
	ModifyTime string `bson:"modifyTime,omitempty"` // 更新时间 yyyy-mm-dd hh:mm:ss
	Method     string `bson:"method,omitempty"`     // 执行方法
}

// CheckAndNotify 检查并通知
type CheckAndNotify struct {
	BizType string `bson:"bizType"`
	CurTag  string `bson:"curTag"`
	PrevTag string `bson:"prevTag"`
	App1Tag string `bson:"app1Tag"`
	App2Tag string `bson:"app2Tag"`
	App3Tag string `bson:"app3Tag"`
	App4Tag string `bson:"app4Tag"`
}

// SpTransLogs 扫码交易日志
type SpTransLogs struct {
	ReqId        string      `bson:"reqId" json:"reqId"`
	Direction    string      `bson:"direction" json:"direction"`
	MerId        string      `bson:"merId" json:"merId"`
	OrderNum     string      `bson:"orderNum,omitempty" json:"orderNum"`
	OrigOrderNum string      `bson:"origOrderNum,omitempty" json:"origOrderNum"`
	TransTime    string      `bson:"transTime" json:"transTime"`
	TransType    string      `bson:"transType" json:"transType"`
	MsgType      int         `bson:"msgType" json:"msgType"`
	Msg          interface{} `bson:"msg" json:"msg"`
}

// Task 任务
type Task struct {
	D          time.Duration `bson:"-"`
	Name       string        `bson:"name"`
	IsDoing    bool          `bson:"isDoing"`
	F          func()        `bson:"-"`
	CreateTime string        `bson:"createTime"`
	UpdateTime string        `bson:"updateTime"`
}
