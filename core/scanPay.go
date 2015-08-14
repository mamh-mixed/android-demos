package core

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/quickpay/weixin"
	"github.com/omigo/log"
)

// PublicPay 公众号页面支付
func PublicPay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	jsPayInfo := &model.PayJson{}

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
		AgentCode:  req.AgentCode,
		Terminalid: req.Terminalid,
		TransAmt:   req.IntTxamt,
		ChanCode:   channel.ChanCodeWeixin,
		VeriCode:   req.VeriCode,
	}

	// 网页授权获取token和openid
	token, err := weixin.GetAuthAccessToken(req.Code)
	if err != nil {
		log.Errorf("get accessToken error: %s", err)
		return adaptor.LogicErrorHandler(t, "AUTH_CODE_ERROR")
	}
	openId := token.OpenId
	req.OpenId = openId

	// 走预下单渠道
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

	// 如果下单成功
	if ret.Respcd == adaptor.InprocessCode {

		if req.NeedUserInfo == "YES" {
			// 用token换取用户信息
			userInfo, err := weixin.GetAuthUserInfo(token.AccessToken, token.OpenId)
			if err != nil {
				log.Errorf("unable to get userInfo by accessToken: %s", err)
			}
			jsPayInfo.UserInfo = userInfo
		}

		// 包装返回值
		config := &model.JsConfig{
			AppID:     req.AppID,
			NonceStr:  util.Nonce(32),
			Timestamp: util.Millisecond(),
		}
		wxpPay := &model.JsWxpPay{
			AppID:     req.AppID,
			NonceStr:  util.Nonce(32),
			TimeStamp: util.Millisecond(),
			SignType:  "MD5",
			Package:   fmt.Sprintf("prepay_id=%s", ret.PrePayId),
		}

		// 签名
		config.Signature = signWithMD5(config, req.SignCert)
		wxpPay.PaySign = signWithMD5(wxpPay, req.SignCert)
		jsPayInfo.Config = config
		jsPayInfo.WxpPay = wxpPay

		bytes, err := json.Marshal(jsPayInfo)
		if err != nil {
			log.Errorf("marshal jsPayInfo error:%s", err)
		}
		ret.PayJsonStr = string(bytes) // 该这段签名时使用
		ret.PayJson = jsPayInfo
	}

	// 预支付凭证
	t.PrePayId = ret.PrePayId
	ret.ConsumerAccount = openId

	// 更新交易信息
	updateTrans(t, ret)

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
		AgentCode:     req.AgentCode,
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
		AgentCode:  req.AgentCode,
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

	// 下单时忽略渠道，以免误送渠道导致交易失败
	// 上送渠道与付款码不符
	// if req.Chcd != "" && req.Chcd != shouldChcd {
	// 	return adaptor.LogicErrorHandler(t, "CODE_CHAN_NOT_MATCH")
	// }
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
		AgentCode:  req.AgentCode,
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
		AgentCode:    req.AgentCode,
		Terminalid:   req.Terminalid,
		TransAmt:     req.IntTxamt,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(refund, "TRADE_NOT_EXIST")
	}
	refund.ChanCode = orig.ChanCode
	refund.ChanMerId = orig.ChanMerId
	refund.SubChanMerId = orig.SubChanMerId

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
			orig.RefundAmt = refundAmt
		} else {
			orig.RefundStatus = model.TransPartRefunded
			orig.RefundAmt = refundAmt // 这个字段的作用主要是为了方便报表时计算部分退款，位了一致性，撤销，取消接口也都统一加上，虽然并没啥作用
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

		// 来到这一般是系统故障，没有更新成功
		// 先拿到原交易
		orig, err := mongo.SpTransColl.FindOne(t.MerId, t.OrigOrderNum)
		if err != nil {
			return returnWithErrorCode("TRADE_NOT_EXIST")
		}

		// 查看原交易的状态
		ret = &model.ScanPayResponse{}
		switch orig.TransStatus {
		// 撤销、退款、取消成功时，订单状态都是已关闭
		case model.TransClosed:
			ret.Respcd, ret.ErrorDetail = adaptor.SuccessCode, adaptor.SuccessMsg
		// 原交易状态没变
		case model.TransSuccess:
			// 如果是退款
			if t.TransType == model.RefundTrans {
				// 部分退款的标识
				if orig.RefundStatus == model.TransPartRefunded {
					ret.Respcd, ret.ErrorDetail = adaptor.SuccessCode, adaptor.SuccessMsg
					break
				}
				// 如果是退款，且原交易状态没变，这时需要去查询
				// 微信可以查退款
				if t.ChanCode == channel.ChanCodeWeixin {
					req.OrderNum = t.OrderNum
					req.OrigOrderNum = t.OrigOrderNum
					ret = adaptor.ProcessWxpRefundQuery(t, req)
					if ret.Respcd == adaptor.SuccessCode {
						// 更新原交易状态
						// TODO
					}
				}
			}

		default:
			ret.Respcd, ret.ErrorDetail = adaptor.FailCode, adaptor.FailMsg
		}

		// 更新
		updateTrans(t, ret)

	// 交易失败
	case model.TransFail:
		// TODO 系统超时
		// 这里只处理渠道应答错误的情况，如渠道系统超时，应答码为91-外部系统错误
		// if t.RespCode == adaptor.UnKnownCode {

		// }

		// 其他情况跳过
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
		AgentCode:    req.AgentCode,
		Terminalid:   req.Terminalid,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(cancel, "TRADE_NOT_EXIST")
	}
	cancel.ChanCode = orig.ChanCode
	cancel.ChanMerId = orig.ChanMerId
	cancel.SubChanMerId = orig.SubChanMerId
	cancel.TransAmt = orig.TransAmt

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
		orig.RefundAmt = orig.TransAmt
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
		AgentCode:    req.AgentCode,
		Terminalid:   req.Terminalid,
	}

	// 判断是否存在该订单
	orig, err := mongo.SpTransColl.FindOne(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(closed, "TRADE_NOT_EXIST")
	}
	closed.ChanCode = orig.ChanCode
	closed.ChanMerId = orig.ChanMerId
	closed.SubChanMerId = orig.SubChanMerId

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
		// 如果原交易成功
		if orig.TransStatus == model.TransSuccess {
			// 这样做方便于报表导出计算
			closed.TransAmt = orig.TransAmt
			orig.RefundAmt = orig.TransAmt
		}
		// 更新原交易信息
		orig.TransStatus = model.TransClosed
		orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
		orig.ErrorDetail = adaptor.CloseMsg
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

func signWithMD5(s interface{}, key string) string {
	buf, err := util.Query(s)
	if err != nil {
		log.Errorf("gen query params error:%s", err)
		return ""
	}
	sign := buf.String() + "&key=" + key
	md5Bytes := md5.Sum([]byte(sign))
	return fmt.Sprintf("%X", md5Bytes[:])
}
