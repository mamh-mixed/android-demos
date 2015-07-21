package scanpay

import (
	"crypto/rsa"
	"encoding/json"
)

// PreCreateReq 撤销订单
type PreCreateReq struct {
	CommonParams

	Subject     string `json:"subject" validate:"nonzero"`      // 订单标题
	OutTradeNo  string `json:"out_trade_no" validate:"nonzero"` // 原支付请求的商户订单号
	TotalAmount string `json:"total_amount" validate:"nonzero"` // 原支付请求的商户订单号
}

// NewPreCreateReq 建议调用这个构造函数，以免漏掉默认参数
func NewPreCreateReq(appID string, privateKey *rsa.PrivateKey) *PreCreateReq {
	return &PreCreateReq{
		CommonParams: NewCommonParams("alipay.trade.precreate", appID, privateKey),
	}
}

// PreCreateBody 撤销订单报文
type PreCreateBody struct {
	Sign string          `json:"sign"`                            // 签名
	Raw  json.RawMessage `json:"alipay_trade_precreate_response"` // 返回消息体
}

// GetSign 撤销订单报文签名
func (d *PreCreateBody) GetSign() string {
	return d.Sign
}

// GetRaw 撤销订单报文内容
func (d *PreCreateBody) GetRaw() []byte {
	return []byte(d.Raw)
}

// PreCreateResp 撤销订单
type PreCreateResp struct {
	Code       string `json:"code"`                   // 结果码
	Msg        string `json:"msg"`                    // 结果码描述
	SubCode    string `json:"sub_code,omitempty"`     // 错误子代码
	SubMsg     string `json:"sub_msg,omitempty" `     // 错误子代码描述
	TradeNo    string `json:"trade_no,omitempty"`     // 支付宝交易号
	OutTradeNo string `json:"out_trade_no,omitempty"` // 商户订单号
	RetryFlag  string `json:"retry_flag,omitempty"`   // 商户订单号
	Action     string `json:"action,omitempty"`       // 商户订单号
	Close      string `json:"close,omitempty"`        // 商户订单号
}
