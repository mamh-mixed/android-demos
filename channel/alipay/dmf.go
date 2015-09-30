package alipay

import (
	"errors"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
)

const NotifyURL = "/scanpay/upNotify/alipay"

var DefaultClient alp
var alipayNotifyUrl = goconf.Config.AlipayScanPay.NotifyUrl + NotifyURL

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
func (a *alp) ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:       req.ChanMerId,
		Service:       createAndPay,
		NotifyUrl:     req.NotifyUrl,
		OutTradeNo:    req.OrderNum, // 送的是原订单号，不转换
		Subject:       req.Subject,
		GoodsDetail:   req.AlpMarshalGoods(),
		ProductCode:   "BARCODE_PAY_OFFLINE",
		TotalFee:      req.ActTxamt,
		ExtendParams:  req.ExtendParams, //...
		ItBPay:        "40m",            // 超时时间
		DynamicIdType: "bar_code",
		DynamicId:     req.ScanCodeId,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignKey)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, createAndPay)
		return nil, err
	}

	// 处理结果返回
	return transform(alpReq.Service, alpResp)
}

// ProcessQrCodeOfflinePay 扫码支付/预下单
func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:        req.ChanMerId,
		Service:        preCreate,
		NotifyUrl:      alipayNotifyUrl,
		OutTradeNo:     req.OrderNum, // 送的是原订单号，不转换,
		Subject:        req.Subject,
		GoodsDetail:    req.AlpMarshalGoods(),
		PassbackParams: req.SysOrderNum, // 传系统订单号，异步通知时可用
		ProductCode:    "QR_CODE_OFFLINE",
		TotalFee:       req.ActTxamt,
		ExtendParams:   req.ExtendParams,
		ItBPay:         "40m", // 超时时间
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignKey)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, preCreate)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:      req.ChanMerId,
		Service:      refund,
		NotifyUrl:    req.NotifyUrl,
		OutTradeNo:   req.OrigOrderNum,
		RefundAmount: req.ActTxamt,
		OutRequestNo: req.OrderNum, //该字段上送才能部分退款，如果不送则只能全额退款
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignKey)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, refund)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessEnquiry 查询，包含支付、退款
func (a *alp) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:    req.ChanMerId,
		Service:    query,
		OutTradeNo: req.OrigOrderNum, // 送的是原订单号，不转换
	}
	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignKey)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, query)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessCancel 撤销
func (a *alp) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:    req.ChanMerId,
		Service:    cancel,
		NotifyUrl:  req.NotifyUrl,
		OutTradeNo: req.OrigOrderNum,
	}

	// req to map
	dict := toMap(alpReq)

	alpResp, err := sendRequest(dict, req.SignKey)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, cancel)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessClose 关闭接口即撤销接口
func (a *alp) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return a.ProcessCancel(req)
}

// ProcessRefundQuery 退款查询
func (a *alp) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, errors.New("not support now")
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
	dict["it_b_pay"] = req.ItBPay
	dict["dynamic_id_type"] = req.DynamicIdType
	dict["dynamic_id"] = req.DynamicId
	dict["refund_amount"] = req.RefundAmount
	dict["out_request_no"] = req.OutRequestNo
	dict["passback_parameters"] = req.PassbackParams

	// utf-8 -> gbk
	e := mahonia.NewEncoder("gbk")
	if req.Subject != "" {
		dict["subject"] = e.ConvertString(req.Subject)
	}
	if req.GoodsDetail != "" {
		dict["goods_detail"] = e.ConvertString(req.GoodsDetail)
	}
	if req.ExtendParams != "" {
		dict["extend_params"] = e.ConvertString(req.ExtendParams)
	}

	return dict
}
