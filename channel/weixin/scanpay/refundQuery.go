package scanpay

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_5
// 应用场景
// 提交退款申请后，通过调用该接口查询退款状态。退款有一定延时，用零钱支付的退款20分钟内到账，银行卡支付的退款3个工作日后重新查询退款状态。
// 接口链接：https://api.mch.weixin.qq.com/pay/refundquery
// 是否需要证书: 不需要

// RefundQueryReq 查询退款
type RefundQueryReq struct {
	CommonParams

	DeviceInfo    string `xml:"device_info,omitempty" url:"device_info,omitempty"`       // 设备号
	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`      // 商户订单号
	OutRefundNo   string `xml:"out_refund_no,omitempty" url:"out_refund_no,omitempty"`   // 商户退款单号
	RefundId      string `xml:"refund_id,omitempty" url:"refund_id,omitempty"`           // 微信退款单号
}

// RefundQueryResp 查询退款
type RefundQueryResp struct {
	CommonBody

	DeviceInfo      string `xml:"device_info,omitempty" url:"device_info,omitempty"`             // 设备号
	TransactionId   string `xml:"transaction_id" url:"transaction_id,omitempty"`                 // 微信订单号
	OutTradeNo      string `xml:"out_trade_no" url:"out_trade_no,omitempty"`                     // 商户订单号
	TotalFee        int    `xml:"total_fee" url:"total_fee,omitempty"`                           // 订单总金额
	FeeType         string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`                   // 订单金额货币种类
	CashFee         int    `xml:"cash_fee" url:"cash_fee,omitempty"`                             // 现金支付金额
	CashFeeType     string `xml:"cash_fee_type,omitempty" url:"cash_fee_type,omitempty"`         // 货币种类
	RefundFee       int    `xml:"refund_fee" url:"refund_fee,omitempty"`                         // 退款金额
	CouponRefundFee int    `xml:"coupon_refund_fee,omitempty" url:"coupon_refund_fee,omitempty"` // 代金券或立减优惠退款金额
	RefundCount     int    `xml:"refund_count" url:"refund_count,omitempty"`                     // 退款笔数

	OutRefundNo0           string `xml:"out_refund_no_0" url:"out_refund_no_0,omitempty"`
	RefundID0              string `xml:"refund_id_0" url:"refund_id_0,omitempty"`
	RefundChannel0         string `xml:"refund_channel_0,omitempty" url:"refund_channel_0,omitempty"`
	RefundFee0             int    `xml:"refund_fee_0" url:"refund_fee_0,omitempty"`
	FeeType0               string `xml:"fee_type_0,omitempty" url:"fee_type_0,omitempty"`
	CouponRefundFee1       int    `xml:"coupon_refund_fee_1,omitempty" url:"coupon_refund_fee_1,omitempty"`
	CouponRefundCount1     int    `xml:"coupon_refund_count_1,omitempty" url:"coupon_refund_count_1,omitempty"`
	CouponRefundBatchID1_1 string `xml:"coupon_refund_batch_id_1_1,omitempty" url:"coupon_refund_batch_id_1_1,omitempty"`
	CouponRefundID1_1      string `xml:"coupon_refund_id_1_1,omitempty" url:"coupon_refund_id_1_1,omitempty"`
	CouponRefundFee1_1     int    `xml:"coupon_refund_fee_1_1,omitempty" url:"coupon_refund_fee_1_1,omitempty"`
	RefundStatus0          string `xml:"refund_status_0,omitempty" url:"refund_status_0,omitempty"`
}
