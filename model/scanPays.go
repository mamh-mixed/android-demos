package model

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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

// ScanPay 扫码支付
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

	// 微信需要的字段
	AppID      string // 公众号ID
	DeviceInfo string // 设备号
	GoodsDesc  string // 商品描述
	Attach     string // 附加数据
	CurrType   string // 货币类型
	GoodsGag   string // 商品标记
	SubMchId   string // 子商户
	TotalTxamt string // 订单总金额

	// 辅助字段
	Subject     string `json:"-"` //  商品名称
	SysOrderNum string `json:"-"` //  渠道交易号
	ActTxamt    string `json:"-"` //  实际交易金额 不同渠道单位不同
	ChanMerId   string `json:"-"` // 渠道商户Id
	SignCert    string `json:"-"` // 可能表示md5key等
	IntTxamt    int64  `json:"-"`
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

// DictSortMsg 字典排序报文
func (s *ScanPayRequest) SignMsg() string {
	return genSignMsg(s)
}

// DictSortMsg 字典排序报文
func (s *ScanPayResponse) SignMsg() string {
	return genSignMsg(s)
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

// WeixinNotifyReq 支付完成后，微信会把相关支付结果和用户信息发送给商户，商户需要接收处理，并返回应答
type WeixinNotifyReq struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息

	// 当return_code为SUCCESS的时候，还会包括以下字段：
	Appid      string `xml:"appid"`                  // 公众账号ID
	MchID      string `xml:"mch_id"`                 // 商户号
	SubMchId   string `xml:"sub_mch_id"`             // 子商户号（文档没有该字段）
	NonceStr   string `xml:"nonce_str"`              // 随机字符串
	Sign       string `xml:"sign"`                   // 签名
	ResultCode string `xml:"result_code"`            // 业务结果
	ErrCode    string `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"` // 错误代码描述

	// 以上为微信接口公共字段

	// 当return_code 和result_code都为SUCCESS的时，还会包括以下字段：
	DeviceInfo    string `xml:"device_info,omitempty"` // 设备号
	OpenID        string `xml:"openid"`                // 用户标识
	IsSubscribe   string `xml:"is_subscribe"`          // 是否关注公众账号
	TradeType     string `xml:"trade_type"`            // 交易类型
	BankType      string `xml:"bank_type"`             // 付款银行
	FeeType       string `xml:"fee_type"`              // 货币类型
	TotalFee      string `xml:"total_fee"`             // 总金额
	CashFeeType   string `xml:"cash_fee_type"`         // 现金支付货币类型
	CashFee       string `xml:"cash_fee"`              // 现金支付金额
	CouponFee     string `xml:"coupon_fee"`            // 代金券或立减优惠金额
	CouponCount   string `xml:"coupon_count"`          // 代金券或立减优惠使用数量
	TransactionId string `xml:"transaction_id"`        // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`          // 商户订单号
	Attach        string `xml:"attach"`                // 商家数据包
	TimeEnd       string `xml:"time_end"`              // 支付完成时间

}

// WeixinNotifyResp 商户需要接收处理，并返回应答
type WeixinNotifyResp struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息
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
	RespCode      string `bson:"respCode"`
	RespMsg       string `bson:"respMsg"`
	ISO8583Code   string `bson:"ISO8583Code"`
	ISO8583Msg    string `bson:"ISO8583Msg"`
	IsUseChanDesc bool   `bson:"isUseChanDesc"`

	Alp []*SpChanCSV `bson:"alp,omitempty"`
	Wxp []*SpChanCSV `bson:"wxp,omitempty"`
	//...
}
