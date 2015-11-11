package scanpay2

import (
	"encoding/json"
	"net/url"
)

// 参考文档https://app.alipay.com/market/document.htm?name=saomazhifu#page-17
// 当交易发生之后一段时间内，由于买家或者卖家的原因需退款，卖家可通过退 款接口将支付款退还给买家，支付宝将在收到退款请求并验证成功后，按退款规则 将支付款按原路退到买家帐号上。
// 交易超过可退款时间(签约时设置的可退款时间)的订单无法进行退款。
// 支付宝退款支持单笔交易分多次退款，多次退款需要提交支付宝交易号并 设置不同的退款单号；总退款金额不能超过用户实际支付金额。
// 分多笔退款时，若一笔退款失败需重新提交，要采用原来的退款单号。

// RefundReq 申请退款
type RefundReq struct {
	CommonParams

	TradeNo       string `json:"trade_no" validate:"nonzero"`      // 支付宝交易号
	RefundAmount  string `json:"refund_amount" validate:"nonzero"` // 退款金额
	OutRequestNo  string `json:"out_request_no,omitempty"`         // 商户退款请求号
	RefundReason  string `json:"refund_reason,omitempty"`          // 退款原因
	StoreID       string `json:"store_id,omitempty"`               // 商户的门店编号
	AlipayStoreID string `json:"alipay_store_id,omitempty"`        // 支付宝店铺编号
	TerminalID    string `json:"terminal_id,omitempty"`            // 商户的终端编号

}

// Values 组装公共参数
func (c *RefundReq) Values() (v url.Values) {
	c.CommonParams.Method = "alipay.trade.refund"
	return c.CommonParams.Values()
}

// RefundResp 申请退款
type RefundResp struct {
	CommonBody
	Raw json.RawMessage `json:"alipay_trade_refund_response"` // 返回消息体

	TradeNo              string `json:"trade_no,omitempty"`       // 支付宝交易号
	OutTradeNo           string `json:"out_trade_no,omitempty"`   // 商户订单号
	OpenID               string `json:"open_id,omitempty"`        // 买家支付宝用户号
	BuyerLogonID         string `json:"buyer_logon_id,omitempty"` // 买家支付宝账号
	FundChange           string `json:"fund_change"`              // 本次退款请求是否发生资金变动
	RefundFee            string `json:"fund_change,omitempty"`    // 累计退款金额
	GmtRefundPay         string `json:"gmt_refund_pay,omitempty"` // 退款时间
	RefundDetailItemList []struct {
		FundChannel string `json:"fund_channel,omitempty"` // 支付渠道,例 COUPON、DISCOUNT
		Amount      string `json:"amount,omitempty"`       // 支付金额
	} `json:"refund_detail_item_list"` // 退款资金明细信息集合
}

// GetRaw 报文内容
func (c *RefundResp) GetRaw() []byte {
	return []byte(c.Raw)
}
