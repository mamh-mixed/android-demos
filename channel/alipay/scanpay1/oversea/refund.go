package oversea

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
)

func NewRefundReq() *RefundReq {
	return &RefundReq{
		Service: "alipay.acquire.overseas.spot.refund",
	}
}

type RefundReq struct {
	scanpay1.CommonReq

	Service         string `url:"service"`
	PartnerTransId  string `url:"partner_trans_id,omitempty"`                // M
	PartnerRefundId string `url:"partner_refund_id" xml:"partner_refund_id"` // M
	Currency        string `url:"currency,omitempty"`                        // M
	RefundAmount    string `url:"refund_amount,omitempty"`                   // M
	RefundReson     string `url:"refund_reson,omitempty"`
}

type RefundResp struct {
	scanpay1.CommonResp

	Response struct {
		Alipay struct {
			ResultCode      string `url:"result_code" xml:"result_code"`
			Error           string `url:"error" xml:"error"`
			PartnerTransId  string `url:"partner_trans_id" xml:"partner_trans_id"`
			PartnerRefundId string `url:"partner_refund_id" xml:"partner_refund_id"` // M
			AlipayTransId   string `url:"alipay_trans_id" xml:"alipay_trans_id"`
			Currency        string `url:"currency" xml:"currency"`
			RefundAmount    string `url:"refund_amount" xml:"refund_amount"`
			ExchangeRate    string `url:"exchange_rate" xml:"exchange_rate"`
			RefundAmountCny string `url:"refund_amount_cny" xml:"refund_amount_cny"`
		} `xml:"alipay,omitempty"`
	} `xml:"response,omitempty" bson:"response,omitempty"`
}

func (p *RefundResp) ResultCode() string {
	return p.Response.Alipay.ResultCode
}

func (p *RefundResp) ErrorCode() string {
	if p.CommonResp.ReqFlag() {
		return p.Response.Alipay.Error
	}
	return p.CommonResp.Error
}
