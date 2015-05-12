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
func BarcodePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	ret = req.Response
	log.Debugf("%+v", ret)
	// 判断订单是否存在
	count, err := mongo.TransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		ret.ErrorDetail = "SYSTEM_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		return ret
	}
	if count > 0 {
		// 没有订单号重复代码
		ret.ErrorDetail = "AUTH_NO_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		return ret
	}

	// 记录该笔交易
	t := &model.Trans{
		MerId:     req.Mchntid,
		OrderNum:  req.OrderNum,
		TransType: model.PayTrans,
		Busicd:    ret.Busicd,
		Inscd:     ret.Inscd,
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		ret.ErrorDetail = "SYSTEM_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		t.RespCode = "000001"
		mongo.TransColl.Add(t)
		return ret
	}
	t.TransAmt = int64(f * 100)

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
		ret.ErrorDetail = "SYSTEM_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		t.RespCode = "000001"
		mongo.TransColl.Add(t)
		return ret
	}
	t.ChanCode = req.Chcd

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, cardBrand)
	if rp == nil {
		// TODO check error code
		ret.ErrorDetail = "SYSTEM_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		t.RespCode = "000001"
		mongo.TransColl.Update(t)
		return ret
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		// TODO check error code
		ret.ErrorDetail = "SYSTEM_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		t.RespCode = "000001"
		mongo.TransColl.Add(t)
		return ret
	}

	// 上送参数
	req.SysOrderNum = tools.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.Key = c.SignCert
	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.TransColl.Add(t)
	if err != nil {
		ret.ErrorDetail = "SYSTEM_ERROR"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		return ret
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	ret = sp.ProcessBarcodePay(req)

	// 渠道
	ret.Chcd = req.Chcd

	// 根据请求结果更新
	t.ChanRespCode = ret.ChanRespCode
	// t.RespCode = ret.RespCode
	t.ChanOrderNum = ret.ChannelOrderNum
	t.ChanDiscount = ret.ChcdDiscount
	t.MerDiscount = ret.MerDiscount
	t.ConsumerAccount = ret.ConsumerAccount
	t.ConsumerId = ret.ConsumerId

	switch ret.Respcd {
	case "00":
		t.TransStatus = model.TransSuccess
		t.RespCode = "000000"
	case "09":
		t.TransStatus = model.TransHandling
		t.RespCode = "000009"
	default:
		t.TransStatus = model.TransFail
		t.RespCode = "100070"
	}
	mongo.TransColl.Update(t)

	return ret
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// TODO 判断订单是否存在

	// TODO 记录该笔交易

	// TODO 判断渠道商户

	// TODO 转换参数

	// TODO 获得渠道实例，请求

	// TODO 根据请求结果更新

	return
}

// Refund 退款
func Refund(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}

// Enquiry 查询
func Enquiry(req *model.ScanPay) (ret *model.ScanPayResponse) {

	ret = req.Response

	// 判断是否存在该订单
	t, err := mongo.TransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		ret.ErrorDetail = "TRADE_NOT_EXIST"
		ret.Respcd = offLineRespCd(ret.ErrorDetail)
		return ret
	}
	log.Debugf("trans:(%+v)", t)

	// 判断订单的状态

	switch t.TransStatus {
	// 如果是处理中或者得不到响应的向渠道发起查询
	case model.TransHandling, "":
		// 获取渠道商户
		c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
		if err != nil {
			// TODO check error code
			ret.ErrorDetail = "SYSTEM_ERROR"
			ret.Respcd = offLineRespCd(ret.ErrorDetail)
			return ret
		}
		// 原订单号
		req.SysOrderNum = t.SysOrderNum
		req.Key = c.SignCert
		// 原订单的性质
		req.Response.Busicd = t.Busicd

		// 向渠道查询
		sp := channel.GetScanPayChan(t.ChanCode)
		ret = sp.ProcessEnquiry(req)

		// 更新交易结果
		t.ChanRespCode = ret.ErrorDetail
		// t.RespCode = ret.RespCode
		t.ChanOrderNum = ret.ChannelOrderNum
		t.ChanDiscount = ret.ChcdDiscount
		t.MerDiscount = ret.MerDiscount
		t.ConsumerAccount = ret.ConsumerAccount
		t.ConsumerId = ret.ConsumerId

		switch ret.Respcd {
		case "00":
			t.TransStatus = model.TransSuccess
			t.RespCode = "000000"
		case "09":
			t.TransStatus = model.TransHandling
			t.RespCode = "000009"
		default:
			t.RespCode = "100070"
			t.TransStatus = model.TransFail
		}
		mongo.TransColl.Update(t)

		fallthrough

	// 直接返回
	default:
		ret.Busicd = t.Busicd
		ret.Chcd = t.ChanCode
		ret.ErrorDetail = t.ChanRespCode
		ret.ChannelOrderNum = t.ChanOrderNum
		ret.ConsumerAccount = t.ConsumerAccount
		ret.ConsumerId = t.ConsumerId
		ret.ChcdDiscount = t.ChanDiscount
		ret.MerDiscount = t.MerDiscount
	}

	return ret
}

// Cancel 撤销
func Cancel(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// TODO 判断是否存在该订单

	// TODO 判断订单的状态

	// TODO 获得渠道商户

	// TODO 请求退款

	// TODO 更新订单状态

	return
}
