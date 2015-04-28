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
		// TODO check the err response message
		return []byte("params invalid")
	}

	var resp interface{}
	switch {
	// TODO
	case req.Busicd == "purc":
		resp = BarcodePay(req)
	case req.Busicd == "paut":
		resp = QrCodeOfflinePay(req)
	case req.Busicd == "inqy":
		resp = Enquiry(req)
	case req.Busicd == "refd":
		resp = Refund(req)
	case req.Busicd == "void":
		resp = Cancel(req)
	default:
		return []byte(fmt.Sprintf("no busicd: %s", req.Busicd))
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", resp, err)
		// TODO retrun system error string
		return []byte("system error")
	}
	return respBytes
}

// BarcodePay 条码下单
func BarcodePay(req *model.ScanPay) (resp *model.QrCodePayResponse) {
	log.Debugf("request body: %+v", req)

	// TODO validite field

	// process
	return core.BarcodePay(req)
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPay) (resp *model.QrCodePrePayResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.QrCodeOfflinePay(req)
}

// Refund 退款
func Refund(req *model.ScanPay) (resp *model.QrCodeRefundResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Refund(req)
}

// Enquiry 查询
func Enquiry(req *model.ScanPay) (resp *model.QrCodeEnquiryResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Enquiry(req)
}

// Cancel 撤销
func Cancel(req *model.ScanPay) (resp *model.QrCodeCancelResponse) {

	log.Debugf("request body: %+v", req)

	// TODO validite field
	return core.Cancel(req)
}
