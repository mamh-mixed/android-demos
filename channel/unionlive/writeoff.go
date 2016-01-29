package unionlive

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CardInfoLink/quickpay/channel/unionlive/coupon"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
)

func Demo() (a interface{}, err error) {
	req := &coupon.PurchaseCouponsReq{
		Header: coupon.PurchaseCouponsReqHeader{
			Version:       Version,                                              // 报文版本号  15 M 当前版本 1.0
			TransDirect:   TransDirectQ,                                         // 交易方向  1 M Q:请求
			TransType:     "W412",                                               // 交易类型 8 M 本交易固定值 W412
			MerchantId:    "182000001000000",                                    // 商户编号 15 M 由优麦圈后台分配给商户的编号
			SubmitTime:    time.Now().Format("20060102150405"),                  // 交易提交时间 14 M 固定格式:yyyyMMddHHmmss
			ClientTraceNo: fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()), // 客户端交易流水号 40 M 客户端的唯一交易流水号
		},
		Body: coupon.PurchaseCouponsReqBody{
			CouponsNo: "1816086060100100", // 优麦圈电子券号 50 M 优麦圈电子券号
			TermId:    "00000667",         // 终端编号 8 M 由优麦圈后台分配给该终端的编号
			// TermSn:    "9e908a255b3e5989", // 终端唯一序列号 100 M 商户终端对应的硬件唯一序列号
			Amount:      1, // 要验证的次数  10 M 要验证该券码的次数,次数必须大于0
			ExtMercId:   "100000000010001",
			ExtTermId:   "1000134",
			TransAmount: 21000,
		},
		SpReq: &model.ScanPayRequest{},
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
			TransDirect:   TransDirectQ,
			TransType:     "W394",
			MerchantId:    "182000001000000",
			SubmitTime:    time.Now().Format("20060102150405"),
			ClientTraceNo: fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		},
		Body: coupon.QueryPurchaseCouponsResultReqBody{
			CouponsNo: "1818303006004106",
			TermId:    "00000667",
			// TermSn:           "9e908a255b3e5989",
			ExtMercId:        "100000000010001",
			ExtTermId:        "1000134",
			Amount:           1,
			OldClientTraceNo: "14494757911298498081",
			OldSubmitTime:    "20151207160951",
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
			TransDirect:   TransDirectQ,
			TransType:     "W395",
			MerchantId:    "182000001000000",
			SubmitTime:    time.Now().Format("20060102150405"),
			ClientTraceNo: fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		},
		Body: coupon.QueryPurchaseLogReqBody{
			TermId: "00000667",
			// Termsn:    "9e908a255b3e5989",
			ExtMercId: "100000000010001",
			ExtTermId: "1000134",
			PageIndex: 0,
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

// PurchaseActCouponsDemo W452-刷卡活动券验证
func PurchaseActCouponsDemo() (a interface{}, err error) {
	req := &coupon.PurchaseActCouponsReq{
		Header: coupon.PurchaseActCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W452",
			MerchantId:    "182000001000000",
			SubmitTime:    time.Now().Format("20060102150405"),
			ClientTraceNo: fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		},
		Body: coupon.PurchaseActCouponsReqBody{
			CouponsNo:      "1818303006004106",
			OldHostTraceNo: "39d2eda6-7ee3-453a-97bb-882238d1b446",
			TermId:         "00000667",
			// TermSn:           "9e908a255b3e5989",
			ExtMercId:   "100000000010001",
			ExtTermId:   "1000134",
			Amount:      1,
			Cardbin:     "622525",
			TransAmount: 100,
			PayType:     2,
		},
		SpReq: &model.ScanPayRequest{},
	}
	resp := &coupon.PurchaseActCouponsResp{}
	if err = Execute(req, resp); err != nil {
		log.Errorf("sendRequest fail, service=PurchaseActCoupons, channel=UNIONLIVE,%s", err)
		return nil, err
	}
	log.Debugf("%#v", resp)
	return
}

// UndoPurchaseActCouponsDemo W492-刷卡活动券验证撤销
func UndoPurchaseActCouponsDemo() (a interface{}, err error) {
	req := &coupon.UndoPurchaseActCouponsReq{
		Header: coupon.UndoPurchaseActCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W492",
			MerchantId:    "182000001000000",
			SubmitTime:    time.Now().Format("20060102150405"),
			ClientTraceNo: fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		},
		Body: coupon.UndoPurchaseActCouponsReqBody{
			CouponsNo: "1818303006004106",
			TermId:    "00000667",
			// TermSn:           "9e908a255b3e5989",
			ExtMercId:        "100000000010001",
			ExtTermId:        "1000134",
			OldTransAmount:   1,
			OldSubmitTime:    "20151207160951",
			OldClientTraceNo: "14494757911298498081",
			OldHostTraceNo:   "805024",
		},
		SpReq: &model.ScanPayRequest{},
	}
	resp := &coupon.UndoPurchaseActCouponsResp{}
	if err = Execute(req, resp); err != nil {
		log.Errorf("sendRequest fail, service=UndoPurchaseActCoupons, channel=UNIONLIVE,%s", err)
		return nil, err
	}
	log.Debugf("%#v", resp)
	return
}

// QueryCouponsPackageDemo W396-礼包券查询
func QueryCouponsPackageDemo() (a interface{}, err error) {
	req := &coupon.QueryCouponsPackageReq{
		Header: coupon.QueryCouponsPackageReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W394",
			MerchantId:    "182000001000000",
			SubmitTime:    time.Now().Format("20060102150405"),
			ClientTraceNo: fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		},
		Body: coupon.QueryCouponsPackageReqBody{
			CouponsNo:      "1802702004000305",
			OldHostTraceNo: "",
			TermId:         "00000667",
			// TermSn:           "9e908a255b3e5989",
			ExtMercId: "100000000010001",
			ExtTermId: "1000134",
			Amount:    1,
		},
		SpReq: &model.ScanPayRequest{},
	}
	resp := &coupon.QueryCouponsPackageResp{}
	if err = Execute(req, resp); err != nil {
		log.Errorf("sendRequest fail, service=QueryCouponsPackage, channel=UNIONLIVE,%s", err)
		return nil, err
	}
	log.Debugf("%#v", resp)

	return
}
