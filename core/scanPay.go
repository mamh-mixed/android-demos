package core

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"strconv"
)

// BarcodePay 条码下单
// TODO 根据errorDetail找到respcd
func BarcodePay(req *model.ScanPay) (resp *model.QrCodePayResponse) {

	resp = new(model.QrCodePayResponse)
	// 判断订单是否存在
	count, err := mongo.TransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		resp.ErrorDetail = "SYSTEM_ERROR"
		return resp
	}
	if count > 0 {
		// 没有订单号重复代码
		resp.ErrorDetail = "AUTH_NO_ERROR"
		return resp
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		resp.ErrorDetail = "SYSTEM_ERROR"
		return
	}

	// 渠道选择
	if req.Chcd != "" {
		// TODO 根据扫码Id判断走哪个渠道
	}

	// 记录该笔交易
	txamt := int64(f * 100)
	t := &model.Trans{
		MerId:     req.Mchntid,
		OrderNum:  req.OrderNum,
		TransType: model.PayTrans,
		ChanCode:  req.Chcd,
		TransAmt:  txamt,
	}
	log.Debug(t)
	// TODO 判断渠道商户

	// TODO 转换参数

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	resp = sp.ProcessBarcodePay(req)

	// 根据请求结果更新

	return
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPay) (resp *model.QrCodePrePayResponse) {

	// TODO 判断订单是否存在

	// TODO 记录该笔交易

	// TODO 判断渠道商户

	// TODO 转换参数

	// TODO 获得渠道实例，请求

	// TODO 根据请求结果更新

	return
}

// Refund 退款
func Refund(req *model.ScanPay) (resp *model.QrCodeRefundResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}

// Enquiry 查询
func Enquiry(req *model.ScanPay) (resp *model.QrCodeEnquiryResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 如果订单状态位处理中，则向渠道查询

	// TODO 更新订单状态

	return
}

// Cancel 撤销
func Cancel(req *model.ScanPay) (resp *model.QrCodeCancelResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}
