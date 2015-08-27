package weixin

import (
	"crypto/tls"
	"encoding/xml"
	"net/http"

	"github.com/omigo/log"
)

const NotifyURL = "/scanpay/upNotify/weixin"

// BaseReq 只是为了注入签名方便
type BaseReq interface {
	SetSign(sign string) // 设置签名 setter
	GetSignKey() string  // 取商户（可能是大商户）签名密钥
	GetURI() string      // GetURI 取接口地址
	GetHTTPClient() *http.Client
}

// CommonParams 微信接口请求公共参数
type CommonParams struct {
	XMLName xml.Name `xml:"xml" url:"-"`

	// 公共字段
	Appid    string `xml:"appid,omitempty" url:"appid,omitempty"`           // 微信分配的公众账号ID validate:"len=18
	SubAppid string `xml:"sub_appid,omitempty" url:"sub_appid,omitempty"`   // 微信分配的子商户公众账号ID
	MchID    string `xml:"mch_id,omitempty" url:"mch_id,omitempty"`         // 微信支付分配的商户号
	SubMchId string `xml:"sub_mch_id,omitempty" url:"sub_mch_id,omitempty"` // 微信支付分配的子商户号，开发者模式下必填
	NonceStr string `xml:"nonce_str" url:"nonce_str" validate:"nonzero"`    // 随机字符串
	Sign     string `xml:"sign" url:"-"`                                    // 签名

	WeixinMD5Key string `xml:"-" url:"-" validate:"nonzero"`

	ClientCert []byte `xml:"-" url:"-"` // HTTPS 双向认证证书
	ClientKey  []byte `xml:"-" url:"-"` // HTTPS 双向认证密钥
}

// SetSign sign setter
func (c *CommonParams) SetSign(sign string) {
	c.Sign = sign
}

// GetSignKey signKey getter
func (c *CommonParams) GetSignKey() string {
	return c.WeixinMD5Key
}

// GetHTTPClient 如果组合结构体不重写这个方法，表示不使用双向 HTTPS 认证，
// 如果使用双向 HTTPS 认证，重写此方法`return GetHTTPSClient()` 即可
func (c *CommonParams) GetHTTPClient() *http.Client {
	// return GetHTTPSClient()
	return http.DefaultClient
}

// GetHTTPSClient 使用双向 HTTPS 认证
func (c *CommonParams) GetHTTPSClient() (cli *http.Client) {
	if len(c.ClientCert) == 0 || len(c.ClientKey) == 0 {
		log.Error("client cert and key must not blank")
		return http.DefaultClient
	}

	cliCrt, err := tls.X509KeyPair(c.ClientCert, c.ClientKey)
	if err != nil {
		log.Errorf("X509KeyPair err: %s", err)
		return http.DefaultClient
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// InsecureSkipVerify: true, // only for testing
			Certificates: []tls.Certificate{cliCrt}},
	}
	cli = &http.Client{Transport: tr}

	return cli
}

// BaseResp 只是为了传参方便
type BaseResp interface {
	GetSign() string
}

// CommonBody 微信接口返回公共字段
type CommonBody struct {
	XMLName xml.Name `xml:"xml" url:"-"`

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

	DeviceInfo     string `xml:"device_info,omitempty"`                             // 设备号
	OpenID         string `xml:"openid"`                                            // 用户标识
	IsSubscribe    string `xml:"is_subscribe"`                                      // 是否关注公众账号
	TradeType      string `xml:"trade_type"`                                        // 交易类型
	BankType       string `xml:"bank_type"`                                         // 付款银行
	FeeType        string `xml:"fee_type"`                                          // 货币类型
	TotalFee       string `xml:"total_fee"`                                         // 总金额
	CashFeeType    string `xml:"cash_fee_type"`                                     // 现金支付货币类型
	CashFee        string `xml:"cash_fee"`                                          // 现金支付金额
	CouponFee      string `xml:"coupon_fee"`                                        // 代金券或立减优惠金额
	CouponCount    string `xml:"coupon_count"`                                      // 代金券或立减优惠使用数量
	TransactionId  string `xml:"transaction_id"`                                    // 微信支付订单号
	OutTradeNo     string `xml:"out_trade_no"`                                      // 商户订单号
	Attach         string `xml:"attach"`                                            // 商家数据包
	TimeEnd        string `xml:"time_end"`                                          // 支付完成时间
	SubOpenid      string `xml:"sub_openid,omitempty" url:"sub_openid,omitempty"`   // 子商户 Open ID
	SubIsSubscribe string `xml:"sub_is_subscribe" url:"sub_is_subscribe,omitempty"` // 是否关注子商户公众账号
}

// WeixinNotifyResp 商户需要接收处理，并返回应答
type WeixinNotifyResp struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息
}
