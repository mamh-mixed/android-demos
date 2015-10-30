package unionlive

import (
	"time"

	"github.com/CardInfoLink/quickpay/channel/unionlive/coupon"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

func Demo() (a interface{}, err error) {
	req := &coupon.PurchaseCouponsReq{
		Header: coupon.PurchaseCouponsReqHeader{
			Version:       Version,           // 报文版本号  15 M 当前版本 1.0
			TransDirect:   TransDirectQ,      // 交易方向  1 M Q:请求
			TransType:     "W412",            // 交易类型 8 M 本交易固定值 W412
			MerchantId:    "182000001000000", // 商户编号 15 M 由优麦圈后台分配给商户的编号
			SubmitTime:    "20130501201012",  // 交易提交时间 14 M 固定格式:yyyyMMddHHmmss
			ClientTraceNo: "497540",          // 客户端交易流水号 40 M 客户端的唯一交易流水号
		},
		Body: coupon.PurchaseCouponsReqBody{
			CouponsNo: "1809706004000705", // 优麦圈电子券号 50 M 优麦圈电子券号
			TermId:    "00000667",         // 终端编号 8 M 由优麦圈后台分配给该终端的编号
			TermSn:    "9e908a255b3e5989", // 终端唯一序列号 100 M 商户终端对应的硬件唯一序列号
			Amount:    1,                  // 要验证的次数  10 M 要验证该券码的次数,次数必须大于0
		},
	}
	resp := &coupon.PurchaseCouponsResp{}
	if err = Execute(req, resp); err != nil {
		return nil, err
	}

	log.Debugf("%#v", resp)
	// process resp code, return
	return
}

// W394-电子券验证结果查询
func QueryPurchaseCouponsResultDemo() (a interface{}, err error) {
	req := &coupon.QueryPurchaseCouponsResultReq{
		Header: coupon.QueryPurchaseCouponsResultReqHeader{
			Version:       Version,
			Transdirect:   TransDirectQ,
			Transtype:     "W394",
			Merchantid:    "182000001000000",
			Submittime:    time.Now().Format("20060102150405"),
			Clienttraceno: "1446102368374",
		},
		Body: coupon.QueryPurchaseCouponsResultReqBody{
			Couponsno:        "1802702004000305",
			Termid:           "00000667",
			Termsn:           "9e908a255b3e5989",
			Amount:           1,
			Oldclienttraceno: "1446109183201",
			Oldsubmittime:    "20151029170019",
			// Couponsno:        "1808700004000875",
			// Oldclienttraceno: "1446029945052",
			// Oldsubmittime:    "20151028185926",
		},
		SpReq: &model.ScanPayRequest{},
	}
	resp := &coupon.QueryPurchaseCouponsResultResp{}
	if err = Execute(req, resp); err != nil {
		log.Errorf("sendRequest fail, service=QueryPurchaseCouponsResult, channel=UNIONLIVE,%s", err)
		return nil, err
	}
	log.Debugf("%#v", resp)
	// process resp code, return
	return
}

func QueryPurchaseLogDemo() (a interface{}, err error) {
	req := &coupon.QueryPurchaseLogReq{
		Header: coupon.QueryPurchaseLogReqHeader{
			Version:       Version,
			Transdirect:   TransDirectQ,
			Transtype:     "W395",
			Merchantid:    "182000001000000",
			Submittime:    time.Now().Format("20060102150405"),
			Clienttraceno: "144611269437",
		},
		Body: coupon.QueryPurchaseLogReqBody{
			Termid:    "00000667",
			Termsn:    "9e908a255b3e5989",
			Pageindex: 0,
		},
		SpReq: &model.ScanPayRequest{},
	}
	resp := &coupon.QueryPurchaseLogResp{}
	if err = Execute(req, resp); err != nil {
		log.Errorf("sendRequest fail, service=QueryPurchaseCouponsResult, channel=UNIONLIVE,%s", err)
		return nil, err
	}
	log.Debugf("%#v", resp)
	// process resp code, return
	return
}
