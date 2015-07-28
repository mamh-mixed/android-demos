package scanpay

import (
	"encoding/json"
	"net/url"
)

// 参考文档https://app.alipay.com/market/document.htm?name=saomazhifu#page-15

// 本接口提供支付宝支付订单的查询的功能，商户可以通过本接口主动查询订单状态，完成下一步的业务逻辑。 需要调用查询接口的情况：
// 当商户后台、网络、服务器等出现异常，商户系统最终未接收到支付通知；
// 调用扫码支付支付接口后，返回系统错误或未知交易状态情况
// 调用扫码支付请求后，如果结果返回处理中（返回结果中的code等于10003）的状态；
// 调用撤销接口API之前，需确认该笔交易目前支付状态；

// QueryReq 查询订单
type QueryReq struct {
	CommonParams

	TradeNo    string `json:"trade_no,omitempty"`     // 支付宝交易号
	OutTradeNo string `json:"out_trade_no,omitempty"` // 原支付请求的商户订单号

}

// Values 组装公共参数
func (c *QueryReq) Values() (v url.Values) {
	c.CommonParams.Method = "alipay.trade.query"
	return c.CommonParams.Values()
}

// QueryResp 查询订单
type QueryResp struct {
	CommonBody

	Raw json.RawMessage `json:"alipay_trade_query_response"` // 返回消息体

	TradeNo        string `json:"trade_no"`                   // 支付宝交易号
	OutTradeNo     string `json:"out_trade_no"`               // 商户订单号
	OpenID         string `json:"open_id"`                    // 买家支付宝用户号
	BuyerLogonID   string `json:"buyer_logon_id"`             // 买家支付宝账号，将用*号屏蔽部分内容
	TradeStatus    string `json:"trade_status"`               // 交易状态
	TotalAmount    string `json:"total_amount"`               // 订单金额
	ReceiptAmount  string `json:"receipt_amount,omitempty"`   // 商家实收金额
	InvoiceAmount  string `json:"invoice_amount,omitempty"`   // 开票金额
	BuyerPayAmount string `json:"buyer_pay_amount,omitempty"` // 付款金额
	PointAmount    string `json:"point_amount,omitempty"`     // 积分宝金额
	SendPayDate    string `json:"send_pay_date,omitempty"`    // 本次交易打款给卖家的时间,格式为 yyyy-MM-dd HH:mm:ss
	TerminalID     string `json:"terminal_id,omitempty"`      // 商户机具终端编号
	AlipayStoreID  string `json:"alipay_store_id,omitempty"`  // 支付宝店铺编号
	StoreID        string `json:"store_id,omitempty"`         // 商户门店编号
	FundBillList   []struct {
		FundChannel string `json:"fund_channel,omitempty"` // 支付渠道,例 COUPON、DISCOUNT
		Amount      string `json:"amount,omitempty"`       // 支付金额
	} `json:"store_id,omitempty"` // 资金单据信息的集合
}

// GetRaw 报文内容
func (c *QueryResp) GetRaw() []byte {
	return []byte(c.Raw)
}
