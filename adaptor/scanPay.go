package adaptor

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/log"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"strings"
	"time"
)

var agentId = goconf.Config.AlipayScanPay.AgentId

// ProcessEnterprisePay 企业支付
func ProcessEnterprisePay(t *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 上送参数
	req.SysOrderNum = t.SysOrderNum
	req.SignKey = c.SignKey
	req.ChanMerId = c.ChanMerId

	// 交易参数
	t.SysOrderNum = req.SysOrderNum

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
	ret, err := ep.ProcessPay(req)
	if err != nil {
		log.Errorf("process BarcodePay error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessBarcodePay 扫条码下单
func ProcessBarcodePay(t *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 选择送往渠道的商户
	chanMer, subMchId, err := chooseChanMer(c)
	if err != nil {
		log.Errorf("chanMer(%s): %s", c.ChanMerId, err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 上送参数
	req.SysOrderNum = t.SysOrderNum
	req.Subject = req.M.Detail.CommodityName
	req.SignKey = chanMer.SignKey
	req.ChanMerId = chanMer.ChanMerId

	// 不同渠道参数转换
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		if channel.Oversea == c.AreaType {
			req.ExtendParams, err = genOverseaExtendInfo(req.M)
			if err != nil {
				return ReturnWithErrorCode("SYSTEM_ERROR")
			}
		} else {
			req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
			req.ExtendParams = genExtendParams(req.M, chanMer)
		}
	case channel.ChanCodeWeixin:
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = chanMer.WxpAppId
		req.SubMchId = subMchId
		req.GoodsTag = req.M.Detail.GoodsTag
	default:
		req.ActTxamt = req.Txamt
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd, c.AreaType)
	if sp == nil {
		return ReturnWithErrorCode("NO_CHANNEL")
	}
	ret, err = sp.ProcessBarcodePay(req)
	if err != nil {
		log.Errorf("process BarcodePay error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 海外接口没有异步通知，需要设置payTime，默认为createTime
	if channel.Oversea == c.AreaType {
		t.PayTime = t.CreateTime
	}

	return ret
}

// ProcessQrCodeOfflinePay 二维码预下单
func ProcessQrCodeOfflinePay(t *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 预下单赞不支持海外
	if c.AreaType == channel.Oversea {
		return ReturnWithErrorCode("NO_ROUTERPOLICY")
	}

	// 选择送往渠道的商户
	chanMer, subMchId, err := chooseChanMer(c)
	if err != nil {
		log.Errorf("chanMer(%s): %s", c.ChanMerId, err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 不同渠道参数转换
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
		req.ExtendParams = genExtendParams(req.M, chanMer)
	case channel.ChanCodeWeixin:
		req.ActTxamt = fmt.Sprintf("%d", t.TransAmt)
		req.AppID = chanMer.WxpAppId
		req.SubMchId = subMchId
		req.GoodsTag = req.M.Detail.GoodsTag
	default:
		req.ActTxamt = req.Txamt
	}

	// 上送参数
	req.SysOrderNum = t.SysOrderNum
	req.Subject = req.M.Detail.CommodityName
	req.SignKey = chanMer.SignKey
	req.ChanMerId = chanMer.ChanMerId

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd, c.AreaType)
	if sp == nil {
		return ReturnWithErrorCode("NO_CHANNEL")
	}
	ret, err = sp.ProcessQrCodeOfflinePay(req)
	if err != nil {
		log.Errorf("process QrCodeOfflinePay error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessRefund 请求渠道退款，不做逻辑处理
func ProcessRefund(orig *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 选择送往渠道的商户
	chanMer, subMchId, err := chooseChanMer(c)
	if err != nil {
		log.Errorf("chanMer(%s): %s", c.ChanMerId, err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 渠道参数
	req.SignKey = chanMer.SignKey
	req.ChanMerId = chanMer.ChanMerId

	// 不同渠道参数转换
	switch orig.ChanCode {
	case channel.ChanCodeAlipay:
		req.ActTxamt = fmt.Sprintf("%0.2f", float64(req.IntTxamt)/100)
	case channel.ChanCodeWeixin:
		req.AppID = chanMer.WxpAppId
		req.ActTxamt = fmt.Sprintf("%d", req.IntTxamt)
		req.TotalTxamt = fmt.Sprintf("%d", orig.TransAmt)
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(chanMer.HttpCert)
		req.WeixinClientKey = []byte(chanMer.HttpKey)
	default:
		req.ActTxamt = req.Txamt
	}

	// 请求退款
	sp := channel.GetScanPayChan(orig.ChanCode, c.AreaType)
	if sp == nil {
		return ReturnWithErrorCode("NO_CHANNEL")
	}

	// 如果是海外支付宝，做一个退款接口到撤销接口的转换，以抵消手续费
	if c.AreaType == channel.Oversea && orig.ChanCode == channel.ChanCodeAlipay {
		// 全额退款
		if orig.TransAmt == req.IntTxamt {
			// 当天
			if strings.HasPrefix(orig.PayTime, time.Now().Format("2006-01-02")) {
				ret, err = sp.ProcessCancel(req)
				if err != nil {
					log.Errorf("process cancel error:%s", err)
					return ReturnWithErrorCode("SYSTEM_ERROR")
				}
				return ret
			}
		}
	}

	ret, err = sp.ProcessRefund(req)
	if err != nil {
		log.Errorf("process refund error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessEnquiry 查询
func ProcessEnquiry(t *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 选择送往渠道的商户
	chanMer, subMchId, err := chooseChanMer(c)
	if err != nil {
		log.Errorf("chanMer(%s): %s", c.ChanMerId, err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 上送参数
	req.SignKey = chanMer.SignKey
	req.ChanMerId = chanMer.ChanMerId

	// 不同渠道参数转换
	switch t.ChanCode {
	case channel.ChanCodeAlipay:
		// do nothing...
	case channel.ChanCodeWeixin:
		req.AppID = chanMer.WxpAppId
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(chanMer.HttpCert)
		req.WeixinClientKey = []byte(chanMer.HttpKey)
	default:
	}

	// 向渠道查询
	sp := channel.GetScanPayChan(t.ChanCode, c.AreaType)
	if sp == nil {
		return ReturnWithErrorCode("NO_CHANNEL")
	}

	ret, err = sp.ProcessEnquiry(req)
	if err != nil {
		log.Errorf("process enquiry error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
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
func ProcessCancel(orig *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 选择送往渠道的商户
	chanMer, subMchId, err := chooseChanMer(c)
	if err != nil {
		log.Errorf("chanMer(%s): %s", c.ChanMerId, err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 渠道参数
	req.SignKey = chanMer.SignKey
	req.ChanMerId = chanMer.ChanMerId

	// 请求撤销
	sp := channel.GetScanPayChan(orig.ChanCode, c.AreaType)

	switch orig.ChanCode {
	case channel.ChanCodeWeixin:
		// 微信用退款接口
		req.AppID = chanMer.WxpAppId
		req.TotalTxamt = fmt.Sprintf("%d", orig.TransAmt)
		req.ActTxamt = req.TotalTxamt
		req.SubMchId = subMchId
		req.WeixinClientCert = []byte(chanMer.HttpCert)
		req.WeixinClientKey = []byte(chanMer.HttpKey)
		ret, err = sp.ProcessRefund(req)
	case channel.ChanCodeAlipay:
		ret, err = sp.ProcessCancel(req)
	default:
		err = fmt.Errorf("unknown scan pay channel `%s`", orig.ChanCode)
	}

	if err != nil {
		log.Errorf("process cancel error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessClose 关闭订单
func ProcessClose(orig *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 支付交易（下单、预下单）
	switch orig.ChanCode {
	case channel.ChanCodeAlipay:
		// 成功支付的交易标记已退款
		if orig.TransStatus == model.TransSuccess {
			orig.RefundStatus = model.TransRefunded
		}
		// 执行撤销流程
		return ProcessCancel(orig, c, req)

	case channel.ChanCodeWeixin:
		// 下单，微信叫做刷卡支付，即被扫，收银员使用扫码设备读取微信用户刷卡授权码
		if orig.Busicd == model.Purc {

			// 走微信撤销接口
			return ProcessWxpCancel(orig, c, req)
		}
		// 预下单，微信叫做扫码支付，即主扫，统一下单，商户系统先调用该接口在微信支付服务后台生成预支付交易单
		if orig.Busicd == model.Paut {
			// 支付成功，调用退款接口
		Tag:
			switch orig.TransStatus {
			case model.TransSuccess:
				// 预下单全额退款
				req.IntTxamt = orig.TransAmt
				return ProcessRefund(orig, c, req)
			case model.TransHandling:
				// 发起查询请求，确认订单状态
				log.Info("query order status before close ...")
				orderStatus := ProcessEnquiry(orig, c, &model.ScanPayRequest{OrigOrderNum: orig.OrderNum})
				if orderStatus.Respcd == SuccessCode {
					orig.TransStatus = model.TransSuccess
					orig.ChanRespCode = orderStatus.ChanRespCode
					orig.ChanOrderNum = orderStatus.ChannelOrderNum
					orig.ConsumerAccount = orderStatus.ConsumerAccount
					orig.RespCode = SuccessCode
					orig.ErrorDetail = SuccessMsg
					goto Tag
				}
				fallthrough
			default:
				// 直接关单 不用判断时间
				// 这里可能出现重新查询时状态为09，但是这时用户马上支付成功
				// 那么实际上是对已支付的订单进行关单，会报错
				return ProcessWxpClose(orig, c, req)
			}
		}
		return ReturnWithErrorCode("NOT_SUPPORT_TYPE")

	default:
		return ReturnWithErrorCode("NO_CHANNEL")
	}
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
func ProcessWxpRefundQuery(t *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	var err error
	errResp := prepareWxpReqData(t, c, req)
	if errResp != nil {
		return errResp
	}

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin, c.AreaType)
	ret, err = sp.ProcessRefundQuery(req)
	if err != nil {
		log.Errorf("process weixin refundQuery error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessWxpCancel 微信撤销接口
func ProcessWxpCancel(orig *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	var err error
	errResp := prepareWxpReqData(orig, c, req)
	if errResp != nil {
		return errResp
	}

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin, c.AreaType)
	ret, err = sp.ProcessCancel(req)
	if err != nil {
		log.Errorf("process weixin cancel error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

// ProcessWxpClose 微信关闭接口
func ProcessWxpClose(orig *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	var err error
	errResp := prepareWxpReqData(orig, c, req)
	if errResp != nil {
		return errResp
	}

	// 指定微信
	sp := channel.GetScanPayChan(channel.ChanCodeWeixin, c.AreaType)
	ret, err = sp.ProcessClose(req)
	if err != nil {
		log.Errorf("process weixin Close error:%s", err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	return ret
}

func prepareWxpReqData(orig *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 选择送往渠道的商户
	chanMer, subMchId, err := chooseChanMer(c)
	if err != nil {
		log.Errorf("chanMer(%s): %s", c.ChanMerId, err)
		return ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 渠道参数
	req.SignKey = chanMer.SignKey
	req.ChanMerId = chanMer.ChanMerId
	req.AppID = chanMer.WxpAppId
	req.SubMchId = subMchId
	req.WeixinClientCert = []byte(chanMer.HttpCert)
	req.WeixinClientKey = []byte(chanMer.HttpKey)

	return nil
}

// chooseChanMer 选择往渠道送的商户
func chooseChanMer(c *model.ChanMer) (chanMer *model.ChanMer, subMchId string, err error) {
	// 受理商模式
	if c.IsAgentMode {
		if c.AgentMer == nil {
			err = fmt.Errorf("%s", "use agentMode but not supply agentMer,please check.")
			return
		}
		subMchId = c.ChanMerId
		chanMer = c.AgentMer
		return
	}
	chanMer = c
	return
}

func genOverseaExtendInfo(mer model.Merchant) (string, error) {
	if mer.Options == nil {
		return "", fmt.Errorf("%s", "no options params found")
	}
	bytes, _ := json.Marshal(mer.Options)
	return string(bytes), nil
}

func genExtendParams(mer model.Merchant, c *model.ChanMer) string {
	var agentCode = agentId
	if c.AgentCode != "" {
		agentCode = c.AgentCode
	}
	var shopInfo = &struct {
		AGENT_ID   string `json:",omitempty"`
		STORE_ID   string `json:",omitempty"`
		STORE_TYPE string `json:",omitempty"`
		SHOP_ID    string `json:",omitempty"`
	}{agentCode, mer.Detail.ShopID, mer.Detail.ShopType, mer.Detail.BrandNum}

	bytes, _ := json.Marshal(shopInfo)
	return string(bytes)
}

// ProcessPurchaseCoupons 卡券核销
// func ProcessPurchaseCoupons(t *model.Trans, c *model.ChanMer, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
//
// 	// 上送参数
// 	req.SysOrderNum = t.SysOrderNum
// 	// req.Subject = mer.Detail.CommodityName
// 	req.ChanMerId = c.ChanMerId
// 	req.Terminalsn = req.Terminalid
// 	req.Terminalid = c.TerminalId
//
// 	// 获得渠道实例，请求
// 	client := unionlive.DefaultClient
// 	ret, err := client.ProcessPurchaseCoupons(req)
// 	if err != nil {
// 		log.Errorf("process PurchaseCoupons error:%s", err)
// 		return ReturnWithErrorCode("SYSTEM_ERROR")
// 	}
//
// 	return ret
// }
