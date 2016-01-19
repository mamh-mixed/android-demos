package core

import (
	"strconv"
	"time"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel/unionlive"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

// PurchaseCoupons 卡券核销
func PurchaseCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isCouponOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}
	// 核销次数不填默认为1
	processVeriTime(req)

	// 如果渠道号为空，默认设置为ULIVE
	processCouponChcd(req)

	// 记录该笔交易
	t := &model.Trans{
		MerId:       req.Mchntid,
		SysOrderNum: util.SerialNumber(),
		OrderNum:    req.OrderNum,
		TransType:   model.PurchaseCoupons,
		Busicd:      req.Busicd,
		AgentCode:   req.AgentCode,
		ChanCode:    req.Chcd,
		Terminalid:  req.Terminalid,
		TradeFrom:   req.TradeFrom,
		CouponsNo:   req.ScanCodeId,
		VeriTime:    req.VeriTime,
		TransAmt:    req.IntTxamt,
	}

	// 补充关联字段
	addRelatedProperties(t, req.M)

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	// 记录交易
	err = mongo.CouTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	submitTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format submitTime err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	req.CreateTime = submitTime.Format("20060102150405")
	req.SysOrderNum = t.SysOrderNum
	req.ChanMerId = c.ChanMerId
	req.Terminalsn = req.Terminalid
	req.Terminalid = c.TerminalId

	// 获得渠道实例，请求
	client := unionlive.DefaultClient
	ret, err = client.ProcessPurchaseCoupons(req)
	if err != nil {
		log.Errorf("process PurchaseCoupons error:%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 更新交易信息
	updateCouponTrans(t, ret)

	return ret
}

// PurchaseActCoupons 刷卡活动券验证
func PurchaseActCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isCouponOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}
	// 核销次数不填默认为1
	processVeriTime(req)

	// 如果渠道号为空，默认设置为ULIVE
	processCouponChcd(req)

	// 记录该笔交易
	t := &model.Trans{
		MerId:       req.Mchntid,
		SysOrderNum: util.SerialNumber(),
		OrderNum:    req.OrderNum,
		// TransType:   model.PurchaseCoupons,
		Busicd:       req.Busicd,
		AgentCode:    req.AgentCode,
		ChanCode:     req.Chcd,
		Terminalid:   req.Terminalid,
		CouponsNo:    req.ScanCodeId,
		VeriTime:     req.VeriTime,
		OrigOrderNum: req.OrigOrderNum,
		Cardbin:      req.Cardbin,
		TransAmt:     req.IntTxamt,
		PayType:      req.PayType,
	}
	// 补充关联字段
	addRelatedProperties(t, req.M)

	// 判断是否存在该订单
	orig, err := mongo.CouTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "TRADE_NOT_EXIST")
	}

	//从原始交易中获取订单号，赋值给该请求的原始订单号字段。
	t.OrigOrderNum = orig.OrderNum

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	// 记录交易
	// t.TransStatus = model.TransNotPay
	err = mongo.CouTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	submitTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format submitTime err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	req.CreateTime = submitTime.Format("20060102150405")
	req.SysOrderNum = t.SysOrderNum
	req.ChanMerId = c.ChanMerId
	req.Terminalsn = req.Terminalid
	req.Terminalid = c.TerminalId
	req.OrigChanOrderNum = orig.ChanOrderNum

	// 获得渠道实例，请求
	client := unionlive.DefaultClient
	ret, err = client.ProcessPurchaseActCoupons(req)
	if err != nil {
		log.Errorf("process PurchaseActCoupons error:%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 更新交易信息
	updateCouponTrans(t, ret)

	return ret
}

// QueryPurchaseCouponsResult 电子券验证结果查询
func QueryPurchaseCouponsResult(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isCouponOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}
	// 核销次数不填默认为1
	processVeriTime(req)

	// 如果渠道号为空，默认设置为ULIVE
	processCouponChcd(req)

	// 记录该笔交易
	t := &model.Trans{
		MerId:       req.Mchntid,
		SysOrderNum: util.SerialNumber(),
		OrderNum:    req.OrderNum,
		// TransType:   model.PurchaseCoupons,
		Busicd:       req.Busicd,
		AgentCode:    req.AgentCode,
		ChanCode:     req.Chcd,
		Terminalid:   req.Terminalid,
		CouponsNo:    req.ScanCodeId,
		VeriTime:     req.VeriTime,
		OrigOrderNum: req.OrigOrderNum,
	}
	// 补充关联字段
	addRelatedProperties(t, req.M)

	// 判断是否存在该订单
	orig, err := mongo.CouTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "TRADE_NOT_EXIST")
	}

	//从原始交易中获取订单号，赋值给该请求的原始订单号字段。
	t.OrigOrderNum = orig.OrderNum

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	// 记录交易
	// t.TransStatus = model.TransNotPay
	err = mongo.CouTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	submitTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format submitTime err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	origSubmitTime, err := time.ParseInLocation("2006-01-02 15:04:05", orig.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format origSubmitTime err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	req.CreateTime = submitTime.Format("20060102150405")
	req.SysOrderNum = t.SysOrderNum
	req.ChanMerId = c.ChanMerId
	req.Terminalsn = req.Terminalid
	req.Terminalid = c.TerminalId
	req.OrigSubmitTime = origSubmitTime.Format("20060102150405")
	req.IntTxamt = orig.TransAmt
	if orig.PayType != "" {
		intPayType, err := strconv.Atoi(orig.PayType)
		if err != nil {
			log.Errorf("format payType to int err,%s", err)
			return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
		}
		req.IntPayType = intPayType
	}

	// 获得渠道实例，请求
	client := unionlive.DefaultClient
	ret, err = client.ProcessQueryPurchaseCouponsResult(req)
	if err != nil {
		log.Errorf("process PurchaseActCoupons error:%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 更新交易信息
	updateCouponTrans(t, ret)

	return ret
}

// UndoPurchaseActCoupons 电子券验证撤销
func UndoPurchaseActCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isCouponOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}
	// 如果渠道号为空，默认设置为ULIVE
	processCouponChcd(req)

	// 记录该笔交易
	t := &model.Trans{
		MerId:       req.Mchntid,
		SysOrderNum: util.SerialNumber(),
		OrderNum:    req.OrderNum,
		// TransType:   model.PurchaseCoupons,
		Busicd:       req.Busicd,
		AgentCode:    req.AgentCode,
		ChanCode:     req.Chcd,
		Terminalid:   req.Terminalid,
		CouponsNo:    req.ScanCodeId,
		OrigOrderNum: req.OrigOrderNum,
	}
	// 补充关联字段
	addRelatedProperties(t, req.M)

	// 判断是否存在该订单
	orig, err := mongo.CouTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "TRADE_NOT_EXIST")
	}

	//从原始交易中获取订单号，赋值给该请求的原始订单号字段。
	t.OrigOrderNum = orig.OrderNum

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	// 记录交易
	// t.TransStatus = model.TransNotPay
	err = mongo.CouTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	submitTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format submitTime err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}
	origSubmitTime, err := time.ParseInLocation("2006-01-02 15:04:05", orig.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format origSubmitTime err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	origVeriTime, err := strconv.Atoi(orig.VeriTime)
	if err != nil {
		log.Errorf("format veriTime to int err,%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	req.CreateTime = submitTime.Format("20060102150405")
	req.SysOrderNum = t.SysOrderNum
	req.ChanMerId = c.ChanMerId
	req.Terminalsn = req.Terminalid
	req.Terminalid = c.TerminalId
	req.OrigSubmitTime = origSubmitTime.Format("20060102150405")
	req.OrigChanOrderNum = orig.ChanOrderNum
	req.OrigVeriTime = origVeriTime

	// 获得渠道实例，请求
	client := unionlive.DefaultClient
	ret, err = client.ProcessUndoPurchaseActCoupons(req)
	if err != nil {
		log.Errorf("process UndoPurchaseActCoupons error:%s", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 更新交易信息
	updateCouponTrans(t, ret)

	return ret
}

// isCouponOrderDuplicate 判断卡券订单号是否重复
func isCouponOrderDuplicate(mchId, orderNum string) (*model.ScanPayResponse, bool) {
	count, err := mongo.CouTransColl.Count(mchId, orderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR"), true
	}
	if count > 0 {
		// 订单号重复
		return adaptor.ReturnWithErrorCode("ORDER_DUPLICATE"), true
	}
	return nil, false
}

// updateTrans 更新卡券交易信息
func updateCouponTrans(t *model.Trans, ret *model.ScanPayResponse) error {
	// 根据请求结果更新
	t.ChanRespCode = ret.ChanRespCode
	t.ChanOrderNum = ret.ChannelOrderNum
	ret.ChannelOrderNum = ""
	t.RespCode = ret.Respcd
	t.ErrorDetail = ret.ErrorDetail
	t.Prodname = ret.CardId
	t.CardInfo = ret.CardInfo
	t.AvailCount = ret.AvailCount
	t.Authcode = ret.Authcode
	t.VoucherType = ret.VoucherType
	t.SaleMinAmount = ret.SaleMinAmount
	t.SaleDiscount = ret.SaleDiscount
	t.ActualPayAmount = ret.ActualPayAmount
	t.OrigRespCode = ret.OrigRespcd
	t.OrigErrorDetail = ret.OrigErrorDetail

	//更新核销状态
	if ret.Respcd == "00" {
		t.WriteoffStatus = model.COUPON_WO_SUCCESS
	} else if ret.Respcd == "09" {
		t.WriteoffStatus = model.COUPON_WO_PROCESS
	} else {
		t.WriteoffStatus = model.COUPON_WO_ERROR
	}

	if ret.ExpDate != "" {
		expDate, err := time.ParseInLocation("20060102", ret.ExpDate, time.Local)
		if err != nil {
			log.Errorf("format ret.ExpDate err,%s", err)
			return err
		}
		t.ExpDate = expDate.Format("2006-01-02")
	} else {
		t.ExpDate = ret.ExpDate
	}

	if ret.ChannelTime != "" {
		ChannelTime, err := time.ParseInLocation("20060102150405", ret.ChannelTime, time.Local)
		if err != nil {
			log.Errorf("format ret.ExpDate err,%s", err)
			return err
		}
		t.ChannelTime = ChannelTime.Format("2006-01-02 15:04:05")
		ret.ChannelTime = ""
	}

	return mongo.CouTransColl.UpdateAndUnlock(t)
}

func processVeriTime(req *model.ScanPayRequest) {
	// 核销次数不填默认为1
	if req.VeriTime == "" {
		req.VeriTime = "1"
	}
}

func processCouponChcd(req *model.ScanPayRequest) {
	// 如果渠道号为空，默认设置为ULIVE
	if req.Chcd == "" {
		req.Chcd = "ULIVE"
	}
}
