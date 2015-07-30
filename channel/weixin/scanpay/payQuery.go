package scanpay

// 参考文档https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_2
// 应用场景
// 该接口提供所有微信支付订单的查询，商户可以通过该接口主动查询订单状态，完成下一步的业务逻辑。
// 需要调用查询接口的情况：
// ◆ 当商户后台、网络、服务器等出现异常，商户系统最终未接收到支付通知；
// ◆ 调用支付接口后，返回系统错误或未知交易状态情况；
// ◆ 调用被扫支付API，返回USERPAYING的状态；
// ◆ 调用关单或撤销接口API之前，需确认支付状态；
// 接口链接 https://api.mch.weixin.qq.com/pay/orderquery
// 是否需要证书 :不需要

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/omigo/log"
)

// PayQueryReq 请求被扫支付API需要提交的数据
type PayQueryReq struct {
	CommonParams

	TransactionId string `xml:"transaction_id" url:"transaction_id" validate:"nonzero"` // 微信的订单号，优先使用
	OutTradeNo    string `xml:"out_trade_no,omitempty" url:"out_trade_no,omitempty"`    // 商户系统内部的订单号，当没提供transaction_id时需要传这个
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
	CommonBody

	// 当return_code 和result_code都为SUCCESS的时，还会包括以下字段：
	DeviceInfo     string `xml:"device_info,omitempty" url:"device_info,omitempty"`     // 设备号
	OpenID         string `xml:"openid" url:"openid"`                                   // 用户标识
	IsSubscribe    string `xml:"is_subscribe" url:"is_subscribe"`                       // 是否关注公众账号
	TradeType      string `xml:"trade_type" url:"trade_type"`                           // 交易类型
	TradeState     string `xml:"trade_state" url:"trade_state"`                         // 交易状态
	BankType       string `xml:"bank_type" url:"bank_type"`                             // 付款银行
	FeeType        string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`           // 货币类型
	TotalFee       string `xml:"total_fee" url:"total_fee"`                             // 总金额
	CashFeeType    string `xml:"cash_fee_type,omitempty" url:"cash_fee_type,omitempty"` // 现金支付货币类型
	CashFee        string `xml:"cash_fee" url:"cash_fee"`                               // 现金支付金额
	CouponFee      string `xml:"coupon_fee,omitempty" url:"coupon_fee,omitempty"`       // 代金券或立减优惠金额
	CouponCount    string `xml:"coupon_count,omitempty" url:"coupon_count,omitempty"`   // 代金券或立减优惠使用数量
	TransactionId  string `xml:"transaction_id" url:"transaction_id"`                   // 微信支付订单号
	OutTradeNo     string `xml:"out_trade_no" url:"out_trade_no"`                       // 商户订单号
	Attach         string `xml:"attach,omitempty" url:"attach,omitempty"`               // 商家数据包
	TimeEnd        string `xml:"time_end" url:"time_end"`                               // 支付完成时间
	TradeStateDesc string `xml:"trade_state_desc" url:"trade_state_desc"`               // 交易状态描述
	SubOpenid      string `xml:"sub_openid,omitempty" url:"sub_openid,omitempty"`       // 子商户 Open ID
	SubIsSubscribe string `xml:"sub_is_subscribe" url:"sub_is_subscribe"`               // 是否关注子商户公众账号
}
