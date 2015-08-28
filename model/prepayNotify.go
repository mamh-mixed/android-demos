package model

// NotifyRecord 存储异步消息通知记录
type NotifyRecord struct {
	MerId       string `bson:"merId"`
	OrderNum    string `bson:"orderNum"`
	FromChanMsg string `bson:"fromChanMsg"`        // 渠道异步消息内容
	ToMerMsg    string `bson:"toMerMsg,omitempty"` // 系统发给商户异步消息内容
	IsToMerFail bool   `bson:"isToMerFail"`        // 是否发送失败
	CreateTime  string `bson:"createTime"`
	UpdateTime  string `bson:"updateTime,omitempty"`
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
