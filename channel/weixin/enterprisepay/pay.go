package enterprisepay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// PayReq 请求被扫支付API需要提交的数据
type EnterprisePayReq struct {
	XMLName xml.Name `xml:"xml" url:"-"`

	// 公共字段
	MchAappid string `xml:"mch_appid" url:"mch_appid"`            // 商户appid
	MchID     string `xml:"mchid" url:"mchid" validate:"nonzero"` // 商户号
	// SubMchId  string `xml:"sub_mch_id" validate:"nonzero"` // 子商户号（文档没有该字段）
	NonceStr       string `xml:"nonce_str" url:"nonce_str" validate:"nonzero"` // 随机字符串
	Sign           string `xml:"sign" url:"sign,omitempty"`                    // 签名
	PartnerTradeNo string `xml:"partner_trade_no" url:"partner_trade_no"`
	OpenId         string `xml:"openid" url:"openid"`
	CheckName      string `xml:"check_name" url:"check_name"`
	ReUserName     string `xml:"re_user_name" url:"re_user_name,omitempty"`
	Amount         string `xml:"amount" url:"amount"`
	Desc           string `xml:"desc" url:"desc"`
	SpbillCreateIp string `xml:"spbill_create_ip" url:"spbill_create_ip,omitempty"`
	WeixinMD5Key   string `xml:"-" validate:"nonzero" url:"-" validate:"nonzero"`
}

// GenSign 计算签名 （写一个 marshal 方法，类似 json 和 xml ，作为工具类，一次搞定 拼串）
func (d *EnterprisePayReq) GenSign() {
	buf, err := util.Query(d)
	if err != nil {
		log.Error(err)
	}
	buf.WriteString("&key=" + d.WeixinMD5Key)
	log.Debug(buf.String())

	sign := md5.Sum(buf.Bytes())
	d.Sign = strings.ToUpper(hex.EncodeToString(sign[:]))
}

// PayResp 被扫支付提交Post数据给到API之后，API会返回XML格式的数据，这个类用来装这些数据
type EnterprisePayResp struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	MchAappid  string `xml:"mch_appid"`              // 公众账号ID
	MchID      string `xml:"mchid"`                  // 商户号
	NonceStr   string `xml:"nonce_str"`              // 随机字符串
	Sign       string `xml:"sign"`                   // 签名
	ResultCode string `xml:"result_code"`            // 业务结果
	ErrCode    string `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"` // 错误代码描述
	DeviceInfo string `xml:"device_info,omitempty"`  // 设备号
	// 以上为微信接口公共字段

	// 当 return_code 和 result_code 都为 SUCCESS 的时，还会包括以下字段：
	PartnerTradeNo string `xml:"partner_trade_no,omitempty"` // 商户订单号，需保持唯一性
	PaymentNo      string `xml:"payment_no,omitempty"`       // 企业付款成功，返回的微信订单号
	PaymentTime    string `xml:"payment_time,omitempty"`     // 企业付款成功时间
}
