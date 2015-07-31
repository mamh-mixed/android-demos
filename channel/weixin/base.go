package weixin

import (
	"crypto/tls"
	"encoding/xml"
	"net/http"

	"github.com/omigo/log"
)

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
	Appid    string `xml:"appid" url:"appid" validate:"len=18"`            // 微信分配的公众账号ID
	SubAppid string `xml:"sub_appid,omitempty" url:"sub_appid,omitempty"`  // 微信分配的子商户公众账号ID
	MchID    string `xml:"mch_id" url:"mch_id" validate:"nonzero"`         // 微信支付分配的商户号
	SubMchId string `xml:"sub_mch_id" url:"sub_mch_id" validate:"nonzero"` // 微信支付分配的子商户号，开发者模式下必填
	NonceStr string `xml:"nonce_str" url:"nonce_str" validate:"nonzero"`   // 随机字符串
	Sign     string `xml:"sign" url:"-"`                                   // 签名

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
	Appid      string `xml:"appid" url:"appid"`                                   // 公众账号ID
	MchID      string `xml:"mch_id" url:"mch_id"`                                 // 商户号
	SubMchId   string `xml:"sub_mch_id" url:"sub_mch_id"`                         // 子商户号
	SubAppid   string `xml:"sub_appid" url:"sub_appid"`                           // 子商户公众账号 ID
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
