package weixin

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// params
const (
	partner  = "2088811767473826"
	charSet  = "utf-8"
	currency = "156"
)

var DefaultClient wxp

const (
	md5        = "12sdffjjguddddd2widousldadi9o0i1"
	mch_id     = "1236593202"
	appid      = "wx25ac886b6dac7dd2"
	acqfee     = "0.02"
	merfee     = "0.03"
	fee        = "0.01"
	sub_mch_id = "1247075201"
)

type wxp struct{}

// ProcessBarcodePay 扫条码下单
func (w *wxp) ProcessBarcodePay(req *model.ScanPay) *model.ScanPayResponse {
	weixinReq := &weiXinRequest{
		Appid:          appid,
		MchId:          mch_id,
		NonceStr:       getRandomStr(),
		TotalFee:       req.Txamt,
		OutTradeNo:     req.OrderNum,
		FeeType:        "CNY",
		SpbillCreateIp: "10.10.10.1",
		Body:           "sdfdfsdds",
		AuthCode:       req.ScanCodeId,

		/*
			DeviceInfo :req.
			GoodsTag       :req.
			Detail     :req.
			Attach     :req.
			Sign           :req.
		*/
	}

	dict := toMap(weixinReq)

	weixinRep := sendRequest(dict, req.Key)

	log.Debugf("weixin response: %+v", weixinRep)

	return transform(weixinReq.Service, weixinRep, req.Response)

}

func toMap(req *weiXinRequest) map[string]string {

	dict := make(map[string]string)

	// 固定参数
	dict["_input_charset"] = charSet
	dict["partner"] = partner
	dict["currency"] = currency
	dict["seller_id"] = partner
	// 参数转换
	dict["appid"] = req.AppId
	dict["mch_id"] = req.MchId
	dict["nonce_str"] = req.NonceStr
	dict["body"] = req.Body
	dict["out_trade_no"] = req.OutTradeNo
	dict["total_fee"] = req.TotalFee
	dict["spbill_create_ip"] = req.SpbillCreateIp
	dict["auto_code"] = req.AuthCode
	// ...

	return dict
}

func getRandomStr() string {
	return "sdfsfdsdfsfds"
}

// ProcessQrCodeOfflinePay 扫二维码预下单
func (w *wxp) ProcessQrCodeOfflinePay(req *model.ScanPay) *model.ScanPayResponse {
	return nil

}

// ProcessRefund 退款
func (w *wxp) ProcessRefund(req *model.ScanPay) *model.ScanPayResponse {
	return nil

}

// ProcessEnquiry 查询
func (w *wxp) ProcessEnquiry(req *model.ScanPay) *model.ScanPayResponse {
	return nil

}
