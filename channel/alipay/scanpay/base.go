package scanpay

import (
	"crypto/rsa"
	"net/url"
	"time"

	"github.com/omigo/log"
	"github.com/omigo/validator"
)

// BaseReq  统一签名和组报文
type BaseReq interface {
	Values() url.Values          // 组装公共参数
	PrivateKey() *rsa.PrivateKey // 商户 RSA 私钥
}

// CommonParams 组装公共参数
type CommonParams struct {
	appID     string // 支付宝服务窗 APPID
	method    string // 接口名称
	charset   string // 请求使用的编码格式，如utf-8,gbk,gb2312等
	signType  string // 商户生成签名字符串所使用的签名算法类型，目前支持RSA,DSA
	sign      string // 商户请求参数的签名串，详见安全规范中的签名生成算法
	timestamp string // 发送请求的时间，格式“yyyy-MM-dd HH:mm:ss”
	version   string // 调用的接口版本，固定为:1.0

	privateKey *rsa.PrivateKey // 商户 RSA 私钥
}

// NewCommonParams 外部只能调用这个构造函数
func NewCommonParams(method, appID string, privateKey *rsa.PrivateKey) CommonParams {
	return CommonParams{
		appID:      appID,
		method:     method,
		charset:    "utf-8",
		signType:   "RSA",
		timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		version:    "1.0",
		privateKey: privateKey,
	}
}

// Values 组装公共参数
func (c *CommonParams) Values() (v url.Values) {
	v = url.Values{}
	v.Set("app_id", c.appID)
	v.Set("charset", c.charset)
	v.Set("method", c.method)
	v.Set("sign_type", c.signType)
	v.Set("timestamp", c.timestamp)
	v.Set("version", c.version)
	return v
}

// PrivateKey 商户 RSA 私钥
func (c *CommonParams) PrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

// BaseBody 应答报文
type BaseBody interface {
	GetSign() string
	GetRaw() []byte
}

// BaseResp 应答报文内容
type BaseResp interface{}

func base(req BaseReq, body BaseBody, resp BaseResp) (err error) {
	if err := validator.Validate(req); err != nil {
		log.Errorf("validate error, %s", err)
		return err
	}

	err = sendRequest(req, body, resp)
	if err != nil {
		log.Errorf("alipay request error: %s", err)
		return err
	}
	return nil
}
