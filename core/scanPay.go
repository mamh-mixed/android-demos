package core

import (
	"encoding/json"
	"fmt"

	"net/url"
	"strconv"
	"strings"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/goconf"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

var alipayNotifyUrl = goconf.Config.AlipayScanPay.NotifyUrl + "/quickpay/back/alp"

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

	// 金额单位转换 txamt:000000000010分
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		ret = mongo.OffLineRespCd("SYSTEM_ERROR")
		t.RespCode = ret.Respcd
		mongo.SpTransColl.Add(t)
		return ret
	}
	t.TransAmt = int64(f)

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

	// 参数处理
	switch rp.ChanCode {
	case "ALP":
		req.ActTxamt = fmt.Sprintf("%0.2f", f/100)
	case "WXP":
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
	default:
		req.ActTxamt = req.Txamt
	}

	// 上送参数
	req.SysOrderNum = tools.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	// req.NotifyUrl = alipayNotifyUrl + "?schema=" + req.SysOrderNum

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
	ret, err = sp.ProcessBarcodePay(req)
	if err != nil {
		log.Errorf("process barcodePay error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

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
	t.TransAmt = int64(f)

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

	// 转换金额单位
	switch rp.ChanCode {
	case "ALP":
		req.ActTxamt = fmt.Sprintf("%0.2f", f/100)
	case "WXP":
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
	default:
		req.ActTxamt = req.Txamt
	}

	// 上送参数
	req.SysOrderNum = tools.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	// req.NotifyUrl = alipayNotifyUrl + "?schema=" + req.SysOrderNum

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
	ret, err = sp.ProcessQrCodeOfflinePay(req)
	if err != nil {
		log.Errorf("process QrCodeOfflinePay error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

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
		MerId:          req.Mchntid,
		OrderNum:       req.OrderNum,
		RefundOrderNum: req.OrigOrderNum,
		TransType:      model.RefundTrans,
		Busicd:         req.Busicd,
		Inscd:          req.Inscd,
		ChanCode:       req.Chcd,
	}

	// 金额单位转换
	f, err := strconv.ParseFloat(req.Txamt, 64)
	if err != nil {
		log.Errorf("转换金额错误，txamt = %s", req.Txamt)
		return logicErrorHandler(refund, "SYSTEM_ERROR")
	}
	refund.TransAmt = int64(f)

	// 判断是否存在该订单
	t, err := mongo.SpTransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST")
	}

	// 是否是支付交易
	if t.TransType != model.PayTrans {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
	}

	// 交易状态是否正常
	if t.TransStatus != model.TransSuccess {
		return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
	}

	var refundStatus int8
	refundAmt := refund.TransAmt
	// 退款状态是否可退
	switch t.RefundStatus {
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
		if refundAmt > t.TransAmt {
			return logicErrorHandler(refund, "TRADE_NOT_EXIST") // TODO check error code
		} else if refundAmt == t.TransAmt {
			refundStatus = model.TransRefunded
		} else {
			refundStatus = model.TransPartRefunded
		}
	}

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return logicErrorHandler(refund, "SYSTEM_ERROR")
	}

	// TODO 这个转换交给各个渠道自己管，如果再加一个渠道，不用改核心逻辑
	// 转换金额单位
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", f/100)
	case channel.ChanCodeWeixin:
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
	default:
		req.ActTxamt = req.Txamt
	}

	// 渠道参数
	req.SysOrderNum = tools.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	// req.NotifyUrl = alipayNotifyUrl + "?schema=" + req.SysOrderNum

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(refund)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 请求退款
	sp := channel.GetScanPayChan(t.ChanCode)
	ret, err = sp.ProcessRefund(req)
	if err != nil {
		log.Errorf("process refund error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 更新交易状态
	if ret.Respcd == "00" {
		t.RefundStatus = refundStatus
		mongo.SpTransColl.Update(t)
	}
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
	case model.TransHandling, model.TransSuccess:
		// 获取渠道商户
		c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
		if err != nil {
			// TODO check error code
			return mongo.OffLineRespCd("SYSTEM_ERROR")
		}
		// 上送参数
		req.AppID = c.WxpAppId    // 微信需要
		req.SubMchId = c.SubMchId // 微信需要
		req.OrderNum = t.OrderNum
		req.SignCert = c.SignCert
		req.ChanMerId = c.ChanMerId

		// 向渠道查询
		sp := channel.GetScanPayChan(t.ChanCode)
		ret, err = sp.ProcessEnquiry(req)
		if err != nil {
			log.Errorf("process enquiry error:%s", err)
			return mongo.OffLineRespCd("SYSTEM_ERROR")
		}

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
		MerId:          req.Mchntid,
		OrderNum:       req.OrderNum,
		RefundOrderNum: req.OrigOrderNum,
		TransType:      model.CancelTrans,
		Busicd:         req.Busicd,
		Inscd:          req.Inscd,
		ChanCode:       req.Chcd,
	}

	// 判断是否存在该订单
	t, err := mongo.SpTransColl.Find(req.Mchntid, req.OrigOrderNum)
	if err != nil {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST")
	}

	// 是否是支付交易
	if t.TransType != model.PayTrans {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST") // TODO check error code
	}

	// 交易状态是否正常
	if t.TransStatus != model.TransSuccess {
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST") // TODO check error code
	}

	var refundStatus int8
	// 判断原交易是否有退款交易
	switch t.RefundStatus {
	// 已被退款
	case model.TransRefunded:
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST") // TODO check error code
	// 存在退款交易
	case model.TransPartRefunded:
		return logicErrorHandler(cancel, "TRADE_NOT_EXIST") // TODO check error code
	default:
		// 可撤销
		refundStatus = model.TransRefunded
	}

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return logicErrorHandler(cancel, "SYSTEM_ERROR")
	}

	// 转换金额单位
	switch t.ChanCode {
	case "ALP":
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
	case "WXP":
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
	}

	// 渠道参数
	req.SysOrderNum = tools.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	// req.NotifyUrl = alipayNotifyUrl + "?schema=" + req.SysOrderNum

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(cancel)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 请求撤销
	sp := channel.GetScanPayChan(t.ChanCode)
	ret, err = sp.ProcessCancel(req)
	if err != nil {
		log.Errorf("process cancel error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 原交易状态更新
	if ret.Respcd == "00" {
		t.RefundStatus = refundStatus
		mongo.SpTransColl.Update(t)
	}

	// 更新交易状态
	updateTrans(cancel, ret)

	return ret
}

// AlpAsyncNotify 支付宝异步通知处理
func AlpAsyncNotify(params url.Values) {

	// 通知动作类型
	notifyAction := params.Get("notify_action_type")
	// 交易订单号
	orderNum := params.Get("out_trade_no")
	// 系统订单号
	sysOrderNum := params.Get("schema")

	// 系统订单号是全局唯一
	t, err := mongo.SpTransColl.FindByOrderNum(sysOrderNum)
	if err != nil {
		log.Errorf("fail to find trans by sysOrderNum=%s", sysOrderNum)
		return
	}

	// 判断是否是原订单
	if t.OrderNum != orderNum {
		log.Errorf("orderNum not match, expect %s, but get %s", t.OrderNum, orderNum)
		return
	}

	switch notifyAction {
	// 退款
	case "refundFPAction":
		// 将优惠信息更新为0.00，貌似为了打单用
		mongo.SpTransColl.UpdateFields(&model.Trans{
			Id:           t.Id,
			MerDiscount:  "0.00",
			ChanDiscount: "0.00",
		})
	// 其他
	default:
		// TODO 是否需要校验
		bills := params.Get("paytools_pay_amount")
		if bills != "" {
			var merDiscount float64
			var arrayBills []map[string]string
			if err := json.Unmarshal([]byte(bills), &arrayBills); err == nil {
				for _, bill := range arrayBills {
					for k, v := range bill {
						if k == "MCOUPON" || k == "MDISCOUNT" {
							f, _ := strconv.ParseFloat(v, 64)
							merDiscount += f
						}
					}
				}
			}
			// 更新指定字段，注意，这里不能全部更新
			// 否则可能会覆盖同步返回的结果
			mongo.SpTransColl.UpdateFields(&model.Trans{
				Id:          t.Id,
				MerDiscount: fmt.Sprintf("%0.2f", merDiscount),
			})
		}
	}

}

// ProcessWeixinNotify 微信异步通知处理
func ProcessWeixinNotify(req *model.WeixinNotifyReq) (resp *model.WeixinNotifyResp) {
	log.Errorf("unimplement method: %#v", req)

	// TODO ...

	return resp
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
