package scanpay

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type HandleFuc func(req *model.ScanPay) (ret *model.ScanPayResponse)

// AsyncNotifyRouter 异步通知处理分发
func AsyncNotifyRouter(values url.Values) {

	// 渠道类型
	chcd := values.Get("scanpay_chcd")
	switch chcd {

	case "ALP":
		core.AlpAsyncNotify(values)
	case "WXP":
		core.WxpAsyncNotify(values)
	default:
		// do nothing
	}
}

// ScanPayHandle 执行扫码支付逻辑
func ScanPayHandle(reqBytes []byte) []byte {

	log.Debugf("request body: %s", string(reqBytes))
	// 解析请求内容
	req := new(model.ScanPay)
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal jsonStr(%s): %s", reqBytes, err)
		return errorResponse(req, "INVALID_PARAMETER")
	}

	// 具体业务
	ret := router(req)

	// 应答
	retBytes, err := json.Marshal(ret)
	log.Debugf("handled body: %s", retBytes)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		return errorResponse(req, "SYSTEM_ERROR")
	}
	return retBytes
}

// router 分发业务逻辑
func router(req *model.ScanPay) (ret *model.ScanPayResponse) {

	switch {
	case req.Busicd == "purc":
		// ret = barcodePay(req)
		ret = doScanPay(validateBarcodePay, core.BarcodePay, req)
	case req.Busicd == "paut":
		// ret = qrCodeOfflinePay(req)
		ret = doScanPay(validateQrCodeOfflinePay, core.QrCodeOfflinePay, req)
	case req.Busicd == "inqy":
		// ret = enquiry(req)
		ret = doScanPay(validateEnquiry, core.Enquiry, req)
	case req.Busicd == "refd":
		// ret = refund(req)
		ret = doScanPay(validateRefund, core.Refund, req)
	case req.Busicd == "void":
		// ret = cancel(req)
		ret = doScanPay(validateCancel, core.Cancel, req)
	default:
		ret = mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// 补充原信息返回
	fillResponseInfo(req, ret)

	return ret
}

// doScanPay 执行业务逻辑
func doScanPay(validateFuc, processFuc HandleFuc, req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证字段
	if validateFuc(req); ret != nil {
		return ret
	}

	mer, err := mongo.MerchantColl.Find(req.Mchntid)
	if err != nil {
		return mongo.OffLineRespCd("NO_MERCHANT_MATCH") // todo check error code
	}

	// TODO valid sign

	// 验证接口权限
	if !strings.Contains(strings.Join(mer.Permission, ","), req.Busicd) {
		log.Errorf("merchant %s request %s interface without permission!", req.Mchntid, req.Busicd)
		return mongo.OffLineRespCd("NO_PERMISSION") // todo check error code
	}

	// process
	ret = processFuc(req)

	// TODO sign

	return ret

}

func fillResponseInfo(req *model.ScanPay, ret *model.ScanPayResponse) {

	// 如果空白，默认将原信息返回
	if ret.Busicd == "" {
		ret.Busicd = req.Busicd
	}
	if ret.Inscd == "" {
		ret.Inscd = req.Inscd
	}
	if ret.Mchntid == "" {
		ret.Mchntid = req.Mchntid
	}
	if ret.Txamt == "" {
		ret.Txamt = req.Txamt
	}
	if ret.OrigOrderNum == "" {
		ret.OrigOrderNum = req.OrigOrderNum
	}
	if ret.OrderNum == "" {
		ret.OrderNum = req.OrderNum
	}
	// TODO
	if ret.Sign == "" {
		ret.Sign = req.Sign
	}
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
