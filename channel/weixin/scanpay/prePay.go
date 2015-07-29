package scanpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

var weixinNotifyURL = goconf.Config.WeixinScanPay.NotifyURL + "/quickpay/back/weixin"

// PrePayReq 请求被扫支付API需要提交的数据
type PrePayReq struct {
	CommonParams

	DeviceInfo     string `xml:"device_info,omitempty"`               // 设备号
	Body           string `xml:"body" validate:"nonzero"`             // 商品描述
	Detail         string `xml:"detail,omitempty"`                    // 商品详情
	Attach         string `xml:"attach,omitempty"`                    // 附加数据
	OutTradeNo     string `xml:"out_trade_no" validate:"nonzero"`     // 商户订单号
	TotalFee       string `xml:"total_fee" validate:"nonzero"`        // 总金额
	FeeType        string `xml:"fee_type,omitempty"`                  // 货币类型
	SpbillCreateIP string `xml:"spbill_create_ip" validate:"nonzero"` // 终端IP
	TimeStart      string `xml:"time_start,omitempty"`                // 交易起始时间
	TimeExpire     string `xml:"time_expire,omitempty"`               // 交易结束时间
	GoodsGag       string `xml:"goods_tag,omitempty"`                 // 商品标记
	NotifyURL      string `xml:"notify_url" validate:"nonzero"`       // 通知地址
	TradeType      string `xml:"trade_type" validate:"nonzero"`       // 交易类型
	ProductID      string `xml:"product_id,omitempty"`                // 商品ID
	Openid         string `xml:"openid,omitempty"`                    // 用户标识
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
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	Appid      string `xml:"appid"`                  // 公众账号ID
	MchID      string `xml:"mch_id"`                 // 商户号
	SubMchId   string `xml:"sub_mch_id"`             // 子商户号（文档没有该字段）
	SubAppid   string `xml:"sub_appid"`              // 子商户公众账号 ID
	NonceStr   string `xml:"nonce_str"`              // 随机字符串
	Sign       string `xml:"sign"`                   // 签名
	ResultCode string `xml:"result_code"`            // 业务结果
	ErrCode    string `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"` // 错误代码描述

	// 以上为微信接口公共字段

	DeviceInfo string `xml:"device_info,omitempty"` // 设备号

	// 当 return_code 和 result_code 都为 SUCCESS 的时，还会包括以下字段：
	TradeType string `xml:"trade_type"` // 交易类型
	PrepayID  string `xml:"prepay_id"`  // 预支付交易会话标识
	CodeURL   string `xml:"code_url"`   // 二维码链接
}
