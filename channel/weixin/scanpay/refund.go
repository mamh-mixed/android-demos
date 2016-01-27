package scanpay

import (
	"net/http"

	"github.com/CardInfoLink/quickpay/channel/weixin"
)

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_4
// 应用场景
// 当交易发生之后一段时间内，由于买家或者卖家的原因需要退款时，卖家可以通过退款接口将支付款退还给买家，微信支付将在收到退款请求并且验证成功之后，按照退款规则将支付款按原路退到买家帐号上。
// 注意：
// 1.交易时间超过半年的订单无法提交退款；
// 2.微信支付退款支持单笔交易分多次退款，多次退款需要提交原支付订单的商户订单号和设置不同的退款单号。一笔退款失败后重新提交，要采用原来的退款单号。总退款金额不能超过用户实际支付金额。
// 接口链接：https://api.mch.weixin.qq.com/secapi/pay/refund
// 是否需要证书
// 请求需要双向证书。

// RefundReq 申请退款
type RefundReq struct {
	weixin.CommonParams

	DeviceInfo    string `xml:"device_info,omitempty" url:"device_info,omitempty"`         // 设备号
	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"`   // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`        // 商户订单号
	OutRefundNo   string `xml:"out_refund_no" url:"out_refund_no" validate:"nonzero"`      // 商户退款单号
	TotalFee      string `xml:"total_fee" url:"total_fee" validate:"nonzero"`              // 总金额
	FeeType       string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`               // 标价币种
	RefundFee     string `xml:"refund_fee" url:"refund_fee" validate:"nonzero"`            // 退款金额
	RefundFeeType string `xml:"refund_fee_type,omitempty" url:"refund_fee_type,omitempty"` // 货币种类
	OpUserId      string `xml:"op_user_id" url:"op_user_id" validate:"nonzero"`            // 操作员
}

// GetURI 取接口地址
func (r *RefundReq) GetURI() string {
	return "/secapi/pay/refund"
}

// GetHTTPClient 使用双向 HTTPS 认证
func (r *RefundReq) GetHTTPClient() *http.Client {
	return r.GetHTTPSClient()
}

// RefundResp 申请退款
type RefundResp struct {
	weixin.CommonBody

	DeviceInfo        string `xml:"device_info,omitempty" url:"device_info,omitempty"`                 // 设备号
	TransactionId     string `xml:"transaction_id" url:"transaction_id,omitempty"`                     // 微信订单号
	OutTradeNo        string `xml:"out_trade_no" url:"out_trade_no,omitempty"`                         // 商户订单号
	OutRefundNo       string `xml:"out_refund_no" url:"out_refund_no,omitempty"`                       // 商户退款单号
	RefundId          string `xml:"refund_id" url:"refund_id,omitempty"`                               // 微信退款单号
	RefundChannel     string `xml:"refund_channel,omitempty" url:"refund_channel,omitempty"`           // 退款渠道
	RefundFee         string `xml:"refund_fee" url:"refund_fee,omitempty"`                             // 退款金额
	TotalFee          string `xml:"total_fee" url:"total_fee,omitempty"`                               // 订单总金额
	FeeType           string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`                       // 订单金额货币种类
	CashFee           string `xml:"cash_fee" url:"cash_fee,omitempty"`                                 // 现金支付金额
	CashRefundFee     string `xml:"cash_refund_fee,omitempty" url:"cash_refund_fee,omitempty"`         // 现金退款金额
	CouponRefundFee   string `xml:"coupon_refund_fee,omitempty" url:"coupon_refund_fee,omitempty"`     // 代金券或立减优惠退款金额
	CouponRefundCount string `xml:"coupon_refund_count,omitempty" url:"coupon_refund_count,omitempty"` // 代金券或立减优惠使用数量
	CouponRefundId0   string `xml:"coupon_refund_id_0,omitempty" url:"coupon_refund_id_0,omitempty"`   // 代金券或立减优惠ID
	CouponRefundFee0  string `xml:"coupon_refund_fee_0,omitempty" url:"coupon_refund_fee_0,omitempty"` // 代金券或立减优惠退款金额
	CouponRefundId1   string `xml:"coupon_refund_id_1,omitempty" url:"coupon_refund_id_1,omitempty"`   // 代金券或立减优惠ID
	CouponRefundFee1  string `xml:"coupon_refund_fee_1,omitempty" url:"coupon_refund_fee_1,omitempty"` // 代金券或立减优惠退款金额
	CouponRefundId2   string `xml:"coupon_refund_id_2,omitempty" url:"coupon_refund_id_2,omitempty"`   // 代金券或立减优惠ID
	CouponRefundFee2  string `xml:"coupon_refund_fee_2,omitempty" url:"coupon_refund_fee_2,omitempty"` // 代金券或立减优惠退款金额
}
