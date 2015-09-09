package core

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/quickpay/weixin"
	"github.com/omigo/log"
	"math"
	"strings"
	"time"
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
		TradeFrom:  req.TradeFrom,
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

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	// 计算费率 四舍五入
	t.Fee = int64(math.Floor(float64(t.TransAmt)*float64(c.MerFee) + 0.5))
	t.NetFee = t.Fee // 净手续费，会在退款时更新

	// 请求渠道
	ret = adaptor.ProcessQrCodeOfflinePay(t, c, req)

	// 如果下单成功
	if ret.Respcd == adaptor.InprocessCode {

		if req.NeedUserInfo == "YES" {
			// 用token换取用户信息
			userInfo, err := weixin.GetAuthUserInfo(token.AccessToken, token.OpenId)
			if err != nil {
				log.Errorf("unable to get userInfo by accessToken: %s", err)
			}
			jsPayInfo.UserInfo = userInfo
			if userInfo != nil {
				t.NickName = userInfo.Nickname
				t.HeadImgUrl = userInfo.Headimgurl
			}
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
		config.Signature = signWithMD5(config, req.SignKey)
		wxpPay.PaySign = signWithMD5(wxpPay, req.SignKey)
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
		TransType:     model.EnterpriseTrans,
		Busicd:        req.Busicd,
		AgentCode:     req.AgentCode,
		Terminalid:    req.Terminalid,
		TransAmt:      req.IntTxamt,
		Remark:        req.Desc,
		GatheringId:   req.OpenId,
		GatheringName: req.UserName,
		ChanCode:      req.Chcd,
		TradeFrom:     req.TradeFrom,
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

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	ret = adaptor.ProcessEnterprisePay(t, c, req)

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
		GoodsInfo:  req.GoodsInfo,
		TradeFrom:  req.TradeFrom,
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

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}

	// 计算费率 四舍五入
	t.Fee = int64(math.Floor(float64(t.TransAmt)*float64(c.MerFee) + 0.5))
	t.NetFee = t.Fee // 净手续费，会在退款时更新

	ret = adaptor.ProcessBarcodePay(t, c, req)

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
		GoodsInfo:  req.GoodsInfo,
		NotifyUrl:  req.NotifyUrl,
		TradeFrom:  req.TradeFrom,
	}

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

	// 计算费率 四舍五入
	t.Fee = int64(math.Floor(float64(t.TransAmt)*float64(c.MerFee) + 0.5))
	t.NetFee = t.Fee // 净手续费，会在退款时更新

	// 将openId参数设置为空，防止tradeType为JSAPI
	req.OpenId = ""

	// 请求渠道
	ret = adaptor.ProcessQrCodeOfflinePay(t, c, req)

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
		TradeFrom:    req.TradeFrom,
	}

	// 判断是否存在该订单
	orig, err := findAndLockOrigTrans(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(refund, err.Error())
	}

	// 解锁
	defer func() {
		// 如果是逻辑错误等导致原交易没解锁
		if orig.LockFlag == 1 {
			mongo.SpTransColl.Unlock(orig.MerId, orig.OrderNum)
		}
	}()

	copyProperties(refund, orig)
	// refund.SubChanMerId = orig.SubChanMerId

	// TODO 退款只能隔天退，按需求投产后先缓冲一段时间再开启。
	// if strings.HasPrefix(orig.CreateTime, time.Now().Format("2006-01-02")) {
	// 	return adaptor.LogicErrorHandler(refund, "REFUND_TIME_ERROR")
	// }

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
		refundAmt += orig.RefundAmt
		fallthrough
	default:
		// 金额过大
		if refundAmt > orig.TransAmt {
			return adaptor.LogicErrorHandler(refund, "TRADE_AMT_INCONSISTENT")
		} else if refundAmt == orig.TransAmt {
			orig.RefundStatus = model.TransRefunded
			orig.TransStatus = model.TransClosed
			orig.RefundAmt = refundAmt
			// 不更新原订单应答码。在查询时做处理，只更新交易的状态。
			// orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
			// orig.ErrorDetail = adaptor.CloseMsg
			// orig.Fee = 0 // 被全额退款时，手续费清零
		} else {
			orig.RefundStatus = model.TransPartRefunded
			orig.RefundAmt = refundAmt // 这个字段的作用主要是为了方便报表时计算部分退款，为了一致性，撤销，取消接口也都统一加上，虽然并没啥作用
		}
	}

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(refund, "NO_CHANMER")
	}
	// 重新计算手续费
	// if orig.RefundStatus == model.TransPartRefunded {
	// 	orig.Fee = int64(math.Floor(float64(orig.TransAmt-orig.RefundAmt))*float64(c.MerFee) + 0.5)
	// }
	// 退款算退款部分的手续费，出报表时，将原订单的跟退款的相减
	refund.Fee = int64(math.Floor(float64(refund.TransAmt)*float64(c.MerFee) + 0.5))
	orig.NetFee = orig.NetFee - refund.Fee // 重新计算原订单的手续费

	ret = adaptor.ProcessRefund(orig, refund, c, req)

	// 更新原交易状态
	if ret.Respcd == adaptor.SuccessCode {
		mongo.SpTransColl.UpdateAndUnlock(orig)
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
	t, err := findAndLockOrigTrans(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.ReturnWithErrorCode(err.Error())
	}

	// 解锁
	defer func() {
		// 如果是逻辑错误等导致原交易没解锁
		if t.LockFlag == 1 {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
		}
	}()

	// 判断订单的状态
	switch t.TransStatus {
	// 如果是处理中或者得不到响应的向渠道发起查询
	case model.TransHandling:

		// 获取渠道商户
		c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
		if err != nil {
			return adaptor.ReturnWithErrorCode("NO_CHANMER")
		}

		// 支付交易则向渠道查询
		if t.TransType == model.PayTrans {
			ret = adaptor.ProcessEnquiry(t, c, req)
			// 更新交易结果
			updateTrans(t, ret)
			break
		}

		// 来到这一般是系统故障，没有更新成功
		// 先拿到原交易
		orig, err := mongo.SpTransColl.FindOne(t.MerId, t.OrigOrderNum)
		if err != nil {
			return adaptor.ReturnWithErrorCode("TRADE_NOT_EXIST")
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
					ret = adaptor.ProcessWxpRefundQuery(t, c, req)
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
	case model.TransClosed:
		// 订单被关闭时，返回关闭的应答码。
		t.RespCode = adaptor.CloseCode
		t.ErrorDetail = adaptor.CloseMsg
		fallthrough
	// 交易失败
	// case model.TransFail:
	// TODO 系统超时
	// 这里只处理渠道应答错误的情况，如渠道系统超时，应答码为91-外部系统错误
	// if t.RespCode == adaptor.UnKnownCode {

	// }
	// 其他情况跳过
	// fallthrough
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
		TradeFrom:    req.TradeFrom,
	}

	// 判断是否存在该订单
	orig, err := findAndLockOrigTrans(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(cancel, err.Error())
	}

	// 解锁
	defer func() {
		// 如果是逻辑错误等导致原交易没解锁
		if orig.LockFlag == 1 {
			mongo.SpTransColl.Unlock(orig.MerId, orig.OrderNum)
		}
	}()

	copyProperties(cancel, orig)
	// 撤销交易，金额=原交易金额
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
		orig.RefundAmt = orig.TransAmt
		// orig.Fee = 0 // 手续费清零
		// orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
		// orig.ErrorDetail = adaptor.CloseMsg
	}

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(cancel, "NO_CHANMER")
	}

	// 对这笔撤销计算手续费，不然会对应不上，出现多扣少退。
	cancel.Fee = int64(math.Floor(float64(cancel.TransAmt)*float64(c.MerFee) + 0.5))
	orig.NetFee = orig.NetFee - cancel.Fee // 重新计算原订单的手续费

	ret = adaptor.ProcessCancel(orig, cancel, c, req)

	// 原交易状态更新
	if ret.Respcd == adaptor.SuccessCode {
		mongo.SpTransColl.UpdateAndUnlock(orig)
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
		TradeFrom:    req.TradeFrom,
	}

	// 判断是否存在该订单
	orig, err := findAndLockOrigTrans(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return adaptor.LogicErrorHandler(closed, err.Error())
	}

	// 解锁
	defer func() {
		// 如果是逻辑错误等导致原交易没解锁
		if orig.LockFlag == 1 {
			mongo.SpTransColl.Unlock(orig.MerId, orig.OrderNum)
		}
	}()

	copyProperties(closed, orig)

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

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(closed, "NO_CHANMER")
	}

	ret = adaptor.ProcessClose(orig, closed, c, req)

	// 成功应答
	if ret.Respcd == adaptor.SuccessCode {
		// 判断原交易是否成功必须在应答回来后判断
		// 因为在执行关单时，还有可能查询订单以明确状态
		// 原交易成功，那么计算这笔取消的手续费
		if orig.TransStatus == model.TransSuccess {
			// 这样做方便于报表导出计算
			closed.TransAmt = orig.TransAmt
			orig.RefundAmt = orig.TransAmt
			closed.Fee = int64(math.Floor(float64(closed.TransAmt)*float64(c.MerFee) + 0.5))
			orig.NetFee = orig.NetFee - closed.Fee // 重新计算原订单的手续费
		}

		// 更新原交易信息
		orig.TransStatus = model.TransClosed
		mongo.SpTransColl.UpdateAndUnlock(orig)
		// orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
		// orig.ErrorDetail = adaptor.CloseMsg
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
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR"), true
	}
	if count > 0 {
		// 订单号重复
		return adaptor.ReturnWithErrorCode("ORDER_DUPLICATE"), true
	}
	return nil, false
}

// updateTrans 更新交易信息
func updateTrans(t *model.Trans, ret *model.ScanPayResponse) error {
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
	return mongo.SpTransColl.UpdateAndUnlock(t)
}

// findAndLockOrigTrans 查找原交易记录
// 如果找到原交易，那么对原交易加锁。
func findAndLockOrigTrans(merId, orderNum string) (orig *model.Trans, err error) {
	var retry int
	for {

		// 判断是否有此订单
		orig, err = mongo.SpTransColl.FindOne(merId, orderNum)
		if err != nil {
			return nil, errors.New("TRADE_NOT_EXIST")
		}
		retry++
		// 判断订单状态主要是保证下单、预下单新增、修改的事务完整性
		if orig.TransStatus == model.TransNotPay {

			// 最多延迟5s，如果这笔交易还是没处理完，报交易超时
			// TODO:时间待商榷
			if retry == 5 {
				log.Errorf("trans(%s,%s) spent long time to update.", merId, orderNum)
				return nil, errors.New("TRADE_OVERTIME")
			}
			// 说明该笔交易在支付接口时还没完成更新，此时数据是脏数据
			log.Info("find trans sleep 500ms ...")
			time.Sleep(500 * time.Millisecond * time.Duration(retry))
			// 等待500ms，继续循环
			continue
		}

		// 如果此时交易被锁住
		if orig.LockFlag == 1 {
			now := time.Now()
			lockTime, err := time.ParseInLocation("2006-01-02 15:04:05", orig.UpdateTime, time.Local)
			if err != nil {
				log.Errorf("fail to parse time : %s", err)
				return nil, errors.New("SYSTEM_ERROR")
			}
			// TODO:被锁时间 1分钟？
			if now.Sub(lockTime) > 1*time.Minute {
				// 直接返回该原交易
				return orig, err
			}
			// 休眠一段时间
			time.Sleep(200 * time.Millisecond * time.Duration(retry))
		}

		// 锁住交易，并且此时交易是最新的
		orig, err = findAndLockTrans(merId, orderNum)
		if err != nil {
			return nil, err
		}
		break
	}
	return
}

// findAndLockTrans
func findAndLockTrans(merId, orderNum string) (orig *model.Trans, err error) {
	var retry int
	for {
		// 锁住该原订单
		orig, err = mongo.SpTransColl.FindAndLock(merId, orderNum)
		if err != nil {
			// 该订单已被读取
			retry++
			// TODO:时间待商榷
			if retry == 10 {
				log.Errorf("trans(%s,%s) waiting for long to lock.", merId, orderNum)
				return nil, errors.New("TRADE_OVERTIME")
			}
			time.Sleep(500 * time.Millisecond * time.Duration(retry))
			continue
		}
		break
	}
	return
}

// copyProperties 从原交易拷贝属性
func copyProperties(current *model.Trans, orig *model.Trans) {
	current.ChanCode = orig.ChanCode
	current.ChanMerId = orig.ChanMerId
	current.MerName = orig.MerName
	current.AgentName = orig.AgentName
	current.GroupCode = orig.GroupCode
	current.GroupName = orig.GroupName
	current.ShortName = orig.ShortName
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
