package core

import (
	// "github.com/omigo/log"
	"github.com/CardInfoLink/quickpay/model"
)

// BarcodePay 条码下单
func BarcodePay(req *model.QrCodePay) (resp *model.QrCodePayResponse) {
	return
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.QrCodePay) (resp *model.QrCodePrePayResponse) {

	return
}

// Refund 退款
func Refund(req *model.QrCodePay) (resp *model.QrCodeRefundResponse) {

	return
}

// Enquiry 查询
func Enquiry(req *model.QrCodePay) (resp *model.QrCodeEnquiryResponse) {

	return
}
