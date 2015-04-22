package core

import (
	"github.com/CardInfoLink/quickpay/channel/cil"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// ProcessApplyPay 执行apple Pay相关的支付
// TODO 默认为消费，预授权尚未处理
func ProcessApplePay(ap *model.ApplePay) (ret *model.BindingReturn) {
	log.Debugf("Apple Pay请求数据: %+v", ap)
	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")
	// 系统唯一的序列号
	chanOrderNum := mongo.SnColl.GetSysSN()

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(ap.MerId, ap.MerOrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return
	}
	if count > 0 {
		return mongo.RespCodeColl.Get("200081")
	}

	//只要订单号不重复就记录这笔交易
	errorTrans := &model.Trans{
		MerId:        ap.MerId,
		OrderNum:     ap.MerOrderNum,
		ChanOrderNum: chanOrderNum,
		TransAmt:     ap.ApplePayData.TransactionAmount,
		TransCurr:    ap.ApplePayData.CurrencyCode,
		TransType:    1, //TODO 预授权不属于支付
		RespCode:     "000001",
		Remark:       "Apple Pay",
		SubMerId:     ap.SubMerId,
	}

	// 如果是预授权交易，先返回不支持此类交易
	if ap.TransType == "AUTH" {
		errorTrans.RespCode = "100030"
		return mongo.RespCodeColl.Get("100030")
	}

	// 根据卡号查找卡属性，然后匹配路由查找
	// 获取卡属性
	appleAcctNum := ap.ApplePayData.ApplicationPrimaryAccountNumber
	bin := tree.match(appleAcctNum)
	log.Debugf("cardNum=%s, cardBin=%s", appleAcctNum, bin)
	// 获取卡bin详情
	cardBin, err := mongo.CardBinColl.Find(bin, len(appleAcctNum))

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
	// 补充信息
	ap.Chcd = chanMer.InsCode
	ap.Mchntid = chanMer.ChanMerId
	ap.CliSN = mongo.SnColl.GetDaySN(chanMer.ChanMerId, chanMer.TerminalId)
	ap.SysSN = chanOrderNum

	// 记录这笔交易，入库
	trans := &model.Trans{
		MerId:        ap.MerId,
		OrderNum:     ap.MerOrderNum,
		ChanOrderNum: ap.SysSN,
		AcctNum:      ap.ApplePayData.ApplicationPrimaryAccountNumber,
		TransAmt:     ap.ApplePayData.TransactionAmount,
		ChanMerId:    chanMer.ChanMerId,
		ChanCode:     chanMer.ChanCode,
		TransType:    1, //TODO 预授权不属于支付
		Remark:       "Apple Pay",
		SubMerId:     ap.SubMerId,
	}
	if err := mongo.TransColl.Add(trans); err != nil {
		log.Errorf("add trans fail: (%s)", err)
		return
	}

	ret = cil.ConsumeByApplePay(ap)

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
		log.Error("update trans status fail ", err)
	}

	return ret
}

// saveErrorTran 保存失败的交易信息
func saveErrorTran(et *model.Trans) {
	if err := mongo.TransColl.Add(et); err != nil {
		log.Error("add errorTrans fail: ", err)
	}
}
