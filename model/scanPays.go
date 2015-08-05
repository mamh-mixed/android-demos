package model

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/util"
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
	Qyfk = "QYFK" // 企业付款
)

// QueryCondition 扫码交易查询
type QueryCondition struct {
	Mchntid      string  `json:"mchntid,omitempty"`
	StartTime    string  `json:"startTime,omitempty"`
	EndTime      string  `json:"endTime,omitempty"`
	Busicd       string  `json:"busicd,omitempty"`
	OrderNum     string  `json:"orderNum,omitempty"`
	OrigOrderNum string  `json:"origOrderNum,omitempty"`
	NextOrderNum string  `json:"nextOrderNum,omitempty"`
	RespCode     string  `json:"respCode,omitempty"`
	RespMsg      string  `json:"respMsg,omitempty"`
	Count        int     `json:"count,omitempty"`
	Page         int     `json:"page,omitempty"`
	Total        int     `json:"total,omitempty"`
	Size         int     `json:"size,omitempty"`
	Rec          []Trans `json:"rec,omitempty"`
}

// ScanPayRequest 扫码支付
type ScanPayRequest struct {
	Txndir       string `json:"txndir,omitempty" url:"txndir,omitempty"`             // 交易方向
	Busicd       string `json:"busicd,omitempty" url:"busicd,omitempty"`             // 交易类型
	Inscd        string `json:"inscd,omitempty" url:"inscd,omitempty"`               // 机构号
	Chcd         string `json:"chcd,omitempty" url:"chcd,omitempty"`                 // 渠道机构
	Mchntid      string `json:"mchntid,omitempty" url:"mchntid,omitempty"`           // 商户号
	Terminalid   string `json:"terminalid,omitempty" url:"terminalid,omitempty"`     // 终端号
	Txamt        string `json:"txamt,omitempty" url:"txamt,omitempty"`               // 订单金额
	Currency     string `json:"currency,omitempty" url:"currency,omitempty"`         // 币种
	GoodsInfo    string `json:"goodsInfo,omitempty" url:"goodsInfo,omitempty"`       // 商品详情
	OrderNum     string `json:"orderNum,omitempty" url:"orderNum,omitempty"`         // 订单号
	OrigOrderNum string `json:"origOrderNum,omitempty" url:"origOrderNum,omitempty"` // 原订单号
	ScanCodeId   string `json:"scanCodeId,omitempty" url:"scanCodeId,omitempty"`     // 扫码号
	Sign         string `json:"sign,omitempty" url:"-"`                              // 签名
	NotifyUrl    string `json:"notifyUrl,omitempty"url:"notifyUrl,omitempty" `       // 异步通知地址
	OpenId       string `json:"openid,omitempty" url:"openid,omitempty" `            // openid
	CheckName    string `json:"checkName,omitempty" url:"checkName,omitempty"`       // 校验用户姓名选项
	UserName     string `json:"userName,omitempty" url:"userName,omitempty"`         // 用户名
	Desc         string `json:"desc,omitempty" url:"desc,omitempty"`                 // 描述

	// 微信需要的字段
	AppID      string `json:"-" url:"-"` // 公众号ID
	DeviceInfo string `json:"-" url:"-"` // 设备号
	SubMchId   string `json:"-" url:"-"` // 子商户
	TotalTxamt string `json:"-" url:"-"` // 订单总金额

	// 辅助字段
	Subject          string `json:"-" url:"-"` // 商品名称
	SysOrderNum      string `json:"-" url:"-"` // 渠道交易号
	ActTxamt         string `json:"-" url:"-"` // 实际交易金额 不同渠道单位不同
	IntTxamt         int64  `json:"-" url:"-"` // 以分为单位的交易金额
	ChanMerId        string `json:"-" url:"-"` // 渠道商户Id
	SignCert         string `json:"-" url:"-"` // 可能表示md5key等
	WeixinClientCert []byte `json:"-" url:"-"` // 商户双向认证证书，如果是大商户模式，用大商户的证书
	WeixinClientKey  []byte `json:"-" url:"-"` // 商户双向认证密钥，如果是大商户模式，用大商户的密钥
}

// FillWithRequest 如果空白，默认将原信息返回
func (ret *ScanPayResponse) FillWithRequest(req *ScanPayRequest) {
	ret.Txndir = "A"

	if ret.Busicd == "" {
		ret.Busicd = req.Busicd
	}
	if ret.Inscd == "" {
		ret.Inscd = req.Inscd
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
}

// ScanPayResponse 下单支付返回体
// M:返回时必须带上
// C:可选
type ScanPayResponse struct {
	Txndir          string `json:"txndir" url:"txndir"`                                       // 交易方向 M M
	Busicd          string `json:"busicd" url:"busicd"`                                       // 交易类型 M M
	Respcd          string `json:"respcd" url:"respcd"`                                       // 交易结果  M
	Inscd           string `json:"inscd,omitempty" url:"inscd,omitempty"`                     // 机构号 M M
	Chcd            string `json:"chcd,omitempty" url:"chcd,omitempty"`                       // 渠道 C C
	Mchntid         string `json:"mchntid" url:"mchntid"`                                     // 商户号 M M
	Terminalid      string `json:"terminalid,omitempty" url:"terminalid,omitempty"`           // 终端号
	Txamt           string `json:"txamt,omitempty" url:"txamt,omitempty"`                     // 订单金额 M M
	ChannelOrderNum string `json:"channelOrderNum,omitempty" url:"channelOrderNum,omitempty"` // 渠道交易号 C
	ConsumerAccount string `json:"consumerAccount,omitempty" url:"consumerAccount,omitempty"` // 渠道账号  C
	ConsumerId      string `json:"consumerId,omitempty" url:"consumerId,omitempty"`           // 渠道账号ID   C
	ErrorDetail     string `json:"errorDetail,omitempty" url:"errorDetail,omitempty"`         // 错误信息   C
	OrderNum        string `json:"orderNum,omitempty" url:"orderNum,omitempty"`               // 订单号 M C
	OrigOrderNum    string `json:"origOrderNum,omitempty" url:"origOrderNum,omitempty"`       // 源订单号 M C
	Sign            string `json:"sign" url:"-"`                                              // 签名 M M
	ChcdDiscount    string `json:"chcdDiscount,omitempty" url:"chcdDiscount,omitempty"`       // 渠道优惠  C
	MerDiscount     string `json:"merDiscount,omitempty" url:"merDiscount,omitempty"`         // 商户优惠  C
	QrCode          string `json:"qrcode,omitempty" url:"qrcode,omitempty"`                   // 二维码 C
	PrePayId        string `json:"-" url:"-"`
	// 辅助字段
	ChanRespCode string `json:"-" url:"-"` // 渠道详细应答码
}

// SignMsg 字典排序报文
func (s *ScanPayRequest) SignMsg() string {
	return genSignMsg(s)
}

// SignMsg 字典排序报文
func (ret *ScanPayResponse) SignMsg() string {
	return genSignMsg(ret)
}

// MarshalGoods 将商品详情解析成字符json字符串
// 格式: 商品名称,价格,数量;商品名称,价格,数量;...
func (s *ScanPayRequest) MarshalGoods() string {

	if s.GoodsInfo == "" {
		return ""
	}

	goods := strings.Split(s.GoodsInfo, ";")
	gs := make([]interface{}, 0, len(goods))

	for i, v := range goods {
		good := strings.Split(v, ",")
		if len(good) != 3 {
			return ""
		}
		g := &struct {
			GoodsId   int    `json:"goodsId"`
			GoodsName string `json:"goodsName"`
			Price     string `json:"price"`
			Quantity  string `json:"quantity"`
		}{
			i, good[0], good[1], good[2],
		}
		gs = append(gs, g)
	}
	formated, err := json.Marshal(gs)
	if err != nil {
		return ""
	}
	return string(formated)
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
