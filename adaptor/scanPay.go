package adaptor

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"math"
)

var agentId = goconf.Config.AlipayScanPay.AgentId

// ProcessEnterprisePay 企业支付
func ProcessEnterprisePay(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_CHANMER")
	}

	// 上送参数
	req.SysOrderNum = util.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	// 不同渠道参数转换
	switch t.ChanCode {
	// 目前暂时不支付支付宝
	// case channel.ChanCodeAlipay:
	// 	req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
	case channel.ChanCodeWeixin:
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = c.WxpAppId
		req.WeixinClientCert = []byte(c.HttpCert)
		req.WeixinClientKey = []byte(c.HttpKey)
		// req.SubMchId = c.SubMchId // remark:暂不支持受理商模式
	}

	ep := channel.GetEnterprisePayChan(t.ChanCode)
	ret, err = ep.ProcessPay(req)
	if err != nil {
		log.Errorf("process BarcodePay error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessBarcodePay 扫条码下单
func ProcessBarcodePay(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_CHANMER")
	}

	// 计算费率 四舍五入
	t.Fee = int64(math.Floor(float64(t.TransAmt)*float64(c.MerFee) + 0.5))

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(t, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
	}

	mer, err := mongo.MerchantColl.Find(t.MerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_MERCHANT")
	}

	// 上送参数
	req.SysOrderNum = util.SerialNumber()
	req.Subject = mer.Detail.CommodityName
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	// 不同渠道参数转换
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
		req.ExtendParams = genExtendParams(mer)
	case channel.ChanCodeWeixin:
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = c.WxpAppId
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(c.HttpCert)
		req.WeixinClientKey = []byte(c.HttpKey)
	default:
		req.ActTxamt = req.Txamt
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	if sp == nil {
		return returnWithErrorCode("NO_CHANNEL")
	}
	ret, err = sp.ProcessBarcodePay(req)
	if err != nil {
		log.Errorf("process BarcodePay error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessQrCodeOfflinePay 二维码预下单
func ProcessQrCodeOfflinePay(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_CHANMER")
	}

	// 计算费率 四舍五入
	t.Fee = int64(math.Floor(float64(t.TransAmt)*float64(c.MerFee) + 0.5))

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(t, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
	}

	mer, err := mongo.MerchantColl.Find(t.MerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_MERCHANT")
	}

	// 不同渠道参数转换
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
		req.ExtendParams = genExtendParams(mer)
	case channel.ChanCodeWeixin:
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = c.WxpAppId
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(c.HttpCert)
		req.WeixinClientKey = []byte(c.HttpKey)
	default:
		req.ActTxamt = req.Txamt
	}

	// 上送参数
	req.SysOrderNum = util.SerialNumber()
	req.Subject = mer.Detail.CommodityName
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(t)
	if err != nil {
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	if sp == nil {
		return returnWithErrorCode("NO_CHANNEL")
	}
	ret, err = sp.ProcessQrCodeOfflinePay(req)
	if err != nil {
		log.Errorf("process QrCodeOfflinePay error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessRefund 请求渠道退款，不做逻辑处理
func ProcessRefund(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return LogicErrorHandler(current, "NO_CHANMER")
	}

	// 重新计算手续费
	if orig.RefundStatus == model.TransPartRefunded {
		orig.Fee = int64(math.Floor(float64(orig.TransAmt-orig.RefundAmt))*float64(c.MerFee) + 0.5)
	}

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(current, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
	}

	// 渠道参数
	req.SysOrderNum = util.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 交易参数
	current.SysOrderNum = req.SysOrderNum

	// 不同渠道参数转换
	switch orig.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(current.TransAmt)/100)
	case channel.ChanCodeWeixin:
		req.AppID = c.WxpAppId
		req.ActTxamt = fmt.Sprintf("%d", current.TransAmt)
		req.TotalTxamt = fmt.Sprintf("%d", orig.TransAmt)
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(c.HttpCert)
		req.WeixinClientKey = []byte(c.HttpKey)
	default:
		req.ActTxamt = req.Txamt
	}

	// 记录交易
	err = mongo.SpTransColl.Add(current)
	if err != nil {
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	// 请求退款
	sp := channel.GetScanPayChan(orig.ChanCode)
	if sp == nil {
		return returnWithErrorCode("NO_CHANNEL")
	}

	ret, err = sp.ProcessRefund(req)
	if err != nil {
		log.Errorf("process refund error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessEnquiry 查询
func ProcessEnquiry(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_CHANMER")
	}

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(t, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
	}

	// 上送参数
	req.OrderNum = t.OrderNum
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId

	// 不同渠道参数转换
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		// do nothing...
	case channel.ChanCodeWeixin:
		req.AppID = c.WxpAppId
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(c.HttpCert)
		req.WeixinClientKey = []byte(c.HttpKey)
	default:
	}

	// 向渠道查询
	sp := channel.GetScanPayChan(t.ChanCode)
	if sp == nil {
		return returnWithErrorCode("NO_CHANNEL")
	}

	ret, err = sp.ProcessEnquiry(req)
	if err != nil {
		log.Errorf("process enquiry error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	// 特殊处理
	// 原交易为支付宝预下单并且返回值为交易不存在时，自动处理为09
	// 已在应答码中转换
	// if t.ChanCode == channel.ChanCodeAlipay && t.Busicd == model.Paut && ret.ChanRespCode == "TRADE_NOT_EXIST" {
	// 	inporcess := mongo.ScanPayRespCol.Get("INPROCESS")
	// 	ret.Respcd = inporcess.ISO8583Code
	// 	ret.ErrorDetail = inporcess.ISO8583Msg
	// }

	return ret
}

// ProcessCancel 请求渠道撤销，不做逻辑处理
func ProcessCancel(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return LogicErrorHandler(current, "NO_CHANMER")
	}

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(current, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
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
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	// 请求撤销
	sp := channel.GetScanPayChan(orig.ChanCode)

	switch orig.ChanCode {
	case channel.ChanCodeWeixin:
		// 微信用退款接口
		req.AppID = c.WxpAppId
		req.TotalTxamt = fmt.Sprintf("%d", orig.TransAmt)
		req.ActTxamt = req.TotalTxamt
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(c.HttpCert)
		req.WeixinClientKey = []byte(c.HttpKey)
		ret, err = sp.ProcessRefund(req)
	case channel.ChanCodeAlipay:
		ret, err = sp.ProcessCancel(req)
	default:
		err = fmt.Errorf("unknown scan pay channel `%s`", orig.ChanCode)
	}

	if err != nil {
		log.Errorf("process cancel error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessClose 关闭订单
func ProcessClose(orig, closed *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 支付交易（下单、预下单）
	if orig.ChanCode == channel.ChanCodeAlipay {
		// 成功支付的交易标记已退款
		if orig.TransStatus == model.TransSuccess {
			orig.RefundStatus = model.TransRefunded
		}
		// 执行撤销流程
		return ProcessCancel(orig, closed, req)
	}

	if orig.ChanCode == channel.ChanCodeWeixin {
		// 下单，微信叫做刷卡支付，即被扫，收银员使用扫码设备读取微信用户刷卡授权码
		if orig.Busicd == model.Purc {

			// 走微信撤销接口
			return ProcessWxpCancel(orig, closed, req)
		}

		// 预下单，微信叫做扫码支付，即主扫，统一下单，商户系统先调用该接口在微信支付服务后台生成预支付交易单
		if orig.Busicd == model.Paut {
			// 支付成功，调用退款接口
		Tag:
			switch orig.TransStatus {
			case model.TransSuccess:
				// 预下单全额退款
				closed.TransAmt = orig.TransAmt
				orig.RefundStatus = model.TransRefunded
				return ProcessRefund(orig, closed, req)
			case model.TransHandling:
				// 发起查询请求，确认订单状态
				log.Info("query order status before close ...")
				orderStatus := ProcessEnquiry(orig, &model.ScanPayRequest{OrderNum: orig.OrderNum})
				if orderStatus.Respcd == SuccessCode {
					orig.TransStatus = model.TransSuccess
					goto Tag
				}
				fallthrough
			default:
				// 直接关单 不用判断时间
				return ProcessWxpClose(orig, closed, req)
			}
		}
		return LogicErrorHandler(closed, "NOT_SUPPORT_TYPE")
	}
	return LogicErrorHandler(closed, "NO_CHANNEL")
}

// func weixinCloseOrder(orig, closed *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
// 	// 以下情况需要调用关单接口：
// 	// 商户订单支付失败需要生成新单号重新发起支付，要对原订单号调用关单，避免重复支付；
// 	// 系统下单后，用户支付超时，系统退出不再受理，避免用户继续，请调用关单接口。
// 	// 注意：订单生成后不能马上调用关单接口，最短调用时间间隔为5分钟。
// 	transTime, err := time.ParseInLocation("2006-01-02 15:04:05", orig.CreateTime, time.Local)
// 	if err != nil {
// 		log.Errorf("parse time error: creatTime=%s, mchntid=%s, origOrderNum=%s",
// 			orig.CreateTime, req.Mchntid, req.OrigOrderNum)
// 		return LogicErrorHandler(closed, "SYSTERM_ERROR")
// 	}

// 	interval := time.Now().Sub(transTime)
// 	// 超过5分钟
// 	if interval >= 5*time.Minute {
// 		return ProcessWxpClose(orig, closed, req)
// 	}

// 	// 系统落地，异步执行关单
// 	time.AfterFunc(5*time.Minute-interval, func() {
// 		ProcessWxpClose(orig, closed, req)
// 	})

// 	// TODO 直接返回 ？？？
// 	return &model.ScanPayResponse{
// 		Respcd:      SuccessCode,
// 		ErrorDetail: SuccessMsg,
// 	}
// }

// ProcessWxpRefundQuery 微信查询退款接口
func ProcessWxpRefundQuery(t *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 获取渠道商户
	c, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		return LogicErrorHandler(t, "NO_CHANMER")
	}

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(t, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
	}

	// 上送参数
	req.OrderNum = t.OrderNum
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	req.AppID = c.WxpAppId
	req.SubMchId = subMchId
	req.WeixinClientCert = []byte(c.HttpCert)
	req.WeixinClientKey = []byte(c.HttpKey)

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin)
	ret, err = sp.ProcessRefundQuery(req)
	if err != nil {
		log.Errorf("process weixin refundQuery error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessWxpCancel 微信撤销接口
func ProcessWxpCancel(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	var err error
	errResp := prepareWxpReqData(orig, current, req)
	if errResp != nil {
		return errResp
	}

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin)
	ret, err = sp.ProcessCancel(req)
	if err != nil {
		log.Errorf("process weixin cancel error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessWxpClose 微信关闭接口
func ProcessWxpClose(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	var err error
	errResp := prepareWxpReqData(orig, current, req)
	if errResp != nil {
		return errResp
	}

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin)
	ret, err = sp.ProcessClose(req)
	if err != nil {
		log.Errorf("process weixin Close error:%s", err)
		return returnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

func prepareWxpReqData(orig, current *model.Trans, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 获得渠道商户
	c, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		return LogicErrorHandler(current, "NO_CHANMER")
	}

	var subMchId string
	// 代理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			log.Error("use agentMode but not supply agentMer,please check.")
			return LogicErrorHandler(current, "SYSTEM_ERROR")
		}
		subMchId = c.ChanMerId
		c = c.AgentMer
	}

	// 渠道参数
	req.SysOrderNum = util.SerialNumber()
	req.SignCert = c.SignCert
	req.ChanMerId = c.ChanMerId
	req.AppID = c.WxpAppId
	req.SubMchId = subMchId
	req.WeixinClientCert = []byte(c.HttpCert)
	req.WeixinClientKey = []byte(c.HttpKey)

	// 系统订单号
	current.SysOrderNum = req.SysOrderNum

	// 记录交易
	err = mongo.SpTransColl.Add(current)
	if err != nil {
		return returnWithErrorCode("SYSTEM_ERROR")
	}
	return nil
}

func genExtendParams(mer *model.Merchant) string {
	var shopInfo = &struct {
		AGENT_ID   string `json:",omitempty"`
		STORE_ID   string `json:",omitempty"`
		STORE_TYPE string `json:",omitempty"`
		SHOP_ID    string `json:",omitempty"`
	}{agentId, mer.Detail.ShopID, mer.Detail.ShopType, mer.Detail.BrandNum}
	bytes, _ := json.Marshal(shopInfo)
	return string(bytes)
}
