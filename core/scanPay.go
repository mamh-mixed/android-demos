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
func BarcodePay(req *model.ScanPay) (resp *model.ScanPayResponse) {

	resp = new(model.ScanPayResponse)
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
		Busicd:    resp.Busicd,
		Inscd:     resp.Inscd,
	}
	// log.Debug(t)
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

	// 渠道
	resp.Chcd = req.Chcd

	// 根据请求结果更新
	t.ChanRespCode = resp.ErrorDetail
	t.RespCode = resp.RespCode
	t.ChanOrderNum = resp.ChannelOrderNum
	t.ChanDiscount = resp.ChcdDiscount
	t.MerDiscount = resp.MerDiscount
	t.ConsumerAccount = resp.ConsumerAccount
	t.ConsumerId = resp.ConsumerId

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

	return resp
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPay) (resp *model.ScanPayResponse) {

	// TODO 判断订单是否存在

	// TODO 记录该笔交易

	// TODO 判断渠道商户

	// TODO 转换参数

	// TODO 获得渠道实例，请求

	// TODO 根据请求结果更新

	return
}

// Refund 退款
func Refund(req *model.ScanPay) (resp *model.ScanPayResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}

// Enquiry 查询
func Enquiry(req *model.ScanPay) (resp *model.ScanPayResponse) {

	resp = new(model.ScanPayResponse)

	// 判断是否存在该订单
	t, err := mongo.TransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		resp.ErrorDetail = "TRADE_NOT_EXIST"
		return resp
	}
	log.Debugf("trans:(%+v)", t)

	// 判断订单的状态

	switch t.TransStatus {
	// 如果是处理中或者得不到响应的向渠道发起查询
	case model.TransHandling, "":
		// 获取渠道商户
		_, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
		if err != nil {
			// TODO check error code
			log.Errorf("not found any chanMer(%s,%s): %s", t.ChanCode, t.ChanMerId, err)
			resp.ErrorDetail = "SYSTEM_ERROR"
			return resp
		}
		// 源订单号
		req.SysOrderNum = t.SysOrderNum

		// 向渠道查询
		sp := channel.GetScanPayChan(t.ChanCode)
		resp = sp.ProcessEnquiry(req)

		// 更新交易结果
		t.ChanRespCode = resp.ErrorDetail
		t.RespCode = resp.RespCode
		t.ChanOrderNum = resp.ChannelOrderNum
		t.ChanDiscount = resp.ChcdDiscount
		t.MerDiscount = resp.MerDiscount
		t.ConsumerAccount = resp.ConsumerAccount
		t.ConsumerId = resp.ConsumerId

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

		fallthrough

	// 直接返回
	default:
		resp.Busicd = t.Busicd
		resp.ChannelOrderNum = t.ChanOrderNum
		resp.ConsumerAccount = t.ConsumerAccount
		resp.ConsumerId = t.ConsumerId
		resp.ChcdDiscount = t.ChanDiscount
		resp.MerDiscount = t.MerDiscount
	}

	return resp
}

// Cancel 撤销
func Cancel(req *model.ScanPay) (resp *model.ScanPayResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}
