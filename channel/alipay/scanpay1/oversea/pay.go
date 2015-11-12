package oversea

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
)

func NewPayReq() *PayReq {
	return &PayReq{
		Service:          "alipay.acquire.overseas.spot.pay",
		BizProduct:       "OVERSEAS_MBARCODE_PAY",
		IdentityCodeType: "barcode",
	}
}

type PayReq struct {
	scanpay1.CommonReq

	Service           string `url:"service"`
	AlipaySellerId    string `url:"alipay_seller_id,omitempty"` // M
	Quantity          int    `url:"quantity,omitempty"`
	TransName         string `url:"trans_name,omitempty"`          // M
	PartnerTransId    string `url:"partner_trans_id,omitempty"`    // M
	Currency          string `url:"currency,omitempty"`            // M
	TransAmount       string `url:"trans_amount,omitempty"`        // M
	BuyerIdentityCode string `url:"buyer_identity_code,omitempty"` // M
	IdentityCodeType  string `url:"identity_code_type,omitempty"`  // M
	TransCreateTime   string `url:"trans_create_time,omitempty"`   // M
	Memo              string `url:"memo,omitempty"`
	BizProduct        string `url:"biz_product,omitempty"` // M
	ExtendInfo        string `url:"extend_info,omitempty"` // M
}

type PayResp struct {
	scanpay1.CommonResp

	Response struct {
		Alipay struct {
			ResultCode         string `url:"result_code" xml:"result_code"`
			Error              string `url:"error" xml:"error"`
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

func (p *PayResp) ResultCode() string {
	return p.Response.Alipay.ResultCode
}

func (p *PayResp) ErrorCode() string {
	if p.CommonResp.ReqFlag() {
		return p.Response.Alipay.Error
	}
	return p.CommonResp.Error
}
