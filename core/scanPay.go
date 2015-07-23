package core

import (
	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"strings"
	"time"
)

// 使用8583应答
var (
	closeCode, closeMsg         = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
	inprocessCode, inprocessMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("INPROCESS")
	successCode, successMsg     = mongo.ScanPayRespCol.Get8583CodeAndMsg("SUCCESS")
)

// TransQuery 交易查询
func TransQuery(q *model.QueryCondition) (ret *model.QueryCondition) {

	now := time.Now().Format("2006-01-02")
	// 默认当天开始
	if q.StartTime == "" {
		q.StartTime = now + " 00:00:00"
	}
	// 默认当天结束
	if q.EndTime == "" {
		q.EndTime = now + " 23:59:59"
	}

	// mongo统计
	trans, total, err := mongo.SpTransColl.Find(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
	}

	size := len(trans)
	ret = &model.QueryCondition{
		Page:     q.Page,
		Size:     size,
		Total:    total,
		RespCode: "000000",
		RespMsg:  "成功",
		Rec:      trans,
		Count:    size,
	}

	return ret
}

// BarcodePay 条码下单
func BarcodePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复
		return returnWithErrorCode("ORDER_DUPLICATE")
	}

	// 记录该笔交易
	t := &model.Trans{
		MerId:      req.Mchntid,
		OrderNum:   req.OrderNum,
		TransType:  model.PayTrans,
		Busicd:     req.Busicd,
		Inscd:      req.Inscd,
		Terminalid: req.Terminalid,
		TransAmt:   req.IntTxamt,
	}

	// 根据扫码Id判断走哪个渠道
	shouldChcd := ""
	switch req.ScanCodeId[0:1] {
	case "1":
		shouldChcd = "WXP"
	case "2":
		shouldChcd = "ALP"
	default:
		return logicErrorHandler(t, "NO_CHANNEL")
	}

	// 上送渠道与付款码不符
	if req.Chcd != "" && req.Chcd != shouldChcd {
		return logicErrorHandler(t, "CODE_CHAN_NOT_MATCH")
	}
	req.Chcd = shouldChcd
	t.ChanCode = shouldChcd

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return logicErrorHandler(t, "NO_ROUTERPOLICY")
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
func QrCodeOfflinePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复
		return returnWithErrorCode("ORDER_DUPLICATE")
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
		TransAmt:   req.IntTxamt,
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		// TODO check error code
		return logicErrorHandler(t, "NO_ROUTERPOLICY")
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
func Refund(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}
	if count > 0 {
		return returnWithErrorCode("ORDER_DUPLICATE")
	}

	// 记录这笔退款
	refund := &model.Trans{
		MerId:        req.Mchntid,
		OrderNum:     req.OrderNum,
		OrigOrderNum: req.OrigOrderNum,
		TransType:    model.RefundTrans,
		Busicd:       req.Busicd,
		Inscd:        req.Inscd,
		Terminalid:   req.Terminalid,
		TransAmt:     req.IntTxamt,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST")
	}
	refund.ChanCode = orig.ChanCode

	// 退款只能隔天退
	if strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
		return logicErrorHandler(refund, "REFUND_TIME_ERROR")
	}

	// 是否是支付交易
	if orig.TransType != model.PayTrans {
		return logicErrorHandler(refund, "NOT_PAYTRADE")
	}

	// 交易状态是否正常
	if orig.TransStatus != model.TransSuccess {
		return logicErrorHandler(refund, "NOT_SUCESS_TRADE")
	}

	refundAmt := refund.TransAmt
	// 退款状态是否可退
	switch orig.RefundStatus {
	// 已退款
	case model.TransRefunded:
		return logicErrorHandler(refund, "TRADE_REFUNDED")
	// 部分退款
	case model.TransPartRefunded:
		refunded, err := mongo.SpTransColl.FindTransRefundAmt(req.Mchntid, req.OrigOrderNum)
		if err != nil {
			return logicErrorHandler(refund, "SYSTEM_ERROR")
		}
		refundAmt += refunded
		fallthrough
	default:
		// 金额过大
		if refundAmt > orig.TransAmt {
			return logicErrorHandler(refund, "TRADE_AMT_INCONSISTENT")
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
	if ret.Respcd == successCode {
		mongo.SpTransColl.Update(orig)
	}

	// 更新这笔交易
	updateTrans(refund, ret)

	return
}

// Enquiry 查询
func Enquiry(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断是否存在该订单
	t, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return returnWithErrorCode("TRADE_NOT_EXIST")
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
		ret.ChannelOrderNum = t.ChanOrderNum
		ret.ConsumerAccount = t.ConsumerAccount
		ret.ConsumerId = t.ConsumerId
		ret.ChcdDiscount = t.ChanDiscount
		ret.MerDiscount = t.MerDiscount
		ret.Respcd = t.RespCode
		ret.ErrorDetail = t.ErrorDetail
	}

	// 渠道
	ret.Chcd = t.ChanCode
	// 请求业务类型，非原业务类型
	ret.Busicd = req.Busicd
	return ret
}

// Cancel 撤销
func Cancel(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	ret = new(model.ScanPayResponse)
	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}
	if count > 0 {
		return returnWithErrorCode("ORDER_DUPLICATE")
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
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST")
	}
	cancel.ChanCode = orig.ChanCode

	// 撤销只能撤当天交易
	if !strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
		return logicErrorHandler(cancel, "CANCEL_TIME_ERROR")
	}

	// 是否是支付交易
	if orig.TransType != model.PayTrans {
		return logicErrorHandler(cancel, "NOT_PAYTRADE")
	}

	// 存在部分退款交易
	if orig.RefundStatus == model.TransPartRefunded {
		return logicErrorHandler(cancel, "TRADE_HAS_REFUND")
	}

	// 判断交易状态
	switch orig.TransStatus {
	case model.TransFail:
		return logicErrorHandler(cancel, "FAIL")
	case model.TransClosed:
		return logicErrorHandler(cancel, "ORDER_CLOSED")
	case model.TransHandling:
		return logicErrorHandler(cancel, "INPROCESS")
	default:
		orig.RefundStatus = model.TransRefunded // 撤销，全部退款
		orig.TransStatus = model.TransClosed
		orig.RespCode = closeCode // 订单已关闭或取消
		orig.ErrorDetail = closeMsg
	}

	ret = adaptor.ProcessCancel(orig, cancel, req)

	// 原交易状态更新
	if ret.Respcd == successCode {
		mongo.SpTransColl.Update(orig)
	}

	// 更新交易状态
	updateTrans(cancel, ret)

	return ret
}

// Close 关闭订单
func Close(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	count, err := mongo.SpTransColl.Count(req.Mchntid, req.OrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}
	if count > 0 {
		// 订单号重复
		return returnWithErrorCode("ORDER_DUPLICATE")
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
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(closed, "TRADE_NOT_EXIST")
	}
	closed.ChanCode = orig.ChanCode

	// 不支持退款、撤销等其他类型交易
	if orig.TransType != model.PayTrans {
		return logicErrorHandler(closed, "NOT_SUPPORT_TYPE")
	}

	// 存在部分退款交易
	if orig.RefundStatus == model.TransPartRefunded {
		return logicErrorHandler(closed, "TRADE_HAS_REFUND")
	}

	// 交易已关闭
	if orig.TransStatus == model.TransClosed {
		return logicErrorHandler(closed, "ORDER_CLOSED")
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
					return logicErrorHandler(closed, "SYSTERM_ERROR")
				}
				interval := time.Now().Sub(transTime)
				// 超过5分钟
				if interval >= 5*time.Minute {
					ret = adaptor.ProcessWxpClose(orig, closed, req)
				} else {
					// 系统落地,异步执行关单
					ret = &model.ScanPayResponse{
						Respcd:      successCode,
						ErrorDetail: successMsg,
					}
					time.AfterFunc(5*time.Minute-interval, func() {
						adaptor.ProcessWxpClose(orig, closed, req)
					})
				}
			}

		default:
			return logicErrorHandler(closed, "NOT_SUPPORT_TYPE")
		}
	default:
		return logicErrorHandler(closed, "NO_CHANNEL")
	}

	// 成功应答
	if ret.Respcd == successCode {
		orig.TransStatus = model.TransClosed
		orig.RespCode = closeCode // 订单已关闭或取消
		orig.ErrorDetail = closeMsg
		// 更新原交易信息
		mongo.SpTransColl.Update(orig)
	}

	// 更新交易状态
	updateTrans(closed, ret)

	return ret
}

// returnWithErrorCode 使用错误码直接返回
func returnWithErrorCode(errorCode string) *model.ScanPayResponse {
	spResp := mongo.ScanPayRespCol.Get(errorCode)
	return &model.ScanPayResponse{
		Respcd:      spResp.ISO8583Code,
		ErrorDetail: spResp.ISO8583Msg,
	}
}

// logicErrorHandler 逻辑错误处理
func logicErrorHandler(t *model.Trans, errorCode string) *model.ScanPayResponse {

	spResp := mongo.ScanPayRespCol.Get(errorCode)
	// 8583应答
	code, msg := spResp.ISO8583Code, spResp.ISO8583Msg

	// 交易保存
	t.RespCode = code
	t.ErrorDetail = msg
	mongo.SpTransColl.Add(t)

	return &model.ScanPayResponse{
		Respcd:      code,
		ErrorDetail: msg,
	}
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
	t.ErrorDetail = ret.ErrorDetail

	// 根据应答码判断交易状态
	switch ret.Respcd {
	case successCode:
		t.TransStatus = model.TransSuccess
	case inprocessCode:
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
	mongo.SpTransColl.Update(t)
}
