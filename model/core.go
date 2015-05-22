package model

import (
	"gopkg.in/mgo.v2/bson"
)

// status
const (
	TransRefunded     = 1 //已退款
	TransPartRefunded = 2 //部分退款
	PayTrans          = 1 //支付交易
	RefundTrans       = 2 //退款交易

	SettSuccess    = 1 //勾兑成功
	SettSysRemain  = 2 //系统多出的
	SettChanRemain = 3 //渠道多出的

	TransHandling = "10" //交易处理中
	TransFail     = "20" //交易失败
	TransSuccess  = "30" //交易成功

	BindingHandling = "10" //绑定处理中
	BindingFail     = "20" //绑定失败
	BindingSuccess  = "30" //绑定成功
	BindingRemoved  = "40" //已解绑（绑定成功过，后续解绑也成功）
)

// cache name
const (
	Cache_RespCode          = "respCode"
	Cache_CardBin           = "cardBin"
	Cache_Merchant          = "merchant"
	Cache_ChanMer           = "chanMer"
	Cache_CfcaBankMap       = "cfcaBankMap"
	Cache_ChanMerRSAPrivKey = "chanMerRSAPrivKey"
)

// RouterPolicy 路由策略
type RouterPolicy struct {
	MerId     string `json:"merId" bson:"merId,omitempty"`         // 商户号
	CardBrand string `json:"cardBrand" bson:"cardBrand,omitempty"` // 卡品牌
	CardType  string `json:"cardType" bson:"cardType,omitempty"`   // 卡类型
	TransType string `json:"transType" bson:"transType,omitempty"` // 交易类型
	BinGroup  string `json:"binGroup" bson:"binGroup,omitempty"`   // 卡Bin组
	InputWay  string `json:"inputWay" bson:"inputWay,omitempty"`   // 输入方式
	MinAmount string `json:"minAmount" bson:"minAmount,omitempty"` // 起始金额
	MaxAmount string `json:"maxAmount" bson:"maxAmount,omitempty"` // 最大金额（与起始金额配套使用，该金额范围）
	ChanCode  string `json:"chanCode" bson:"chanCode,omitempty"`   // 渠道代码
	ChanMerId string `json:"chanMerId" bson:"chanMerId,omitempty"` // 渠道商户号
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
}

// Merchant 商户基本信息
type Merchant struct {
	MerId      string    `bson:"merId,omitempty" json:"merId,omitempty"`           //商户号
	MerStatus  string    `bson:"merStatus,omitempty" json:"merStatus,omitempty"`   //商户状态（Normal，Deleted，Test）
	TransCurr  string    `bson:"transCurr,omitempty" json:"transCurr,omitempty"`   //商户交易币种
	SignKey    string    `bson:"signKey,omitempty" json:"signKey,omitempty"`       //商户签名密钥
	EncryptKey string    `bson:"encryptKey,omitempty" json:"encryptKey,omitempty"` //商户加密密钥
	Remark     string    `bson:"remark,omitempty" json:"remark,omitempty"`         //备注信息
	Detail     MerDetail `bson:"merDetail,omitempty" json:"detail,omitempty"`      //商户详细信息
}

// MerDetail 商户详细信息
type MerDetail struct {
	// MerId         string `bson:"merId,omitempty"`         //商户号
	MerName       string `bson:"merName,omitempty"`       //商户名称
	ShortName     string `bson:"shortName,omitempty"`     //商户简称
	City          string `bson:"city,omitempty"`          //商户城市
	Nation        string `bson:"nation,omitempty"`        //商户国家
	MerType       string `bson:"merType,omitempty"`       //商户类型
	BillingScheme string `bson:"billingScheme,omitempty"` //商户计费方案代码
	SettCurr      string `bson:"SettCurr,omitempty"`      //商户清算币种
	AcctName      string `bson:"acctName,omitempty"`      //商户账户名称
	AcctNum       string `bson:"acctNum,omitempty"`       //商户账户
	Corp          string `bson:"corp,omitempty"`          //法人代表
	Master        string `bson:"master,omitempty"`        //商户负责人
	Contact       string `bson:"contact,omitempty"`       //商户联系人
	ContactTel    string `bson:"contactTel,omitempty"`    //商户联系电话
	Fax           string `bson:"fax,omitempty"`           //商户传真
	Email         string `bson:"email,omitempty"`         //商户邮箱
	Addr          string `bson:"addr,omitempty"`          //商户地址
	Postcode      string `bson:"postcode,omitempty"`      //商户邮编
	Password      string `bson:"password,omitempty"`      //商户密码
}

// ChanMer 渠道商户
type ChanMer struct {
	ChanCode       string `bson:"chanCode,omitempty" json:"chanCode,omitempty"`             //渠道代码
	ChanMerId      string `bson:"chanMerId,omitempty" json:"chanMerId,omitempty"`           //商户号
	ChanMerName    string `bson:"chanMerName,omitempty" json:"chanMerName,omitempty"`       //商户名称
	SettFlag       string `bson:"settFlag,omitempty" json:"settFlag,omitempty"`             //清算标识
	SettRole       string `bson:"settRole,omitempty" json:"settRole,omitempty"`             //清算角色
	SignCert       string `bson:"signCert,omitempty" json:"signCert,omitempty"`             //签名证书
	CheckSignCert  string `bson:"checkSignCert,omitempty" json:"checkSignCert,omitempty"`   //验签证书
	AlpMd5Key      string `bson:"alpMd5Key,omitempty" json:"alpMd5Key,omitempty"`           //支付宝 MD5 Key
	WxpAppId       string `bson:"wxpAppId,omitempty" json:"wxpAppId,omitempty"`             //微信支付App Id
	WxpPartnerKey  string `bson:"wxpPartnerKey,omitempty" json:"wxpPartnerKey,omitempty"`   //微信支付Partner Key
	WxpEncryptCert string `bson:"wxpEncryptCert,omitempty" json:"wxpEncryptCert,omitempty"` //微信支付加密证书
	InsCode        string `bson:"insCode,omitempty" json:"insCode,omitempty"`               //机构号，Apple Pay支付需要把该字段对应到线下网关的chcd域
	TerminalId     string `bson:"terminalId,omitempty" json:"terminalId,omitempty"`         //终端号，Apple Pay支付需要把该字段对应到线下网关的terminalid域
	//...
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
	Id              bson.ObjectId `bson:"_id" json:",omitempty"`
	OrderNum        string        `bson:"orderNum,omitempty"`                         //商户订单流水号、退款流水号
	SysOrderNum     string        `bson:"sysOrderNum,omitempty"`                      //系统订单流水号、退款流水号
	ChanOrderNum    string        `bson:"chanOrderNum,omitempty"`                     //渠道返回订单流水号
	RefundOrderNum  string        `bson:"refundOrderNum,omitempty"`                   //退款订单号 当交易类型为退款时
	BindingId       string        `bson:"bindingId,omitempty"`                        //商户绑定ID
	ChanBindingId   string        `bson:"chanBindingId,omitempty"`                    //渠道绑定ID
	AcctNum         string        `bson:"acctNum,omitempty"`                          //交易账户
	RespCode        string        `bson:"respCode,omitempty"`                         //网关应答码
	MerId           string        `bson:"merId,omitempty"`                            //商户号
	TransAmt        int64         `bson:"transAmt,omitempty"`                         //交易金额
	TransCurr       string        `bson:"transCurr,omitempty"`                        //交易币种
	TransStatus     string        `bson:"transStatus,omitempty"`                      //交易状态 10-处理中 20-失败 30-成功
	TransType       int8          `bson:"transType,omitempty"`                        //交易类型 1-支付 2-退款 3-预授权
	ChanMerId       string        `bson:"chanMerId,omitempty"`                        //渠道商户号
	ChanCode        string        `bson:"chanCode,omitempty"`                         //渠道代码
	ChanRespCode    string        `bson:"chanRespCode,omitempty"`                     //渠道应答码
	CreateTime      string        `bson:"createTime,omitempty"`                       //交易创建时间 yyyy-mm-dd hh:mm:ss
	UpdateTime      string        `bson:"updateTime,omitempty"`                       //交易更新时间 yyyy-mm-dd hh:mm:ss
	RefundStatus    int8          `bson:"refundStatus,omitempty" json:"refundStatus"` //退款状态 当交易类型为支付时 0-正常 1-已退款 2-部分退款
	SendSmsId       string        `bson:"sendSmsId,omitempty"`                        //短信流水号
	SmsCode         string        `bson:"smsCode,omitempty"`                          //短信验证码
	Remark          string        `bson:"remark,omitempty"`                           //备注
	SubMerId        string        `bson:"subMerId,omitempty"`                         //子商户id
	ChanDiscount    string        `bson:"chanDiscount,omitempty"`                     //渠道折扣 支付宝、微信
	MerDiscount     string        `bson:"merDiscount,omitempty"`                      //商户折扣 支付宝、微信
	ConsumerAccount string        `bson:"consumerAccount,omitempty"`                  //消费帐号 支付宝、微信
	ConsumerId      string        `bson:"consumerId,omitempty"`                       //消费id  支付宝、微信
	Busicd          string        `bson:"busicd,omitempty"`                           //业务id
	Inscd           string        `bson:"inscd,omitempty"`                            //机构号
}

// SummarySettData 交易汇总
type SummarySettData struct {
	TransType     int8  `bson:"transType" json:"transType"`         //交易类型
	TotalTransNum int8  `bson:"totalTransNum" json:"totalTransNum"` //总交易数量
	TotalTransAmt int64 `bson:"totalTransAmt" json:"totalTransAmt"` //总交易金额
	TotalSettAmt  int64 `bson:"totalSettAmt" json:"totalSettAmt"`   //总清算金额
	TotalMerFee   int64 `bson:"totalMerFee" json:"totalMerFee"`     //总手续费
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
	BankId   string `bson:"bankId,omitempty"`   //银行ID
	BankName string `bson:"bankName,omitempty"` //银行名称
}

// NewTransInfo TransInfo 构造方法
func NewTransInfo(t Trans) (info *TransInfo) {
	info = new(TransInfo)
	info.TransType = t.TransType
	switch info.TransType {
	case PayTrans:
		info.TransAmt = t.TransAmt
		info.RefundStatus = t.RefundStatus
		//退款金额暂默认等于支付金额
		if info.RefundStatus == TransRefunded {
			info.RefundAmt = t.TransAmt
		}
	case RefundTrans:
		info.TransAmt = t.TransAmt
		info.PayOrderNum = t.RefundOrderNum
	}
	return
}

// TransSett 清算信息
type TransSett struct {
	Tran        Trans  `bson:",inline"`
	SettFlag    int8   `bson:"settFlag"`    //清算标志
	SettDate    string `bson:"settDate"`    //清算时间
	MerSettAmt  int64  `bson:"merSettAmt"`  //商户清算金额
	MerFee      int64  `bson:"merFee"`      //商户手续费
	ChanSettAmt int64  `bson:"chanSettAmt"` //渠道清算金额
	ChanFee     int64  `bson:"chanFee"`     //渠道手续费
}

// TransSettInfo 清分信息明细
type TransSettInfo struct {
	OrderNum   string `bson:"orderNum" json:"orderNum"`     //订单号
	TransType  int8   `bson:"transType" json:"transType"`   //交易类型
	CreateTime string `bson:"createTime" json:"transTime"`  //交易时间
	TransAmt   int64  `bson:"transAmt" json:"transAmt"`     //交易金额
	MerFee     int64  `bson:"merFee" json:"merFee"`         //商户手续费
	MerSettAmt int64  `bson:"merSettAmt" json:"merSettAmt"` //商户清算金额
	//TODO check 交易日期
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

// ChanCsv 渠道文件csv
type ChanCsv struct {
	Code     string `bson:"code"`
	Msg      string `bson:"msg"`
	RespCode string `bson:",omitempty"`
	RespMsg  string `bson:",omitempty"`
}

// QuickpayCsv 系统应答码
type QuickpayCsv struct {
	RespCode string     `bson:"respCode"`
	RespMsg  string     `bson:"respMsg"`
	Cfca     []*ChanCsv `bson:"cfca,omitempty"`
	Cil      []*ChanCsv `bson:"cil,omitempty"`
	//...
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
