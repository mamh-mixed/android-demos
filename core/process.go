package core

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"strings"
)

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

	// 获取卡属性
	cardBin, err := mongo.CardBinColl.Find(bc.AcctNumDecrypt)
	if err != nil {
		if err.Error() == "not found" {
			return mongo.RespCodeColl.Get("200110")
		}
		return
	}
	log.Debugf("CardBin: %+v", cardBin)

	// 如果是银联卡，验证证件信息
	if strings.EqualFold("CUP", cardBin.CardBrand) || strings.EqualFold("UPI", cardBin.CardBrand) {
		ret = UnionPayCardValidity(bc)
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
		BankId:    bc.BankId,
		IdentType: bc.IdentType,
		IdentNum:  bc.IdentNum,
		PhoneNum:  bc.PhoneNum,
		ValidDate: bc.ValidDate,
		Cvv2:      bc.Cvv2,
	}
	if err := mongo.BindingInfoColl.Insert(bi); err != nil {
		log.Errorf("'InsertBindingInfo' error: (%s)\n 'BindingInfo': %+v", err, bi)
		return
	}

	// 根据商户、卡号、绑定Id、渠道、渠道商户生成一个系统绑定Id(ChanBindingId)，并将绑定关系映射入库
	bm := &model.BindingMap{
		MerId:         bc.MerId,
		BindingId:     bc.BindingId,
		ChanCode:      rp.ChanCode,
		ChanMerId:     rp.ChanMerId,
		ChanBindingId: tools.SerialNumber(),
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
		return
	}

	// bc(BindingCreate)用来向渠道发送请求，增加一些渠道要求的数据。
	bc.ChanMerId = rp.ChanMerId
	bc.ChanBindingId = bm.ChanBindingId
	bc.SignCert = chanMer.SignCert
	// TODO对加密的字段进行解密再送往渠道方

	log.Debugf("'BindingCreate' is: %+v", bc)

	// 如果是中金渠道，到数据库查找中金支持的银行卡的ID，并赋值给bindingCreate
	cm, err := mongo.CfcaBankMapColl.Find(cardBin.InsCode)
	if err != nil {
		log.Errorf("find CfcaBankMap ERROR!error message is: %s", err)
		return
	}
	bc.BankId = cm.BankId

	// todo 根据路由策略里面不同的渠道调用不同的绑定接口，这里为了简单，调用中金的接口
	c := channel.GetChan(bm.ChanCode)
	ret = c.ProcessBindingCreate(bc)
	// ret = cfca.ProcessBindingCreate(bc)

	// 渠道返回后，根据应答码，判断绑定是否成功，如果成功，更新绑定关系映射，绑定关系生效
	switch ret.RespCode {
	case "000000":
		bm.BindingStatus = model.BindingSuccess
	case "000009":
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
	be.SignCert = chanMer.SignCert
	// todo 查找该商户配置的渠道，这里为了简单，到中金查找。
	c := channel.GetChan(bm.ChanCode)
	ret = c.ProcessBindingEnquiry(be)

	// 转换绑定状态
	switch ret.RespCode {
	case "000000":
		bm.BindingStatus = model.BindingSuccess
	case "000009":
		bm.BindingStatus = model.BindingHandling
	default:
		bm.BindingStatus = model.BindingFail
	}

	// 更新绑定关系的状态
	if err = mongo.BindingMapColl.Update(bm); err != nil {
		log.Infof("'UpdateBindingMap' is: %+v", bm)
		log.Error("'UpdateBindingRelation' error: ", err)
	}

	return &model.BindingReturn{
		RespCode:      "000000",
		RespMsg:       "请求处理成功",
		BindingStatus: bm.BindingStatus,
	}
}

// ProcessBindingPayment 绑定支付
func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 默认返回
	ret = mongo.RespCodeColl.Get("000001")

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
		OrderNum:  be.MerOrderNum,
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
		errorTrans.RespCode = "200070"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Error("add errorTrans fail: ", err)
		}
		log.Error("not found any BindingMap: ", err)
		return mongo.RespCodeColl.Get("200070")
	}

	// 如果绑定关系不是成功的状态，返回
	switch bm.BindingStatus {
	case model.BindingHandling:
		errorTrans.RespCode = "200075"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Errorf("add errorTrans fail: (%s)", err)
		}
		return mongo.RespCodeColl.Get("200075")
	case model.BindingFail, model.BindingRemoved:
		errorTrans.RespCode = "200074"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Error("add errorTrans fail: ", err)
		}
		return mongo.RespCodeColl.Get("200074")
	}

	// 查找商家的绑定信息
	bi, err := mongo.BindingInfoColl.Find(be.MerId, be.BindingId)
	if err != nil {
		errorTrans.RespCode = "200070"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Error("add errorTrans fail: ", err)
		}
		log.Error("not found any BindingInfo: ", err)
		return mongo.RespCodeColl.Get("200070")
	}

	// 根据绑定关系得到渠道商户信息
	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
		errorTrans.RespCode = "300030"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Errorf("add errorTrans fail: (%s)", err)
		}
		log.Error("not found any chanMer: ", err)
		return mongo.RespCodeColl.Get("300030")
	}

	// 赋值
	be.SettFlag = chanMer.SettFlag
	be.ChanBindingId = bm.ChanBindingId
	be.ChanMerId = bm.ChanMerId
	be.ChanOrderNum = tools.SerialNumber()
	be.SignCert = chanMer.SignCert

	// 记录这笔交易
	trans := &model.Trans{
		OrderNum:      be.MerOrderNum,
		BindingId:     bi.BindingId,
		ChanOrderNum:  be.ChanOrderNum,
		ChanBindingId: be.ChanBindingId,
		AcctNum:       bi.AcctNum,
		MerId:         be.MerId,
		TransAmt:      be.TransAmt,
		ChanMerId:     be.ChanMerId,
		ChanCode:      bm.ChanCode,
		TransType:     model.PayTrans, //支付
		SendSmsId:     be.SendSmsId,
		SmsCode:       be.SmsCode,
		Remark:        be.Remark,
		SubMerId:      be.SubMerId,
	}
	if err = mongo.TransColl.Add(trans); err != nil {
		log.Errorf("add trans fail: (%s)", err)
		return
	}

	// 支付
	c := channel.GetChan(chanMer.ChanCode)
	ret = c.ProcessBindingPayment(be)

	// 处理结果
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

// ProcessBindingReomve 绑定解除
func ProcessBindingReomve(br *model.BindingRemove) (ret *model.BindingReturn) {
	ret = mongo.RespCodeColl.Get("000001")

	// 本地查询绑定关系
	bm, err := mongo.BindingMapColl.Find(br.MerId, br.BindingId)
	if err != nil {
		log.Error("'FindBindingRelation' error: ", err)
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
		log.Debugf("not found any chanMer (%s)", err)
		return mongo.RespCodeColl.Get("300030")
	}

	// 转换关系，补充信息
	br.ChanMerId = bm.ChanMerId
	br.ChanBindingId = bm.ChanBindingId
	br.TxSNUnBinding = tools.SerialNumber()
	br.SignCert = chanMer.SignCert

	// 到渠道解绑
	c := channel.GetChan(chanMer.ChanCode)
	ret = c.ProcessBindingRemove(br)

	// 如果解绑成功，更新本地数据库
	if ret.RespCode == "000000" {
		bm.BindingStatus = model.BindingRemoved
		if err := mongo.BindingMapColl.Update(bm); err != nil {
			log.Error("'Update BindingMap' error: ", err)
		}
	}

	return ret
}

// ProcessBindingRefund 退款
func ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {

	// default
	ret = mongo.RespCodeColl.Get("000001")

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
		OrderNum:       be.MerOrderNum,
		MerId:          be.MerId,
		RefundOrderNum: be.OrigOrderNum,
		TransAmt:       be.TransAmt,
	}

	// 是否有该源订单号
	orign, err := mongo.TransColl.Find(be.MerId, be.OrigOrderNum)
	switch {
	// 不存在原交易
	case err != nil:
		errorTrans.RespCode = "200091"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Errorf("add errorTrans fail: (%s)", err)
		}
		return mongo.RespCodeColl.Get("100020")
	// 已退款
	case orign.RefundStatus == model.TransRefunded:
		errorTrans.RespCode = "100010"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Errorf("add errorTrans fail: (%s)", err)
		}
		return mongo.RespCodeColl.Get("100010")
	// 退款金额过大
	case be.TransAmt > orign.TransAmt:
		errorTrans.RespCode = "200191"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Errorf("add errorTrans fail: (%s)", err)
		}
		return mongo.RespCodeColl.Get("200191")
	}

	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(orign.ChanCode, orign.ChanMerId)
	if err != nil {
		errorTrans.RespCode = "300030"
		if err = mongo.TransColl.Add(errorTrans); err != nil {
			log.Errorf("add errorTrans fail: (%s)", err)
		}
		log.Error("not found any chanMer: ", err)
		return mongo.RespCodeColl.Get("300030")
	}

	// 赋值
	be.ChanMerId = orign.ChanMerId
	be.ChanOrderNum = tools.SerialNumber()
	be.ChanOrigOrderNum = orign.ChanOrderNum
	be.SignCert = chanMer.SignCert

	// 记录这笔退款
	refund := &model.Trans{
		OrderNum:      be.MerOrderNum,
		ChanOrderNum:  be.ChanOrderNum,
		ChanBindingId: orign.ChanBindingId,
		//记录商户原订单而不是渠道原订单号
		RefundOrderNum: be.OrigOrderNum,
		AcctNum:        orign.AcctNum,
		MerId:          be.MerId,
		TransAmt:       be.TransAmt,
		ChanMerId:      be.ChanMerId,
		ChanCode:       orign.ChanCode,
		TransType:      model.RefundTrans, //退款
	}
	if err = mongo.TransColl.Add(refund); err != nil {
		log.Errorf("add refund trans fail : (%s)", err)
		return
	}

	// 退款
	c := channel.GetChan(chanMer.ChanCode)
	ret = c.ProcessBindingRefund(be)

	// 更新结果
	refund.ChanRespCode = ret.ChanRespCode
	refund.RespCode = ret.RespCode
	switch ret.RespCode {
	case "000000":
		refund.TransStatus = model.TransSuccess
		//更新原交易状态
		orign.RefundStatus = model.TransRefunded
		if err = mongo.TransColl.Update(orign); err != nil {
			log.Errorf("update orign trans RefundStatus fail : (%s)", err)
		}
	//只有超时才会出现000009
	case "000009":
		refund.TransStatus = model.TransHandling
	default:
		refund.TransStatus = model.TransFail
	}
	if err = mongo.TransColl.Update(refund); err != nil {
		log.Errorf("update refund trans status fail : (%s)", err)
	}
	return
}

// ProcessOrderEnquiry 订单查询
func ProcessOrderEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	// 默认返回成功的应答码
	ret = mongo.RespCodeColl.Get("000000")

	// 是否有该订单号
	t, err := mongo.TransColl.Find(be.MerId, be.OrigOrderNum)
	if err != nil {
		return mongo.RespCodeColl.Get("200080")
	}
	log.Debugf("trans:(%+v)", t)
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
	be.SignCert = chanMer.SignCert
	be.ChanMerId = chanMer.ChanMerId
	be.ChanOrderNum = t.ChanOrderNum

	// 原订单为处理中，向渠道发起查询
	result := new(model.BindingReturn)
	c := channel.GetChan(chanMer.ChanCode)
	switch t.TransType {
	//支付
	case model.PayTrans:
		result = c.ProcessPaymentEnquiry(be)
	//退款
	case model.RefundTrans:
		result = c.ProcessRefundEnquiry(be)
	}

	//更新交易状态
	switch result.RespCode {
	case "000000":
		t.TransStatus = model.TransSuccess
	case "000009":
		t.TransStatus = model.TransHandling
	default:
		t.TransStatus = model.TransFail
	}
	if err = mongo.TransColl.Update(t); err != nil {
		log.Errorf("Update trans error : %s", err)
	}

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
	ret = mongo.RespCodeColl.Get("000000")

	//查询
	rec, err := mongo.TransSettColl.Find(be.MerId, be.SettDate, be.NextOrderNum)
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
	ret = mongo.RespCodeColl.Get("000000")

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

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	log.Infof("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
