package oversea

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
)

type PayReq struct {
	scanpay1.CommonParams

	AlipaySellerId    string `url:"alipay_seller_id"`
	Quantity          string `url:"quantity"`
	TransName         string `url:"trans_name"`
	PartnerTransId    string `url:"partner_trans_id"`
	Currency          string `url:"currency"`
	TransAmount       string `url:"trans_amount"`
	BuyerIdentityCode string `url:"buyer_identity_code"`
	IdentityCodeType  string `url:"identity_code_type"`
	TransCreateTime   string `url:"trans_create_time"`
	Memo              string `url:"memo"`
	BizProduct        string `url:"biz_product"`
	ExtendInfo        string `url:"extend_info"`
}

type PayResp struct {
}
