package core

import (
	"strings"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

var (
	successCode  = "000000"
	handlingCode = "000009"
)

// ProcessPaySettlement 支付结算
func ProcessPaySettlement(be *model.PaySettlement) (ret *model.BindingReturn) {

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum)
	if err != nil {
		return mongo.RespCodeColl.Get("000001")
	}
	if count > 0 {
		return mongo.RespCodeColl.Get("200081")
	}

	// 记录这笔结算
	sett := &model.Trans{
		MerId:        be.MerId,
		OrderNum:     be.MerOrderNum,
		TransType:    model.SettTrans,
		SubMerId:     be.SubMerId,
		TransAmt:     be.SettAmt,
		SettOrderNum: be.SettOrderNum,
		AcctNum:      be.SettAccountNum,
		AcctName:     be.SettAccountName,
		Province:     be.Province,
		City:         be.City,
		BranchName:   be.SettBranchName,
	}

	// 判断是否是合法的结算订单号
	settOrder, err := mongo.TransColl.GetBySettOrder(be.MerId, be.SettOrderNum)
	if err != nil {
		return logicErrorHandle(sett, "200252")
	}

	// 获取卡bin详情
	cardBin, err := findCardBin(be.AcctNumDecrypt)
	if err != nil {
		return logicErrorHandle(sett, "200110")
	}
	sett.ChanCode = settOrder.ChanCode
	sett.ChanMerId = settOrder.ChanMerId

	// 渠道商户信息
	chanMer, err := mongo.ChanMerColl.Find(sett.ChanCode, sett.ChanMerId)
	if err != nil {
		return logicErrorHandle(sett, "300030")
	}

	// 暂时判断中金渠道是否支持该银行卡
	cm, err := mongo.CfcaBankMapColl.Find(cardBin.InsCode)
	if err != nil {
		return logicErrorHandle(sett, "400020")
	}

	// 交易参数
	sett.SysOrderNum = util.SerialNumber()
	// 渠道请求参数
	be.ChanMerId = sett.ChanMerId
	be.SysOrderNum = sett.SysOrderNum
	be.PrivateKey = chanMer.PrivateKey
	be.BankCode = cm.BankId
	be.SettOrderNum = be.MerId + be.SettOrderNum

	// 获取渠道接口
	c := channel.GetChan(sett.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return logicErrorHandle(sett, "510010")
	}

	// 记录这笔交易
	if err = mongo.TransColl.Add(sett); err != nil {
		log.Errorf("add trans error: %s", err)
		return mongo.RespCodeColl.Get("000001")
	}

	ret = c.ProcessPaySettlement(be)

	// 更新
	transStatusHandle(ret, sett)
	mongo.TransColl.Update(sett)

	return ret
}

// ProcessGetCardInfo 获取卡片信息
func ProcessGetCardInfo(bc *model.CardInfo) (ret *model.BindingReturn) {

	// 获取卡bin详情
	cardBin, err := findCardBin(bc.CardNum)
	if err != nil {
		log.Error(err)
		return mongo.RespCodeColl.Get("200110")
	}

	// 通过路由策略找到渠道和渠道商户
	// rp := mongo.RouterPolicyColl.Find(bc.MerId, cardBin.CardBrand)
	// if rp == nil {
	// 	return mongo.RespCodeColl.Get("300030")
	// }

	// 返回卡片信息
	ret = mongo.RespCodeColl.Get(successCode)
	ret.CardNum = bc.CardNum
	ret.CardBrand = cardBin.CardBrand
	ret.AcctType = cardBin.AcctType
	ret.IssBankName = cardBin.InsName
	ret.IssBankNum = cardBin.InsCode

	// 暂时判断中金渠道是否支持该银行卡
	cm, err := mongo.CfcaBankMapColl.Find(cardBin.InsCode)
	if err != nil {
		// 没找到说明不支持
		ret.BindingPaySupport = "0"
		return ret
	}

	ret.BindingPaySupport = "1"
	ret.BankCode = cm.BankId
	// ...图片地址

	return ret
}

// ProcessBindingCreate 绑定建立的业务处理
func ProcessBindingCreate(bc *model.BindingCreate) (ret *model.BindingReturn) {
	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")

	// 验证该机构下，该绑定号是否已经绑定了
	count, err := mongo.BindingMapColl.Count(bc.MerId, bc.BindingId)
	if count > 0 {
		return mongo.RespCodeColl.Get("200071")
	}

	// todo 如果需要校验短信，验证短信
	// ret = validateSmsCode(bc.SendSmsId, bc.SmsCode)
	// if ret != nil {
	// 	return ret
	// }

	// 获取卡bin详情
	cardBin, err := findCardBin(bc.AcctNumDecrypt)
	if err != nil {
		log.Error(err)
		return mongo.RespCodeColl.Get("200110")
	}
	// log.Debugf("CardBin: %+v", cardBin)

	// 如果是银联卡，验证证件信息
	if strings.EqualFold("CUP", cardBin.CardBrand) || strings.EqualFold("UPI", cardBin.CardBrand) {
		ret = UnionPayCardCommonValidity(bc.IdentType, bc.IdentNumDecrypt, bc.PhoneNumDecrypt)
		if ret != nil {
			return ret
		}
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(bc.MerId, cardBin.CardBrand)
	if rp == nil {
		return mongo.RespCodeColl.Get("300030")
	}

	// 商家绑定信息和绑定映射入库
	bi := &model.BindingInfo{
		MerId:     bc.MerId,
		BindingId: bc.BindingId,
		CardBrand: rp.CardBrand,
		AcctType:  bc.IdentType,
		AcctName:  bc.AcctName,
		AcctNum:   bc.AcctNum,
		BankId:    bc.BankCode,
		IdentType: bc.IdentType,
		IdentNum:  bc.IdentNum,
		PhoneNum:  bc.PhoneNum,
		ValidDate: bc.ValidDate,
		Cvv2:      bc.Cvv2,
	}
	if err := mongo.BindingInfoColl.Insert(bi); err != nil {
		log.Errorf("'InsertBindingInfo' error: (%s) 'BindingInfo': %+v", err, bi)
		return
	}

	// 根据商户、卡号、绑定Id、渠道、渠道商户生成一个系统绑定Id(ChanBindingId)，并将绑定关系映射入库
	bm := &model.BindingMap{
		MerId:         bc.MerId,
		BindingId:     bc.BindingId,
		ChanCode:      rp.ChanCode,
		ChanMerId:     rp.ChanMerId,
		ChanBindingId: util.SerialNumber(),
		BindingStatus: "",
	}
	if err := mongo.BindingMapColl.Insert(bm); err != nil {
		log.Errorf("'InsertBindingMap' error: (%s)\n 'BindingMap': %+v", err, bm)
		return
	}

	// 根据绑定关系得到渠道商户信息
	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		log.Errorf("not found any chanMer (%s)", err)
		return mongo.RespCodeColl.Get("300030")
	}

	// bc(BindingCreate)用来向渠道发送请求，增加一些渠道要求的数据。
	bc.ChanMerId = rp.ChanMerId
	bc.ChanBindingId = bm.ChanBindingId
	bc.PrivateKey = chanMer.PrivateKey

	log.Tracef("'BindingCreate' is: %+v", bc)

	// 如果接入方没有送，则到渠道获取
	// 如果是中金渠道，到数据库查找中金支持的银行卡的ID，并赋值给bindingCreate
	if bc.BankCode == "" {
		cm, err := mongo.CfcaBankMapColl.Find(cardBin.InsCode)
		if err != nil {
			log.Errorf("find CfcaBankMap ERROR!error message is: %s", err)
			return mongo.RespCodeColl.Get("400020")
		}
		bc.BankCode = cm.BankId
	}

	// 根据路由策略里面不同的渠道调用不同的绑定接口
	c := channel.GetChan(bm.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return mongo.RespCodeColl.Get("510010")
	}

	ret = c.ProcessBindingCreate(bc)

	// 渠道返回后，根据应答码，判断绑定是否成功，如果成功，更新绑定关系映射，绑定关系生效
	switch ret.RespCode {
	case successCode:
		bm.BindingStatus = model.BindingSuccess
	case handlingCode:
		bm.BindingStatus = model.BindingHandling
	default:
		bm.BindingStatus = model.BindingFail
	}
	err = mongo.BindingMapColl.Update(bm)
	if err != nil {
		log.Infof("'BindingMap' is: %+v", bm)
		log.Error("'BindingMapColl update' error: ", err)
	}

	return ret
}

// ProcessBindingEnquiry 绑定关系查询
// 先到本地库去查找，如果本地库查找的结果是正在处理中，就到渠道查找；查找完更新到数据库中。
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")

	// 本地查询绑定关系
	bm, err := mongo.BindingMapColl.Find(be.MerId, be.BindingId)
	if err != nil {
		log.Errorf("'FindBindingMap' error: %s", err)
		return mongo.RespCodeColl.Get("200070")
	}
	log.Debugf("binding result: %#v", bm)

	// 非处理中，直接返回结果
	if bm.BindingStatus != model.BindingHandling {
		return &model.BindingReturn{
			RespCode:      "000000",
			RespMsg:       "请求处理成功",
			BindingStatus: bm.BindingStatus,
		}
	}

	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
		log.Debugf("not found any chanMer (%s)", err)
		return ret
	}

	// 正在处理中，到渠道那边查找
	// 转换绑定关系、请求
	be.ChanMerId = bm.ChanMerId
	be.ChanBindingId = bm.ChanBindingId
	be.PrivateKey = chanMer.PrivateKey

	// 查找该商户配置的渠道。
	c := channel.GetChan(bm.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return mongo.RespCodeColl.Get("510010")
	}

	ret = c.ProcessBindingEnquiry(be)

	// 转换绑定状态
	switch ret.RespCode {
	case successCode:
		bm.BindingStatus = model.BindingSuccess
	case handlingCode:
		bm.BindingStatus = model.BindingHandling
	default:
		bm.BindingStatus = model.BindingFail
	}

	// 更新绑定关系的状态
	if err = mongo.BindingMapColl.Update(bm); err != nil {
		// log.Infof("'UpdateBindingMap' is: %+v", bm)
		log.Error("'UpdateBindingRelation' error: ", err)
	}

	return &model.BindingReturn{
		RespCode:      successCode,
		RespMsg:       "请求处理成功",
		BindingStatus: bm.BindingStatus,
	}
}

// ProcessBindingReomve 绑定解除
func ProcessBindingReomve(br *model.BindingRemove) (ret *model.BindingReturn) {
	ret = mongo.RespCodeColl.Get("000001")

	// 本地查询绑定关系
	bm, err := mongo.BindingMapColl.Find(br.MerId, br.BindingId)
	if err != nil {
		return mongo.RespCodeColl.Get("200070")
	}

	switch bm.BindingStatus {
	// // todo 绑定状态为处理中的话
	// case model.BindingHandling:
	// 	return model.NewBindingReturn("200070", "绑定ID有误")
	// 绑定状态为已解绑的话
	case model.BindingRemoved:
		return mongo.RespCodeColl.Get("200072")
	// 绑定状态为失败的话
	case model.BindingFail:
		return mongo.RespCodeColl.Get("200073")
	}

	// 查找渠道商户信息，获取证书
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
		return mongo.RespCodeColl.Get("300030")
	}

	// 转换关系，补充信息
	br.ChanMerId = bm.ChanMerId
	br.ChanBindingId = bm.ChanBindingId
	br.TxSNUnBinding = util.SerialNumber()
	br.PrivateKey = chanMer.PrivateKey

	// 到渠道解绑
	c := channel.GetChan(chanMer.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return mongo.RespCodeColl.Get("510010")
	}

	ret = c.ProcessBindingRemove(br)

	// 如果解绑成功，更新本地数据库
	if ret.RespCode == successCode {
		bm.BindingStatus = model.BindingRemoved
		mongo.BindingMapColl.Update(bm)
	}

	return ret
}

// ProcessPaymentWithSMS 通过验证码完成交易
func ProcessPaymentWithSMS(be *model.BindingPayment) (ret *model.BindingReturn) {

	// 找到原订单
	orig, err := mongo.TransColl.FindOne(be.MerId, be.MerOrderNum)
	if err != nil {
		return mongo.RespCodeColl.Get("200080")
	}

	// 订单状态是否是待支付
	if orig.TransStatus != model.TransNotPay {
		return mongo.RespCodeColl.Get("200081") // TODO:返回订单号重复。待确认
	}

	// 获取渠道接口
	c := channel.GetChan(orig.ChanCode)
	if c == nil {
		// 来到这里必须是系统错误。
		return mongo.RespCodeColl.Get("000001")
	}

	// 根据绑定关系得到渠道商户信息
	chanMer, err := mongo.ChanMerColl.Find(orig.ChanCode, orig.ChanMerId)
	if err != nil {
		// 来到这里必须是系统错误。
		return mongo.RespCodeColl.Get("000001")
	}

	switch chanMer.TransMode {
	case model.MarketMode:
		if be.SettOrderNum == "" {
			return model.NewBindingReturn("200050", "字段 SettOrderNum 不能为空")
		}
		if be.SettOrderNum != orig.SettOrderNum {
			return mongo.RespCodeColl.Get("200251")
		}
		be.SettOrderNum = be.MerId + be.SettOrderNum
	case model.MerMode:
	default:
		log.Errorf("Unsupport mode %s", chanMer.TransMode)
		return
	}

	// 渠道需要参数
	be.Mode = chanMer.TransMode
	be.SysOrderNum = orig.SysOrderNum
	be.ChanMerId = orig.ChanMerId
	be.PrivateKey = chanMer.PrivateKey

	// 请求支付
	ret = c.ProcessPaymentWithSMS(be)

	// 更新交易
	transStatusHandle(ret, orig)
	mongo.TransColl.Update(orig)
	return ret
}

// ProcessBindingPayment 绑定支付
// isSendSMS:true 将调用支付前发短信接口;false 直接支付，中金暂不支持
func ProcessBindingPayment(be *model.BindingPayment, isSendSMS bool) (ret *model.BindingReturn) {
	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum)
	if err != nil {
		return
	}
	if count > 0 {
		return mongo.RespCodeColl.Get("200081")
	}

	//只要订单号不重复就记录这笔交易
	pay := &model.Trans{
		MerId:     be.MerId,
		OrderNum:  be.MerOrderNum,
		TransType: model.PayTrans,
		BindingId: be.BindingId,
		TransAmt:  be.TransAmt,
		SendSmsId: be.SendSmsId,
		SmsCode:   be.SmsCode,
		Remark:    be.Remark,
		SubMerId:  be.SubMerId,
	}

	// 本地查询绑定关系。查询绑定关系的状态是否成功
	bm, err := mongo.BindingMapColl.Find(be.MerId, be.BindingId)
	if err != nil {
		return logicErrorHandle(pay, "200070")
	}
	pay.ChanCode = bm.ChanCode
	pay.ChanMerId = bm.ChanMerId
	pay.ChanBindingId = bm.ChanBindingId

	// 绑定状态处理中
	if bm.BindingStatus == model.BindingHandling {
		return logicErrorHandle(pay, "200075")
	}

	// 绑定状态失败
	if bm.BindingStatus == model.BindingFail || bm.BindingStatus == model.BindingRemoved {
		return logicErrorHandle(pay, "200074")
	}

	// 查找商家的绑定信息
	bi, err := mongo.BindingInfoColl.Find(be.MerId, be.BindingId)
	if err != nil {
		return logicErrorHandle(pay, "200070")
	}
	pay.AcctNum = bi.AcctNum

	// TODO 金额是否超出最大可支付金额

	// 根据绑定关系得到渠道商户信息
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
		return logicErrorHandle(pay, "300030")
	}

	// 交易模式
	switch chanMer.TransMode {
	case model.MarketMode:
		if be.SettOrderNum == "" {
			return model.NewBindingReturn("200050", "字段 SettOrderNum 不能为空")
		}
		pay.SettOrderNum = be.SettOrderNum
		be.SettOrderNum = be.MerId + be.SettOrderNum
	case model.MerMode:
		be.SettFlag = chanMer.SettFlag
	default:
		log.Errorf("Unsupport mode %s", chanMer.TransMode)
		return
	}

	// 交易参数
	pay.SysOrderNum = util.SerialNumber()
	// 渠道请求参数
	be.Mode = chanMer.TransMode
	be.ChanBindingId = pay.ChanBindingId
	be.ChanMerId = pay.ChanMerId
	be.SysOrderNum = pay.SysOrderNum
	be.PrivateKey = chanMer.PrivateKey

	// 获取渠道接口
	c := channel.GetChan(chanMer.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return logicErrorHandle(pay, "510010")
	}

	// 记录这笔交易
	if err = mongo.TransColl.Add(pay); err != nil {
		log.Errorf("add trans error: %s", err)
		return
	}

	if isSendSMS {
		// 发送支付验证码
		ret = c.ProcessSendBindingPaySMS(be)
		// 发送结果成功，交易待支付
		pay.RespCode = ret.RespCode
		pay.ChanRespCode = ret.ChanRespCode
		pay.ErrorDetail = ret.RespMsg
		if ret.RespCode == successCode {
			pay.TransStatus = model.TransNotPay
		}
	} else {
		// 直接支付
		ret = c.ProcessBindingPayment(be)
		// 处理结果
		transStatusHandle(ret, pay)
	}

	// 更新交易
	mongo.TransColl.Update(pay)
	return ret
}

// ProcessBindingRefund 退款
func ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return
	}
	if count > 0 {
		return mongo.RespCodeColl.Get("200081")
	}
	//只要订单号不重复就记录这笔交易
	refund := &model.Trans{
		OrderNum:     be.MerOrderNum,
		MerId:        be.MerId,
		OrigOrderNum: be.OrigOrderNum,
		TransAmt:     be.TransAmt,
		TransType:    model.RefundTrans,
	}

	// 是否有该源订单号
	orign, err := mongo.TransColl.FindOne(be.MerId, be.OrigOrderNum)
	if err != nil {
		return logicErrorHandle(refund, "200082")
	}
	refund.ChanBindingId = orign.ChanBindingId
	refund.AcctNum = orign.AcctNum
	refund.ChanCode = orign.ChanCode
	refund.ChanMerId = orign.ChanMerId

	// 不能对退款的交易号进行退款
	if orign.TransType != model.PayTrans {
		return logicErrorHandle(refund, "200090")
	}

	// 原交易不成功
	if orign.TransStatus != model.TransSuccess {
		return logicErrorHandle(refund, "100020")
	}

	// 中金不支持部分退款
	if orign.ChanCode == "CFCA" && be.TransAmt != orign.TransAmt {
		return logicErrorHandle(refund, "200190")
	}

	refundAmt := refund.TransAmt
	// 退款状态是否可退
	switch orign.RefundStatus {
	// 已退款
	case model.TransRefunded:
		return logicErrorHandle(refund, "100010")
	// 部分退款
	case model.TransPartRefunded:
		refundAmt += orign.RefundAmt
		fallthrough
	default:
		// 金额过大
		if refundAmt > orign.TransAmt {
			return logicErrorHandle(refund, "200191")
		} else if refundAmt == orign.TransAmt {
			orign.RefundStatus = model.TransRefunded
			orign.RefundAmt = refundAmt
		} else {
			orign.RefundStatus = model.TransPartRefunded
			orign.RefundAmt = refundAmt
		}
	}

	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(orign.ChanCode, orign.ChanMerId)
	if err != nil {
		log.Error("Find channel merchant error,error message is '%s'", err)
		return logicErrorHandle(refund, "300030")
	}

	// 是中金退款的话
	if orign.ChanCode == "CFCA" {
		// 交易模式
		switch chanMer.TransMode {
		case model.MarketMode:
			if be.SettOrderNum == "" {
				return model.NewBindingReturn("200050", "字段 SettOrderNum 不能为空")
			}
			if be.SettOrderNum != orign.SettOrderNum {
				return logicErrorHandle(refund, "200251")
			}
			refund.SettOrderNum = be.SettOrderNum
			be.SettOrderNum = be.MerId + be.SettOrderNum
		case model.MerMode:
			// dothing
		default:
			log.Errorf("Unsupport mode %d", chanMer.TransMode)
			return logicErrorHandle(refund, "000001")
		}
	}

	// 获取渠道接口
	c := channel.GetChan(chanMer.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return logicErrorHandle(refund, "510010")
	}

	// 请求信息
	be.Mode = chanMer.TransMode
	be.ChanMerId = orign.ChanMerId
	be.SysOrderNum = util.SerialNumber()
	be.SysOrigOrderNum = orign.SysOrderNum
	be.PrivateKey = chanMer.PrivateKey
	// 交易信息
	refund.SysOrderNum = be.SysOrderNum

	// 记录这笔退款
	if err = mongo.TransColl.Add(refund); err != nil {
		return mongo.RespCodeColl.Get("000001")
	}

	// 退款
	ret = c.ProcessBindingRefund(be)

	// 更新原交易状态
	if ret.RespCode == successCode {
		mongo.TransColl.Update(orign)
	}

	// 更新结果
	transStatusHandle(ret, refund)
	mongo.TransColl.Update(refund)
	return
}

// ProcessOrderEnquiry 订单查询
func ProcessOrderEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	// 默认返回成功的应答码
	ret = mongo.RespCodeColl.Get(successCode)

	// 是否有该订单号
	t, err := mongo.TransColl.FindOne(be.MerId, be.OrigOrderNum)
	if err != nil {
		return mongo.RespCodeColl.Get("200080")
	}

	// 如果交易状态不是在处理中
	if t.TransStatus != model.TransHandling {
		ret.TransStatus = t.TransStatus
		if be.ShowOrigInfo == "1" {
			ret.OrigTransDetail = model.NewTransInfo(*t)
		}
		return
	}

	// 获得渠道商户信息
	chanMer, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		log.Error("not found any chanMer: ", err)
		return mongo.RespCodeColl.Get("300030")
	}

	//赋值
	be.Mode = chanMer.TransMode
	be.PrivateKey = chanMer.PrivateKey
	be.ChanMerId = chanMer.ChanMerId
	be.SysOrderNum = t.SysOrderNum

	// 原订单为处理中，向渠道发起查询
	result := new(model.BindingReturn)
	c := channel.GetChan(chanMer.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return mongo.RespCodeColl.Get("510010")
	}

	switch t.TransType {
	//支付
	case model.PayTrans:
		result = c.ProcessPaymentEnquiry(be)
	//退款
	case model.RefundTrans:
		result = c.ProcessRefundEnquiry(be)
	}

	//更新交易状态
	transStatusHandle(result, t)
	mongo.TransColl.Update(t)

	//返回结果
	ret.TransStatus = t.TransStatus
	if be.ShowOrigInfo == "1" {
		ret.OrigTransDetail = model.NewTransInfo(*t)
	}

	return
}

// ProcessBillingDetails 交易对账明细查询
func ProcessBillingDetails(be *model.BillingDetails) (ret *model.BindingReturn) {

	//default return
	ret = mongo.RespCodeColl.Get(successCode)

	//查询
	rec, err := mongo.TransSettColl.FindByDate(be.MerId, be.SettDate, be.NextOrderNum)
	if err != nil {
		log.Errorf("Find transSett records error : %s", err)
		return
	}

	//结果处理
	//暂时默认商户一次可取为10条
	//实际查询可取为11条、包含下次查询的第一条
	if len(rec) == 11 {
		ret.Rec = rec[:len(rec)-1]
		ret.NextOrderNum = rec[len(rec)-1].OrderNum
	} else {
		//如果不够11条、直接赋值
		//NextOrderNum为空
		ret.Rec = rec
	}

	//赋值
	ret.Count = len(ret.Rec)
	return
}

// ProcessBillingSummary 交易对账汇总查询
func ProcessBillingSummary(be *model.BillingSummary) (ret *model.BindingReturn) {

	//default return
	ret = mongo.RespCodeColl.Get(successCode)

	//查询
	data, err := mongo.TransSettColl.Summary(be.MerId, be.SettDate)
	if err != nil {
		log.Errorf("summary transSett records error : %s", err)
		return
	}

	//赋值
	ret.SettDate = be.SettDate
	ret.Data = data
	return
}

// ProcessNoTrackPayment 处理无卡直接支付的业务逻辑
func ProcessNoTrackPayment(be *model.NoTrackPayment) (ret *model.BindingReturn) {
	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")

	// 系统唯一的序列号
	sysOrderNum := mongo.SnColl.GetSysSN()

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum)
	if err != nil {
		log.Errorf("find trans fail : (%s)", err)
		return
	}
	if count > 0 {
		return mongo.RespCodeColl.Get("200081")
	}

	//只要订单号不重复就记录这笔交易
	errorTrans := &model.Trans{
		OrderNum:    be.MerOrderNum,
		SysOrderNum: sysOrderNum,
		AcctNum:     be.AcctNum,
		MerId:       be.MerId,
		TransAmt:    be.TransAmt,
		TransCurr:   be.CurrCode,
		SendSmsId:   be.SendSmsId,
		SmsCode:     be.SmsCode,
		SubMerId:    be.SubMerId,
	}

	// 暂时不支持预授权交易
	if be.TransType == "AUTH" {
		log.Error("暂时不支持预授权交易")
		errorTrans.RespCode = "100030"
		saveErrorTran(errorTrans)
		return mongo.RespCodeColl.Get("100030")
	}

	// 暂不支持借记卡
	if be.AcctType == "10" {
		log.Error("暂不支持借记卡")
		errorTrans.RespCode = "300030"
		saveErrorTran(errorTrans)
		return mongo.RespCodeColl.Get("300030")
	}

	// 获取卡bin详情
	cardBin, err := findCardBin(be.AcctNumDecrypt)
	if err != nil {
		log.Errorf("find card bin error: %s", err)
		if err.Error() == "not found" {
			errorTrans.RespCode = "200070"
			saveErrorTran(errorTrans)
			return mongo.RespCodeColl.Get("200110")
		}
		saveErrorTran(errorTrans)
		return
	}
	log.Debugf("CardBin: %+v", cardBin)

	// 银联卡校验
	if strings.EqualFold("CUP", cardBin.CardBrand) || strings.EqualFold("UPI", cardBin.CardBrand) {
		result := UnionPayCardCommonValidity(be.IdentType, be.IdentNumDecrypt, be.PhoneNumDecrypt)
		if result != nil {
			return result
		}
	}

	// 通过路由策略找到渠道和渠道商户
	rp := mongo.RouterPolicyColl.Find(be.MerId, cardBin.CardBrand)
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

	// 交易币种如果为空的话默认为156
	if be.CurrCode == "" {
		be.CurrCode = "156"
	}

	// 下游送来的终端号，如果没有的话，填上渠道商户里面的配置的终端号
	if be.TerminalId == "" {
		be.TerminalId = chanMer.TerminalId
	}
	// 补充信息
	be.Chcd = chanMer.InsCode
	be.Mchntid = chanMer.ChanMerId
	be.CliSN = mongo.SnColl.GetDaySN(chanMer.ChanMerId, chanMer.TerminalId)
	be.SysSN = sysOrderNum

	// 记录这笔交易，入库
	trans := &model.Trans{
		OrderNum:    be.MerOrderNum,
		SysOrderNum: sysOrderNum,
		AcctNum:     be.AcctNum,
		MerId:       be.MerId,
		TransAmt:    be.TransAmt,
		TransCurr:   be.CurrCode,
		SendSmsId:   be.SendSmsId,
		SmsCode:     be.SmsCode,
		SubMerId:    be.SubMerId,
		ChanMerId:   chanMer.ChanMerId,
		ChanCode:    chanMer.ChanCode,
		TransType:   1, //TODO 预授权不属于支付
		Remark:      "Apple Pay",
	}
	if err := mongo.TransColl.Add(trans); err != nil {
		log.Errorf("add trans fail: (%s)", err)
		return
	}

	// 查找配置的渠道入口
	c := channel.GetDirectPayChan(chanMer.ChanCode)
	if c == nil {
		log.Error("Channel interface is unavailable,error message is 'get channel return nil'")
		return mongo.RespCodeColl.Get("510010")
	}

	// 消费
	ret = c.Consume(be)

	// 更新结构
	transStatusHandle(ret, trans)
	mongo.TransColl.Update(trans)

	return ret
}

func transStatusHandle(ret *model.BindingReturn, t *model.Trans) {
	t.ChanRespCode = ret.ChanRespCode
	t.RespCode = ret.RespCode
	t.ErrorDetail = ret.RespMsg
	switch ret.RespCode {
	case successCode:
		t.TransStatus = model.TransSuccess
	case handlingCode:
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	log.Infof("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}

// logicErrorHandle 逻辑错误错误，保存交易
func logicErrorHandle(t *model.Trans, respCode string) (ret *model.BindingReturn) {
	t.RespCode = respCode
	mongo.TransColl.Add(t)
	return mongo.RespCodeColl.Get(respCode)
}
