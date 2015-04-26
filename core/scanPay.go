package core

import (
	// "github.com/omigo/log"
	"github.com/CardInfoLink/quickpay/model"
)

// BarcodePay 条码下单
func BarcodePay(req *model.QrCodePay) (resp *model.QrCodePayResponse) {

	// TODO 判断订单是否存在

	// TODO 记录该笔交易

	// TODO 判断渠道商户

	// TODO 转换参数

	// TODO 获得渠道实例，请求

	// TODO 根据请求结果更新

	return
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.QrCodePay) (resp *model.QrCodePrePayResponse) {

	// TODO 判断订单是否存在

	// TODO 记录该笔交易

	// TODO 判断渠道商户

	// TODO 转换参数

	// TODO 获得渠道实例，请求

	// TODO 根据请求结果更新

	return
}

// Refund 退款
func Refund(req *model.QrCodePay) (resp *model.QrCodeRefundResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}

// Enquiry 查询
func Enquiry(req *model.QrCodePay) (resp *model.QrCodeEnquiryResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 如果订单状态位处理中，则向渠道查询

	// TODO 更新订单状态

	return
}

// Cancel 撤销
func Cancel(req *model.QrCodePay) (resp *model.QrCodeCancelResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}
