package oversea

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
)

func NewReverseReq() *PayReq {
	return &PayReq{
		Service: "alipay.acquire.overseas.spot.reverse",
	}
}

type ReverseReq struct {
	scanpay1.CommonReq

	Service        string `url:"service"`                    // M
	PartnerTransId string `url:"partner_trans_id,omitempty"` // M
}

type ReverseResp struct {
	scanpay1.CommonResp

	Response struct {
		Alipay struct {
			ResultCode        string `url:"result_code" xml:"result_code"`
			Error             string `url:"error" xml:"error"`
			PartnerTransId    string `url:"partner_trans_id" xml:"partner_trans_id"`
			AlipayTransId     string `url:"alipay_trans_id" xml:"alipay_trans_id"`
			AlipayReverseTime string `url:"alipay_reverse_time" xml:"alipay_reverse_time"`
		} `xml:"alipay,omitempty"`
	} `xml:"response,omitempty" bson:"response,omitempty"`
}

func (p *ReverseResp) ResultCode() string {
	return p.Response.Alipay.ResultCode
}

func (p *ReverseResp) ErrorCode() string {
	if p.CommonResp.ReqFlag() {
		return p.Response.Alipay.Error
	}
	return p.CommonResp.Error
}
