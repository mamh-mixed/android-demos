package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

var DefaultClient alp

// alp 当面付，扫码支付
type alp struct{}

// service
const (
	createAndPay = "alipay.acquire.createandpay"
	preCreate    = "precreate"
	refund       = "alipay.acquire.refund"
	query        = "alipay.acquire.query"
	cancel       = "alipay.acquire.cancel"
)

// params
const (
	partner  = "2088811767473826"
	charSet  = "utf-8"
	currency = "156"
)

// ProcessBarcodePay 条码支付/下单
func (a *alp) ProcessBarcodePay(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Service:       createAndPay,
		NotifyUrl:     req.NotifyUrl,
		OutTradeNo:    req.SysOrderNum,
		Subject:       req.Subject,
		GoodsDetail:   req.MarshalGoods(),
		ProductCode:   "BARCODE_PAY_OFFLINE",
		TotalFee:      req.Txamt,
		ExtendParams:  "",   //...
		ItBPay:        "1m", // 超时时间
		DynamicIdType: "bar_code",
		DynamicId:     req.ScanCodeId,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", alpResp)

	// 处理结果返回
	return transform(alpReq.Service, alpResp)
}

// ProcessQrCodeOfflinePay 扫码支付/预下单
func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Service:       preCreate,
		NotifyUrl:     "",
		OutTradeNo:    req.SysOrderNum,
		Subject:       req.Subject,
		GoodsDetail:   req.MarshalGoods(),
		ProductCode:   "BARCODE_PAY_OFFLINE",
		TotalFee:      req.Txamt,
		ExtendParams:  "",
		ItBPay:        "1m", // 超时时间
		DynamicIdType: "bar_code",
		DynamicId:     req.ScanCodeId,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", alpResp)

	return transform(alpReq.Service, alpResp)
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Service:       refund,
		NotifyUrl:     "",
		OutTradeNo:    req.SysOrderNum,
		Subject:       req.Subject,
		GoodsDetail:   req.MarshalGoods(),
		ProductCode:   "BARCODE_PAY_OFFLINE",
		TotalFee:      req.Txamt,
		ExtendParams:  "",
		ItBPay:        "1m", // 超时时间
		DynamicIdType: "bar_code",
		DynamicId:     req.ScanCodeId,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", alpResp)

	return transform(alpReq.Service, alpResp)
}

// ProcessEnquiry 查询，包含支付、退款
func (a *alp) ProcessEnquiry(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Service:    query,
		OutTradeNo: req.SysOrderNum,
	}
	// req to map
	dict := toMap(alpReq)

	alpResp := sendRequest(dict, req.Key)
	log.Infof("alp response: %+v", alpResp)

	return transform(alpReq.Service, alpResp)
}

// ProcessVoid 撤销
func (a *alp) ProcessCancel(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Service:       cancel,
		NotifyUrl:     "",
		OutTradeNo:    req.SysOrderNum,
		Subject:       req.Subject,
		GoodsDetail:   req.MarshalGoods(),
		ProductCode:   "BARCODE_PAY_OFFLINE",
		TotalFee:      req.Txamt,
		ExtendParams:  "",
		ItBPay:        "1m", // 超时时间
		DynamicIdType: "bar_code",
		DynamicId:     req.ScanCodeId,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", alpResp)

	return transform(alpReq.Service, alpResp)
}

func toMap(req *alpRequest) map[string]string {

	dict := make(map[string]string)

	// 固定参数
	dict["_input_charset"] = charSet
	dict["partner"] = partner
	dict["currency"] = currency
	dict["seller_id"] = partner
	// 参数转换
	dict["service"] = req.Service
	dict["notify_url"] = req.NotifyUrl
	dict["product_code"] = req.ProductCode
	dict["out_trade_no"] = req.OutTradeNo
	dict["subject"] = req.Subject
	dict["total_fee"] = req.TotalFee
	dict["extend_params"] = req.ExtendParams
	dict["it_b_pay"] = req.ItBPay
	dict["dynamic_id_type"] = req.DynamicIdType
	dict["dynamic_id"] = req.DynamicId
	dict["goods_detail"] = req.GoodsDetail

	// ...

	return dict
}
