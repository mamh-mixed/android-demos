package weixin

import (
	"crypto/tls"
	"encoding/xml"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
)

// NotifyURL 微信异步通知
const NotifyURL = "/scanpay/upNotify/weixin"

// BaseReq 只是为了注入签名方便
type BaseReq interface {
	SetSign(sign string) // 设置签名 setter
	GetSignKey() string  // 取商户（可能是大商户）签名密钥
	GetURI() string      // GetURI 取接口地址
	GetHTTPClient() *http.Client
	GetSpReq() *model.ScanPayRequest
}

// CommonParams 微信接口请求公共参数
type CommonParams struct {
	XMLName xml.Name `xml:"xml" url:"-" bson:"-"`

	// 公共字段
	Appid    string `xml:"appid,omitempty" url:"appid,omitempty"`           // 微信分配的公众账号ID validate:"len=18
	SubAppid string `xml:"sub_appid,omitempty" url:"sub_appid,omitempty"`   // 微信分配的子商户公众账号ID
	MchID    string `xml:"mch_id,omitempty" url:"mch_id,omitempty"`         // 微信支付分配的商户号
	SubMchId string `xml:"sub_mch_id,omitempty" url:"sub_mch_id,omitempty"` // 微信支付分配的子商户号，开发者模式下必填
	NonceStr string `xml:"nonce_str" url:"nonce_str" validate:"nonzero"`    // 随机字符串
	Sign     string `xml:"sign" url:"-"`                                    // 签名

	WeixinMD5Key string `xml:"-" url:"-" validate:"nonzero" bson:"-"`

	ClientCert []byte                `xml:"-" url:"-" bson:"-"` // HTTPS 双向认证证书
	ClientKey  []byte                `xml:"-" url:"-" bson:"-"` // HTTPS 双向认证密钥
	Req        *model.ScanPayRequest `xml:"-" url:"-" bson:"-"`
}

// SetSign sign setter
func (c *CommonParams) SetSign(sign string) {
	c.Sign = sign
}

// GetSignKey signKey getter
func (c *CommonParams) GetSignKey() string {
	return c.WeixinMD5Key
}

// GetSpReq 请求对象
func (c *CommonParams) GetSpReq() *model.ScanPayRequest {
	return c.Req
}

// GetHTTPClient 如果组合结构体不重写这个方法，表示不使用双向 HTTPS 认证，
// 如果使用双向 HTTPS 认证，重写此方法`return GetHTTPSClient()` 即可
func (c *CommonParams) GetHTTPClient() *http.Client {
	// return c.GetHTTPSClient()
	return getDefaultWeixinClient()
}

// GetHTTPSClient 使用双向 HTTPS 认证
func (c *CommonParams) GetHTTPSClient() (cli *http.Client) {
	if len(c.ClientCert) == 0 || len(c.ClientKey) == 0 {
		log.Error("client cert and key must not blank")
		return getDefaultWeixinClient()
	}

	cliCrt, err := tls.X509KeyPair(c.ClientCert, c.ClientKey)
	if err != nil {
		log.Errorf("X509KeyPair err: %s", err)
		return getDefaultWeixinClient()
	}

	return getPrivateWeixinClient(&cliCrt)
}

// BaseResp 只是为了传参方便
type BaseResp interface {
	GetSign() string
}

// CommonBody 微信接口返回公共字段
type CommonBody struct {
	XMLName xml.Name `xml:"xml" url:"-" bson:"-"`

	ReturnCode string `xml:"return_code" url:"return_code"`                   // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty" url:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	Appid      string `xml:"appid,omitempty" url:"appid,omitempty"`               // 公众账号ID
	MchID      string `xml:"mch_id,omitempty" url:"mch_id,omitempty"`             // 商户号
	SubMchId   string `xml:"sub_mch_id,omitempty" url:"sub_mch_id,omitempty"`     // 子商户号
	SubAppid   string `xml:"sub_appid,omitempty" url:"sub_appid,omitempty"`       // 子商户公众账号 ID
	NonceStr   string `xml:"nonce_str" url:"nonce_str"`                           // 随机字符串
	Sign       string `xml:"sign" url:"-"`                                        // 签名
	ResultCode string `xml:"result_code" url:"result_code"`                       // 业务结果
	ErrCode    string `xml:"err_code,omitempty" url:"err_code,omitempty"`         // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty" url:"err_code_des,omitempty"` // 错误代码描述
}

// GetSign sign getter
func (c *CommonBody) GetSign() string {
	return c.Sign
}

// WeixinNotifyReq 支付完成后，微信会把相关支付结果和用户信息发送给商户，商户需要接收处理，并返回应答
type WeixinNotifyReq struct {
	CommonBody

	DeviceInfo       string `xml:"device_info,omitempty" url:"device_info,omitempty" `                // 设备号
	OpenID           string `xml:"openid" url:"openid"`                                               // 用户标识
	IsSubscribe      string `xml:"is_subscribe" url:"is_subscribe"`                                   // 是否关注公众账号
	TradeType        string `xml:"trade_type" url:"trade_type"`                                       // 交易类型
	BankType         string `xml:"bank_type" url:"bank_type"`                                         // 付款银行
	FeeType          string `xml:"fee_type" url:"fee_type"`                                           // 货币类型
	TotalFee         string `xml:"total_fee" url:"total_fee"`                                         // 总金额
	CashFeeType      string `xml:"cash_fee_type" url:"cash_fee_type"`                                 // 现金支付货币类型
	CashFee          string `xml:"cash_fee" url:"cash_fee"`                                           // 现金支付金额
	CouponFee        string `xml:"coupon_fee" url:"coupon_fee"`                                       // 代金券或立减优惠金额
	CouponCount      string `xml:"coupon_count" url:"coupon_count"`                                   // 代金券或立减优惠使用数量
	TransactionId    string `xml:"transaction_id" url:"transaction_id"`                               // 微信支付订单号
	OutTradeNo       string `xml:"out_trade_no" url:"out_trade_no"`                                   // 商户订单号
	Attach           string `xml:"attach" url:"attach"`                                               // 商家数据包
	TimeEnd          string `xml:"time_end" url:"time_end"`                                           // 支付完成时间
	SubOpenid        string `xml:"sub_openid,omitempty" url:"sub_openid,omitempty"`                   // 子商户 Open ID
	SubIsSubscribe   string `xml:"sub_is_subscribe" url:"sub_is_subscribe,omitempty"`                 // 是否关注子商户公众账号
	CouponRefundId0  string `xml:"coupon_refund_id_0,omitempty" url:"coupon_refund_id_0,omitempty"`   // 代金券或立减优惠ID
	CouponRefundId1  string `xml:"coupon_refund_id_1,omitempty" url:"coupon_refund_id_1,omitempty"`   // 代金券或立减优惠ID
	CouponRefundId2  string `xml:"coupon_refund_id_2,omitempty" url:"coupon_refund_id_2,omitempty"`   // 代金券或立减优惠ID
	CouponRefundFee0 string `xml:"coupon_refund_fee_0,omitempty" url:"coupon_refund_fee_0,omitempty"` // 代金券或立减优惠退款金额
	CouponRefundFee1 string `xml:"coupon_refund_fee_1,omitempty" url:"coupon_refund_fee_1,omitempty"` // 代金券或立减优惠退款金额
	CouponRefundFee2 string `xml:"coupon_refund_fee_2,omitempty" url:"coupon_refund_fee_2,omitempty"` // 代金券或立减优惠退款金额
}

// WeixinNotifyResp 商户需要接收处理，并返回应答
type WeixinNotifyResp struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code" url:"return_code"`                   // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty" url:"return_msg,omitempty"` // 返回信息
}
