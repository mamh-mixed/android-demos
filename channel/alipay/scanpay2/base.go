package scanpay2

import (
	"crypto/rsa"
	"errors"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"github.com/omigo/validator"
	"net/url"
	"time"
)

// 编码、签名算法、版本
const (
	CharsetUTF8 = "utf-8"
	SignTypeRSA = "RSA"
	Version1_0  = "1.0"
)

// BaseReq  统一签名和组报文
type BaseReq interface {
	Values() url.Values             // 组装公共参数
	GetPrivateKey() *rsa.PrivateKey // 商户 RSA 私钥
	GetSpReq() *model.ScanPayRequest
	SaveLog() bool // 是否需要记录日志
}

// CommonParams 组装公共参数
type CommonParams struct {
	AppID        string `json:"-" validate:"nonzero"` // 支付宝服务窗 APPID
	Method       string `json:"-"`                    // 接口名称
	Charset      string `json:"-"`                    // 请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType     string `json:"-"`                    // 商户生成签名字符串所使用的签名算法类型，目前支持RSA,DSA
	Sign         string `json:"-"`                    // 商户请求参数的签名串，详见安全规范中的签名生成算法
	Timestamp    string `json:"-"`                    // 发送请求的时间，格式“yyyy-MM-dd HH:mm:ss”
	Version      string `json:"-"`                    // 调用的接口版本，固定为:1.0
	NotifyUrl    string `json:"-"`                    // 异步通知地址
	AppAuthToken string `json:"-"`                    // 是否授权

	Req        *model.ScanPayRequest `json:"-" bson:"-"`
	PrivateKey *rsa.PrivateKey       `json:"-" bson:"-"` // 商户 RSA 私钥
}

type GoodsDetail struct {
	GoodsId       string `json:"goods_id"`
	GoodsName     string `json:"goods_name"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	AlipayGoodsId string `json:"alipay_goods_id,omitempty"`
	GoodsCategory string `json:"goods_category",omitempty`
	Body          string `json:"body,omitempty"`
}

type Params struct {
	SysServiceProviderId string `json:"sys_service_provider_id,omitempty"`
}

// SaveLog 是否需要记录日志
func (c *CommonParams) SaveLog() bool {
	return true
}

// GetSpReq 前端请求报文
func (c *CommonParams) GetSpReq() *model.ScanPayRequest {
	return c.Req
}

// Values 组装公共参数
func (c *CommonParams) Values() (v url.Values) {
	v = url.Values{}

	v.Set("app_id", c.AppID)
	// 固定为 UTF-8 编码
	v.Set("charset", CharsetUTF8)
	v.Set("method", c.Method)

	if c.SignType == "" {
		v.Set("sign_type", SignTypeRSA)
	} else {
		v.Set("sign_type", c.SignType)
	}

	if c.Timestamp == "" {
		v.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	} else {
		v.Set("timestamp", c.Timestamp)
	}

	if c.Version == "" {
		v.Set("version", Version1_0)
	} else {
		v.Set("version", c.Version)
	}

	if c.NotifyUrl != "" {
		v.Set("notify_url", c.NotifyUrl)
	}

	if c.AppAuthToken != "" {
		v.Set("app_auth_token", c.AppAuthToken)
	}

	return v
}

// GetPrivateKey 商户 RSA 私钥
func (c *CommonParams) GetPrivateKey() *rsa.PrivateKey {
	return c.PrivateKey
}

// BaseResp 应答报文
type BaseResp interface {
	GetSign() string
	GetRaw() []byte
}

// CommonBody 公共返回参数
type CommonBody struct {
	Sign    string `json:"sign"`               // 签名
	Code    string `json:"code"`               // 结果码
	Msg     string `json:"msg"`                // 结果码描述
	SubCode string `json:"sub_code,omitempty"` // 错误子代码
	SubMsg  string `json:"sub_msg,omitempty" ` // 错误子代码描述
}

// GetSign 报文签名
func (c *CommonBody) GetSign() string {
	return c.Sign
}

// Execute 这个是扫码支付入口，所有请求，准备好参数后，调用此方法发送到支付宝
func Execute(req BaseReq, resp BaseResp) (err error) {

	if req.GetPrivateKey() == nil {
		return errors.New("private key is nil")
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("validate error, %s", err)
		return err
	}

	err = sendRequest(req, resp)
	if err != nil {
		log.Errorf("alipay request error: %s", err)
		return err
	}

	return nil
}
