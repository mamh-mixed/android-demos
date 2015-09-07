package scanpay

import "github.com/CardInfoLink/quickpay/channel/weixin"

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_10&index=1
// 应用场景
// 收银员使用扫码设备读取微信用户刷卡授权码以后，二维码或条码信息传送至商户收银台，由商户收银台或者商户后台调用该接口发起支付。
// 提醒1：提交支付请求后微信会同步返回支付结果。当返回结果为“系统错误”时，商户系统等待5秒后调用【查询订单API】，查询支付实际交易结果；当返回结果为“USERPAYING”时，商户系统可设置间隔时间(建议10秒)重新查询支付结果，直到支付成功或超时(建议30秒)；
// 提醒2：在调用查询接口返回后，如果交易状况不明晰，请调用【撤销订单API】，此时如果交易失败则关闭订单，该单不能再支付成功；如果交易成功，则将扣款退回到用户账户。当撤销无返回或错误时，请再次调用。注意：请勿调用扣款后立即调用【撤销订单API】。撤销订单API需要双向证书。
// 接口地址 https://api.mch.weixin.qq.com/pay/micropay
// 是否需要证书: 不需要

// PayReq 请求被扫支付API需要提交的数据
type PayReq struct {
	weixin.CommonParams

	DeviceInfo     string `xml:"device_info,omitempty" url:"device_info,omitempty"`          // 设备号
	Body           string `xml:"body" url:"body" validate:"nonzero"`                         // 商品描述
	Detail         string `xml:"detail,omitempty" url:"detail,omitempty"`                    // 商品详情
	Attach         string `xml:"attach,omitempty" url:"attach,omitempty"`                    // 附加数据
	OutTradeNo     string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`         // 商户订单号
	TotalFee       string `xml:"total_fee" url:"total_fee" validate:"nonzero"`               // 总金额
	FeeType        string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`                // 货币类型
	SpbillCreateIP string `xml:"spbill_create_ip" url:"spbill_create_ip" validate:"nonzero"` // 终端IP
	GoodsGag       string `xml:"goods_tag,omitempty" url:"goods_tag,omitempty"`              // 商品标记
	LimitPay       string `xml:"limit_pay,omitempty" url:"limit_pay,omitempty"`              // 指定支付方式
	AuthCode       string `xml:"auth_code" url:"auth_code" validate:"nonzero"`               // 授权码
	// AuthCode       string `xml:"auth_code" url:"auth_code" validate:"regexp=^1\\d{17}$"` // 授权码
}

// GetURI 取接口地址
func (p *PayReq) GetURI() string {
	return "/pay/micropay"
}

// PayResp 被扫支付提交Post数据给到API之后，API会返回XML格式的数据，这个类用来装这些数据
type PayResp struct {
	weixin.CommonBody

	DeviceInfo     string `xml:"device_info,omitempty" url:"device_info,omitempty"`     // 设备号
	Openid         string `xml:"openid" url:"openid,omitempty"`                         // 用户标识
	IsSubscribe    string `xml:"is_subscribe" url:"is_subscribe,omitempty"`             // 是否关注公众账号
	TradeType      string `xml:"trade_type" url:"trade_type,omitempty"`                 // 交易类型
	BankType       string `xml:"bank_type" url:"bank_type,omitempty"`                   // 付款银行
	FeeType        string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`           // 货币类型
	TotalFee       string `xml:"total_fee" url:"total_fee,omitempty"`                   // 总金额
	CashFeeType    string `xml:"cash_fee_type,omitempty" url:"cash_fee_type,omitempty"` // 现金支付货币类型
	CashFee        string `xml:"cash_fee" url:"cash_fee,omitempty"`                     // 现金支付金额
	CouponFee      string `xml:"coupon_fee,omitempty" url:"coupon_fee,omitempty"`       // 代金券或立减优惠金额
	TransactionId  string `xml:"transaction_id" url:"transaction_id,omitempty"`         // 微信支付订单号
	OutTradeNo     string `xml:"out_trade_no" url:"out_trade_no,omitempty"`             // 商户订单号
	Attach         string `xml:"attach,omitempty" url:"attach,omitempty"`               // 商家数据包
	TimeEnd        string `xml:"time_end" url:"time_end,omitempty"`                     // 支付完成时间
	CouponId0      string `xml:"coupon_id_0,omitempty" url:"coupon_id_0,omitempty"`     // 代金券或立减优惠ID
	CouponFee0     string `xml:"coupon_fee_0,omitempty" url:"coupon_fee_0,omitempty"`   // 代金券或立减优惠退款金额
	CouponId1      string `xml:"coupon_id_1,omitempty" url:"coupon_id_1,omitempty"`     // 代金券或立减优惠ID
	CouponFee1     string `xml:"coupon_fee_1,omitempty" url:"coupon_fee_1,omitempty"`   // 代金券或立减优惠退款金额
	CouponId2      string `xml:"coupon_id_2,omitempty" url:"coupon_id_2,omitempty"`     // 代金券或立减优惠ID
	CouponFee2     string `xml:"coupon_fee_2,omitempty" url:"coupon_fee_2,omitempty"`   // 代金券或立减优惠退款金额
	SubOpenid      string `xml:"sub_openid,omitempty" url:"sub_openid,omitempty"`       // 子商户 Open ID
	SubIsSubscribe string `xml:"sub_is_subscribe" url:"sub_is_subscribe,omitempty"`     // 是否关注子商户公众账号
}
