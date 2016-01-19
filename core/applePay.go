package core

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

// ProcessApplyPay 执行apple Pay相关的支付
// TODO 目前只处理消费类请求，预授权类请求暂时返回无权限
func ProcessApplePay(ap *model.ApplePay) (ret *model.BindingReturn) {
	log.Tracef("before process: %+v", ap)

	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")

	// 系统唯一的序列号
	sysOrderNum := mongo.SnColl.GetSysSN()

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(ap.MerId, ap.MerOrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return
	}
	if count > 0 {
		return mongo.RespCodeColl.Get("200081")
	}

	// 只要订单号不重复就记录这笔交易
	errorTrans := &model.Trans{
		MerId:       ap.MerId,
		OrderNum:    ap.MerOrderNum,
		SysOrderNum: sysOrderNum,
		TransAmt:    ap.ApplePayData.TransactionAmount,
		TransCurr:   ap.ApplePayData.CurrencyCode,
		TransType:   1, //TODO 预授权不属于支付
		RespCode:    "000001",
		Remark:      "Apple Pay",
		SubMerId:    ap.SubMerId,
	}

	// 暂时不支持预授权交易
	if ap.TransType == "AUTH" {
		errorTrans.RespCode = "100030"
		saveErrorTran(errorTrans)
		return mongo.RespCodeColl.Get("100030")
	}

	// 获取卡bin详情
	cardBin, err := findCardBin(ap.ApplePayData.ApplicationPrimaryAccountNumber)
	if err != nil {
		if err.Error() == "not found" {
			errorTrans.RespCode = "200070"
			saveErrorTran(errorTrans)
			return mongo.RespCodeColl.Get("200110")
		}
		saveErrorTran(errorTrans)
		return
	}
	log.Debugf("CardBin: %+v", cardBin)

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(ap.MerId, cardBin.CardBrand)
	if rp == nil {
		errorTrans.RespCode = "300030"
		saveErrorTran(errorTrans)
		return mongo.RespCodeColl.Get("300030")
	}

	// 根据绑定关系得到渠道商户信息
	chanMer, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		errorTrans.RespCode = "300030"
		saveErrorTran(errorTrans)
		log.Errorf("not found any chanMer (%s)", err)
		return mongo.RespCodeColl.Get("300030")
	}

	// 下游送来的终端号，如果没有的话，填上渠道商户里面的配置的终端号
	if ap.TerminalId == "" {
		ap.TerminalId = chanMer.TerminalId
	}

	// 查找配置的渠道入口
	c := channel.GetDirectPayChan(chanMer.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		errorTrans.RespCode = "510010"
		saveErrorTran(errorTrans)
		return mongo.RespCodeColl.Get("510010")
	}

	// 补充信息
	ap.Chcd = chanMer.InsCode
	ap.Mchntid = chanMer.ChanMerId
	ap.CliSN = mongo.SnColl.GetDaySN(chanMer.ChanMerId, chanMer.TerminalId)
	ap.SysSN = sysOrderNum

	// 记录这笔交易，入库
	trans := &model.Trans{
		MerId:       ap.MerId,
		OrderNum:    ap.MerOrderNum,
		SysOrderNum: ap.SysSN,
		AcctNum:     ap.ApplePayData.ApplicationPrimaryAccountNumber,
		TransAmt:    ap.ApplePayData.TransactionAmount,
		ChanMerId:   chanMer.ChanMerId,
		ChanCode:    chanMer.ChanCode,
		TransType:   1, //TODO 预授权不属于支付
		Remark:      "Apple Pay",
		SubMerId:    ap.SubMerId,
	}
	if err := mongo.TransColl.Add(trans); err != nil {
		log.Errorf("add trans fail: (%s)", err)
		return
	}

	// Apple Pay 支付
	ret = c.ConsumeByApplePay(ap)

	trans.ChanRespCode = ret.ChanRespCode
	trans.RespCode = ret.RespCode
	switch ret.RespCode {
	case "000000":
		trans.TransStatus = model.TransSuccess
	case "000009":
		trans.TransStatus = model.TransHandling
	default:
		trans.TransStatus = model.TransFail
	}
	if err = mongo.TransColl.Update(trans); err != nil {
		log.Errorf("update trans status fail: %s", err)
	}

	return ret
}

// saveErrorTran 保存失败的交易信息
func saveErrorTran(et *model.Trans) {
	if err := mongo.TransColl.Add(et); err != nil {
		log.Error("add errorTrans fail: ", err)
	}
}
