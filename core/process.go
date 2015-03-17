package core

import (
	"quickpay/channel/cfca"
	"quickpay/model"
	"quickpay/mongo"
	"quickpay/tools"
	"strings"

	"github.com/omigo/g"
)

// ProcessBindingCreate 绑定建立的业务处理
// todo 先验证是否已经绑定过
func ProcessBindingCreate(bc *model.BindingCreate) (ret *model.BindingReturn) {
	// todo 如果需要校验短信，验证短信
	// ret = validateSmsCode(bc.SendSmsId, bc.SmsCode)
	// if ret != nil {
	// 	return ret
	// }

	// 获取卡属性
	cardBin := mongo.CardBinColl.Find(bc.AcctNum)
	g.Debug("CardBin: %+v", cardBin)

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
		return model.NewBindingReturn("300030", "无此交易权限")
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
		g.Error("'InsertBindingInfo' error: (%s)\n 'BindingInfo': %+v", err, bi)
		return model.NewBindingReturn("000001", "系统内部错误")
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
		g.Error("'InsertBindingMap' error: (%s)\n 'BindingMap': %+v", err, bm)
		return model.NewBindingReturn("000001", "系统内部错误")
	}

	// 根据绑定关系得到渠道商户信息
	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(rp.ChanCode, rp.ChanMerId)
	if err != nil {
		g.Debug("not found any chanMer (%s)", err)
		return ret
	}

	// bc(BindingCreate)用来向渠道发送请求，增加一些渠道要求的数据。
	bc.ChanMerId = rp.ChanMerId
	bc.ChanBindingId = bm.ChanBindingId
	bc.SignCert = chanMer.SignCert
	g.Info("'BindingCreate' is: %+v", bc)
	// todo 根据路由策略里面不同的渠道调用不同的绑定接口，这里为了简单，调用中金的接口。
	ret = cfca.ProcessBindingCreate(bc)

	// 渠道返回后，根据应答码，判断绑定是否成功，如果成功，更新绑定关系映射，绑定关系生效。
	bm.BindingStatus = ret.RespCode
	err = mongo.BindingMapColl.Update(bm)
	if err != nil {
		g.Info("'BindingMap' is: %+v", bm)
		g.Error("'BindingMapColl update' error: ", err.Error())
	}

	return ret
}

// ProcessBindingEnquiry 绑定关系查询。
// 先到本地库去查找，如果本地库查找的结果是正在处理中，就到渠道查找；查找完更新到数据库中。
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 默认返回
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 本地查询绑定关系
	bm, err := mongo.BindingMapColl.Find(be.MerId, be.BindingId)
	if err != nil {
		g.Error("'FindBindingMap' error: ", err.Error())
		return model.NewBindingReturn("200101", "绑定ID不正确")
	}

	// 非处理中，直接返回结果
	if bm.BindingStatus != "000009" {
		return mongo.RespCodeColl.Get(bm.BindingStatus)
	}

	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
		g.Debug("not found any chanMer (%s)", err)
		return ret
	}

	// 正在处理中，到渠道那边查找
	// 转换绑定关系、请求
	be.ChanMerId = bm.ChanMerId
	be.ChanBindingId = bm.ChanBindingId
	be.SignCert = chanMer.SignCert
	// todo 查找该商户配置的渠道，这里为了简单，到中金查找。
	ret = cfca.ProcessBindingEnquiry(be)

	// 更新绑定关系的状态
	bm.BindingStatus = ret.RespCode
	if err = mongo.BindingMapColl.Update(bm); err != nil {
		g.Info("'UpdateBindingMap' is: %+v", bm)
		g.Error("'UpdateBindingRelation' error: ", err.Error())
		return model.NewBindingReturn("000001", "系统内部错误")
	}

	return ret
}

// ProcessBindingPayment 绑定支付
func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 默认返回
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum, 1)
	if err != nil {
		g.Error("find trans fail : (%s)", err)
		return
	}
	if count > 0 {
		return model.NewBindingReturn("200081", "订单号重复")
	}

	// 本地查询绑定关系。查询绑定关系的状态是否成功
	bm, err := mongo.BindingMapColl.Find(be.MerId, be.BindingId)
	if err != nil {
		g.Error("not found any BindingMap: ", err)
		return model.NewBindingReturn("200101", "绑定ID不正确")
	}
	// 如果绑定关系不是成功的状态，返回
	if bm.BindingStatus != "000000" {
		return model.NewBindingReturn(bm.BindingStatus, "绑定中或者绑定失败，请查询绑定关系。")
	}

	// 查找商家的绑定信息
	bi, err := mongo.BindingInfoColl.Find(be.MerId, be.BindingId)
	if err != nil {
		g.Error("not found any BindingInfo: ", err)
		return model.NewBindingReturn("200101", "绑定ID不正确")
	}

	// 根据绑定关系得到渠道商户信息
	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
		g.Error("not found any chanMer: ", err)
		return model.NewBindingReturn("300030", "无此交易权限")
	}

	// 赋值
	be.SettFlag = chanMer.SettFlag
	be.ChanBindingId = bm.ChanBindingId
	be.ChanMerId = bm.ChanMerId
	be.SignCert = chanMer.SignCert

	// 记录这笔交易
	trans := &model.Trans{
		OrderNum:      be.MerOrderNum,
		ChanOrderNum:  tools.SerialNumber(),
		ChanBindingId: be.ChanBindingId,
		AcctNum:       bi.AcctNum,
		MerId:         be.MerId,
		TransAmount:   be.TransAmt,
		ChanMerId:     be.ChanMerId,
		ChanCode:      bm.ChanCode,
		TransType:     1,
	}
	if err = mongo.TransColl.Add(trans); err != nil {
		g.Error("add trans fail: (%s)", err)
		return
	}

	// 支付
	ret = cfca.ProcessBindingPayment(be)

	// 处理结果
	trans.ChanRespCode = ret.ChanRespCode
	trans.RespCode = ret.RespCode
	switch ret.RespCode {
	case "000000":
		trans.TransStatus = 1
	case "000009":
		trans.TransStatus = 2
	default:
		trans.TransStatus = 3
	}
	if err = mongo.TransColl.Update(trans); err != nil {
		g.Error("update trans status fail ", err)
	}
	return ret
}

// ProcessBindingReomve 绑定解除
func ProcessBindingReomve(br *model.BindingRemove) (ret *model.BindingReturn) {
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 本地查询绑定关系
	bm, err := mongo.BindingMapColl.Find(br.MerId, br.BindingId)
	if err != nil {
		g.Error("'FindBindingRelation' error: ", err.Error())
		return model.NewBindingReturn("200101", "绑定ID不正确")
	}

	// 如果绑定状态非成功状态
	if bm.BindingStatus != "000000" {
		return mongo.RespCodeColl.Get(bm.BindingStatus)
	}

	// 查找渠道商户信息，获取证书
	chanMer := &model.ChanMer{
		ChanCode:  bm.ChanCode,
		ChanMerId: bm.ChanMerId,
	}
	if err = mongo.FindChanMer(chanMer); err != nil {
		g.Debug("not found any chanMer (%s)", err)
		return model.NewBindingReturn("300030", "无此交易权限")
	}

	// 转换关系，补充信息
	br.ChanMerId = bm.ChanMerId
	br.ChanBindingId = bm.ChanBindingId
	br.TxSNUnBinding = tools.SerialNumber()
	br.SignCert = chanMer.SignCert

	// 到渠道解绑
	ret = cfca.ProcessBindingRemove(br)

	// 如果解绑成功，更新本地数据库
	if ret.RespCode == "000000" {
		bm.BindingStatus = "100050"
		if err := mongo.BindingMapColl.Update(bm); err != nil {
			g.Error("'Update BindingMap' error: ", err.Error())
		}
	}

	// todo 成功后是否要入库

	return ret
}

// ProcessBindingRefund 退款
func ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {

	// default
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 是否有该订单号
	t, err := mongo.TransColl.Find(be.MerId, be.OrigOrderNum, 1)
	if err != nil {
		return model.NewBindingReturn("100020", "原交易不成功，不能退款")
	}

	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
	if err != nil {
		g.Error("not found any chanMer: ", err)
		// TODO 找不到渠道商户的错误码
		return model.NewBindingReturn("-100000", "找不到渠道商户")
	}

	// 赋值
	be.ChanMerId = t.ChanMerId
	be.ChanOrderNum = tools.SerialNumber()
	be.ChanOrigOrderNum = t.ChanOrderNum
	be.SignCert = chanMer.SignCert

	// 记录这笔退款
	trans := &model.Trans{
		OrderNum:       be.MerOrderNum,
		ChanOrderNum:   be.ChanOrderNum,
		ChanBindingId:  t.ChanBindingId,
		RefundOrderNum: be.ChanOrigOrderNum,
		AcctNum:        t.AcctNum,
		MerId:          be.MerId,
		TransAmount:    be.TransAmt,
		ChanMerId:      be.ChanMerId,
		ChanCode:       t.ChanCode,
		TransType:      2,
	}
	if err = mongo.TransColl.Add(trans); err != nil {
		g.Error("add trans fail : (%s)", err)
		return
	}

	// 退款
	ret = cfca.ProcessBindingRefund(be)

	// 更新结果
	trans.ChanRespCode = ret.ChanRespCode
	trans.RespCode = ret.RespCode
	switch ret.RespCode {
	case "000000":
		trans.TransStatus = 1
	case "000009":
		trans.TransStatus = 2
	default:
		trans.TransStatus = 3
	}
	if err = mongo.TransColl.Update(trans); err != nil {
		g.Error("update trans status fail : (%s)", err)
	}
	return
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	g.Info("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
