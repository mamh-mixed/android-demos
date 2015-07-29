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
}

// AlipayNotifyResp 商户需要接收处理，并返回应答
type AlipayNotifyResp struct {
}
