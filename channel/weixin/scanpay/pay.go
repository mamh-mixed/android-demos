package scanpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/omigo/log"
)

// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
// 统一下单
// 应用场景
// 除被扫支付场景以外，商户系统先调用该接口在微信支付服务后台生成预支付交易单，返回正确的预支付交易回话标识后再按扫码、JSAPI、APP等不同场景生成交易串调起支付。
// 接口链接
// URL地址：https://api.mch.weixin.qq.com/pay/unifiedorder
// 是否需要证书
// 不需要

// PayReq 请求被扫支付API需要提交的数据
type PayReq struct {
	XMLName xml.Name `xml:"xml"`

	// 公共字段
	Appid    string `xml:"appid" validate:"len=18"`       // 公众账号ID
	MchID    string `xml:"mch_id" validate:"nonzero"`     // 商户号
	SubMchId string `xml:"sub_mch_id" validate:"nonzero"` // 子商户号（文档没有该字段）
	NonceStr string `xml:"nonce_str" validate:"nonzero"`  // 随机字符串
	Sign     string `xml:"sign"`                          // 签名

	WeixinMD5Key string `xml:"-" validate:"nonzero"`

	DeviceInfo     string `xml:"device_info,omitempty"`               // 设备号
	Body           string `xml:"body" validate:"nonzero"`             // 商品描述
	Detail         string `xml:"detail,omitempty"`                    // 商品详情
	Attach         string `xml:"attach,omitempty"`                    // 附加数据
	OutTradeNo     string `xml:"out_trade_no" validate:"nonzero"`     // 商户订单号
	TotalFee       string `xml:"total_fee" validate:"nonzero"`        // 总金额
	FeeType        string `xml:"fee_type,omitempty"`                  // 货币类型
	SpbillCreateIP string `xml:"spbill_create_ip" validate:"nonzero"` // 终端IP
	GoodsGag       string `xml:"goods_tag,omitempty"`                 // 商品标记
	AuthCode       string `xml:"auth_code"`                           // 授权码
	// AuthCode       string `xml:"auth_code" validate:"regexp=^1\\d{17}$"` // 授权码
}

// GenSign 计算签名 （写一个 marshal 方法，类似 json 和 xml ，作为工具类，一次搞定 拼串）
func (d *PayReq) GenSign() {
	buf := bytes.Buffer{}

	buf.WriteString("appid=" + d.Appid)
	if d.Attach != "" {
		buf.WriteString("&attach=" + d.Attach)
	}
	buf.WriteString("&auth_code=" + d.AuthCode)
	buf.WriteString("&body=" + d.Body)
	if d.Detail != "" {
		buf.WriteString("&detail=" + d.Detail)
	}
	if d.DeviceInfo != "" {
		buf.WriteString("&device_info=" + d.DeviceInfo)
	}
	if d.FeeType != "" {
		buf.WriteString("&fee_type=" + d.FeeType)
	}
	if d.GoodsGag != "" {
		buf.WriteString("&goods_tag=" + d.GoodsGag)
	}
	buf.WriteString("&mch_id=" + d.MchID)
	buf.WriteString("&nonce_str=" + d.NonceStr)
	buf.WriteString("&out_trade_no=" + d.OutTradeNo)
	buf.WriteString("&spbill_create_ip=" + d.SpbillCreateIP)
	buf.WriteString("&sub_mch_id=" + d.SubMchId)
	buf.WriteString("&total_fee=" + d.TotalFee)
	buf.WriteString("&key=" + d.WeixinMD5Key)

	log.Debug(buf.String())

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))
}

// PayResp 被扫支付提交Post数据给到API之后，API会返回XML格式的数据，这个类用来装这些数据
type PayResp struct {
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

	// 当 return_code 和 result_code 都为 SUCCESS 的时，还会包括以下字段：
	OpenID        string `xml:"openid"`         // 用户标识
	IsSubscribe   string `xml:"is_subscribe"`   // 是否关注公众账号
	TradeType     string `xml:"trade_type"`     // 交易类型
	BankType      string `xml:"bank_type"`      // 付款银行
	FeeType       string `xml:"fee_type"`       // 货币类型
	TotalFee      string `xml:"total_fee"`      // 总金额
	CashFeeType   string `xml:"cash_fee_type"`  // 现金支付货币类型
	CashFee       string `xml:"cash_fee"`       // 现金支付金额
	CouponFee     string `xml:"coupon_fee"`     // 代金券或立减优惠金额
	TransactionId string `xml:"transaction_id"` // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	Attach        string `xml:"attach"`         // 商家数据包
	TimeEnd       string `xml:"time_end"`       // 支付完成时间
}
