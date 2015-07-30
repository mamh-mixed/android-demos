package scanpay

// CloseReq 关闭订单
type CloseReq struct {
	CommonParams

	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`      // 商户订单号
}

// CloseResp 撤销订单
type CloseResp struct {
	CommonBody
}
