package model

import "encoding/xml"

// WeixinNotifyReq 支付完成后，微信会把相关支付结果和用户信息发送给商户，商户需要接收处理，并返回应答
type WeixinNotifyReq struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息

	// 当return_code为SUCCESS的时候，还会包括以下字段：
	Appid      string `xml:"appid"`                  // 公众账号ID
	MchID      string `xml:"mch_id"`                 // 商户号
	SubMchId   string `xml:"sub_mch_id"`             // 子商户号（文档没有该字段）
	NonceStr   string `xml:"nonce_str"`              // 随机字符串
	Sign       string `xml:"sign"`                   // 签名
	ResultCode string `xml:"result_code"`            // 业务结果
	ErrCode    string `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"` // 错误代码描述

	// 以上为微信接口公共字段

	// 当return_code 和result_code都为SUCCESS的时，还会包括以下字段：
	DeviceInfo    string `xml:"device_info,omitempty"` // 设备号
	OpenID        string `xml:"openid"`                // 用户标识
	IsSubscribe   string `xml:"is_subscribe"`          // 是否关注公众账号
	TradeType     string `xml:"trade_type"`            // 交易类型
	BankType      string `xml:"bank_type"`             // 付款银行
	FeeType       string `xml:"fee_type"`              // 货币类型
	TotalFee      string `xml:"total_fee"`             // 总金额
	CashFeeType   string `xml:"cash_fee_type"`         // 现金支付货币类型
	CashFee       string `xml:"cash_fee"`              // 现金支付金额
	CouponFee     string `xml:"coupon_fee"`            // 代金券或立减优惠金额
	CouponCount   string `xml:"coupon_count"`          // 代金券或立减优惠使用数量
	TransactionId string `xml:"transaction_id"`        // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`          // 商户订单号
	Attach        string `xml:"attach"`                // 商家数据包
	TimeEnd       string `xml:"time_end"`              // 支付完成时间

}

// WeixinNotifyResp 商户需要接收处理，并返回应答
type WeixinNotifyResp struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息
}

// AlipayNotifyReq 预下单用户支付完成后，支付宝会把相关支付结果和用户信息发送给商户，商户需要接收处理，并返回应答
type AlipayNotifyReq struct {
	NotifyTime       string `json:"notify_time" validate:"nonzero"`        // 通知时间
	NotifyType       string `json:"notify_type" validate:"nonzero"`        // 通知类型
	NotifyID         string `json:"notify_id" validate:"nonzero"`          // 通知校验ID
	SignType         string `json:"sign_type" validate:"nonzero"`          // 签名类型
	Sign             string `json:"sign" validate:"nonzero"`               // 签名
	NotifyActionType string `json:"notify_action_type" validate:"nonzero"` // 通知动作类型
	TradeNo          string `json:"trade_no" validate:"nonzero"`           // 支付宝交易号
	AppID            string `json:"app_id" validate:"nonzero"`             // 开发者的appid
	OutTradeNo       string `json:"out_trade_no,omitempty"`                // 商户订单号
	OutBizNo         string `json:"out_biz_no,omitempty"`                  // 商户业务号
	OpenID           string `json:"open_id,omitempty"`                     // 买家支付宝用户号
	BuyerLogonID     string `json:"buyer_logon_id,omitempty"`              // 买家支付宝账号
	SellerID         string `json:"seller_id,omitempty"`                   // 卖家支付宝用户号
	SellerEmail      string `json:"seller_email,omitempty"`                // 卖家支付宝账号
	TradeStatus      string `json:"trade_status,omitempty"`                // 交易状态
	TotalAmount      string `json:"total_amount,omitempty"`                // 订单金额
	ReceiptAmount    string `json:"receipt_amount,omitempty"`              // 实收金额
	InvoiceAmount    string `json:"invoice_amount,omitempty"`              // 开票金额
	BuyerPayAmount   string `json:"buyer_pay_amount,omitempty"`            // 付款金额
	PointAmount      string `json:"point_amount,omitempty"`                // 积分宝金额
	RefundFee        string `json:"refund_fee,omitempty"`                  // 退款金额
	Subject          string `json:"subject,omitempty"`                     // 订单标题
	Body             string `json:"body,omitempty"`                        // 商品描述
	GmtCreate        string `json:"gmt_create,omitempty"`                  // 交易创建时间
	GmtPayment       string `json:"gmt_payment,omitempty"`                 // 交易付款时间
	GmtRefund        string `json:"gmt_refund,omitempty"`                  // 交易退款时间
	GmtClose         string `json:"gmt_close,omitempty"`                   // 交易结束时间
	FundBillList     string `json:"fund_bill_list,omitempty"`              // 支付金额信息
}
