package scanpay

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_5
// 应用场景
// 提交退款申请后，通过调用该接口查询退款状态。退款有一定延时，用零钱支付的退款20分钟内到账，银行卡支付的退款3个工作日后重新查询退款状态。
// 接口链接：https://api.mch.weixin.qq.com/pay/refundquery
// 是否需要证书: 不需要

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/omigo/log"
)

// RefundQueryReq 查询退款
type RefundQueryReq struct {
	CommonParams

	DeviceInfo    string `xml:"device_info,omitempty" url:"device_info,omitempty"`       // 设备号
	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`      // 商户订单号
	OutRefundNo   string `xml:"out_refund_no,omitempty" url:"out_refund_no,omitempty"`   // 商户退款单号
	RefundId      string `xml:"refund_id,omitempty" url:"refund_id,omitempty"`           // 微信退款单号
}

// GenSign 计算签名
func (d *RefundQueryReq) GenSign() {
	buf := bytes.Buffer{}

	buf.WriteString("appid=" + d.Appid)
	if d.DeviceInfo != "" {
		buf.WriteString("&device_info=" + d.DeviceInfo)
	}
	buf.WriteString("&mch_id=" + d.MchID)
	buf.WriteString("&nonce_str=" + d.NonceStr)
	if d.OutRefundNo != "" {
		buf.WriteString("&out_refund_no=" + d.OutRefundNo)
	}
	buf.WriteString("&out_trade_no=" + d.OutTradeNo)
	if d.RefundId != "" {
		buf.WriteString("&refund_id=" + d.RefundId)
	}
	buf.WriteString("&sub_mch_id=" + d.SubMchId)
	if d.TransactionId != "" {
		buf.WriteString("&transaction_id=" + d.TransactionId)
	}

	buf.WriteString("&key=" + d.WeixinMD5Key)

	log.Debug(buf.String())

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))
}

// RefundQueryResp 查询退款
type RefundQueryResp struct {
	XMLName xml.Name `xml:"xml" url:"-"`

	ReturnCode string `xml:"return_code" url:"return_code"`                   // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty" url:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	Appid      string `xml:"appid,omitempty" url:"appid,omitempty"`               // 公众账号ID
	MchID      string `xml:"mch_id" url:"mch_id"`                                 // 商户号
	SubMchId   string `xml:"sub_mch_id" url:"sub_mch_id"`                         // 子商户号（文档没有该字段）
	SubAppid   string `xml:"sub_appid,omitempty" url:"sub_appid,omitempty"`       // 子商户公众账号 ID
	NonceStr   string `xml:"nonce_str" url:"nonce_str"`                           // 随机字符串
	Sign       string `xml:"sign" url:"-"`                                        // 签名
	ResultCode string `xml:"result_code" url:"result_code"`                       // 业务结果
	ErrCode    string `xml:"err_code,omitempty" url:"err_code,omitempty"`         // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty" url:"err_code_des,omitempty"` // 错误代码描述

	// 以上为微信接口公共字段

	DeviceInfo string `xml:"device_info,omitempty" url:"device_info,omitempty"` // 设备号

	TransactionId   string `xml:"transaction_id" url:"transaction_id"`                           // 微信订单号
	OutTradeNo      string `xml:"out_trade_no" url:"out_trade_no"`                               // 商户订单号
	TotalFee        int    `xml:"total_fee" url:"total_fee"`                                     // 订单总金额
	FeeType         string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`                   // 订单金额货币种类
	CashFee         int    `xml:"cash_fee" url:"cash_fee"`                                       // 现金支付金额
	CashFeeType     string `xml:"cash_fee_type,omitempty" url:"cash_fee_type,omitempty"`         // 货币种类
	RefundFee       int    `xml:"refund_fee" url:"refund_fee"`                                   // 退款金额
	CouponRefundFee int    `xml:"coupon_refund_fee,omitempty" url:"coupon_refund_fee,omitempty"` // 代金券或立减优惠退款金额
	RefundCount     int    `xml:"refund_count" url:"refund_count"`                               // 退款笔数
}
