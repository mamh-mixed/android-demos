package alp

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

var DefaultClient alp

const (
	partner  = "2088811767473826"
	charSet  = "UTF-8"
	currency = "156"
)

// alp 当面付，扫码支付
type alp struct{}

// ProcessBarcodePay 条码支付/下单
func (a *alp) ProcessBarcodePay(req *model.ScanPay) *model.QrCodePayResponse {

	alpReq := &alpRequest{
		Service:       "alipay.acquire.createandpay",
		NotifyUrl:     req.NotifyUrl,
		OutTradeNo:    req.ChannelOrderNum,
		Subject:       "test", // TODO
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

	resp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", resp)

	// 请求成功
	ret := new(model.QrCodePayResponse)
	if resp.IsSuccess == "T" {
		alipay := resp.Response.Alipay
		ret.ChanRespCode = alipay.ResultCode
		switch alipay.ResultCode {
		case "ORDER_SUCCESS_PAY_SUCCESS":

			ret.ChannelOrderNum = alipay.TradeNo
			ret.ConsumerAccount = alipay.BuyerLogonId
			ret.ConsumerId = alipay.BuyerUserId
			// 计算折扣
			ret.MerDiscount, ret.ChcdDiscount = alipay.DisCount()

		case "ORDER_FAIL":
			// do nothing
		default:
			ret.ChannelOrderNum = alipay.TradeNo
		}
	} else {
		ret.ChanRespCode = resp.Error
		ret.ErrorDetail = "SYSTEM_ERROR"
	}

	return ret
}

// ProcessQrCodeOfflinePay 扫码支付/预下单
func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPay) *model.QrCodePrePayResponse {

	alpReq := &alpRequest{
		Service:       "alipay.acquire.createandpay",
		NotifyUrl:     "",
		OutTradeNo:    req.ChanOrderNum,
		Subject:       "",
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

	resp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", resp)

	return nil
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPay) *model.QrCodeRefundResponse {

	alpReq := &alpRequest{
		Service:       "alipay.acquire.refund",
		NotifyUrl:     "",
		OutTradeNo:    req.ChanOrderNum,
		Subject:       "",
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

	resp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", resp)

	return nil
}

// ProcessEnquiry 查询，包含支付、退款
func (a *alp) ProcessEnquiry(req *model.ScanPay) *model.QrCodeEnquiryResponse {

	alpReq := &alpRequest{
		Service:       "alipay.acquire.query",
		NotifyUrl:     "",
		OutTradeNo:    req.ChanOrderNum,
		Subject:       "",
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

	resp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", resp)

	return nil
}

// ProcessVoid 撤销
func (a *alp) ProcessCancel(req *model.ScanPay) *model.QrCodeCancelResponse {

	alpReq := &alpRequest{
		Service:       "alipay.acquire.cancel",
		NotifyUrl:     "",
		OutTradeNo:    req.ChanOrderNum,
		Subject:       "",
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

	resp := sendRequest(dict, req.Key)
	log.Debugf("alp response: %+v", resp)

	return nil
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
	dict["product_code"] = req.ProductCode
	dict["total_fee"] = req.TotalFee
	dict["extend_params"] = req.ExtendParams
	dict["it_b_pay"] = req.ItBPay
	dict["dynamic_id_type"] = req.DynamicIdType
	dict["dynamic_id"] = req.DynamicId
	dict["goods_detail"] = req.GoodsDetail

	// ...

	return dict
}
