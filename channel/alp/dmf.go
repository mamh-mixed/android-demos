package alp

import (
	_ "github.com/CardInfoLink/quickpay/model"
)

var Obj alp

// alp 当面付，扫码支付
type alp struct{}

// ProcessBarcodePay 条码支付
func (a *alp) ProcessBarcodePay(req *AlpRequest) *AlpResponse {

	// req to map
	dict := toMap(req)

	return sendRequest(dict, req.Key)
}

// ProcessQrCodeOfflinePay 扫码支付
func (a *alp) ProcessQrCodeOfflinePay(req *AlpRequest) *AlpResponse {

	// req to map
	dict := toMap(req)

	return sendRequest(dict, req.Key)
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *AlpRequest) *AlpResponse {
	// req to map
	dict := toMap(req)

	return sendRequest(dict, req.Key)
}

// ProcessEnquiry 查询，包含支付、退款
func (a *alp) ProcessEnquiry(req *AlpRequest) *AlpResponse {
	// req to map
	dict := toMap(req)

	return sendRequest(dict, req.Key)
}

func toMap(req *AlpRequest) map[string]string {

	dict := make(map[string]string)
	// 参数转换
	dict["server"] = req.Service
	dict["_input_charset"] = req.Charset
	dict["currency"] = req.Currency
	dict["notify_url"] = req.NotifyUrl
	dict["partner"] = req.Partner
	dict["product_code"] = req.ProductCode
	dict["out_trade_no"] = req.OutTradeNo
	dict["subject"] = req.Subject
	dict["product_code"] = req.ProductCode
	dict["total_fee"] = req.TotalFee
	dict["seller_id"] = req.SellerId
	dict["extend_params"] = req.ExtendParams
	dict["it_b_pay"] = req.ItBPay
	dict["dynamic_id_type"] = req.DynamicIdType
	dict["dynamic_id"] = req.DynamicId

	// ...

	return dict
}
