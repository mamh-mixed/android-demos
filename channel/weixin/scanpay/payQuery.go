package scanpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/omigo/log"
)

// PayQueryReq 请求被扫支付API需要提交的数据
type PayQueryReq struct {
	XMLName xml.Name `xml:"xml"`

	// 公共字段
	Appid    string `xml:"appid" validate:"len=18"`       // 公众账号ID
	MchID    string `xml:"mch_id" validate:"nonzero"`     // 商户号
	SubMchId string `xml:"sub_mch_id" validate:"nonzero"` // 子商户号（文档没有该字段）
	NonceStr string `xml:"nonce_str" validate:"nonzero"`  // 随机字符串
	Sign     string `xml:"sign"`                          // 签名

	WeixinMD5Key string `xml:"-" validate:"nonzero"`

	TransactionId string `xml:"transaction_id,omitempty"` // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no,omitempty"`   // 商户订单号
}

// GenSign 计算签名 （写一个 marshal 方法，类似 json 和 xml ，作为工具类，一次搞定 拼串）
func (d *PayQueryReq) GenSign() {
	buf := bytes.Buffer{}

	buf.WriteString("appid=" + d.Appid)
	buf.WriteString("&mch_id=" + d.MchID)
	buf.WriteString("&nonce_str=" + d.NonceStr)
	if d.OutTradeNo != "" {
		buf.WriteString("&out_trade_no=" + d.OutTradeNo)
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

// PayQueryResp 被扫支付提交Post数据给到API之后，API会返回XML格式的数据，这个类用来装这些数据
type PayQueryResp struct {
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
	DeviceInfo     string `xml:"device_info,omitempty"` // 设备号
	OpenID         string `xml:"openid"`                // 用户标识
	IsSubscribe    string `xml:"is_subscribe"`          // 是否关注公众账号
	TradeType      string `xml:"trade_type"`            // 交易类型
	TradeState     string `xml:"tradeState"`            // 交易状态
	BankType       string `xml:"bank_type"`             // 付款银行
	FeeType        string `xml:"fee_type"`              // 货币类型
	TotalFee       string `xml:"total_fee"`             // 总金额
	CashFeeType    string `xml:"cash_fee_type"`         // 现金支付货币类型
	CashFee        string `xml:"cash_fee"`              // 现金支付金额
	CouponFee      string `xml:"coupon_fee"`            // 代金券或立减优惠金额
	CouponCount    string `xml:"coupon_count"`          // 代金券或立减优惠使用数量
	TransactionId  string `xml:"transaction_id"`        // 微信支付订单号
	OutTradeNo     string `xml:"out_trade_no"`          // 商户订单号
	Attach         string `xml:"attach"`                // 商家数据包
	TimeEnd        string `xml:"time_end"`              // 支付完成时间
	TradeStateDesc string `xml:"trade_state_desc"`      // 交易状态描述
}
