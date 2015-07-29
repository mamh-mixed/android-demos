package scanpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/omigo/log"
)

// CloseReq 关闭订单
type CloseReq struct {
	CommonParams

	TransactionId string `xml:"transaction_id,omiempty"`         // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" validate:"nonzero"` // 商户订单号
}

// GenSign 计算签名
func (d *CloseReq) GenSign() {
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

// CloseResp 撤销订单
type CloseResp struct {
	CommonBody
}
