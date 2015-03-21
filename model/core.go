package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	TransRefunded = 1  //已退款
	PayTrans      = 1  //支付交易
	RefundTrans   = 2  //退款交易
	TransHandling = 10 //交易处理中
	TransFail     = 20 //交易失败
	TransSuccess  = 30 //交易成功
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
}

// Merchant 商户基本信息
type Merchant struct {
	MerId      string `bson:"merId,omitempty"`      //商户号
	MerStatus  string `bson:"merStatus,omitempty"`  //商户状态（Normal，Deleted）
	TransCurr  string `bson:"transCurr,omitempty"`  //商户交易币种
	SignKey    string `bson:"signKey,omitempty"`    //商户签名密钥
	EncryptKey string `bson:"encryptKey,omitempty"` //商户加密密钥
}

// MerDetail 商户详细信息
type MerDetail struct {
	MerId         string `bson:"merId,omitempty"`         //商户号
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
	ChanCode       string `bson:"chanCode,omitempty"`       //渠道代码
	ChanMerId      string `bson:"chanMerId,omitempty"`      //商户号
	ChanMerName    string `bson:"chanMerName,omitempty"`    //商户名称
	SettFlag       string `bson:"settFlag,omitempty"`       //清算标识
	SettRole       string `bson:"settRole,omitempty"`       //清算角色
	SignCert       string `bson:"signCert,omitempty"`       //签名证书
	CheckSignCert  string `bson:"checkSignCert,omitempty"`  //验签证书
	AlpMd5Key      string `bson:"alpMd5Key,omitempty"`      //支付宝 MD5 Key
	WxpAppId       string `bson:"wxpAppId,omitempty"`       //微信支付App Id
	WxpPartnerKey  string `bson:"wxpPartnerKey,omitempty"`  //微信支付Partner Key
	WxpEncryptCert string `bson:"wxpEncryptCert,omitempty"` //微信支付加密证书
	//...
}

// Trans 支付、退款交易记录
type Trans struct {
	Id             bson.ObjectId `bson:"_id" json:",omitempty"`
	OrderNum       string        `bson:"orderNum"`                                   //商户订单流水号、退款流水号
	ChanOrderNum   string        `bson:"chanOrderNum"`                               //渠道订单流水号、退款流水号
	RefundOrderNum string        `bson:"refundOrderNum,omitempty"`                   //退款订单号 当交易类型为退款时
	ChanBindingId  string        `bson:"chanBindingId"`                              //绑定ID
	AcctNum        string        `bson:"acctNum"`                                    //交易账户
	RespCode       string        `bson:"respCode"`                                   //网关应答码
	MerId          string        `bson:"merId"`                                      //商户号
	TransAmt       int64         `bson:"transAmt"`                                   //交易金额
	TransCurr      string        `bson:"transCurr"`                                  //交易币种
	TransStatus    int8          `bson:"transStatus"`                                //交易状态 10-处理中 20-失败 30-成功
	TransType      int8          `bson:"transType"`                                  //交易类型 1-支付 2-退款
	ChanMerId      string        `bson:"chanMerId"`                                  //渠道商户号
	ChanCode       string        `bson:"chanCode"`                                   //渠道代码
	ChanRespCode   string        `bson:"chanRespCode"`                               //渠道应答码
	CreateTime     int64         `bson:"createTime"`                                 //交易创建时间
	UpdateTime     int64         `bson:"updateTime"`                                 //交易更新时间
	RefundStatus   int8          `bson:"refundStatus,omitempty" json:"refundStatus"` //退款状态 当交易类型为支付时 0-正常 1-已退款
}

// TransInfo 交易明细 对商户
type TransInfo struct {
	TransType    int8   `json:"transtype,omitempty"`
	TransAmt     int64  `json:"transAmt,omitempty"`
	RefundStatus int8   `json:"refundStatus,omitempty"`
	RefundAmt    int64  `json:"refundAmt,omitempty"`
	PayOrderNum  string `json:"payOrderNum,omitempty"`
}

// NerTransInfo TransInfo 构造方法
func NerTransInfo(t Trans) (info *TransInfo) {
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
		info.PayOrderNum = t.OrderNum
	}
	return
}
