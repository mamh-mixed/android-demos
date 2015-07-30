package scanpay

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_4
// 应用场景
// 当交易发生之后一段时间内，由于买家或者卖家的原因需要退款时，卖家可以通过退款接口将支付款退还给买家，微信支付将在收到退款请求并且验证成功之后，按照退款规则将支付款按原路退到买家帐号上。
// 注意：
// 1.交易时间超过半年的订单无法提交退款；
// 2.微信支付退款支持单笔交易分多次退款，多次退款需要提交原支付订单的商户订单号和设置不同的退款单号。一笔退款失败后重新提交，要采用原来的退款单号。总退款金额不能超过用户实际支付金额。
// 接口链接：https://api.mch.weixin.qq.com/secapi/pay/refund
// 是否需要证书
// 请求需要双向证书。

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/omigo/log"
)

// RefundReq 申请退款
type RefundReq struct {
	CommonParams

	DeviceInfo    string `xml:"device_info,omitempty" url:"device_info,omitempty"`         // 设备号
	TransactionId string `xml:"transaction_id" url:"transaction_id" validate:"nonzero"`    // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`        // 商户订单号
	OutRefundNo   string `xml:"out_refund_no" url:"out_refund_no" validate:"nonzero"`      // 商户退款单号
	TotalFee      string `xml:"total_fee" url:"total_fee" validate:"nonzero"`              // 总金额
	RefundFee     string `xml:"refund_fee" url:"refund_fee" validate:"nonzero"`            // 退款金额
	RefundFeeType string `xml:"refund_fee_type,omitempty" url:"refund_fee_type,omitempty"` // 货币种类
	OpUserId      string `xml:"op_user_id" url:"op_user_id" validate:"nonzero"`            // 操作员
}

// GenSign 计算签名
func (d *RefundReq) GenSign() {
	buf := bytes.Buffer{}

	buf.WriteString("appid=" + d.Appid)
	if d.DeviceInfo != "" {
		buf.WriteString("&device_info=" + d.DeviceInfo)
	}
	buf.WriteString("&mch_id=" + d.MchID)
	buf.WriteString("&nonce_str=" + d.NonceStr)
	buf.WriteString("&op_user_id=" + d.OpUserId)
	buf.WriteString("&out_refund_no=" + d.OutRefundNo)
	buf.WriteString("&out_trade_no=" + d.OutTradeNo)
	buf.WriteString("&refund_fee=" + d.RefundFee)
	if d.RefundFeeType != "" {
		buf.WriteString("&refund_fee_type=" + d.RefundFeeType)
	}
	buf.WriteString("&sub_mch_id=" + d.SubMchId)
	buf.WriteString("&total_fee=" + d.TotalFee)
	if d.TransactionId != "" {
		buf.WriteString("&transaction_id=" + d.TransactionId)
	}
	buf.WriteString("&key=" + d.WeixinMD5Key)

	log.Debug(buf.String())

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))
}

// RefundResp 申请退款
type RefundResp struct {
	XMLName xml.Name `xml:"xml" url:"-"`

	ReturnCode string `xml:"return_code" url:"return_code"`                   // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty" url:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	Appid      string `xml:"appid" url:"appid"`                                   // 公众账号ID
	MchID      string `xml:"mch_id" url:"mch_id"`                                 // 商户号
	SubMchId   string `xml:"sub_mch_id" url:"sub_mch_id"`                         // 子商户号（文档没有该字段）
	SubAppid   string `xml:"sub_appid,omitempty" url:"sub_appid,omitempty"`       // 子商户公众账号 ID
	NonceStr   string `xml:"nonce_str" url:"nonce_str"`                           // 随机字符串
	Sign       string `xml:"sign" url:"-"`                                        // 签名
	ResultCode string `xml:"result_code" url:"result_code"`                       // 业务结果
	ErrCode    string `xml:"err_code,omitempty" url:"err_code,omitempty"`         // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty" url:"err_code_des,omitempty"` // 错误代码描述

	// 以上为微信接口公共字段

	DeviceInfo        string `xml:"device_info,omitempty" url:"device_info,omitempty"`                 // 设备号
	TransactionId     string `xml:"transaction_id" url:"transaction_id"`                               // 微信订单号
	OutTradeNo        string `xml:"out_trade_no" url:"out_trade_no"`                                   // 商户订单号
	OutRefundNo       string `xml:"out_refund_no" url:"out_refund_no"`                                 // 商户退款单号
	RefundId          string `xml:"refund_id" url:"refund_id"`                                         // 微信退款单号
	RefundChannel     string `xml:"refund_channel,omitempty" url:"refund_channel,omitempty"`           // 退款渠道
	RefundFee         string `xml:"refund_fee" url:"refund_fee"`                                       // 退款金额
	TotalFee          string `xml:"total_fee" url:"total_fee"`                                         // 订单总金额
	FeeType           string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`                       // 订单金额货币种类
	CashFee           int    `xml:"cash_fee" url:"cash_fee"`                                           // 现金支付金额
	CashRefundFee     int    `xml:"cash_refund_fee,omitempty" url:"cash_refund_fee,omitempty"`         // 现金退款金额
	CouponRefundFee   int    `xml:"coupon_refund_fee,omitempty" url:"coupon_refund_fee,omitempty"`     // 代金券或立减优惠退款金额
	CouponRefundCount int    `xml:"coupon_refund_count,omitempty" url:"coupon_refund_count,omitempty"` // 代金券或立减优惠使用数量
	CouponRefundId    string `xml:"coupon_refund_id,omitempty" url:"coupon_refund_id,omitempty"`       // 代金券或立减优惠ID
}
