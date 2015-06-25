package scanpay

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
	XMLName xml.Name `xml:"xml"`

	// 公共字段
	Appid    string `xml:"appid" validate:"len=18"`       // 公众账号ID
	MchID    string `xml:"mch_id" validate:"nonzero"`     // 商户号
	SubMchId string `xml:"sub_mch_id" validate:"nonzero"` // 子商户号（文档没有该字段）
	NonceStr string `xml:"nonce_str" validate:"nonzero"`  // 随机字符串
	Sign     string `xml:"sign"`                          // 签名

	WeixinMD5Key string `xml:"-" validate:"nonzero"`

	DeviceInfo    string `xml:"device_info,omiempty"`             // 设备号
	TransactionId string `xml:"transaction_id,omiempty"`          // 微信订单号
	OutTradeNo    string `xml:"out_trade_no,omiempty"`            // 商户订单号
	OutRefundNo   string `xml:"out_refund_no" validate:"nonzero"` // 商户退款单号
	TotalFee      string `xml:"total_fee" validate:"nonzero"`     // 总金额
	RefundFee     string `xml:"refund_fee" validate:"nonzero"`    // 退款金额
	RefundFeeType string `xml:"refund_fee_type,omiempty"`         // 货币种类
	OpUserId      string `xml:"op_user_id" validate:"nonzero"`    // 操作员
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
	buf.WriteString("&transaction_id=" + d.TransactionId)

	buf.WriteString("&key=" + d.WeixinMD5Key)

	log.Debug(buf.String())

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))
}

// RefundResp 申请退款
type RefundResp struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	Appid      string `xml:"appid"`                  // 公众账号ID
	MchID      string `xml:"mch_id"`                 // 商户号
	SubMchId   string `xml:"sub_mch_id"`             // 子商户号（文档没有该字段）
	NonceStr   string `xml:"nonce_str"`              // 随机字符串
	Sign       string `xml:"sign"`                   // 签名
	ResultCode string `xml:"result_code"`            // 业务结果
	ErrCode    string `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"` // 错误代码描述

	// 以上为微信接口公共字段

	DeviceInfo string `xml:"device_info,omitempty"` // 设备号

	TransactionId     string `xml:"transaction_id"`               // 微信订单号
	OutTradeNo        string `xml:"out_trade_no"`                 // 商户订单号
	OutRefundNo       string `xml:"out_refund_no"`                // 商户退款单号
	RefundId          string `xml:"refund_id"`                    // 微信退款单号
	RefundChannel     string `xml:"refund_channel,omiempty"`      // 退款渠道
	RefundFee         string `xml:"refund_fee"`                   // 退款金额
	TotalFee          string `xml:"total_fee"`                    // 订单总金额
	FeeType           string `xml:"fee_type,omiempty"`            // 订单金额货币种类
	CashFee           int    `xml:"cash_fee"`                     // 现金支付金额
	CashRefundFee     int    `xml:"cash_refund_fee,omiempty"`     // 现金退款金额
	CouponRefundFee   int    `xml:"coupon_refund_fee,omiempty"`   // 代金券或立减优惠退款金额
	CouponRefundCount int    `xml:"coupon_refund_count,omiempty"` // 代金券或立减优惠使用数量
	CouponRefundId    string `xml:"coupon_refund_id,omiempty"`    // 代金券或立减优惠ID
}
