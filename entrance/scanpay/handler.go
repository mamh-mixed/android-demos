package scanpay

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

func Router(reqBytes []byte) []byte {

	req := new(model.ScanPay)
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal jsonStr(%s): %s", reqBytes, err)
		// TODO check the err retonse message
		return []byte("params invalid")
	}

	// TODO valid sign

	// var ret *ScanPayResponse
	ret := new(model.ScanPayResponse)
	switch {
	// TODO
	case req.Busicd == "purc":
		ret = BarcodePay(req)
	case req.Busicd == "paut":
		ret = QrCodeOfflinePay(req)
	case req.Busicd == "inqy":
		ret = Enquiry(req)
	case req.Busicd == "refd":
		ret = Refund(req)
	case req.Busicd == "void":
		ret = Cancel(req)
	default:
		return []byte(fmt.Sprintf("no busicd: %s", req.Busicd))
	}
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		// TODO retrun system error string
		return []byte("system error")
	}
	return retBytes
}

// BarcodePay 条码下单
func BarcodePay(req *model.ScanPay) (ret *model.ScanPayResponse) {
	log.Debugf("request body: %+v", req)

	// validite field
	if ret = validateBarcodePay(req); ret == nil {
		// process
		ret = core.BarcodePay(req)
	}
	log.Debugf("handled body: %+v", ret)

	// get ret.Respcd
	ret.Respcd = responseCode(ret.ErrorDetail, ret.Chcd)

	// retonse info
	ret.Busicd = req.Busicd
	ret.Inscd = req.Inscd
	ret.Mchntid = req.Mchntid
	ret.Sign = req.Sign
	ret.Txamt = req.Txamt
	ret.Txndir = "A"

	return
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.QrCodeOfflinePay(req)
}

// Refund 退款
func Refund(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Refund(req)
}

// Enquiry 查询
func Enquiry(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Enquiry(req)
}

// Cancel 撤销
func Cancel(req *model.ScanPay) (ret *model.ScanPayResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Cancel(req)
}
