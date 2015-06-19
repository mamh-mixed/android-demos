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

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复
		return mongo.OffLineRespCd("AUTH_NO_ERROR")
	}

	// 记录该笔交易
	t := &model.Trans{
		MerId:     req.Mchntid,
		OrderNum:  req.OrderNum,
		TransType: model.PayTrans,
		Busicd:    req.Busicd,
		Inscd:     req.Inscd,
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}
	t.TransAmt = int64(f * 100)

	// 渠道选择
	// 根据扫码Id判断走哪个渠道
	if strings.HasPrefix(req.ScanCodeId, "1") {
		req.Chcd = "WXP"
	} else if strings.HasPrefix(req.ScanCodeId, "2") {
		req.Chcd = "ALP"
	} else {
		// 不送，返回 TODO check error code
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}
	t.ChanCode = req.Chcd

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		// TODO check error code
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		// TODO check error code
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}

	// 上送参数
	req.SysOrderNum = tools.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.Key = c.SignCert
	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	if sp == nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}
	ret = sp.ProcessBarcodePay(req)

	// 渠道
	ret.Chcd = req.Chcd

	// 更新交易信息
	updateTrans(t, ret)

	return ret
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复
		return mongo.OffLineRespCd("AUTH_NO_ERROR")
	}

	// 记录该笔交易
	t := &model.Trans{
		MerId:     req.Mchntid,
		OrderNum:  req.OrderNum,
		TransType: model.PayTrans,
		Busicd:    req.Busicd,
		Inscd:     req.Inscd,
		ChanCode:  req.Chcd,
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}
	t.TransAmt = int64(f * 100)

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		// TODO check error code
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		// TODO check error code
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}

	// 上送参数
	req.SysOrderNum = tools.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.Key = c.SignCert
	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	if sp == nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}
	ret = sp.ProcessQrCodeOfflinePay(req)

	// 渠道
	ret.Chcd = req.Chcd

	// 更新交易信息
	updateTrans(t, ret)

	return ret
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

	ret = new(model.ScanPayResponse)
	// 判断是否存在该订单
	t, err := mongo.SpTransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return mongo.OffLineRespCd("TRADE_NOT_EXIST")
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
			return mongo.OffLineRespCd("SYSTEM_ERROR")
		}
		// 原订单号
		req.SysOrderNum = t.SysOrderNum
		req.Key = c.SignCert

		// 向渠道查询
		sp := channel.GetScanPayChan(t.ChanCode)
		ret = sp.ProcessEnquiry(req)

		// 更新交易结果
		updateTrans(t, ret)

	default:

		// 原交易信息
		ret.ErrorDetail = t.ChanRespCode
		ret.ChannelOrderNum = t.ChanOrderNum
		ret.ConsumerAccount = t.ConsumerAccount
		ret.ConsumerId = t.ConsumerId
		ret.ChcdDiscount = t.ChanDiscount
		ret.MerDiscount = t.MerDiscount
		ret.Respcd = t.RespCode
	}

	// 渠道
	ret.Chcd = t.ChanCode
	// 请求业务类型，非原业务类型
	ret.Busicd = req.Busicd
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

// updateTrans 更新交易信息
func updateTrans(t *model.Trans, ret *model.ScanPayResponse) {

	// 根据请求结果更新
	t.ChanRespCode = ret.ChanRespCode
	t.ChanOrderNum = ret.ChannelOrderNum
	t.ChanDiscount = ret.ChcdDiscount
	t.MerDiscount = ret.MerDiscount
	t.ConsumerAccount = ret.ConsumerAccount
	t.ConsumerId = ret.ConsumerId
	t.RespCode = ret.Respcd
	// ...

	// 根据应答码判断交易状态
	switch ret.Respcd {
	case "00":
		t.TransStatus = model.TransSuccess
	case "09":
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
	mongo.SpTransColl.Update(t)
}
