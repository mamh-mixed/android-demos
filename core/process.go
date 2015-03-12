package core

import (
	"quickpay/channel/cfca"
	"quickpay/model"
	"quickpay/mongo"

	"github.com/omigo/g"
	"strings"
)

// 绑定建立的业务处理
func ProcessBindingCreate(bc *model.BindingCreate) (ret *model.BindingReturn) {
	// todo 如果需要校验短信，验证短信
	ret = validateSmsCode(bc.SendSmsId, bc.SmsCode)
	if ret != nil {
		return ret
	}
	// 获取卡属性
	cardBin := mongo.FindCardBin(bc.AcctNum)
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
	br := &mongo.BindingRelation{
		CardInfo: *bc,
		Router:   *routerPolicy,
	}
	if err := mongo.InsertOneBindingRelation(br); err != nil {
		// todo 插入绑定关系失败的错误码
		return model.NewBindingReturn("", err.Error())
	}
	// todo 根据路由策略里面不同的渠道调用不同的绑定接口，这里为了简单，调用中金的接口
	ret = cfca.ProcessBindingCreate(bc)
	// 如果返回成功，则更新数据库，将返回的绑定ID存库
	if ret.RespCode != "000000" {
		return ret
	}

	br.ChannelBindingId = ret.BindingId
	err = mongo.UpdateOneBindingRelation(br)
	if err != nil {
		return model.NewBindingReturn("", err.Error())
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

func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 本地查询绑定关系
	// merId = ass.merId
	// 根据绑定关系得到渠道商户信息
	chanMer := mongo.ChanMer{
		ChanCode:  "",
		ChanMerId: "",
	}
	err := chanMer.Init()
	if err != nil {
		//not found
		//return
	}
	// 记录这笔交易

	be.SettlementFlag = chanMer.SettlementFlag
	ret = cfca.ProcessBindingPayment(be)

	// post process
	return ret
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	g.Info("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
