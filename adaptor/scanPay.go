package adaptor

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

var (
	alipayNotifyUrl = goconf.Config.AlipayScanPay.NotifyUrl + "/qp/back/alipay"
	weixinNotifyUrl = goconf.Config.AlipayScanPay.NotifyUrl + "/qp/back/weixin"
)

// ProcessBarcodePay 扫条码下单
func ProcessBarcodePay(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return logicErrorHandler(t, "NO_CHANMER")
	}

	// 上送参数
	req.SysOrderNum = util.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 不同渠道参数转换
	switch t.ChanCode {
	case "ALP":
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
	case "WXP":
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
	default:
		req.ActTxamt = req.Txamt
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	if sp == nil {
		return mongo.OffLineRespCd("NO_CHANNEL")
	}
	ret, err = sp.ProcessBarcodePay(req)
	if err != nil {
		log.Errorf("process BarcodePay error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	return ret
}

// ProcessQrCodeOfflinePay 二维码预下单
func ProcessQrCodeOfflinePay(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return logicErrorHandler(t, "NO_CHANMER")
	}

	// 不同渠道参数转换
	switch t.ChanCode {
	case "ALP":
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
		req.NotifyUrl = alipayNotifyUrl
	case "WXP":
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
		req.NotifyUrl = weixinNotifyUrl
	default:
		req.ActTxamt = req.Txamt
	}

	// 上送参数
	req.SysOrderNum = util.SerialNumber()
	req.Subject = c.ChanMerName // TODO check
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

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
		return mongo.OffLineRespCd("NO_CHANNEL")
	}
	ret, err = sp.ProcessQrCodeOfflinePay(req)
	if err != nil {
		log.Errorf("process QrCodeOfflinePay error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	return ret
}

// ProcessEnquiry 查询
func ProcessEnquiry(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return mongo.OffLineRespCd("NO_CHANMER")
	}
	// 上送参数

	req.OrderNum = t.OrderNum
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 不同渠道参数转换
	switch t.ChanCode {
	case "ALP":
		// do nothing...
	case "WXP":
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
	default:
	}

	// 向渠道查询
	sp := channel.GetScanPayChan(t.ChanCode)
	if sp == nil {
		return mongo.OffLineRespCd("NO_CHANNEL")
	}

	ret, err = sp.ProcessEnquiry(req)
	if err != nil {
		log.Errorf("process enquiry error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 特殊处理
	// 原交易为支付宝预下单并且返回值为交易不存在时，自动处理为09
	if t.ChanCode == "ALP" && t.Busicd == "paut" && ret.ErrorDetail == "TRADE_NOT_EXIST" {
		inporcess := mongo.OffLineRespCd("INPROCESS")
		ret.Respcd = inporcess.Respcd
		ret.ErrorDetail = inporcess.ErrorDetail
	}

	return ret
}

// ProcessRefund 请求渠道退款，不做逻辑处理
func ProcessRefund(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return logicErrorHandler(current, "NO_CHANMER")
	}

	// 不同渠道参数转换
	switch orig.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(current.TransAmt)/100)
	case channel.ChanCodeWeixin:
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
		req.ActTxamt = fmt.Sprintf("%d", current.TransAmt)
		req.TotalTxamt = fmt.Sprintf("%d", orig.TransAmt)
	default:
		req.ActTxamt = req.Txamt
	}

	// 渠道参数
	req.SysOrderNum = util.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	current.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(current)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 请求退款
	sp := channel.GetScanPayChan(orig.ChanCode)
	if sp == nil {
		return mongo.OffLineRespCd("NO_CHANNEL")
	}

	ret, err = sp.ProcessRefund(req)
	if err != nil {
		log.Errorf("process refund error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	return ret
}

// processCancel 请求渠道撤销，不做逻辑处理
func ProcessCancel(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return logicErrorHandler(current, "NO_CHANMER")
	}

	// 渠道参数
	req.SysOrderNum = util.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	current.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(current)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 请求撤销
	sp := channel.GetScanPayChan(orig.ChanCode)

	switch orig.ChanCode {
	case channel.ChanCodeWeixin:
		// 微信用退款接口
		req.AppID = c.WxpAppId
		req.SubMchId = c.SubMchId
		req.TotalTxamt = fmt.Sprintf("%d", orig.TransAmt)
		req.ActTxamt = req.TotalTxamt
		ret, err = sp.ProcessRefund(req)
	case channel.ChanCodeAlipay:
		ret, err = sp.ProcessCancel(req)
	default:
		err = fmt.Errorf("unknown scan pay channel `%s`", orig.ChanCode)
	}

	if err != nil {
		log.Errorf("process cancel error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	return ret
}

// processWxpClose 微信关闭接口
func ProcessWxpClose(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return logicErrorHandler(current, "NO_CHANMER")
	}

	// 渠道参数
	req.SysOrderNum = util.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	req.AppID = c.WxpAppId
	req.SubMchId = c.SubMchId

	// 系统订单号
	current.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(current)
	if err != nil {
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin)
	ret, err = sp.ProcessClose(req)
	if err != nil {
		log.Errorf("process weixin Close error:%s", err)
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	}

	return ret
}

// logicErrorHandler 逻辑错误处理
func logicErrorHandler(t *model.Trans, errorCode string) *model.ScanPayResponse {
	ret := mongo.OffLineRespCd(errorCode)
	t.RespCode = ret.Respcd
	t.ErrorCode = errorCode
	mongo.SpTransColl.Add(t)
	return ret
}
