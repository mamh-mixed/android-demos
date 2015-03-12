package core

import (
	"quickpay/channel/cfca"
	"quickpay/model"
	"quickpay/mongo"
	"strings"

	"github.com/omigo/g"
)

// ProcessBindingCreate 绑定建立的业务处理
func ProcessBindingCreate(bc *model.BindingCreate) (ret *model.BindingReturn) {
	// todo 如果需要校验短信，验证短信
	ret = validateSmsCode(bc.SendSmsId, bc.SmsCode)
	if ret != nil {
		return ret
	}
	// 获取卡属性
	cardBin := mongo.FindCardBin(bc.AcctNum)

	g.Debug("CardBin: %+v", cardBin)
	// 如果是银联卡，验证证件信息
	if strings.EqualFold("CUP", cardBin.CardBrand) || strings.EqualFold("UPI", cardBin.CardBrand) {
		ret = UnionPayCardValidity(bc)
		if ret != nil {
			return ret
		}
	}
	// 通过路由策略找到渠道和渠道商户
	routerPolicy, err := mongo.FindRouter(bc.MerId, cardBin.CardBrand)
	if err != nil {
		// todo 错误返回校验码
		return model.NewBindingReturn("", "找不到路由策略")
	}
	// 根据商户、卡号、绑定Id、渠道、渠道商户生成一个系统绑定Id，并将这些关系入库
	bc.SendSmsId = ""
	bc.SmsCode = ""
	br := &model.BindingRelation{
		*bc,
		*routerPolicy,
		"",
	}
	if err := mongo.InsertBindingRelation(br); err != nil {
		// todo 插入绑定关系失败的错误码
		return model.NewBindingReturn("-100000", err.Error())
	}
	// todo 根据路由策略里面不同的渠道调用不同的绑定接口，这里为了简单，调用中金的接口
	ret = cfca.ProcessBindingCreate(bc)
	// 如果返回成功，则更新数据库，将返回的绑定ID存库
	if ret.RespCode != "000000" {
		return ret
	}

	br.ChanBindingId = ret.BindingId
	err = mongo.UpdateBindingRelation(br)
	if err != nil {
		// todo 更新数据库错误码
		return model.NewBindingReturn("-100000", err.Error())
	}
	return ret
}

// ProcessBindingEnquiry 绑定关系查询
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {

	// 本地查询绑定关系
	// merId = ass.merId

	ret = cfca.ProcessBindingEnquiry(be)

	// post process

	return ret
}

// ProcessBindingPayment 绑定支付
func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 默认返回
	ret = &model.BindingReturn{
		RespCode: "000001",
		RespMsg:  "系统错误",
	}
	// 本地查询绑定关系
	bindRelation, err := mongo.FindBindingRelation(be.MerId, be.BindingId)
	if err != nil {
		g.Debug("not found any bindRelation (%s)", err)
		return
	}
	// 根据绑定关系得到渠道商户信息
	chanMer := mongo.ChanMer{
		ChanCode:  bindRelation.ChanCode,
		ChanMerId: bindRelation.ChanMerId,
	}
	if err = chanMer.Init(); err != nil {
		g.Debug("not found any chanMer (%s)", err)
		return
	}
	// 记录这笔交易
	trans := mongo.Trans{
		Chan:    chanMer,
		Payment: *be,
	}
	if err = trans.Add(); err != nil {
		g.Debug("add trans fail  (%s)", err)
		return
	}
	be.SettlementFlag = chanMer.SettlementFlag
	be.BindingId = bindRelation.ChanBindingId
	be.MerId = bindRelation.ChanMerId
	ret = cfca.ProcessBindingPayment(be)

	// 处理结果
	if ret.RespCode == "000000" {
		trans.Flag = 1
		// 不关心是否更新成功
		trans.Modify()
	}

	return
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	g.Info("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
