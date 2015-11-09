package scanpay1

import (
	"github.com/CardInfoLink/quickpay/model"
)

const (
	Gw = ""
)

type BaseReq interface {
	GetURI() string     // 请求网关
	GetSignKey() string // 签名密钥
	GetSpReq() *model.ScanPayRequest
}

type BaseResp interface {
}

type CommonParams struct {
	Service      string `url:"service"`
	Sign         string `url:"sign"`
	SignType     string `url:"Sign_type"`
	Partner      string `url:"partner"`
	InputCharset string `url:"_input_charset"`

	SignKey string                `url:"-"`
	SpReq   *model.ScanPayRequest `url:"-"`
}

func (c *CommonParams) GetSpReq() *model.ScanPayRequest {
	return c.SpReq
}

func (c *CommonParams) GetSignKey() string {
	return c.SignKey
}

func (c *CommonParams) GetURI() string {
	return Gw + ""
}
