package scanpay1

import (
	"github.com/CardInfoLink/quickpay/model"
)

const (
	Gw = "https://intlmapi.alipay.com/gateway.do?"
)

type BaseReq interface {
	GetURI() string     // 请求网关
	GetSignKey() string // 签名密钥
	GetSpReq() *model.ScanPayRequest
}

type BaseResp interface {
	ReqFlag() bool
	ResultCode() string
	ErrorCode() string
}

type CommonReq struct {
	Sign         string                `url:"sign,omitempty"`
	SignType     string                `url:"Sign_type,omitempty"`
	Partner      string                `url:"partner,omitempty"`
	InputCharset string                `url:"_input_charset,omitempty"`
	SignKey      string                `url:"-"`
	SpReq        *model.ScanPayRequest `url:"-"`
}

type CommonResp struct {
	IsSuccess string  `url:"is_success" xml:"is_success"`
	Error     string  `url:"error" xml:"error"`
	SignType  string  `url:"sign_type" xml:"sign_type"`
	Sign      string  `url:"sign" xml:"sign"`
	Request   []Param `xml:"request>param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}

func (c *CommonReq) GetSpReq() *model.ScanPayRequest {
	return c.SpReq
}

func (c *CommonReq) GetSignKey() string {
	return c.SignKey
}

func (c *CommonReq) GetURI() string {
	return Gw
}

func (c *CommonResp) ReqFlag() bool {
	if c.IsSuccess == "T" {
		return true
	} else {
		return false
	}
}
