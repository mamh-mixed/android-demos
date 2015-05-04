package core

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"strconv"
	"strings"
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
		return resp
	}

	// 渠道选择
	// 根据扫码Id判断走哪个渠道
	cardBrand := ""
	if strings.HasPrefix(req.ScanCodeId, "1") {
		req.Chcd = "Weixin"
		cardBrand = "WXP"
	} else if strings.HasPrefix(req.ScanCodeId, "2") {
		req.Chcd = "Alipay"
		cardBrand = "ALP"
	} else {
		// 不送，返回 TODO check error code
		resp.ErrorDetail = "SYSTEM_ERROR"
		return resp
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
	err = mongo.TransColl.Add(t)
	if err != nil {
		log.Errorf("add trans(%+v) fail: %s", t, err)
		resp.ErrorDetail = "SYSTEM_ERROR"
		return resp
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, cardBrand)
	if rp == nil {
		// TODO check error code
		resp.ErrorDetail = "SYSTEM_ERROR"
		return resp
	}

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		// TODO check error code
		log.Errorf("not found any chanMer(%s,%s): %s", rp.ChanCode, rp.ChanMerId, err)
		resp.ErrorDetail = "SYSTEM_ERROR"
		return resp
	}

	// 上送参数
	req.SysOrderNum = tools.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.Key = c.SignCert
	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	resp = sp.ProcessBarcodePay(req)

	// 根据请求结果更新
	t.ChanRespCode = resp.ChanRespCode
	t.RespCode = resp.RespCode
	switch resp.RespCode {
	case "000000":
		t.TransStatus = model.TransSuccess
	case "000009":
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
	if err = mongo.TransColl.Update(t); err != nil {
		log.Errorf("update trans(%+v) status fail: %s ", t, err)
	}

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
