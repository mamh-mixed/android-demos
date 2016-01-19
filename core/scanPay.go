package core

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/crontab"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/quickpay/weixin"
	"github.com/CardInfoLink/log"
)

var MsgQueue = make(chan *model.Trans, 1e6)

var (
	closeInterval   = time.Duration(goconf.Config.App.OrderCloseTime)
	refreshInterval = time.Duration(goconf.Config.App.OrderRefreshTime)
)

func init() {
	crontab.RegisterTask(refreshInterval, "refreshOrder", RefreshOrder)
	crontab.RegisterTask(5*time.Minute, "closeOrder", CloseOrder)
}

// PublicPay 公众号页面支付
func PublicPay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	jsPayInfo := &model.PayJson{}

	// 判断订单是否存在
	if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
		return err
	}

	// 记录该笔交易
	t := &model.Trans{
		MerId:       req.Mchntid,
		SysOrderNum: req.ReqId,
		OrderNum:    req.OrderNum,
		TransType:   model.PayTrans,
		Busicd:      req.Busicd,
		AgentCode:   req.AgentCode,
		GoodsInfo:   req.GoodsInfo,
		Terminalid:  req.Terminalid,
		TransAmt:    req.IntTxamt,
		ChanCode:    channel.ChanCodeWeixin,
		VeriCode:    req.VeriCode,
		TradeFrom:   req.TradeFrom,
		NotifyUrl:   req.NotifyUrl,
		Attach:      req.Attach,
		Currency:    req.Currency,
		LockFlag:    1, // 锁住
	}

	// 补充关联字段
	addRelatedProperties(t, req.M)

	// 路由策略
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId
	t.SettRole = rp.SettRole

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}
	t.AppID = c.WxpAppId

	var chanMerId string
	if c.IsAgentMode && c.AgentMer != nil {
		chanMerId = c.AgentMer.ChanMerId
	} else {
		chanMerId = c.ChanMerId
	}

	var authClient weixin.AuthClient
	pa, err := mongo.PulicAccountCol.Get(chanMerId)
	if err != nil {
		authClient = weixin.DefaultClient
	} else {
		authClient = weixin.AuthClient{pa.AppID, pa.AppSecret}
	}

	// 网页授权获取token和openid
	token, err := authClient.GetAuthAccessToken(req.Code)
	if err != nil {
		log.Errorf("get accessToken error: %s", err)
		return adaptor.LogicErrorHandler(t, "AUTH_CODE_ERROR")
	}
	openId := token.OpenId
	req.OpenId = openId

	t.Fee = int64(math.Floor(float64(t.TransAmt)*rp.MerFee + 0.5))
	t.NetFee = t.Fee // 净手续费，会在退款时更新

	// 记录交易
	// t.TransStatus = model.TransNotPay
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 请求渠道
	ret = adaptor.ProcessQrCodeOfflinePay(t, c, req)

	// 如果下单成功
	if ret.Respcd == adaptor.InprocessCode {

		if req.NeedUserInfo == "YES" {
			// 用token换取用户信息
			userInfo, err := authClient.GetAuthUserInfo(token.AccessToken, token.OpenId)
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
	// if err, exist := isOrderDuplicate(req.Mchntid, req.OrderNum); exist {
	// 	return err
	// }

	var isNewReq bool // 可重试
	// 判断是否存在该订单
	t, err := findAndLockOrigTrans(req.Mchntid, req.OrderNum)
	if err != nil {
		// 没找到原订单
		isNewReq = true
	}

	// 解锁
	defer func() {
		// 如果是逻辑错误等导致原交易没解锁
		if t.LockFlag == 1 {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
		}
	}()

	// 记录该笔交易
	if isNewReq {
		t = &model.Trans{
			MerId:         req.Mchntid,
			SysOrderNum:   req.ReqId,
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
			Currency:      req.Currency,
			LockFlag:      1,
		}
		// 补充关联字段
		addRelatedProperties(t, req.M)

	} else {
		// 比较数据是否一致
		if t.Remark != req.Desc || t.GatheringId != req.OpenId || t.GatheringName != req.UserName ||
			t.TransAmt != req.IntTxamt || t.ChanCode != req.Chcd {
			return adaptor.ReturnWithErrorCode("ORDER_DUPLICATE") // 不一致时认为两笔交易
		}

		// 如果之前是成功的
		if t.TransStatus == model.TransSuccess {
			return adaptor.ReturnWithErrorCode("SUCCESS")
		}
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
	t.SettRole = rp.SettRole

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}
	t.AppID = c.WxpAppId

	// 记录交易
	if isNewReq {
		err = mongo.SpTransColl.Add(t)
		if err != nil {
			return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
		}
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
		MerId:          req.Mchntid,
		CouponOrderNum: req.CouponOrderNum, // 优惠券核销后的订单号
		DiscountAmt:    req.IntDiscountAmt, // 卡券优惠金额
		PayType:        req.PayType,        // 卡券指定的支付方式
		SysOrderNum:    req.ReqId,
		OrderNum:       req.OrderNum,
		TransType:      model.PayTrans,
		Busicd:         req.Busicd,
		AgentCode:      req.AgentCode,
		Terminalid:     req.Terminalid,
		TransAmt:       req.IntTxamt,
		GoodsInfo:      req.GoodsInfo,
		TradeFrom:      req.TradeFrom,
		Currency:       req.Currency,
		LockFlag:       1,
	}
	// 补充关联字段
	addRelatedProperties(t, req.M)

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

	// 处理卡券相关逻辑
	errCode := couponLogicProcess(req, shouldChcd)
	if errCode != "" {
		return adaptor.LogicErrorHandler(t, errCode)
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
	t.SettRole = rp.SettRole

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}
	t.AppID = c.WxpAppId // 支付宝时会有作用

	// 计算手续费
	t.Fee = int64(math.Floor(float64(t.TransAmt)*rp.MerFee + 0.5))
	t.NetFee = t.Fee

	// 记录交易
	// t.TransStatus = model.TransNotPay
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	ret = adaptor.ProcessBarcodePay(t, c, req)
	if ret.Respcd == adaptor.InprocessCode {
		go refresh(req, c)
	}

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
		MerId:          req.Mchntid,
		SysOrderNum:    req.ReqId,
		OrderNum:       req.OrderNum,
		TransType:      model.PayTrans,
		Busicd:         req.Busicd,
		AgentCode:      req.AgentCode,
		ChanCode:       req.Chcd,
		Terminalid:     req.Terminalid,
		TransAmt:       req.IntTxamt,
		GoodsInfo:      req.GoodsInfo,
		NotifyUrl:      req.NotifyUrl,
		TradeFrom:      req.TradeFrom,
		Attach:         req.Attach,
		Currency:       req.Currency,
		LockFlag:       1,
		CouponOrderNum: req.CouponOrderNum, // 优惠券核销后的订单号
		DiscountAmt:    req.IntDiscountAmt, // 卡券优惠金额
		PayType:        req.PayType,        // 卡券指定的支付方式
	}
	// 补充关联字段
	addRelatedProperties(t, req.M)

	// 处理卡券相关逻辑
	errCode := couponLogicProcess(req, req.Chcd)
	if errCode != "" {
		return adaptor.LogicErrorHandler(t, errCode)
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(req.Mchntid, req.Chcd)
	if rp == nil {
		return adaptor.LogicErrorHandler(t, "NO_ROUTERPOLICY")
	}
	t.ChanMerId = rp.ChanMerId
	t.SettRole = rp.SettRole

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(t, "NO_CHANMER")
	}
	t.AppID = c.WxpAppId // 支付宝时会有作用

	t.Fee = int64(math.Floor(float64(t.TransAmt)*rp.MerFee + 0.5))
	t.NetFee = t.Fee // 净手续费，会在退款时更新

	// 将openId参数设置为空，防止tradeType为JSAPI
	req.OpenId = ""

	// 记录交易
	// t.TransStatus = model.TransNotPay
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

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
		SysOrderNum:  req.ReqId,
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
		// 补充关联字段
		addRelatedProperties(refund, req.M)
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

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(orig.MerId, orig.ChanCode)
	if rp == nil {
		return adaptor.LogicErrorHandler(refund, "NO_ROUTERPOLICY")
	}

	mer, err := mongo.MerchantColl.Find(orig.MerId)
	if err != nil {
		return adaptor.LogicErrorHandler(refund, "NO_MERCHANT")
	}

	// 退款类型
	nowStr := time.Now().Format("2006-01-02")
	switch mer.RefundType {
	case model.NoLimitRefund:
	case model.CurrentDayRefund:
		// 隔天不能退
		if !strings.HasPrefix(orig.CreateTime, nowStr) {
			return adaptor.LogicErrorHandler(refund, "CANCEL_TIME_ERROR")
		}
	case model.OtherDayRefund:
		// 隔天才能退
		if strings.HasPrefix(orig.CreateTime, nowStr) {
			return adaptor.LogicErrorHandler(refund, "REFUND_TIME_ERROR")
		}
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

	// 退款算退款部分的手续费，出报表时，将原订单的跟退款的相减
	refund.Fee = int64(math.Floor(float64(refund.TransAmt)*rp.MerFee + 0.5))
	orig.NetFee = orig.NetFee - refund.Fee // 重新计算原订单的手续费

	// 记录交易
	err = mongo.SpTransColl.Add(refund)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	ret = adaptor.ProcessRefund(orig, c, req)

	// 更新原交易状态
	if ret.Respcd == adaptor.SuccessCode {
		mongo.SpTransColl.UpdateAndUnlock(orig)
	}

	// 补充支付时间
	ret.PayTime = refund.CreateTime

	// 更新这笔交易
	updateTrans(refund, ret)

	// 返回真实渠道
	req.Chcd = refund.ChanCode
	ret.ChannelOrderNum = orig.ChanOrderNum
	ret.ConsumerAccount = orig.ConsumerAccount

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
		// 先简单处理
		if ret.Respcd == adaptor.SuccessCode {
			ret.ErrorCode = "SUCCESS"
		} else {
			// 其他都为失败，处理中上面已处理
			ret.ErrorCode = "FAIL"
		}
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
		SysOrderNum:  req.ReqId,
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
		// 补充关联字段
		addRelatedProperties(cancel, req.M)
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

	rp := mongo.RouterPolicyColl.Find(orig.MerId, orig.ChanCode)
	if rp == nil {
		return adaptor.LogicErrorHandler(cancel, "NO_ROUTERPOLICY")
	}

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
		return adaptor.LogicErrorHandler(cancel, "NOT_SUCESS_TRADE")
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
	cancel.Fee = int64(math.Floor(float64(cancel.TransAmt)*rp.MerFee + 0.5))
	orig.NetFee = orig.NetFee - cancel.Fee // 重新计算原订单的手续费

	// 记录交易
	err = mongo.SpTransColl.Add(cancel)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	ret = adaptor.ProcessCancel(orig, c, req)

	// 原交易状态更新
	if ret.Respcd == adaptor.SuccessCode {
		mongo.SpTransColl.UpdateAndUnlock(orig)
	}

	// 补充支付时间
	ret.PayTime = cancel.CreateTime

	// 更新交易状态
	updateTrans(cancel, ret)

	// 返回真实渠道
	ret.Chcd = cancel.ChanCode
	ret.ChannelOrderNum = orig.ChanOrderNum
	ret.ConsumerAccount = orig.ConsumerAccount

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
		SysOrderNum:  req.ReqId,
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
		// 补充关联字段
		addRelatedProperties(closed, req.M)
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

	rp := mongo.RouterPolicyColl.Find(orig.MerId, orig.ChanCode)
	if rp == nil {
		return adaptor.LogicErrorHandler(closed, "NO_ROUTERPOLICY")
	}

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return adaptor.LogicErrorHandler(closed, "NO_CHANMER")
	}

	// 保存交易
	err = mongo.SpTransColl.Add(closed)
	if err != nil {
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	ret = adaptor.ProcessClose(orig, c, req)

	// 成功应答
	if ret.Respcd == adaptor.SuccessCode {
		// 判断原交易是否成功必须在应答回来后判断
		// 因为在执行关单时，还有可能查询订单以明确状态
		// 原交易成功，那么计算这笔取消的手续费
		if orig.TransStatus == model.TransSuccess {

			// 这样做方便于报表导出计算
			closed.TransAmt = orig.TransAmt
			orig.RefundAmt = orig.TransAmt
			closed.Fee = int64(math.Floor(float64(closed.TransAmt)*rp.MerFee + 0.5))
			orig.NetFee = orig.NetFee - closed.Fee // 重新计算原订单的手续费
			orig.RefundStatus = model.TransRefunded
		}

		// 更新原交易信息
		orig.TransStatus = model.TransClosed
		mongo.SpTransColl.UpdateAndUnlock(orig)
		// orig.RespCode = adaptor.CloseCode // 订单已关闭或取消
		// orig.ErrorDetail = adaptor.CloseMsg
	}

	// 补充支付时间
	ret.PayTime = closed.CreateTime

	// 更新交易状态
	updateTrans(closed, ret)

	// 返回真实渠道
	ret.Chcd = closed.ChanCode
	ret.ChannelOrderNum = orig.ChanOrderNum
	ret.ConsumerAccount = orig.ConsumerAccount

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

// CloseOrder
func CloseOrder() {
	log.Info("process closeOrder ...")
	ts, err := mongo.SpTransColl.FindHandingTrans(closeInterval)
	if err != nil {
		log.Errorf("find handing trans error: %s", err)
		return
	}

	for _, tran := range ts {
		// 锁住该交易
		t, err := mongo.SpTransColl.FindAndLock(tran.MerId, tran.OrderNum)
		if err != nil {
			continue
		}

		// 可能已被更新
		if t.TransStatus != model.TransHandling {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
			continue
		}

		// 获得渠道商户
		c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
		if err != nil {
			log.Errorf("find chanMer error:%s", err)
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
			continue
		}

		// 记录请求日志的request
		req := &model.ScanPayRequest{
			OrigOrderNum: t.OrderNum,
			Mchntid:      t.MerId,
			Busicd:       model.Inqy,
			ReqId:        t.SysOrderNum, // 以系统订单号为日志号
		}

		var closedResult *model.ScanPayResponse
		ret := adaptor.ProcessEnquiry(t, c, req)

		// 还是处理中
		if ret.Respcd == adaptor.InprocessCode {
			switch t.ChanCode {
			case channel.ChanCodeAlipay:
				if ret.ChanRespCode == "TRADE_NOT_EXIST" {
					// 该订单还没被扫，直接取消
					req.Busicd = model.Canc
					closedResult = adaptor.ProcessClose(t, c, req)
				}
				// 假如该ret.ChanRespCode == "WAIT_BUYER_PAY",那么不处理
			case channel.ChanCodeWeixin:
				req.Busicd = model.Canc
				closedResult = adaptor.ProcessWxpClose(t, c, req)
			}
		}

		// 已被关闭
		if ret.Respcd == adaptor.CloseCode {
			closedResult = &model.ScanPayResponse{Respcd: adaptor.SuccessCode}
		}

		// 对结果处理
		if closedResult != nil {
			if closedResult.Respcd == adaptor.SuccessCode {
				// 关闭成功
				t.TransStatus = model.TransClosed
				t.RefundStatus = model.TransOverTimeClosed
				mongo.SpTransColl.UpdateAndUnlock(t)
				continue
			}
		}

		// 其他情况
		updateTrans(t, ret)
	}
}

// 针对09的交易，系统跟踪查询
// 会阻塞，需要单独运行在goruntine里
func refresh(req *model.ScanPayRequest, c *model.ChanMer) {
	var interval = []time.Duration{3, 10, 60, 180}
	for _, d := range interval {
		// 等一等
		time.Sleep(time.Second * d)
		log.Infof("trace inporcess trans, merId=%s, orderNum=%s", req.Mchntid, req.OrderNum)
		// 重新锁住并刷新交易
		t, err := findAndLockTrans(req.Mchntid, req.OrderNum)
		if err != nil {
			// 其他请求在获取该交易
			log.Warnf("refresh warn: %s", err)
			continue
		}

		// 如果交易已经不是处理中，那么直接返回即可
		if t.TransStatus != model.TransHandling {
			return
		}

		// 查询
		ret := adaptor.ProcessEnquiry(t, c, &model.ScanPayRequest{
			ReqId:        req.ReqId, // 关联ReqId，在交易报文里可以看到
			OrigOrderNum: req.OrderNum,
			Busicd:       model.Inqy,
			Mchntid:      req.Mchntid,
		})
		// 处理
		switch ret.Respcd {
		case adaptor.InprocessCode:
			// 记得解锁
			mongo.SpTransColl.Unlock(req.Mchntid, req.OrderNum)
			continue
		case adaptor.SuccessCode:
			// 进入消息推送队列
			if t.TradeFrom == model.IOS || t.TradeFrom == model.Android {
				MsgQueue <- t
			}
			fallthrough
		default:
			// 更新交易状态
			updateTrans(t, ret)
		}
		break
	}
}

// RefreshOrder 刷新支付交易，针对下单需要输入密码的
func RefreshOrder() {
	log.Info("process refreshOrder ...")
	ts, err := mongo.SpTransColl.FindHandingTrans(0, model.Purc)
	if err != nil {
		log.Errorf("find handing trans error: %s", err)
		return
	}
	for _, tran := range ts {
		// 锁住该交易
		t, err := mongo.SpTransColl.FindAndLock(tran.MerId, tran.OrderNum)
		if err != nil {
			continue
		}

		// 可能已被更新
		if t.TransStatus != model.TransHandling {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
			continue
		}

		// 获得渠道商户
		c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
		if err != nil {
			log.Errorf("find chanMer error:%s", err)
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
			continue
		}

		// 记录请求日志的request
		req := &model.ScanPayRequest{
			OrigOrderNum: t.OrderNum,
			Mchntid:      t.MerId,
			Busicd:       model.Inqy,
			ReqId:        t.SysOrderNum, // 以系统订单号为日志号
		}

		ret := adaptor.ProcessEnquiry(t, c, req)

		// 更新
		updateTrans(t, ret)
	}
}

// findAndLockOrigTrans 查找原交易记录
// 如果找到原交易，那么对原交易加锁。
func findAndLockOrigTrans(merId, orderNum string) (orig *model.Trans, err error) {

	// 判断是否有此订单
	orig, err = mongo.SpTransColl.FindOneInMaster(merId, orderNum)
	if err != nil {
		return nil, errors.New("TRADE_NOT_EXIST")
	}

	// 如果此时交易被锁住
	if orig.LockFlag == 1 {
		// 这时交易可能没有被update
		if orig.UpdateTime != "" {
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
		}
		// 休眠一段时间
		time.Sleep(300 * time.Millisecond)
	}

	// 锁住交易，并且此时交易是最新的
	orig, err = findAndLockTrans(merId, orderNum)
	if err != nil {
		return nil, err
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

// updateTrans 更新交易信息
func updateTrans(t *model.Trans, ret *model.ScanPayResponse) error {
	// 根据请求结果更新
	t.ChanRespCode = ret.ChanRespCode
	t.ChanOrderNum = ret.ChannelOrderNum
	t.ChanDiscount = ret.ChcdDiscount
	t.MerDiscount = ret.MerDiscount
	t.RespCode = ret.Respcd
	t.ErrorDetail = ret.ErrorDetail

	if ret.ConsumerAccount != "" {
		t.ConsumerAccount = ret.ConsumerAccount
	}
	if ret.ConsumerId != "" {
		t.ConsumerId = ret.ConsumerId
	}
	if ret.Rate != "" {
		t.ExchangeRate = ret.Rate
	}
	if ret.PayTime != "" {
		t.PayTime = dateFormat(ret.PayTime)
	}

	// 根据应答码判断交易状态
	switch ret.Respcd {
	case adaptor.SuccessCode:
		t.TransStatus = model.TransSuccess
	case adaptor.InprocessCode:
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
	//将支付信息更新到卡券交易中
	updateScanPayTransToCouponTrans(t)
	return mongo.SpTransColl.UpdateAndUnlock(t)
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
	current.Currency = orig.Currency
	current.SubAgentCode = orig.SubAgentCode
	current.SubAgentName = orig.SubAgentName
	current.ConsumerAccount = orig.ConsumerAccount
	current.SettRole = orig.SettRole
	current.ChanOrderNum = orig.ChanOrderNum
	current.AppID = orig.AppID
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

func dateFormat(payTime string) string {
	if payTime == "" {
		return ""
	}
	switch len(payTime) {
	case 14:
		// 20060102150405
		return payTime[0:4] + "-" + payTime[4:6] + "-" + payTime[6:8] + " " + payTime[8:10] + ":" + payTime[10:12] + ":" + payTime[12:14]
	case 19:
		// 2015-01-02 15:04:05
		return payTime
	default:
		// unknown
		log.Errorf("payTime format error, unknown length,get length=%d, patTime=%s", len(payTime), payTime)
		return ""
	}
}

// 为交易关联商户属性
func addRelatedProperties(current *model.Trans, m model.Merchant) {
	current.MerName = m.Detail.MerName
	current.AgentName = m.AgentName
	current.GroupCode = m.GroupCode
	current.GroupName = m.GroupName
	current.ShortName = m.Detail.ShortName
	current.SubAgentCode = m.SubAgentCode
	current.SubAgentName = m.SubAgentName
}

func updateScanPayTransToCouponTrans(t *model.Trans) {
	// 如果卡券订单号为空，则返回
	if t.CouponOrderNum == "" {
		return
	}
	scanPayCoupon := &model.ScanPayCoupon{
		OrderNum:     t.OrderNum,
		RespCode:     t.RespCode,
		TransAmt:     t.TransAmt,
		TransStatus:  t.TransStatus,
		TransType:    t.TransType,
		ChanCode:     t.ChanCode,
		CreateTime:   t.CreateTime,
		UpdateTime:   t.UpdateTime,
		TradeFrom:    t.TradeFrom,
		PayTime:      t.PayTime,
		Currency:     t.Currency,
		ExchangeRate: t.ExchangeRate,
		DiscountAmt:  t.DiscountAmt,
		PayType:      t.PayType,
		MerId:        t.MerId,
		Busicd:       t.Busicd,
		AgentCode:    t.AgentCode,
		Terminalid:   t.Terminalid,
	}
	update := []interface{}{"scanPayCoupon", scanPayCoupon}
	err := mongo.CouTransColl.UpdateFields(t.MerId, t.CouponOrderNum, update...)
	if err != nil {
		log.Errorf("save scanPayTrans to couponTrans fail,scanPayOrderNum:%s,%s", scanPayCoupon.OrderNum, err)
	}
}

func couponLogicProcess(req *model.ScanPayRequest, chanCode string) string {
	// 如果卡券订单号不为空，则查询出卡券订单
	if req.CouponOrderNum != "" {
		payType := ""
		// 判断是否存在该订单
		couponTrans, err := mongo.CouTransColl.FindOne(req.Mchntid, req.CouponOrderNum)
		if err != nil {
			return "COUPON_TRADE_NOT_EXIST"
		}
		// 卡券核销失败或者卡券已被撤销,则不能再支付
		if couponTrans.RespCode != "00" || couponTrans.TransStatus != model.TransSuccess {
			return "COUPON_VERI_ERROR_OR_CANCEL"
		}
		// 如果该卡券订单已经支付成功,则不能再支付
		if couponTrans.ScanPayCoupon != nil && couponTrans.ScanPayCoupon.RespCode == "00" {
			return "COUPON_ALREADY_PAY"
		}
		if len(couponTrans.VoucherType) > 1 {
			payType = couponTrans.VoucherType[0:1]
		}

		// 实际支付方式与卡券指定支付方式不符，则拒掉交易
		if payType != "" {
			if payType == "2" {
				return "CODE_PAYTYPE_NOT_MATCH"
			} else if payType == "4" && chanCode != channel.ChanCodeWeixin {
				return "CODE_PAYTYPE_NOT_MATCH"
			} else if payType == "5" && chanCode != channel.ChanCodeAlipay {
				return "CODE_PAYTYPE_NOT_MATCH"
			}
		}

	}
	return ""
}
