package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/quickpay/weixin"
	"github.com/CardInfoLink/log"
)

// busiType
const (
	Purc = "PURC"
	Paut = "PAUT"
	Inqy = "INQY"
	Refd = "REFD"
	Void = "VOID"
	Canc = "CANC"
	Qyzf = "QYZF" // 企业付款
	Jszf = "JSZF"
	Veri = "VERI" // 电子券验证/刷卡活动券查询
	Crve = "CRVE" // 刷卡活动券验证
	Quve = "QUVE" // 电子券验证结果查询
	Cave = "CAVE" // 电子券验证撤销
	List = "LIST"
	// 卡券核销状态
	COUPON_WO_SUCCESS = "SUCCESS"
	COUPON_WO_ERROR   = "ERROR"
	COUPON_WO_PROCESS = "PROCESS"
)

// QueryCondition 扫码交易查询字段
type QueryCondition struct {
	MerName            string   `json:"mchntName,omitempty"` // 可用于商户名称、商户简称模糊查询
	MerId              string   `json:"mchntid,omitempty"`   // 可用于商户号模糊查询
	MerIds             []string `json:"-"`
	UserType           string
	Col                string   `json:"-"`
	BindingId          string   `json:"bindingId"`
	AgentCode          string   `json:"agentCode,omitempty"`
	SubAgentCode       string   `json:"subAgentCode,omitempty"`
	GroupCode          string   `json:"groupCode,omitempty"`
	TransStatus        []string `json:"transStatus,omitempty"`
	RefundStatus       int      `json:"refundStatus,omitempty"`
	TransType          int      `json:"transType,omitempty"`
	StartTime          string   `json:"startTime,omitempty"`
	EndTime            string   `json:"endTime,omitempty"`
	Date               string   `json:"date,omitempty"`
	Busicd             string   `json:"busicd,omitempty"`
	OrderNum           string   `json:"orderNum,omitempty"`
	OrigOrderNum       string   `json:"origOrderNum,omitempty"`
	NextOrderNum       string   `json:"nextOrderNum,omitempty"`
	Count              int      `json:"count,omitempty"`
	Page               int      `json:"page,omitempty"`
	Total              int      `json:"total,omitempty"`
	Size               int      `json:"size,omitempty"`
	IsForReport        bool     `json:"-"`
	Respcd             string   `json:"respcd" url:"respcd"`
	RespcdNotIn        string   `json:"respcdNotIn"`
	TradeFrom          string   `json:"tradeFrom,omitempty"`
	Skip               int      `json:"skip,omitempty"`
	ChanCode           string   `json:"chanCode,omitempty"`
	Direction          string
	ReqIds             []string
	TimeType           string
	SettRole           string `json:"settRole,omitempty"`
	Locale             string
	Currency           string
	UtcOffset          int
	IsAggregateByGroup bool   `json:"isAggregateByGroup,omitempty"`                             // 是否按照商户号汇总
	CouponsNo          string `bson:"couponsNo,omitempty" json:"couponsNo,omitempty"`           // 卡券号
	WriteoffStatus     string `bson:"writeoffStatus,omitempty" json:"writeoffStatus,omitempty"` // 核销状态
	Terminalid         string `bson:"terminalid,omitempty" json:"terminalid,omitempty"`         // 终端代码
}

// QueryResult 查询结果值
type QueryResult struct {
	Rec          interface{} `json:"rec,omitempty"` // 交易明细
	Page         int         `json:"page,omitempty"`
	Total        int         `json:"total,omitempty"`
	Size         int         `json:"size,omitempty"`
	Count        int         `json:"count,omitempty"`
	RespCode     string      `json:"respCode,omitempty"`
	RespMsg      string      `json:"respMsg,omitempty"`
	NextOrderNum string      `json:"-"`
}

// Summary 商户交易汇总
type Summary struct {
	MerId         string `json:"merId,omitempty"`
	MerName       string `json:"merName,omitempty"`
	AgentName     string `json:"agentName,omitempty"`
	GroupName     string `json:"groupName,omitempty"`
	CompanyName   string `json:"-"`
	TotalTransAmt int64  `json:"totalTransAmt"`
	TotalTransNum int    `json:"totalTransNum"`
	TotalFee      int64  `json:"-"`
	Wxp           struct {
		TransAmt int64 `json:"transAmt"`
		TransNum int   `json:"transNum"`
		Fee      int64 `json:"-"`
	} `json:"wxp"`
	Alp struct {
		TransAmt int64 `json:"transAmt"`
		TransNum int   `json:"transNum"`
		Fee      int64 `json:"-"`
	} `json:"alp"`
	Data []Summary `json:"data"` // 包含每个商户单独数据
}

// TransGroup 按商户号和渠道号统计分组
type TransGroup struct {
	MerId       string    `bson:"_id"`
	MerName     string    `bson:"merName,omitempty"`
	GroupCode   string    `bson:"groupCode"`
	GroupName   string    `bson:"groupName"`
	CompanyName string    `bson:"companyName"`
	AgentName   string    `bson:"agentName"`
	TransAmt    int64     `bson:"transAmt"`
	RefundAmt   int64     `bson:"refundAmt"`
	TransNum    int       `bson:"transNum"`
	Detail      []Channel `bson:"detail"`
	Fee         int64     `bson:"fee" json:"-"`
}

// Channel 按渠道类型分组
type Channel struct {
	ChanCode  string `bson:"chanCode"`
	TransAmt  int64  `bson:"transAmt"`
	RefundAmt int64  `bson:"refundAmt"`
	TransNum  int    `bson:"transNum"`
	Fee       int64  `bson:"fee" json:"-"`
}

// Mer 按商户分组
type MerGroup struct {
	MerId    string `bson:"merId"`
	TransAmt int64  `bson:"transAmt"`
	Fee      int64  `bson:"fee"`
}

// TransTypeGroup 按单个商户交易类型分组
type TransTypeGroup struct {
	TransType int   `bson:"transType"`
	TransAmt  int64 `bson:"transAmt"`
	TransNum  int   `bson:"transNum"`
}

// SettRoleGroup 按清算角色分组
type SettRoleGroup struct {
	SettRole  string     `bson:"settRole"`
	MerGroups []MerGroup `bson:"mers"`
}

// NewScanPayRequest 带请求id的request对象
func NewScanPayRequest() *ScanPayRequest {
	return &ScanPayRequest{
		ReqId: util.SerialNumber(),
	}
}

// ScanPayRequest 扫码支付
type ScanPayRequest struct {
	Txndir       string `json:"txndir,omitempty" url:"txndir,omitempty" bson:"txndir,omitempty"`                   // 交易方向
	Busicd       string `json:"busicd,omitempty" url:"busicd,omitempty" bson:"busicd,omitempty"`                   // 交易类型
	AgentCode    string `json:"inscd,omitempty" url:"inscd,omitempty" bson:"inscd,omitempty"`                      // 代理/机构号
	Chcd         string `json:"chcd,omitempty" url:"chcd,omitempty" bson:"chcd,omitempty"`                         // 渠道机构
	Mchntid      string `json:"mchntid,omitempty" url:"mchntid,omitempty" bson:"mchntid,omitempty"`                // 商户号
	Terminalid   string `json:"terminalid,omitempty" url:"terminalid,omitempty" bson:"terminalid,omitempty"`       // 终端号
	Txamt        string `json:"txamt,omitempty" url:"txamt,omitempty" bson:"txamt,omitempty"`                      // 订单金额
	Currency     string `json:"currency,omitempty" url:"currency,omitempty" bson:"currency,omitempty"`             // 币种
	GoodsInfo    string `json:"goodsInfo,omitempty" url:"goodsInfo,omitempty" bson:"goodsInfo,omitempty"`          // 商品详情
	OrderNum     string `json:"orderNum,omitempty" url:"orderNum,omitempty" bson:"orderNum,omitempty"`             // 订单号
	OrigOrderNum string `json:"origOrderNum,omitempty" url:"origOrderNum,omitempty" bson:"origOrderNum,omitempty"` // 原订单号
	ScanCodeId   string `json:"scanCodeId,omitempty" url:"scanCodeId,omitempty" bson:"scanCodeId,omitempty"`       // 扫码号
	Sign         string `json:"sign,omitempty" url:"-" bson:"sign,omitempty" `                                     // 签名
	NotifyUrl    string `json:"backUrl,omitempty" url:"backUrl,omitempty" bson:"backUrl,omitempty"`                // 异步通知地址
	OpenId       string `json:"openid,omitempty" url:"openid,omitempty" bson:"openid,omitempty"`                   // openid
	CheckName    string `json:"checkName,omitempty" url:"checkName,omitempty" bson:"checkName,omitempty"`          // 校验用户姓名选项
	UserName     string `json:"userName,omitempty" url:"userName,omitempty" bson:"userName,omitempty"`             // 用户名
	Desc         string `json:"desc,omitempty" url:"desc,omitempty" bson:"desc,omitempty"`                         // 描述
	Code         string `json:"code,omitempty" url:"code,omitempty" bson:"code,omitempty"`                         // 认证码
	NeedUserInfo string `json:"needUserInfo,omitempty" url:"needUserInfo,omitempty" bson:"needUserInfo,omitempty"` // 是否需要获取用户信息
	VeriCode     string `json:"veriCode,omitempty" url:"veriCode,omitempty" bson:"veriCode,omitempty"`             // js支付用到的凭证
	Attach       string `json:"attach,omitempty" url:"attach,omitempty" bson:"attach,omitempty"`
	TimeExpire   string `json:"timeExpire,omitempty" url:"timeExpire,omitempty" bson:"timeExpire,omitempty"` // 过期时间

	TradeFrom    string `json:"tradeFrom,omitempty" url:"tradeFrom,omitempty" bson:"tradeFrom,omitempty"` // 交易来源
	SettDate     string `json:"settDate,omitempty" url:"settDate,omitempty" bson:"settDate,omitempty"`
	NextOrderNum string `json:"nextOrderNum,omitempty" url:"nextOrderNum,omitempty" bson:"nextOrderNum,omitempty"`

	CreateTime string `json:"-" url:"-" bson:"-"` // 卡券交易创建时间
	// 卡券相关字段
	VeriTime         string `json:"veriTime,omitempty" url:"veriTime,omitempty" bson:"veriTime,omitempty"`       // 核销次数 C
	Terminalsn       string `json:"terminalsn,omitempty" url:"terminalsn,omitempty" bson:"terminalsn,omitempty"` // 终端号
	Cardbin          string `json:"cardbin,omitempty" url:"cardbin,omitempty" bson:"cardbin,omitempty"`          // 银行卡cardbin或者用户标识等 C
	PayType          string `json:"payType,omitempty" url:"payType,omitempty" bson:"payType,omitempty"`          // 支付方式 M
	OrigChanOrderNum string `json:"-" url:"-" bson:"-"`                                                          // 辅助字段 原渠道订单号
	OrigSubmitTime   string `json:"-" url:"-" bson:"-"`                                                          // 辅助字段原交易提交时间
	OrigVeriTime     int    `json:"-" url:"-" bson:"-"`                                                          // 辅助字段 原交易验证时间
	IntPayType       int    `json:"-" url:"-" bson:"-"`                                                          // 辅助字段 核销次数
	IntVeriTime      int    `json:"-" url:"-" bson:"-"`

	// 微信需要的字段
	AppID      string `json:"-" url:"-" bson:"-"` // 公众号ID
	DeviceInfo string `json:"-" url:"-" bson:"-"` // 设备号
	SubMchId   string `json:"-" url:"-" bson:"-"` // 子商户
	TotalTxamt string `json:"-" url:"-" bson:"-"` // 订单总金额
	GoodsTag   string `json:"-" url:"-" bson:"-"` // 商品标识

	// 辅助字段
	Subject          string `json:"-" url:"-" bson:"-"` // 商品名称
	SysOrderNum      string `json:"-" url:"-" bson:"-"` // 渠道交易号
	ActTxamt         string `json:"-" url:"-" bson:"-"` // 实际交易金额 不同渠道单位不同
	IntTxamt         int64  `json:"-" url:"-" bson:"-"` // 以分为单位的交易金额
	ChanMerId        string `json:"-" url:"-" bson:"-"` // 渠道商户Id
	SignKey          string `json:"-" url:"-" bson:"-"` // 可能表示md5key等
	ExtendParams     string `json:"-" url:"-" bson:"-"` // 业务扩展参数
	WeixinClientCert []byte `json:"-" url:"-" bson:"-"` // 商户双向认证证书，如果是大商户模式，用大商户的证书
	WeixinClientKey  []byte `json:"-" url:"-" bson:"-"` // 商户双向认证密钥，如果是大商户模式，用大商户的密钥
	ReqId            string `json:"-" url:"-" bson:"-"`

	// 访问方式
	IsGBK bool     `json:"-" url:"-" bson:"-"`
	M     Merchant `json:"-" url:"-" bson:"-"`

	// //对账
	// SettleDate string `json:"settleDate,omitempty" url:"settleDate,omitempty" bson:"settleDate,omitempty"` // 对账日期 微信
	// StartTime  string `json:"startTime,omitempty" url:"startTime,omitempty" bson:"startTime,omitempty"`    // 对账开始时间 支付宝
	// EndTime    string `json:"endTime,omitempty" url:"endTime,omitempty" bson:"endTime,omitempty"`          // 对账结束时间 支付宝
}

// FillWithRequest 如果空白，默认将原信息返回
func (ret *ScanPayResponse) FillWithRequest(req *ScanPayRequest) {
	ret.Txndir = "A"

	if ret.Busicd == "" {
		ret.Busicd = req.Busicd
	}
	if ret.AgentCode == "" {
		ret.AgentCode = req.AgentCode
	}
	if ret.Chcd == "" {
		ret.Chcd = req.Chcd
	}
	if ret.Mchntid == "" {
		ret.Mchntid = req.Mchntid
	}
	if ret.Terminalid == "" {
		ret.Terminalid = req.Terminalid
	}
	if ret.Txamt == "" {
		ret.Txamt = req.Txamt
	}
	if ret.OrigOrderNum == "" {
		ret.OrigOrderNum = req.OrigOrderNum
	}
	if ret.OrderNum == "" {
		ret.OrderNum = req.OrderNum
	}
	if ret.VeriCode == "" {
		ret.VeriCode = req.VeriCode
	}
	// 如果是卡券核销，应答报文设置ScanCodeId
	if req.Busicd == Veri {
		ret.ScanCodeId = req.ScanCodeId
		ret.Terminalid = ""
	}
	if ret.PayType == "" {
		ret.PayType = req.PayType
	}
	if ret.PayType == "" {
		ret.PayType = req.PayType
	}
	if ret.Cardbin == "" {
		ret.Cardbin = req.Cardbin
	}

	ret.Attach = req.Attach
}

// ScanPayResponse 下单支付返回体
// M:返回时必须带上
// C:可选
type ScanPayResponse struct {
	Txndir          string   `json:"txndir,omitempty" url:"txndir,omitempty" bson:"txndir,omitempty"`                            // 交易方向 M M
	Busicd          string   `json:"busicd" url:"busicd" bson:"busicd"`                                                          // 交易类型 M M
	Respcd          string   `json:"respcd" url:"respcd" bson:"respcd"`                                                          // 交易结果  M
	AgentCode       string   `json:"inscd,omitempty" url:"inscd,omitempty" bson:"inscd,omitempty"`                               // 代理/机构号 M M
	Chcd            string   `json:"chcd,omitempty" url:"chcd,omitempty" bson:"chcd,omitempty"`                                  // 渠道 C C
	Mchntid         string   `json:"mchntid" url:"mchntid" bson:"mchntid"`                                                       // 商户号 M M
	Terminalid      string   `json:"terminalid,omitempty" url:"terminalid,omitempty" bson:"terminalid,omitempty"`                // 终端号
	Txamt           string   `json:"txamt,omitempty" url:"txamt,omitempty" bson:"txamt,omitempty"`                               // 订单金额 M M
	ChannelOrderNum string   `json:"channelOrderNum,omitempty" url:"channelOrderNum,omitempty" bson:"channelOrderNum,omitempty"` // 渠道交易号 C
	ConsumerAccount string   `json:"consumerAccount,omitempty" url:"consumerAccount,omitempty" bson:"consumerAccount,omitempty"` // 渠道账号  C
	ConsumerId      string   `json:"consumerId,omitempty" url:"consumerId,omitempty" bson:"consumerId,omitempty"`                // 渠道账号ID   C
	ErrorDetail     string   `json:"errorDetail,omitempty" url:"errorDetail,omitempty" bson:"errorDetail,omitempty"`             // 错误信息   C
	OrderNum        string   `json:"orderNum,omitempty" url:"orderNum,omitempty" bson:"orderNum,omitempty"`                      // 订单号 M C
	OrigOrderNum    string   `json:"origOrderNum,omitempty" url:"origOrderNum,omitempty" bson:"origOrderNum,omitempty"`          // 源订单号 M C
	Sign            string   `json:"sign" url:"-" bson:"sign"`                                                                   // 签名 M M
	ChcdDiscount    string   `json:"chcdDiscount,omitempty" url:"chcdDiscount,omitempty" bson:"chcdDiscount,omitempty"`          // 渠道优惠  C
	MerDiscount     string   `json:"merDiscount,omitempty" url:"merDiscount,omitempty" bson:"merDiscount,omitempty"`             // 商户优惠  C
	QrCode          string   `json:"qrcode,omitempty" url:"qrcode,omitempty" bson:"qrcode,omitempty"`                            // 二维码 C
	PayJson         *PayJson `json:"payjson,omitempty" url:"-" bson:"payjson,omitempty"`                                         // json字符串
	PayJsonStr      string   `json:"-" url:"payjson,omitempty" bson:"-"`                                                         // 签名时用
	VeriCode        string   `json:"veriCode,omitempty" url:"veriCode,omitempty" bson:"veriCode,omitempty"`
	GoodsInfo       string   `json:"goodsInfo,omitempty" url:"goodsInfo,omitempty" bson:"goodsInfo,omitempty"`
	Attach          string   `json:"attach,omitempty" url:"attach,omitempty" bson:"attach,omitempty"`

	Count        string      `json:"count,omitempty" url:"count,omitempty" bson:"count,omitempty"`
	Rec          interface{} `json:"rec,omitempty" url:"-" bson:"-"`
	RecStr       string      `json:"-" url:"rec,omitempty" bson:"-"`
	NextOrderNum string      `json:"nextOrderNum,omitempty" url:"nextOrderNum,omitempty" bson:"-"`

	ScanCodeId      string `json:"scanCodeId,omitempty" url:"scanCodeId,omitempty" bson:"scanCodeId,omitempty"`                // 扫码号 卡券核销M
	VeriTime        string `json:"veriTime,omitempty" url:"veriTime,omitempty" bson:"veriTime,omitempty"`                      // 核销次数 C
	CardId          string `json:"cardId,omitempty" url:"cardId,omitempty" bson:"cardId,omitempty"`                            // 卡券类型 C
	CardInfo        string `json:"cardInfo,omitempty" url:"cardInfo,omitempty" bson:"cardInfo,omitempty"`                      // 卡券详情 C
	AvailCount      string `json:"availCount,omitempty" url:"availCount,omitempty" bson:"availCount,omitempty"`                // 卡券剩余可用次数C
	ExpDate         string `json:"expDate,omitempty" url:"expDate,omitempty" bson:"expDate,omitempty"`                         // 卡券有效期 C
	Authcode        int    `json:"-" url:"-" bson:"-"`                                                                         // 授权码
	VoucherType     string `json:"voucherType,omitempty" url:"voucherType,omitempty" bson:"voucherType,omitempty"`             // 券类型 C
	SaleMinAmount   string `json:"saleMinAmount,omitempty" url:"saleMinAmount,omitempty" bson:"saleMinAmount,omitempty"`       // 满足优惠条件的最小金额         C
	SaleDiscount    string `json:"saleDiscount,omitempty" url:"saleDiscount,omitempty" bson:"saleDiscount,omitempty"`          // 抵扣值 C
	TransAmount     string `json:"transAmount,omitempty" url:"transAmount,omitempty" bson:"transAmount,omitempty"`             // 交易原始金额 M
	PayType         string `json:"payType,omitempty" url:"payType,omitempty" bson:"payType,omitempty"`                         // 支付方式 M
	ActualPayAmount string `json:"actualPayAmount,omitempty" url:"actualPayAmount,omitempty" bson:"actualPayAmount,omitempty"` // 实际支付金额 M
	ChannelTime     string `json:"-" url:"-" bson:"-"`                                                                         // 渠道处理时间
	Cardbin         string `json:"cardbin,omitempty" url:"cardbin,omitempty" bson:"cardbin,omitempty"`                         // 银行卡cardbin或者用户标识等 C
	OrigRespcd      string `json:"origRespcd,omitempty" url:"origRespcd,omitempty" bson:"origRespcd,omitempty"`                // 原交易结果 C
	OrigErrorDetail string `json:"origErrorDetail,omitempty" url:"origErrorDetail,omitempty" bson:"origErrorDetail,omitempty"` // 原错误信息   C

	// 辅助字段
	ChanRespCode string `json:"-" url:"-" bson:"-"` // 渠道详细应答码
	PrePayId     string `json:"-" url:"-" bson:"-"`
	ErrorCode    string `json:"-" url:"-" bson:"-"`
	PayTime      string `json:"-" url:"-" bson:"-"`
	Rate         string `json:"-" url:"-" bson:"-"` // 币种费率
}

// PayJson 公众号支付字段
type PayJson struct {
	UserInfo *weixin.AuthUserInfoResp `json:"userinfo,omitempty"`
	Config   *JsConfig                `json:"config,omitempty"`
	WxpPay   *JsWxpPay                `json:"chooseWXPay,omitempty"`
}

type JsConfig struct {
	AppID     string `json:"appId" url:"appId"`
	NonceStr  string `json:"nonceStr" url:"nonceStr"`
	Signature string `json:"signature" url:"-"`
	Timestamp string `json:"timestamp" url:"timestamp"` // 中间不用大写
}

type JsWxpPay struct {
	AppID     string `json:"appId" url:"appId"`
	NonceStr  string `json:"nonceStr" url:"nonceStr"`
	Package   string `json:"package" url:"package"`
	PaySign   string `json:"paySign" url:"-"`
	TimeStamp string `json:"timeStamp" url:"timeStamp"` //支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
	SignType  string `json:"signType" url:"signType"`
}

// SignMsg 字典排序报文
func (s *ScanPayRequest) SignMsg() string {
	return genSignMsg(s)
}

// SignMsg 字典排序报文
func (ret *ScanPayResponse) SignMsg() string {
	return genSignMsg(ret)
}

// GetMerReqLogs 商户请求日志体
func (s *ScanPayRequest) GetMerReqLogs() *SpTransLogs {
	return s.newTransLogs("in", 1, s)
}

// GetMerRetLogs 返回商户报文体
func (s *ScanPayRequest) GetMerRetLogs(ret interface{}) *SpTransLogs {
	return s.newTransLogs("out", 1, ret)
}

// GetChanReqLogs 请求渠道日志体
func (s *ScanPayRequest) GetChanReqLogs(reqData interface{}) *SpTransLogs {
	return s.newTransLogs("out", 2, reqData)
}

// GetChanRetLogs 请求渠道日志体
func (s *ScanPayRequest) GetChanRetLogs(reqData interface{}) *SpTransLogs {
	return s.newTransLogs("in", 2, reqData)
}

func (s *ScanPayRequest) newTransLogs(direct string, mt int, data interface{}) *SpTransLogs {
	return &SpTransLogs{
		Direction: direct, MsgType: mt, Msg: data,
		ReqId: s.ReqId, MerId: s.Mchntid, OrderNum: s.OrderNum,
		OrigOrderNum: s.OrigOrderNum, TransType: s.Busicd,
	}
}

func (s *ScanPayRequest) WxpMarshalGoods() string {

	goods, err := marshalGoods(s.GoodsInfo)
	if err != nil {
		// 格式不对，送配置的商品名称，防止商户送的内容过长
		return s.Subject
	}

	var goodsName []string
	if len(goods) > 0 {
		for _, v := range goods {
			goodsName = append(goodsName, v.GoodsName)
		}

		body := strings.Join(goodsName, ",")
		bodySizes := []rune(body)
		if len(bodySizes) > 20 {
			body = string(bodySizes[:20]) + "..."
		}
		return body
	}

	// 假如商品详细为空，送配置的商品名称
	return s.Subject
}

func (s *ScanPayRequest) AlpMarshalGoods() string {

	goods, err := marshalGoods(s.GoodsInfo)
	if err != nil {
		return ""
	}

	if len(goods) > 0 {
		gbs, err := json.Marshal(goods)
		if err != nil {
			log.Errorf("goodsInfo marshal error:%s", err)
			return ""
		}
		return string(gbs)
	}

	return ""
}

// marshalGoods 将商品详情解析
// 格式: 商品名称,价格,数量;商品名称,价格,数量;...
func marshalGoods(goodsInfo string) ([]goodsDetail, error) {

	var gs []goodsDetail

	if strings.TrimSpace(goodsInfo) == "" {
		return gs, nil
	}

	goods := strings.Split(goodsInfo, ";")
	for i, v := range goods {
		if i == len(goods)-1 && v == "" {
			continue
		}
		good := strings.Split(v, ",")
		if len(good) != 3 {
			return nil, fmt.Errorf("%s", "goodsInfo format error")
		}
		g := goodsDetail{
			GoodsId: i, GoodsName: good[0], Price: good[1], Quantity: good[2],
		}
		gs = append(gs, g)
	}

	return gs, nil
}

type goodsDetail struct {
	GoodsId   int    `json:"goodsId"`
	GoodsName string `json:"goodsName"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
}

// genSignMsg 获取字符串签名字段
func genSignMsg(o interface{}) string {

	buf, err := util.Query(o)
	if err != nil {
		log.Errorf("gen sign msg error: %s", err)
		return ""
	}
	return buf.String()
}

// NewScanPayResponse 构造方法
// 默认使用8583应答
func NewScanPayResponse(s ScanPayRespCode) *ScanPayResponse {
	return &ScanPayResponse{
		Respcd:      s.ISO8583Code,
		ErrorDetail: s.ISO8583Msg,
		ErrorCode:   s.ErrorCode,
	}
}

// ScanPayRespCode 扫码支付应答
type ScanPayRespCode struct {
	RespCode    string `bson:"respCode"`
	RespMsg     string `bson:"respMsg"`
	ISO8583Code string `bson:"ISO8583Code"`
	ISO8583Msg  string `bson:"ISO8583Msg"`
	IsUseISO    bool   `bson:"isUseISO"`
	ErrorCode   string `bson:"errorCode"`
}

// PublicAccount 公众号参数
type PublicAccount struct {
	ChanMerId string `bson:"chanMerId"`
	AppID     string `bson:"appID"`
	AppSecret string `bson:"appSecret"`
}

/* only use for import respCode */

// SpChanCSV 渠道文件csv
type SpChanCSV struct {
	Code    string `bson:"code"`
	Msg     string `bson:"msg"`
	Busicd  string `bson:"busicd"`
	ISOCode string `bson:"-"`
	ISOMsg  string `bson:"-"`
}

// ScanPayCSV 扫码支付应答码
type ScanPayCSV struct {
	RespCode    string `bson:"respCode"`
	RespMsg     string `bson:"respMsg"`
	ISO8583Code string `bson:"ISO8583Code"`
	ISO8583Msg  string `bson:"ISO8583Msg"`
	IsUseISO    bool   `bson:"isUseISO"`
	ErrorCode   string `bson:"errorCode"`

	Alp []*SpChanCSV `bson:"alp,omitempty"`
	Wxp []*SpChanCSV `bson:"wxp,omitempty"`
	//...
}

// MerSettStatus 商户清算状态
type MerSettStatus struct {
	MerId  string `bson:"merId" json:"merId"`
	Status int    `bson:"status" json:"status"`
}

// RoleSett 清算角色级别信息
type RoleSett struct {
	SettRole   string `json:"settRole" bson:"settRole"`
	SettDate   string `json:"settDate" bson:"settDate"`
	ReportName string `json:"reportName" bson:"reportName"`
	UpdateTime string `json:"updateTime" bson:"updateTime"`
	ReportType int    `json:"reportType" bson:"reportType"`
}

// ChanBlendMap 渠道勾兑数据集合
// 外部key为渠道商户号，内部key为渠道订单号
type ChanBlendMap map[string]map[string][]BlendElement

// LocalBlendMap 系统本地勾兑数据集合
// 外部key为渠道商户号，内部key为渠道订单号
type LocalBlendMap map[string]map[string][]TransSett

// 勾兑结构体
type BlendElement struct {
	Chcd      string //渠道编号
	ChcdName  string //渠道名称
	MerID     string //商户号
	ChanMerID string //渠道商户号
	MerName   string //商户名称
	LocalID   string //系统订单号
	OrderID   string //渠道订单号
	OrderTime string //交易时间
	OrderType string //交易类型
	OrderAct  string //交易金额
	IsBlend   bool   //对账标识
}
