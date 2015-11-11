package oversea

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
)

func NewQueryReq() *QueryReq {
	return &QueryReq{
		Service: "alipay.acquire.overseas.query",
	}
}

type QueryReq struct {
	scanpay1.CommonReq

	Service        string `url:"service"`
	PartnerTransId string `url:"partner_trans_id,omitempty"` // M
	AlipayTransId  string `url:"alipay_trans_id,omitempty" xml:"alipay_trans_id"`
}

type QueryResp struct {
	scanpay1.CommonResp

	Response struct {
		Alipay struct {
			ResultCode         string `url:"result_code" xml:"result_code"`
			Error              string `url:"error" xml:"error"`
			AlipayTransStatus  string `url:"alipay_trans_status" xml:"alipay_trans_status"`
			AlipayBuyerLoginId string `url:"alipay_buyer_login_id" xml:"alipay_buyer_login_id"`
			AlipayBuyerUserId  string `url:"alipay_buyer_user_id" xml:"alipay_buyer_user_id"`
			PartnerTransId     string `url:"partner_trans_id" xml:"partner_trans_id"`
			AlipayTransId      string `url:"alipay_trans_id" xml:"alipay_trans_id"`
			AlipayPayTime      string `url:"alipay_pay_time" xml:"alipay_pay_time"`
			Currency           string `url:"currency" xml:"currency"`
			TransAmount        string `url:"trans_amount" xml:"trans_amount"`
			ExchangeRate       string `url:"exchange_rate" xml:"exchange_rate"`
			TransAmountCny     string `url:"trans_amount_cny" xml:"trans_amount_cny"`
		} `xml:"alipay,omitempty"`
	} `xml:"response,omitempty" bson:"response,omitempty"`
}

func (p *QueryResp) ResultCode() string {
	return p.Response.Alipay.ResultCode
}

func (p *QueryResp) ErrorCode() string {
	if p.CommonResp.ReqFlag() {
		return p.Response.Alipay.Error
	}
	return p.CommonResp.Error
}
