package scanpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/omigo/log"
)

// ReverseReq 撤销订单
type ReverseReq struct {
	CommonParams

	TransactionId string `xml:"transaction_id,omiempty"`         // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" validate:"nonzero"` // 商户订单号
}

// GenSign 计算签名
func (d *ReverseReq) GenSign() {
	buf := bytes.Buffer{}

	buf.WriteString("appid=" + d.Appid)
	buf.WriteString("&mch_id=" + d.MchID)
	buf.WriteString("&nonce_str=" + d.NonceStr)
	buf.WriteString("&out_trade_no=" + d.OutTradeNo)
	buf.WriteString("&sub_mch_id=" + d.SubMchId)
	if d.TransactionId != "" {
		buf.WriteString("&transaction_id=" + d.TransactionId)
	}

	buf.WriteString("&key=" + d.WeixinMD5Key)

	log.Debug(buf.String())

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))
}

// ReverseResp 撤销订单
type ReverseResp struct {
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

	Recall string `xml:"recall"` // 是否重调
}
