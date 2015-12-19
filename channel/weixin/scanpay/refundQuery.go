package scanpay

import "github.com/CardInfoLink/quickpay/channel/weixin"

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_5
// 应用场景
// 提交退款申请后，通过调用该接口查询退款状态。退款有一定延时，用零钱支付的退款20分钟内到账，银行卡支付的退款3个工作日后重新查询退款状态。
// 接口链接：https://api.mch.weixin.qq.com/pay/refundquery
// 是否需要证书: 不需要

// RefundQueryReq 查询退款
type RefundQueryReq struct {
	weixin.CommonParams

	DeviceInfo    string `xml:"device_info,omitempty" url:"device_info,omitempty"`       // 设备号
	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`      // 商户订单号
	OutRefundNo   string `xml:"out_refund_no,omitempty" url:"out_refund_no,omitempty"`   // 商户退款单号
	RefundId      string `xml:"refund_id,omitempty" url:"refund_id,omitempty"`           // 微信退款单号
}

// GetURI 取接口地址
func (p *RefundQueryReq) GetURI() string {
	return "/pay/refundquery"
}

// RefundQueryResp 查询退款
type RefundQueryResp struct {
	weixin.CommonBody

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

	OutRefundNo0           string `xml:"out_refund_no_0" url:"out_refund_no_0,omitempty"`                                 // 商户退款单号
	OutRefundNo1           string `xml:"out_refund_no_1" url:"out_refund_no_1,omitempty"`                                 // 商户退款单号
	OutRefundNo2           string `xml:"out_refund_no_2" url:"out_refund_no_2,omitempty"`                                 // 商户退款单号
	RefundID0              string `xml:"refund_id_0" url:"refund_id_0,omitempty"`                                         // 微信退款单号
	RefundID1              string `xml:"refund_id_1" url:"refund_id_1,omitempty"`                                         // 微信退款单号
	RefundID2              string `xml:"refund_id_2" url:"refund_id_2,omitempty"`                                         // 微信退款单号
	RefundChannel0         string `xml:"refund_channel_0,omitempty" url:"refund_channel_0,omitempty"`                     // 退款渠道
	RefundChannel1         string `xml:"refund_channel_1,omitempty" url:"refund_channel_1,omitempty"`                     // 退款渠道
	RefundChannel2         string `xml:"refund_channel_2,omitempty" url:"refund_channel_2,omitempty"`                     // 退款渠道
	RefundFee0             string `xml:"refund_fee_0" url:"refund_fee_0,omitempty"`                                       // 退款金额
	RefundFee1             string `xml:"refund_fee_1" url:"refund_fee_1,omitempty"`                                       // 退款金额
	RefundFee2             string `xml:"refund_fee_2" url:"refund_fee_2,omitempty"`                                       // 退款金额
	FeeType0               string `xml:"fee_type_0,omitempty" url:"fee_type_0,omitempty"`                                 // 货币种类
	FeeType1               string `xml:"fee_type_1,omitempty" url:"fee_type_1,omitempty"`                                 // 货币种类
	FeeType2               string `xml:"fee_type_2,omitempty" url:"fee_type_2,omitempty"`                                 // 货币种类
	CouponRefundFee1       string `xml:"coupon_refund_fee_1,omitempty" url:"coupon_refund_fee_1,omitempty"`               // 代金券或立减优惠退款金额
	CouponRefundFee2       string `xml:"coupon_refund_fee_2,omitempty" url:"coupon_refund_fee_2,omitempty"`               // 代金券或立减优惠退款金额
	CouponRefundFee3       string `xml:"coupon_refund_fee_3,omitempty" url:"coupon_refund_fee_3,omitempty"`               // 代金券或立减优惠退款金额
	CouponRefundCount1     string `xml:"coupon_refund_count_1,omitempty" url:"coupon_refund_count_1,omitempty"`           // 代金券或立减优惠使用数量
	CouponRefundCount2     string `xml:"coupon_refund_count_2,omitempty" url:"coupon_refund_count_2,omitempty"`           // 代金券或立减优惠使用数量
	CouponRefundCount3     string `xml:"coupon_refund_count_3,omitempty" url:"coupon_refund_count_3,omitempty"`           // 代金券或立减优惠使用数量
	CouponRefundBatchID1_1 string `xml:"coupon_refund_batch_id_1_1,omitempty" url:"coupon_refund_batch_id_1_1,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID1_2 string `xml:"coupon_refund_batch_id_1_2,omitempty" url:"coupon_refund_batch_id_1_2,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID1_3 string `xml:"coupon_refund_batch_id_1_3,omitempty" url:"coupon_refund_batch_id_1_3,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID2_1 string `xml:"coupon_refund_batch_id_2_1,omitempty" url:"coupon_refund_batch_id_2_1,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID2_2 string `xml:"coupon_refund_batch_id_2_2,omitempty" url:"coupon_refund_batch_id_2_2,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID2_3 string `xml:"coupon_refund_batch_id_2_3,omitempty" url:"coupon_refund_batch_id_2_3,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID3_1 string `xml:"coupon_refund_batch_id_3_1,omitempty" url:"coupon_refund_batch_id_3_1,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID3_2 string `xml:"coupon_refund_batch_id_3_2,omitempty" url:"coupon_refund_batch_id_3_2,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundBatchID3_3 string `xml:"coupon_refund_batch_id_3_3,omitempty" url:"coupon_refund_batch_id_3_3,omitempty"` // 代金券或立减优惠批次ID
	CouponRefundID1_1      string `xml:"coupon_refund_id_1_1,omitempty" url:"coupon_refund_id_1_1,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID1_2      string `xml:"coupon_refund_id_1_2,omitempty" url:"coupon_refund_id_1_2,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID1_3      string `xml:"coupon_refund_id_1_3,omitempty" url:"coupon_refund_id_1_3,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID2_1      string `xml:"coupon_refund_id_2_1,omitempty" url:"coupon_refund_id_2_1,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID2_2      string `xml:"coupon_refund_id_2_2,omitempty" url:"coupon_refund_id_2_2,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID2_3      string `xml:"coupon_refund_id_2_3,omitempty" url:"coupon_refund_id_2_3,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID3_1      string `xml:"coupon_refund_id_3_1,omitempty" url:"coupon_refund_id_3_1,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID3_2      string `xml:"coupon_refund_id_3_2,omitempty" url:"coupon_refund_id_3_2,omitempty"`             // 代金券或立减优惠ID
	CouponRefundID3_3      string `xml:"coupon_refund_id_3_3,omitempty" url:"coupon_refund_id_3_3,omitempty"`             // 代金券或立减优惠ID
	CouponRefundFee1_1     string `xml:"coupon_refund_fee_1_1,omitempty" url:"coupon_refund_fee_1_1,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee1_2     string `xml:"coupon_refund_fee_1_2,omitempty" url:"coupon_refund_fee_1_1,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee1_3     string `xml:"coupon_refund_fee_1_3,omitempty" url:"coupon_refund_fee_1_1,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee2_1     string `xml:"coupon_refund_fee_2_1,omitempty" url:"coupon_refund_fee_2_1,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee2_2     string `xml:"coupon_refund_fee_2_2,omitempty" url:"coupon_refund_fee_2_2,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee2_3     string `xml:"coupon_refund_fee_2_3,omitempty" url:"coupon_refund_fee_2_3,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee3_1     string `xml:"coupon_refund_fee_3_1,omitempty" url:"coupon_refund_fee_3_1,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee3_2     string `xml:"coupon_refund_fee_3_2,omitempty" url:"coupon_refund_fee_3_2,omitempty"`           // 单个代金券或立减优惠支付金额
	CouponRefundFee3_3     string `xml:"coupon_refund_fee_3_3,omitempty" url:"coupon_refund_fee_3_3,omitempty"`           // 单个代金券或立减优惠支付金额
	RefundStatus0          string `xml:"refund_status_0,omitempty" url:"refund_status_0,omitempty"`                       // 退款状态
	RefundStatus1          string `xml:"refund_status_1,omitempty" url:"refund_status_1,omitempty"`                       // 退款状态
	RefundStatus2          string `xml:"refund_status_2,omitempty" url:"refund_status_2,omitempty"`                       // 退款状态

}
