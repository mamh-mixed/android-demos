package scanpay

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"net/url"
)

// AsyncNotifyHandle 异步通知处理
func AsyncNotifyHandle(values url.Values) {

	// 渠道类型
	chcd := values.Get("chcd")
	switch chcd {

	case "ALP":
		// TODO check sign
		result := true
		if result {
			core.AlpAsyncNotify(values)
		} else {
			log.Errorf("check sign error, chcd=%s, request data=(%+v)", chcd, values)
		}
	case "WXP":
		// TODO check sign
		result := true
		if result {
			core.WxpAsyncNotify(values)
		} else {
			log.Errorf("check sign error, chcd=%s, request data=(%+v)", chcd, values)
		}
	default:
		// do nothing
	}
}

// Router tcp请求路由
func Router(reqBytes []byte) []byte {

	req := new(model.ScanPay)
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal jsonStr(%s): %s", reqBytes, err)
		return errorResponse(req, "INVALID_PARAMETER")
	}

	// TODO valid sign

	ret := new(model.ScanPayResponse)
	switch {
	case req.Busicd == "purc":
		ret = barcodePay(req)
	case req.Busicd == "paut":
		ret = qrCodeOfflinePay(req)
	case req.Busicd == "inqy":
		ret = enquiry(req)
	case req.Busicd == "refd":
		ret = refund(req)
	case req.Busicd == "void":
		ret = cancel(req)
	default:
		return errorResponse(req, "INVALID_PARAMETER")
	}

	// TODO sign

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		return errorResponse(req, "SYSTEM_ERROR")
	}
	return retBytes
}

// barcodePay 条码下单
func barcodePay(req *model.ScanPay) (ret *model.ScanPayResponse) {
	log.Debugf("request body: %+v", req)

	// validate field
	if ret = validateBarcodePay(req); ret == nil {
		// process
		ret = core.BarcodePay(req)
	}

	// 补充原信息返回
	fillResponseInfo(req, ret)

	log.Debugf("handled body: %+v", ret)

	return ret
}

// qrCodeOfflinePay 扫二维码预下单
func qrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// validate field
	if ret = validateQrCodeOfflinePay(req); ret == nil {
		// process
		ret = core.QrCodeOfflinePay(req)
	}

	// 补充原信息返回
	fillResponseInfo(req, ret)

	log.Debugf("handled body: %+v", ret)

	return ret
}

// refund 退款
func refund(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// validate field
	if ret = validateRefund(req); ret == nil {
		// process
		ret = core.Refund(req)
	}

	// 补充原信息返回
	fillResponseInfo(req, ret)

	log.Debugf("handled body: %+v", ret)

	return ret
}

// enquiry 查询
func enquiry(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	if ret = validateEnquiry(req); ret == nil {
		// process
		ret = core.Enquiry(req)
		// 直接返回，查询得到的是原交易信息，不需要补充返回信息
		return ret
	}

	// 错误信息补充完整
	fillResponseInfo(req, ret)

	return ret

}

// cancel 撤销
func cancel(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Cancel(req)
}

func fillResponseInfo(req *model.ScanPay, ret *model.ScanPayResponse) {

	// 默认将原信息返回
	ret.Busicd = req.Busicd
	ret.Inscd = req.Inscd
	ret.Mchntid = req.Mchntid
	ret.Sign = req.Sign // TODO
	ret.Txamt = req.Txamt
	ret.OrigOrderNum = req.OrigOrderNum
	ret.OrderNum = req.OrderNum
	ret.Txndir = "A"
}

// errorResponse 返回错误信息
func errorResponse(req *model.ScanPay, errorCode string) []byte {

	ret := mongo.OffLineRespCd(errorCode)
	ret.Busicd = req.Busicd
	ret.Txndir = "A"

	bytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
	}
	return bytes
}
