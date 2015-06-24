package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
)

var DefaultClient alp

// alp 当面付，扫码支付
type alp struct{}

// service
const (
	createAndPay = "alipay.acquire.createandpay"
	preCreate    = "alipay.acquire.precreate"
	refund       = "alipay.acquire.refund"
	query        = "alipay.acquire.query"
	cancel       = "alipay.acquire.cancel"
)

// params
const (
	// partner  = "2088811767473826" // only for test
	charSet  = "utf-8"
	currency = "156"
)

// ProcessBarcodePay 条码支付/下单
func (a *alp) ProcessBarcodePay(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Partner:       req.ChanMerId,
		Service:       createAndPay,
		NotifyUrl:     req.NotifyUrl,
		OutTradeNo:    req.OrderNum, // 送的是原订单号，不转换
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

	alpResp, err := sendRequest(dict, req.SignCert)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, createAndPay)
	}

	// 处理结果返回
	return transform(alpReq.Service, alpResp, err)
}

// ProcessQrCodeOfflinePay 扫码支付/预下单
func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Partner:      req.ChanMerId,
		Service:      preCreate,
		NotifyUrl:    "",
		OutTradeNo:   req.OrderNum, // 送的是原订单号，不转换,
		Subject:      req.Subject,
		GoodsDetail:  req.MarshalGoods(),
		ProductCode:  "QR_CODE_OFFLINE",
		TotalFee:     req.Txamt,
		ExtendParams: "",
		ItBPay:       "1m", // 超时时间
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignCert)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, preCreate)
	}

	return transform(alpReq.Service, alpResp, err)
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Partner:      req.ChanMerId,
		Service:      refund,
		NotifyUrl:    "",
		OutTradeNo:   req.OrigOrderNum,
		RefundAmount: req.Txamt,
		OutRequestNo: req.OrderNum, // 送的是原订单号，不转换
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignCert)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, refund)
	}

	return transform(alpReq.Service, alpResp, err)
}

// ProcessEnquiry 查询，包含支付、退款
func (a *alp) ProcessEnquiry(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Partner:    req.ChanMerId,
		Service:    query,
		OutTradeNo: req.OrderNum, // 送的是原订单号，不转换
	}
	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignCert)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, query)
	}

	return transform(alpReq.Service, alpResp, err)
}

// ProcessCancel 撤销
func (a *alp) ProcessCancel(req *model.ScanPay) *model.ScanPayResponse {

	alpReq := &alpRequest{
		Partner:    req.ChanMerId,
		Service:    cancel,
		NotifyUrl:  "",
		OutTradeNo: req.OrigOrderNum,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignCert)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, cancel)
	}

	return transform(alpReq.Service, alpResp, err)
}

func toMap(req *alpRequest) map[string]string {

	dict := make(map[string]string)

	// 固定参数
	dict["_input_charset"] = charSet
	dict["partner"] = req.Partner
	dict["currency"] = currency
	dict["seller_id"] = req.Partner
	// 参数转换
	dict["service"] = req.Service
	dict["notify_url"] = req.NotifyUrl
	dict["product_code"] = req.ProductCode
	dict["out_trade_no"] = req.OutTradeNo
	dict["total_fee"] = req.TotalFee
	dict["extend_params"] = req.ExtendParams
	dict["it_b_pay"] = req.ItBPay
	dict["dynamic_id_type"] = req.DynamicIdType
	dict["dynamic_id"] = req.DynamicId
	dict["refund_amount"] = req.RefundAmount

	// utf-8 -> gbk
	e := mahonia.NewEncoder("gbk")
	if req.Subject != "" {
		dict["subject"] = e.ConvertString(req.Subject)
	}
	if req.GoodsDetail != "" {
		dict["goods_detail"] = e.ConvertString(req.GoodsDetail)
	}

	return dict
}
