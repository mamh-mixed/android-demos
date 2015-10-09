package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"net/http"
)

type EnterpriseQueryReq struct {
	weixin.CommonParams
	AppId          string `xml:"appid,omitempty" url:"appid,omitempty"`
	MchId          string `xml:"mch_id,omitempty" url:"mch_id,omitempty"`
	PartnerTradeNo string `xml:"partner_trade_no,omitempty" url:"partner_trade_no,omitempty"`
}

// GetURI 取接口地址
func (p *EnterpriseQueryReq) GetURI() string {
	return "/mmpaymkttransfers/gettransferinfo"
}

// GetHTTPClient 使用双向 HTTPS 认证
func (p *EnterpriseQueryReq) GetHTTPClient() *http.Client {
	return p.GetHTTPSClient()
}

type EnterpriseQueryResp struct {
	weixin.CommonBody
	PartnerTradeNo  string `xml:"partner_trade_no,omitempty" url:"partner_trade_no,omitempty"`
	DetailId        string `xml:"detail_id,omitempty" url:"detail_id,omitempty"`
	Status          string `xml:"status,omitempty" url:"status,omitempty"`
	Reason          string `xml:"reason,omitempty" url:"reason,omitempty"`
	OpenID          string `xml:"openid,omitempty" url:"openid,omitempty"`
	TransferName    string `xml:"transfer_name,omitempty" url:"transfer_name,omitempty"`
	PaymentAmount   string `xml:"payment_amount,omitempty" url:"payment_amount,omitempty"`
	TransferTime    string `xml:"transfer_time,omitempty" url:"transfer_time,omitempty"`
	Desc            string `xml:"desc,omitempty" url:"desc,omitempty"`
	CheckName       string `xml:"check_name,omitempty" url:"check_name,omitempty"`
	CheckNameResult string `xml:"check_name_result,omitempty" url:"check_name_result,omitempty"`
}
