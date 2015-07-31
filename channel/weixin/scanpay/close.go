package scanpay

import "github.com/CardInfoLink/quickpay/channel/weixin"

// CloseReq 关闭订单
type CloseReq struct {
	weixin.CommonParams

	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`      // 商户订单号
}

// GetURI 取接口地址
func (p *CloseReq) GetURI() string {
	return "/pay/closeorder"
}

// CloseResp 撤销订单
type CloseResp struct {
	weixin.CommonBody
}
