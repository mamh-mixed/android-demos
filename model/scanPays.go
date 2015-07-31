package model

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
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
	Txndir       string `json:"txndir,omitempty"`       // 交易方向
	Busicd       string `json:"busicd,omitempty"`       // 交易类型
	Inscd        string `json:"inscd,omitempty"`        // 机构号
	Chcd         string `json:"chcd,omitempty"`         // 渠道机构
	Mchntid      string `json:"mchntid,omitempty"`      // 商户号
	Terminalid   string `json:"terminalid,omitempty"`   // 终端号
	Txamt        string `json:"txamt,omitempty"`        // 订单金额
	Currency     string `json:"currency,omitempty"`     // 币种
	GoodsInfo    string `json:"goodsInfo,omitempty"`    // 商品详情
	OrderNum     string `json:"orderNum,omitempty"`     // 订单号
	OrigOrderNum string `json:"origOrderNum,omitempty"` // 原订单号
	ScanCodeId   string `json:"scanCodeId,omitempty"`   // 扫码号
	Sign         string `json:"sign,omitempty"`         // 签名
	NotifyUrl    string `json:"notifyUrl,omitempty"`    // 异步通知地址
	OpenId       string `json:"openid,omitempty"`       // openid
	CheckName    string `json:"checkName,omitempty"`    // 校验用户姓名选项
	UserName     string `json:"userName,omitempty"`     // 用户名
	Desc         string `json:"desc,omitempty"`         // 描述

	// 微信需要的字段
	AppID      string `json:"-"` // 公众号ID
	DeviceInfo string `json:"-"` // 设备号
	GoodsDesc  string `json:"-"` // 商品描述
	Attach     string `json:"-"` // 附加数据
	CurrType   string `json:"-"` // 货币类型
	GoodsGag   string `json:"-"` // 商品标记
	SubMchId   string `json:"-"` // 子商户
	TotalTxamt string `json:"-"` // 订单总金额

	// 辅助字段
	Subject          string `json:"-"` // 商品名称
	SysOrderNum      string `json:"-"` // 渠道交易号
	ActTxamt         string `json:"-"` // 实际交易金额 不同渠道单位不同
	IntTxamt         int64  `json:"-"` // 以分为单位的交易金额
	ChanMerId        string `json:"-"` // 渠道商户Id
	SignCert         string `json:"-"` // 可能表示md5key等
	WeixinClientCert []byte `json:"-"` // 商户双向认证证书，如果是大商户模式，用大商户的证书
	WeixinClientKey  []byte `json:"-"` // 商户双向认证密钥，如果是大商户模式，用大商户的密钥
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
	Txndir          string `json:"txndir"`                    // 交易方向 M M
	Busicd          string `json:"busicd"`                    // 交易类型 M M
	Respcd          string `json:"respcd"`                    // 交易结果  M
	Inscd           string `json:"inscd,omitempty"`           // 机构号 M M
	Chcd            string `json:"chcd,omitempty"`            // 渠道 C C
	Mchntid         string `json:"mchntid"`                   // 商户号 M M
	Terminalid      string `json:"terminalid,omitempty"`      // 终端号
	Txamt           string `json:"txamt,omitempty"`           // 订单金额 M M
	ChannelOrderNum string `json:"channelOrderNum,omitempty"` // 渠道交易号 C
	ConsumerAccount string `json:"consumerAccount,omitempty"` // 渠道账号  C
	ConsumerId      string `json:"consumerId,omitempty"`      // 渠道账号ID   C
	ErrorDetail     string `json:"errorDetail,omitempty"`     // 错误信息   C
	OrderNum        string `json:"orderNum,omitempty"`        // 订单号 M C
	OrigOrderNum    string `json:"origOrderNum,omitempty"`    // 源订单号 M C
	Sign            string `json:"sign"`                      // 签名 M M
	ChcdDiscount    string `json:"chcdDiscount,omitempty"`    // 渠道优惠  C
	MerDiscount     string `json:"merDiscount,omitempty"`     // 商户优惠  C
	QrCode          string `json:"qrcode,omitempty"`          // 二维码 C
	// 辅助字段
	ChanRespCode string `json:"-"` // 渠道详细应答码
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

	var mFields []string
	sv := reflect.ValueOf(o)
	if sv.Kind() != reflect.Ptr || sv.IsNil() {
		return ""
	}
	t := sv.Type().Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Name == "Sign" {
			continue
		}
		mFields = append(mFields, f.Name)
	}

	// 排序
	sort.Strings(mFields)
	var buf bytes.Buffer
	for _, field := range mFields {

		// jsonTag
		v := sv.Elem().FieldByName(field)
		f, _ := t.FieldByName(field)
		jsonTag := f.Tag.Get("json")

		// fieldName
		jn := ""
		jn = strings.Split(jsonTag, ",")[0]
		if jn == "" {
			jn = field
		}

		// 组装报文
		if jsonTag != "-" && v.CanSet() {
			if v.Kind() == reflect.String {
				fv := v.String()
				if fv != "" {
					if buf.Len() > 0 {
						buf.WriteByte('&')
					}
					buf.WriteString(jn + "=" + fv)
				}
			}
		}
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
