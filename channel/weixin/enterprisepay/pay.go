package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
)

// PayReq 请求被扫支付API需要提交的数据
type EnterprisePayReq struct {
	weixin.CommonParams

	MchAappid      string `xml:"mch_appid" url:"mch_appid"` // 商户appid
	PartnerTradeNo string `xml:"partner_trade_no" url:"partner_trade_no"`
	OpenId         string `xml:"openid" url:"openid"`
	CheckName      string `xml:"check_name" url:"check_name"`
	ReUserName     string `xml:"re_user_name" url:"re_user_name,omitempty"`
	Amount         string `xml:"amount" url:"amount"`
	Desc           string `xml:"desc" url:"desc"`
	SpbillCreateIp string `xml:"spbill_create_ip" url:"spbill_create_ip,omitempty"`
}

// GetURI 取接口地址
func (p *EnterprisePayReq) GetURI() string {
	return "/mmpaymkttransfers/promotion/transfers"
}

// PayResp 被扫支付提交Post数据给到API之后，API会返回XML格式的数据，这个类用来装这些数据
type EnterprisePayResp struct {
	weixin.CommonBody

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	MchAappid  string `xml:"mch_appid"`             // 公众账号ID
	DeviceInfo string `xml:"device_info,omitempty"` // 设备号
	// 以上为微信接口公共字段

	// 当 return_code 和 result_code 都为 SUCCESS 的时，还会包括以下字段：
	PartnerTradeNo string `xml:"partner_trade_no,omitempty"` // 商户订单号，需保持唯一性
	PaymentNo      string `xml:"payment_no,omitempty"`       // 企业付款成功，返回的微信订单号
	PaymentTime    string `xml:"payment_time,omitempty"`     // 企业付款成功时间
}
