package scanpay

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
// 应用场景
// 除被扫支付场景以外，商户系统先调用该接口在微信支付服务后台生成预支付交易单，返回正确的预支付交易回话标识后再按扫码、JSAPI、APP等不同场景生成交易串调起支付。
// 接口链接
// URL地址：https://api.mch.weixin.qq.com/pay/unifiedorder
// 是否需要证书: 不需要

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

var weixinNotifyURL = goconf.Config.WeixinScanPay.NotifyURL + "/quickpay/back/weixin"

// PrePayReq 请求被扫支付API需要提交的数据
type PrePayReq struct {
	CommonParams

	DeviceInfo     string `xml:"device_info,omitempty" url:"device_info,omitempty"`          // 设备号
	Body           string `xml:"body" url:"body" validate:"nonzero"`                         // 商品描述
	Detail         string `xml:"detail,omitempty" url:"detail,omitempty"`                    // 商品详情
	Attach         string `xml:"attach,omitempty" url:"attach,omitempty"`                    // 附加数据
	OutTradeNo     string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`         // 商户订单号
	TotalFee       string `xml:"total_fee" url:"total_fee" validate:"nonzero"`               // 总金额
	FeeType        string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`                // 货币类型
	SpbillCreateIP string `xml:"spbill_create_ip" url:"spbill_create_ip" validate:"nonzero"` // 终端IP
	TimeStart      string `xml:"time_start,omitempty" url:"time_start,omitempty"`            // 交易起始时间
	TimeExpire     string `xml:"time_expire,omitempty" url:"time_expire,omitempty"`          // 交易结束时间
	GoodsGag       string `xml:"goods_tag,omitempty" url:"goods_tag,omitempty"`              // 商品标记
	NotifyURL      string `xml:"notify_url" url:"notify_url" validate:"nonzero"`             // 通知地址
	TradeType      string `xml:"trade_type" url:"trade_type" validate:"nonzero"`             // 交易类型
	ProductID      string `xml:"product_id,omitempty" url:"product_id,omitempty"`            // 商品ID
	Openid         string `xml:"openid,omitempty" url:"openid,omitempty"`                    // 用户标识
	SubOpenid      string `xml:"sub_openid,omitempty" url:"sub_openid,omitempty"`            // 子商户用户标识
}

// GenSign 计算签名 （写一个 marshal 方法，类似 json 和 xml ，作为工具类，一次搞定 拼串）
func (d *PrePayReq) GenSign() {
	buf := bytes.Buffer{}

	buf.WriteString("appid=" + d.Appid)
	if d.Attach != "" {
		buf.WriteString("&attach=" + d.Attach)
	}
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
	buf.WriteString("&notify_url=" + d.NotifyURL)
	if d.Openid != "" {
		buf.WriteString("&openid=" + d.Openid)
	}
	buf.WriteString("&out_trade_no=" + d.OutTradeNo)
	if d.ProductID != "" {
		buf.WriteString("&product_id=" + d.ProductID)
	}
	buf.WriteString("&spbill_create_ip=" + d.SpbillCreateIP)
	buf.WriteString("&sub_mch_id=" + d.SubMchId)
	if d.SubOpenid != "" {
		buf.WriteString("&sub_openid=" + d.SubOpenid)
	}
	if d.TimeStart != "" {
		buf.WriteString("&time_start=" + d.TimeStart)
	}
	if d.TimeExpire != "" {
		buf.WriteString("&time_expire=" + d.TimeExpire)
	}
	buf.WriteString("&total_fee=" + d.TotalFee)
	buf.WriteString("&trade_type=" + d.TradeType)
	buf.WriteString("&key=" + d.WeixinMD5Key)

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))

	log.Debugf("%s = md5( %s )", d.Sign, buf.String())
}

// PrePayResp 被扫支付提交Post数据给到API之后，API会返回XML格式的数据，这个类用来装这些数据
type PrePayResp struct {
	CommonBody

	DeviceInfo string `xml:"device_info,omitempty" url:"device_info,omitempty"` // 设备号

	// 当 return_code 和 result_code 都为 SUCCESS 的时，还会包括以下字段：
	TradeType string `xml:"trade_type" url:"trade_type"`                 // 交易类型
	PrepayID  string `xml:"prepay_id" url:"prepay_id"`                   // 预支付交易会话标识
	CodeURL   string `xml:"code_url,omitempty" url:"code_url,omitempty"` // 二维码链接
}
