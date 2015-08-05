package core

import (
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// PublicPay 公众号页面支付
func PublicPay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
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

	// TODO: token换取openid
	openId := ""
	t.ConsumerAccount = openId
	req.OpenId = openId

	// TODO: 判断是否取用户信息

	// TODO: 走预下单渠道
	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}

	t.ChanMerId = rp.ChanMerId

	// 是否是代理商模式
	if rp.IsAgent {
		req.SubMchId = rp.SubMerId
		t.SubChanMerId = rp.SubMerId
	}

	// 请求渠道
	ret = adaptor.ProcessQrCodeOfflinePay(t, req)

	// 预支付凭证
	t.PrePayId = ret.PrePayId

	// 更新交易信息
	updateTrans(t, ret)

	// TODO:包装返回值

	return ret
}

// EnterprisePay 企业支付接口
func EnterprisePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}

	// 记录该笔交易
	t := &model.Trans{
		MerId:         req.Mchntid,
		OrderNum:      req.OrderNum,
		TransType:     model.PayTrans,
		Busicd:        req.Busicd,
		Inscd:         req.Inscd,
		Terminalid:    req.Terminalid,
		TransAmt:      req.IntTxamt,
		Remark:        req.Desc,
		GatheringId:   req.OpenId,
		GatheringName: req.UserName,
		ChanCode:      req.Chcd,
	}

	// 渠道是否合法
	switch req.Chcd {
	case channel.ChanCodeWeixin:
		// ok

	default:
		// alipay not support now
		return adaptor.LogicErrorHandler(t, "NO_CHANNEL")
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 是否是代理商模式
	if rp.IsAgent {
		req.SubMchId = rp.SubMerId
		t.SubChanMerId = rp.SubMerId
	}

	ret = adaptor.ProcessEnterprisePay(t, req)

	// 更新交易信息
	updateTrans(t, ret)

	return ret
}

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

	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
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
		Remark:     req.GoodsInfo,
	}

	// 根据扫码Id判断走哪个渠道
	shouldChcd := ""
	switch req.ScanCodeId[0:1] {
	case "1":
		shouldChcd = channel.ChanCodeWeixin
	case "2":
		shouldChcd = channel.ChanCodeAlipay
	default:
		return adaptor.LogicErrorHandler(t, "NO_CHANNEL")
	}

	// 上送渠道与付款码不符
	if req.Chcd != "" && req.Chcd != shouldChcd {
		return adaptor.LogicErrorHandler(t, "CODE_CHAN_NOT_MATCH")
	}
	req.Chcd = shouldChcd
	t.ChanCode = shouldChcd

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 是否是代理商模式
	if rp.IsAgent {
		req.SubMchId = rp.SubMerId
		t.SubChanMerId = rp.SubMerId
	}

	ret = adaptor.ProcessBarcodePay(t, req)

	// 渠道
	ret.Chcd = req.Chcd

	// 更新交易信息
	updateTrans(t, ret)

	return ret
}

// QrCodeOfflinePay 扫二维码预下单
func QrCodeOfflinePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
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
		Remark:     req.GoodsInfo,
		NotifyUrl:  req.NotifyUrl,
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId

	// 是否是代理商模式
	if rp.IsAgent {
		req.SubMchId = rp.SubMerId
		t.SubChanMerId = rp.SubMerId
	}

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
	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
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
		return adaptor.LogicErrorHandler(refund, "TRADE_NOT_EXIST")
	}
	refund.ChanCode = orig.ChanCode

	// 退款只能隔天退
	if strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
		return adaptor.LogicErrorHandler(refund, "REFUND_TIME_ERROR")
	}

	// 是否是支付交易
	if orig.TransType != model.PayTrans {
		return adaptor.LogicErrorHandler(refund, "NOT_PAYTRADE")
	}

	// 交易状态是否正常
	if orig.TransStatus == model.TransFail || orig.TransStatus == model.TransHandling {
		return adaptor.LogicErrorHandler(refund, "NOT_SUCESS_TRADE")
	}

	refundAmt := refund.TransAmt
	// 退款状态是否可退
	switch orig.RefundStatus {
	// 已退款
	case model.TransRefunded:
		return adaptor.LogicErrorHandler(refund, "TRADE_REFUNDED")
	// 部分退款
	case model.TransPartRefunded:
		refunded, err := mongo.SpTransColl.FindTransRefundAmt(req.Mchntid, req.OrigOrderNum)
		if err != nil {
			return adaptor.LogicErrorHandler(refund, "SYSTEM_ERROR")
		}
		refundAmt += refunded
		fallthrough
	default:
		// 金额过大
		if refundAmt > orig.TransAmt {
			return adaptor.LogicErrorHandler(refund, "TRADE_AMT_INCONSISTENT")
		} else if refundAmt == orig.TransAmt {
			orig.RefundStatus = model.TransRefunded
			orig.TransStatus = model.TransClosed
			orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
			orig.ErrorDetail = adaptor.CloseMsg
		} else {
			orig.RefundStatus = model.TransPartRefunded
		}
	}

	ret = adaptor.ProcessRefund(orig, refund, req)

	// 更新原交易状态
	if ret.Respcd == adaptor.SuccessCode {
		mongo.SpTransColl.Update(orig)
	}

	// 更新这笔交易
	updateTrans(refund, ret)

	// 返回真实渠道
	req.Chcd = refund.ChanCode

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
	// 原订单非支付交易
	// if t.TransType != model.PayTrans {
	// 	return returnWithErrorCode("CAN_NOT_QUERY_NOT_PAYTRANS")
	// }

	// 判断订单的状态
	switch t.TransStatus {
	// 如果是处理中或者得不到响应的向渠道发起查询
	case model.TransHandling, "":
		// 支付交易则向渠道查询
		if t.TransType == model.PayTrans {
			ret = adaptor.ProcessEnquiry(t, req)
			// 更新交易结果
			updateTrans(t, ret)
			break
		}

		// TODO 待确定
		// 查看原交易的状态
		orig, err := mongo.SpTransColl.FindOne(t.MerId, t.OrigOrderNum)
		if err != nil {
			return returnWithErrorCode("TRADE_NOT_EXIST")
		}

		ret = &model.ScanPayResponse{}
		switch orig.TransStatus {
		// 撤销、退款、取消成功时，订单状态都是已关闭
		case model.TransClosed:
			ret.Respcd, ret.ErrorDetail = adaptor.SuccessCode, adaptor.SuccessMsg
		// TODO 如果是部分退款呢？
		default:
			ret.Respcd, ret.ErrorDetail = adaptor.FailCode, adaptor.FailMsg
		}

		// 更新
		updateTrans(t, ret)

		fallthrough

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
	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}

	// 记录这笔撤销
	cancel := &model.Trans{
		MerId:        req.Mchntid,
		OrderNum:     req.OrderNum,
		OrigOrderNum: req.OrigOrderNum,
		TransType:    model.CancelTrans,
		Busicd:       req.Busicd,
		Inscd:        req.Inscd,
		Terminalid:   req.Terminalid,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(cancel, "TRADE_NOT_EXIST")
	}
	cancel.ChanCode = orig.ChanCode

	// 撤销只能撤当天交易
	if !strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
		return adaptor.LogicErrorHandler(cancel, "CANCEL_TIME_ERROR")
	}

	// 是否是支付交易
	if orig.TransType != model.PayTrans {
		return adaptor.LogicErrorHandler(cancel, "NOT_PAYTRADE")
	}

	// 存在部分退款交易
	if orig.RefundStatus == model.TransPartRefunded {
		return adaptor.LogicErrorHandler(cancel, "TRADE_REFUNDED")
	}

	// 判断交易状态
	switch orig.TransStatus {
	case model.TransFail:
		return adaptor.LogicErrorHandler(cancel, "FAIL")
	case model.TransClosed:
		return adaptor.LogicErrorHandler(cancel, "ORDER_CLOSED")
	case model.TransHandling:
		return adaptor.LogicErrorHandler(cancel, "INPROCESS")
	default:
		orig.RefundStatus = model.TransRefunded // 撤销，全部退款
		orig.TransStatus = model.TransClosed
		orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
		orig.ErrorDetail = adaptor.CloseMsg
	}

	ret = adaptor.ProcessCancel(orig, cancel, req)

	// 原交易状态更新
	if ret.Respcd == adaptor.SuccessCode {
		mongo.SpTransColl.Update(orig)
	}

	// 更新交易状态
	updateTrans(cancel, ret)

	// 返回真实渠道
	req.Chcd = cancel.ChanCode

	return ret
}

// Close 关闭订单
func Close(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
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
		return adaptor.LogicErrorHandler(closed, "TRADE_NOT_EXIST")
	}
	closed.ChanCode = orig.ChanCode

	// 不支持退款、撤销等其他类型交易
	if orig.TransType != model.PayTrans {
		return adaptor.LogicErrorHandler(closed, "NOT_SUPPORT_TYPE")
	}

	// 存在部分退款交易
	if orig.RefundStatus == model.TransPartRefunded {
		return adaptor.LogicErrorHandler(closed, "TRADE_REFUNDED")
	}

	// 交易已关闭
	if orig.TransStatus == model.TransClosed {
		return adaptor.LogicErrorHandler(closed, "ORDER_CLOSED")
	}

	ret = adaptor.ProcessClose(orig, closed, req)

	// 成功应答
	if ret.Respcd == adaptor.SuccessCode {
		orig.TransStatus = model.TransClosed
		orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
		orig.ErrorDetail = adaptor.CloseMsg
		// 更新原交易信息
		mongo.SpTransColl.Update(orig)
	}

	// 更新交易状态
	updateTrans(closed, ret)

	// 返回真实渠道
	req.Chcd = closed.ChanCode

	return ret
}

// isOrderDuplicate 判断订单号是否重复
func isOrderDuplicate(mchId, orderNum string) (*model.ScanPayResponse, bool) {
	count, err := mongo.SpTransColl.Count(mchId, orderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return returnWithErrorCode("SYSTEM_ERROR"), true
	}
	if count > 0 {
		// 订单号重复
		return returnWithErrorCode("ORDER_DUPLICATE"), true
	}
	return nil, false
}

// returnWithErrorCode 使用错误码直接返回
func returnWithErrorCode(errorCode string) *model.ScanPayResponse {
	spResp := mongo.ScanPayRespCol.Get(errorCode)
	return &model.ScanPayResponse{
		Respcd:      spResp.ISO8583Code,
		ErrorDetail: spResp.ISO8583Msg,
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
	case adaptor.SuccessCode:
		t.TransStatus = model.TransSuccess
	case adaptor.InprocessCode:
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
	mongo.SpTransColl.Update(t)
}
