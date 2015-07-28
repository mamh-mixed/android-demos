package scanpay

import (
	"encoding/json"
	"net/url"
)

// https://app.alipay.com/market/document.htm?name=saomazhifu#page-16
// 应用场景实例
// 调用支付宝支付接口时未返回明确的返回结果(如由于系统错误或网络异常导 致无返回结果)，可使用本接口将交易进行撤销。
// 如果用户支付失败，支付宝会将此订单关闭；如果用户支付成功，支付宝会将 支付的资金退还给用户。

// CancelReq 撤销订单
type CancelReq struct {
	CommonParams

	OutTradeNo string `json:"out_trade_no" validate:"nonzero"` // 原支付请求的商户订单号
}

// Values 组装公共参数
func (c *CancelReq) Values() (v url.Values) {
	c.CommonParams.Method = "alipay.trade.cancel"
	return c.CommonParams.Values()
}

// CancelResp 撤销订单
type CancelResp struct {
	CommonBody

	Raw json.RawMessage `json:"alipay_trade_cancel_response"` // 返回消息体

	TradeNo    string `json:"trade_no"`         // 支付宝交易号
	OutTradeNo string `json:"out_trade_no"`     // 原支付请求的商户订单号
	RetryFlag  string `json:"retry_flag"`       // 撤销已成功，无需重试
	Action     string `json:"action,omitempty"` // 撤销执行的动作，close：直接撤销，无退款；refund：撤销，有退款。
}

// GetRaw 报文内容
func (c *CancelResp) GetRaw() []byte {
	return []byte(c.Raw)
}
