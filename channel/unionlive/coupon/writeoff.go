package coupon

import (
	"github.com/CardInfoLink/quickpay/channel/unionlive"
	"github.com/omigo/log"
)

func Demo() (a interface{}, err error) {
	req := &PurchaseCouponsReq{
		Header: PurchaseCouponsReqHeader{
			Version:       unionlive.Version,      // 报文版本号  15 M 当前版本 1.0
			TransDirect:   unionlive.TransDirectQ, // 交易方向  1 M Q:请求
			TransType:     "W412",                 // 交易类型 8 M 本交易固定值 W412
			MerchantId:    "182000001000000",      // 商户编号 15 M 由优麦圈后台分配给商户的编号
			SubmitTime:    "20130501201012",       // 交易提交时间 14 M 固定格式:yyyyMMddHHmmss
			ClientTraceNo: "497540",               // 客户端交易流水号 40 M 客户端的唯一交易流水号
		},
		Body: PurchaseCouponsReqBody{
			CouponsNo: "1809706004000705", // 优麦圈电子券号 50 M 优麦圈电子券号
			TermId:    "00000667",         // 终端编号 8 M 由优麦圈后台分配给该终端的编号
			TermSn:    "9e908a255b3e5989", // 终端唯一序列号 100 M 商户终端对应的硬件唯一序列号
			Amount:    1,                  // 要验证的次数  10 M 要验证该券码的次数,次数必须大于0
		},
	}
	resp := &PurchaseCouponsResp{}
	if err = unionlive.Execute(req, resp); err != nil {
		return nil, err
	}

	log.Debugf("%#v", resp)
	// process resp code, return
	return
}
