package core

import (
	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"strconv"
	"strings"
	"time"
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
		MerId:      req.Mchntid,
		OrderNum:   req.OrderNum,
		TransType:  model.PayTrans,
		Busicd:     req.Busicd,
		Inscd:      req.Inscd,
		Terminalid: req.Terminalid,
	}

	// 金额单位转换 txamt:000000000010分
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		return logicErrorHandler(t, "SYSTEM_ERROR")
	}
	t.TransAmt = int64(f)

	// 根据扫码Id判断走哪个渠道
	if req.Chcd == "" {
		if strings.HasPrefix(req.ScanCodeId, "1") {
			req.Chcd = "WXP"
		} else if strings.HasPrefix(req.ScanCodeId, "2") {
			req.Chcd = "ALP"
		} else {
			// 不送，返回 TODO check error code
			return logicErrorHandler(t, "SYSTEM_ERROR")
		}
	}
	t.ChanCode = req.Chcd

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		// TODO check error code
		return logicErrorHandler(t, "SYSTEM_ERROR")
	}
	t.ChanMerId = rp.ChanMerId

	ret = adaptor.ProcessBarcodePay(t, req)

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
		MerId:      req.Mchntid,
		OrderNum:   req.OrderNum,
		TransType:  model.PayTrans,
		Busicd:     req.Busicd,
		Inscd:      req.Inscd,
		ChanCode:   req.Chcd,
		Terminalid: req.Terminalid,
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		return logicErrorHandler(t, "SYSTEM_ERROR")
	}
	t.TransAmt = int64(f)

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		// TODO check error code
		return logicErrorHandler(t, "SYSTEM_ERROR")
	}
	t.ChanMerId = rp.ChanMerId

	// 请求渠道
	ret = adaptor.ProcessQrCodeOfflinePay(t, req)

	// 二维码
	t.QrCode = ret.QrCode

	// 更新交易信息
	updateTrans(t, ret)

	return ret
}

// Refund 退款
func Refund(req *model.ScanPay) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复 check error code
		return mongo.OffLineRespCd("AUTH_NO_ERROR")
	}

	// 记录这笔退款
	refund := &model.Trans{
		MerId:        req.Mchntid,
		OrderNum:     req.OrderNum,
		OrigOrderNum: req.OrigOrderNum,
		TransType:    model.RefundTrans,
		Busicd:       req.Busicd,
		Inscd:        req.Inscd,
		ChanCode:     req.Chcd,
		Terminalid:   req.Terminalid,
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		log.Errorf("转换金额错误，txamt = %s", req.Txamt)
		return logicErrorHandler(refund, "SYSTEM_ERROR")
	}
	refund.TransAmt = int64(f)

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST")
	}
	refund.ChanCode = orig.ChanCode

	// 退款只能隔天退
	if strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
		return logicErrorHandler(refund, "REFUND_SHOULD_BE_NEXT_DAY") // TODO check error code
	}

	// 是否是支付交易
	if orig.TransType != model.PayTrans {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
	}

	// 交易状态是否正常
	if orig.TransStatus != model.TransSuccess {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
	}

	refundAmt := refund.TransAmt
	// 退款状态是否可退
	switch orig.RefundStatus {
	// 已退款
	case model.TransRefunded:
		return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
	// 部分退款
	case model.TransPartRefunded:
		refunded, err := mongo.SpTransColl.FindTransRefundAmt(req.Mchntid, req.OrigOrderNum)
		if err != nil {
			return logicErrorHandler(refund, "SYSTEM_ERROR") // TODO check error code
		}
		refundAmt += refunded
		fallthrough
	default:
		// 金额过大
		if refundAmt > orig.TransAmt {
			return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
		} else if refundAmt == orig.TransAmt {
			orig.RefundStatus = model.TransRefunded
			orig.TransStatus = model.TransClosed
			orig.RespCode = "54" // 订单已关闭或取消
		} else {
			orig.RefundStatus = model.TransPartRefunded
		}
	}

	ret = adaptor.ProcessRefund(orig, refund, req)

	// 更新原交易状态
	if ret.Respcd == "00" {
		mongo.SpTransColl.Update(orig)
	}

	// 更新这笔交易
	updateTrans(refund, ret)

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

		ret = adaptor.ProcessEnquiry(t, req)
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

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复 check error code
		return mongo.OffLineRespCd("AUTH_NO_ERROR")
	}

	// 记录这笔撤销
	cancel := &model.Trans{
		MerId:        req.Mchntid,
		OrderNum:     req.OrderNum,
		OrigOrderNum: req.OrigOrderNum,
		TransType:    model.CancelTrans,
		Busicd:       req.Busicd,
		Inscd:        req.Inscd,
		ChanCode:     req.Chcd,
		Terminalid:   req.Terminalid,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST")
	}
	cancel.ChanCode = orig.ChanCode

	// 撤销只能撤当天交易
	if !strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
		return logicErrorHandler(cancel, "NOT_CURRENT_DAY_TRAN") // TODO check error code
	}

	// 是否是支付交易
	if orig.TransType != model.PayTrans {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST") // TODO check error code
	}

	// 存在部分退款交易
	if orig.RefundStatus == model.TransPartRefunded {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST") // TODO check error code
	}

	// 判断交易状态
	switch orig.TransStatus {
	case model.TransFail:
		return logicErrorHandler(cancel, "TRADE_FAIL") // TODO check error code
	case model.TransClosed:
		return logicErrorHandler(cancel, "TRADE_CLOSED") // TODO check error code
	case model.TransHandling:
		return logicErrorHandler(cancel, "TRADE_HANDLING") // TODO check error code
	default:
		orig.RefundStatus = model.TransRefunded // 撤销，全部退款
		orig.TransStatus = model.TransClosed
		orig.RespCode = "54" // 订单已关闭或取消
	}

	ret = adaptor.ProcessCancel(orig, cancel, req)

	// 原交易状态更新
	if ret.Respcd == "00" {
		mongo.SpTransColl.Update(orig)
	}

	// 更新交易状态
	updateTrans(cancel, ret)

	return ret
}

// Close 关闭订单
func Close(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR") // todo check error code
	}
	if count > 0 {
		// 订单号重复
		return mongo.OffLineRespCd("AUTH_NO_ERROR") // todo check error code
	}

	// 记录这笔关单
	closed := &model.Trans{
		MerId:        req.Mchntid,
		OrderNum:     req.OrderNum,
		OrigOrderNum: req.OrigOrderNum,
		TransType:    model.CloseTrans,
		Busicd:       req.Busicd,
		Inscd:        req.Inscd,
		Terminalid:   req.Terminalid,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(closed, "ORER_NUM_ERROR")
	}
	closed.ChanCode = orig.ChanCode

	// 不支持退款、撤销等其他类型交易
	if orig.TransType != model.PayTrans {
		return logicErrorHandler(closed, "ORER_NUM_ERROR")
	}

	// 存在部分退款交易
	if orig.RefundStatus == model.TransPartRefunded {
		return logicErrorHandler(closed, "TRADE_HAS_REFUND") // TODO check error code
	}

	// 交易已关闭
	if orig.TransStatus == model.TransClosed {
		return logicErrorHandler(closed, "TRADE_HAS_CLOSED") // TODO check error code
	}

	// 支付交易（下单、预下单）
	switch orig.ChanCode {
	case channel.ChanCodeAlipay:
		// 成功支付的交易标记已退款
		if orig.TransStatus == model.TransSuccess {
			orig.RefundStatus = model.TransRefunded
		}
		// 执行撤销流程
		ret = adaptor.ProcessCancel(orig, closed, req)
	case channel.ChanCodeWeixin:
		switch orig.Busicd {
		// 下单
		case model.Purc:
			ret = adaptor.ProcessCancel(orig, closed, req)
		// 预下单
		case model.Paut:

			if orig.TransStatus == model.TransSuccess {
				// 预下单全额退款
				closed.TransAmt = orig.TransAmt
				orig.RefundStatus = model.TransRefunded
				ret = adaptor.ProcessRefund(orig, closed, req)
			} else {
				transTime, err := time.ParseInLocation("2006-01-02 15:04:05", orig.CreateTime, time.Local)
				if err != nil {
					log.Errorf("parse time error, creatTime=%s, origOrderNum=%s", orig.CreateTime, req.OrigOrderNum)
					return logicErrorHandler(closed, "SYSTERM_ERROR") // TODO check error code
				}
				interval := time.Now().Sub(transTime)
				// 超过5分钟
				if interval >= 5*time.Minute {
					ret = adaptor.ProcessWxpClose(orig, closed, req)
				} else {
					// 系统落地,异步执行关单
					ret = mongo.OffLineRespCd("SUCCESS")
					time.AfterFunc(5*time.Minute-interval, func() {
						adaptor.ProcessWxpClose(orig, closed, req)
					})
				}
			}

		default:
			return logicErrorHandler(closed, "ORER_NUM_ERROR")
		}
	default:
		return logicErrorHandler(closed, "ORER_NUM_ERROR")
	}

	// 成功应答
	if ret.Respcd == "00" {
		orig.TransStatus = model.TransClosed
		orig.RespCode = "54" // 订单已关闭或取消
		// 更新原交易信息
		mongo.SpTransColl.Update(orig)
	}

	// 更新交易状态
	updateTrans(closed, ret)

	return ret
}

// logicErrorHandler 逻辑错误处理
func logicErrorHandler(t *model.Trans, errorDetail string) *model.ScanPayResponse {
	ret := mongo.OffLineRespCd(errorDetail)
	t.RespCode = ret.Respcd
	mongo.SpTransColl.Add(t)
	return ret
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
