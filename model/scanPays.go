package model

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/quickpay/weixin"
	"github.com/omigo/log"
	"strings"
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
)

// QueryCondition 扫码交易查询字段
type QueryCondition struct {
	MerName      string   `json:"mchntName,omitempty"` // 可用于商户名称、商户简称模糊查询
	MerId        string   `json:"mchntid,omitempty"`   // 可用于商户号模糊查询
	MerIds       []string `json:"-"`
	Col          string   `json:"-"`
	BindingId    string   `json:"bindingId"`
	AgentCode    string   `json:"agentCode,omitempty"`
	TransStatus  []string `json:"transStatus,omitempty"`
	RefundStatus int      `json:"refundStatus,omitempty"`
	TransType    int      `json:"transType,omitempty"`
	StartTime    string   `json:"startTime,omitempty"`
	EndTime      string   `json:"endTime,omitempty"`
	Busicd       string   `json:"busicd,omitempty"`
	OrderNum     string   `json:"orderNum,omitempty"`
	OrigOrderNum string   `json:"origOrderNum,omitempty"`
	NextOrderNum string   `json:"nextOrderNum,omitempty"`
	Count        int      `json:"count,omitempty"`
	Page         int      `json:"page,omitempty"`
	Total        int      `json:"total,omitempty"`
	Size         int      `json:"size,omitempty"`
	IsForReport  bool     `json:"-"`
	GroupCode    string   `json:"groupCode,omitempty"`
	Respcd       string   `json:"respcd" url:"respcd"`
	RespcdNotIn  string   `json:"respcdNotIn"`
	TradeFrom    string   `json:"tradeFrom,omitempty"`
	Skip         int      `json:"skip,omitempty"`
	Direction    string
	ReqIds       []string
}

// QueryResult 查询结果值
type QueryResult struct {
	Rec      interface{} `json:"rec,omitempty"` // 交易明细
	Page     int         `json:"page,omitempty"`
	Total    int         `json:"total,omitempty"`
	Size     int         `json:"size,omitempty"`
	Count    int         `json:"count,omitempty"`
	RespCode string      `json:"respCode,omitempty"`
	RespMsg  string      `json:"respMsg,omitempty"`
}

// Summary 商户交易汇总
type Summary struct {
	MerId         string  `json:"merId,omitempty"`
	MerName       string  `json:"merName,omitempty"`
	AgentName     string  `json:"agentName,omitempty"`
	TotalTransAmt float32 `json:"totalTransAmt"`
	TotalTransNum int     `json:"totalTransNum"`
	TotalFee      float32 `json:"-"`
	Wxp           struct {
		TransAmt float32 `json:"transAmt"`
		TransNum int     `json:"transNum"`
		Fee      float32 `json:"-"`
	} `json:"wxp"`
	Alp struct {
		TransAmt float32 `json:"transAmt"`
		TransNum int     `json:"transNum"`
		Fee      float32 `json:"-"`
	} `json:"alp"`
	Data []Summary `json:"data"` // 包含每个商户单独数据
}

// TransGroup 按商户号和渠道号统计分组
type TransGroup struct {
	MerName   string    `bson:"merName"`
	AgentName string    `bson:"agentName"`
	TransAmt  int64     `bson:"transAmt"`
	RefundAmt int64     `bson:"refundAmt"`
	TransNum  int       `bson:"transNum"`
	Detail    []Channel `bson:"detail"`
	MerId     string    `bson:"_id"`
	Fee       int64     `bson:"fee" json:"-"`
}

// Channel 按渠道类型分组
type Channel struct {
	ChanCode  string `bson:"chanCode"`
	TransAmt  int64  `bson:"transAmt"`
	RefundAmt int64  `bson:"refundAmt"`
	TransNum  int    `bson:"transNum"`
	Fee       int64  `bson:"fee" json:"-"`
}

// TransTypeGroup 按单个商户交易类型分组
type TransTypeGroup struct {
	TransType int   `bson:"transType"`
	TransAmt  int64 `bson:"transAmt"`
	TransNum  int   `bson:"transNum"`
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
	IsGBK     bool   `json:"-" url:"-" bson:"-"`
	TradeFrom string `json:"tradeFrom,omitempty" url:"tradeFrom,omitempty" bson:"tradeFrom,omitempty"` // 交易来源
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
	// 辅助字段
	ChanRespCode string `json:"-" url:"-" bson:"-"` // 渠道详细应答码
	PrePayId     string `json:"-" url:"-" bson:"-"`
	ErrorCode    string `json:"-" url:"-" bson:"-"`
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
		return s.GoodsInfo // 假如格式错误，直接透传给微信
	}

	var goodsName []string
	if len(goods) > 0 {
		for _, v := range goods {
			goodsName = append(goodsName, v.GoodsName)
		}
		return strings.Join(goodsName, ",")
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
