package core

import (
	"strings"

	"github.com/CardInfoLink/quickpay/channel/cfca"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"

	"github.com/omigo/g"
)

// ProcessBindingCreate 绑定建立的业务处理
func ProcessBindingCreate(bc *model.BindingCreate) (ret *model.BindingReturn) {
	// 验证该机构下，该绑定号是否已经绑定了
	count, err := mongo.BindingMapColl.Count(bc.MerId, bc.BindingId)
	if count > 0 {
		return model.NewBindingReturn("200071", "绑定ID重复")
	}

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
		return model.NewBindingReturn("000001", "系统内部错误")
	}

	// bc(BindingCreate)用来向渠道发送请求，增加一些渠道要求的数据。
	bc.ChanMerId = rp.ChanMerId
	bc.ChanBindingId = bm.ChanBindingId
	bc.SignCert = chanMer.SignCert
	g.Trace("'BindingCreate' is: %+v", bc)
	// todo 根据路由策略里面不同的渠道调用不同的绑定接口，这里为了简单，调用中金的接口。
	ret = cfca.ProcessBindingCreate(bc)

	// 渠道返回后，根据应答码，判断绑定是否成功，如果成功，更新绑定关系映射，绑定关系生效。
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
		g.Error("'FindBindingMap' error: %s", err.Error())
		return model.NewBindingReturn("200101", "绑定ID不正确")
	}
	g.Debug("binding result: %#v", bm)

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
		g.Info("'UpdateBindingMap' is: %+v", bm)
		g.Error("'UpdateBindingRelation' error: ", err.Error())
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
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum)
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
	if bm.BindingStatus != model.BindingSuccess {
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
	be.ChanOrderNum = tools.SerialNumber()
	be.SignCert = chanMer.SignCert

	// 记录这笔交易
	trans := &model.Trans{
		OrderNum:      be.MerOrderNum,
		ChanOrderNum:  be.ChanOrderNum,
		ChanBindingId: be.ChanBindingId,
		AcctNum:       bi.AcctNum,
		MerId:         be.MerId,
		TransAmt:      be.TransAmt,
		ChanMerId:     be.ChanMerId,
		ChanCode:      bm.ChanCode,
		TransType:     model.PayTrans, //支付
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
		trans.TransStatus = model.TransSuccess
	case "000009":
		trans.TransStatus = model.TransHandling
	default:
		trans.TransStatus = model.TransFail
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
	if bm.BindingStatus != model.BindingSuccess {
		return mongo.RespCodeColl.Get(bm.BindingStatus)
	}

	// 查找渠道商户信息，获取证书
	chanMer, err := mongo.ChanMerColl.Find(bm.ChanCode, bm.ChanMerId)
	if err != nil {
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
		bm.BindingStatus = model.BindingRemoved
		if err := mongo.BindingMapColl.Update(bm); err != nil {
			g.Error("'Update BindingMap' error: ", err.Error())
		}
	}

	return ret
}

// ProcessBindingRefund 退款
func ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {

	// default
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 检查同一个商户的订单号是否重复
	count, err := mongo.TransColl.Count(be.MerId, be.MerOrderNum)
	if err != nil {
		g.Error("find trans fail : (%s)", err)
		return
	}
	if count > 0 {
		return model.NewBindingReturn("200081", "订单号重复")
	}

	// 是否有该源订单号
	orign, err := mongo.TransColl.Find(be.MerId, be.OrigOrderNum)
	switch {
	// 不存在原交易
	case err != nil:
		return model.NewBindingReturn("100020", "原交易不成功，不能退款")
	// 已退款
	case orign.RefundStatus == model.TransRefunded:
		return model.NewBindingReturn("100010", "该笔订单已经存在退款交易，不能再次退款")
	// 退款金额过大
	case be.TransAmt > orign.TransAmt:
		return model.NewBindingReturn("200191", "退款金额（累计）大于可退金额")
	}

	// 获得渠道商户
	chanMer, err := mongo.ChanMerColl.Find(orign.ChanCode, orign.ChanMerId)
	if err != nil {
		g.Error("not found any chanMer: ", err)
		return model.NewBindingReturn("300030", "无此交易权限")
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
		g.Error("add refund trans fail : (%s)", err)
		return
	}

	// 退款
	ret = cfca.ProcessBindingRefund(be)

	// 更新结果
	refund.ChanRespCode = ret.ChanRespCode
	refund.RespCode = ret.RespCode
	switch ret.RespCode {
	case "000000":
		refund.TransStatus = model.TransSuccess
		//更新原交易状态
		orign.RefundStatus = model.TransRefunded
		if err = mongo.TransColl.Update(orign); err != nil {
			g.Error("update orign trans RefundStatus fail : (%s)", err)
		}
	//只有超时才会出现000009
	case "000009":
		refund.TransStatus = model.TransHandling
	default:
		refund.TransStatus = model.TransFail
	}
	if err = mongo.TransColl.Update(refund); err != nil {
		g.Error("update refund trans status fail : (%s)", err)
	}
	return
}

// ProcessOrderEnquiry 订单查询
func ProcessOrderEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	// 默认返回成功的应答码
	ret = &model.BindingReturn{
		RespCode: "000000",
		RespMsg:  "success",
	}

	// 是否有该订单号
	t, err := mongo.TransColl.Find(be.MerId, be.OrigOrderNum)
	if err != nil {
		return model.NewBindingReturn("200082", "订单号不存在")
	}
	g.Debug("trans:(%+v)", t)
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
		g.Error("not found any chanMer: ", err)
		return model.NewBindingReturn("300030", "无此交易权限")
	}

	//赋值
	be.SignCert = chanMer.SignCert
	be.ChanMerId = chanMer.ChanMerId
	be.ChanOrderNum = t.ChanOrderNum

	// 原订单为处理中，向渠道发起查询
	result := new(model.BindingReturn)
	switch t.TransType {
	//支付
	case model.PayTrans:
		result = cfca.ProcessPaymentEnquiry(be)
	//退款
	case model.RefundTrans:
		result = cfca.ProcessRefundEnquiry(be)
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
		g.Error("Update trans error : %s", err)
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
	ret = model.NewBindingReturn("000001", "系统内部错误")

	//查询
	rec, err := mongo.TransSettColl.Find(be.MerId, be.SettDate, be.NextOrderNum)
	if err != nil {
		g.Error("Find transSett records error : %s", err)
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
	ret.RespCode = "000000"
	ret.RespMsg = "success"
	ret.Count = len(ret.Rec)
	return
}

// ProcessBillingSummary 交易对账汇总查询
func ProcessBillingSummary(be *model.BillingSummary) (ret *model.BindingReturn) {

	//default return
	ret = model.NewBindingReturn("000001", "系统内部错误")

	//查询
	data, err := mongo.TransSettColl.Summary(be.MerId, be.SettDate)
	if err != nil {
		g.Error("summary transSett records error : %s", err)
		return
	}

	//赋值
	ret.RespCode = "000000"
	ret.RespMsg = "success"
	ret.SettDate = be.SettDate
	ret.Data = data
	return
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	g.Info("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
