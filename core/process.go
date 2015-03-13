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
	rp := mongo.FindRouterPolicy(bc.MerId, cardBin.CardBrand)
	if rp == nil {
		// todo 错误返回校验码
		return model.NewBindingReturn("", "找不到路由策略")
	}
	// 根据商户、卡号、绑定Id、渠道、渠道商户生成一个系统绑定Id(ChanBindingId)，并将这些关系入库
	bc.SendSmsId = ""
	bc.SmsCode = ""
	// br(BindingRelation)用来入库
	br := &model.BindingRelation{
		BindingId:     bc.BindingId,
		MerId:         bc.MerId,
		AcctName:      bc.AcctName,
		AcctNum:       bc.AcctNum,
		IdentType:     bc.IdentType,
		IdentNum:      bc.IdentNum,
		PhoneNum:      bc.PhoneNum,
		AcctType:      bc.AcctType,
		ValidDate:     bc.ValidDate,
		Cvv2:          bc.Cvv2,
		BankId:        bc.BankId,
		CardBrand:     rp.CardBrand,
		ChanCode:      rp.ChanCode,
		ChanMerId:     rp.ChanMerId,
		SysBindingId:  tools.SerialNumber(),
		BindingStatus: "",
	}
	// bc(BindingCreate)用来向渠道发送请求，增加一些渠道要求的数据。
	bc.ChanMerId = rp.ChanMerId
	bc.ChanBindingId = br.SysBindingId
	g.Info("'BindingCreate' is: %+v", bc)
	g.Info("'BindingRelation' is: %+v", br)
	if err := mongo.InsertBindingRelation(br); err != nil {
		// todo 插入绑定关系失败的错误码
		return model.NewBindingReturn("-100000", err.Error())
	}
	// todo 根据路由策略里面不同的渠道调用不同的绑定接口，这里为了简单，调用中金的接口。
	ret = cfca.ProcessBindingCreate(bc)
	// 渠道返回后，根据应答码，判断绑定是否成功，如果成功，更新数据库，绑定关系生效。
	br.BindingStatus = ret.RespCode
	err := mongo.UpdateBindingRelation(br)
	if err != nil {
		// todo 更新数据库错误码
		return model.NewBindingReturn("-100000", err.Error())
	}
	return ret
}

// ProcessBindingEnquiry 绑定关系查询
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 默认返回
	ret = model.NewBindingReturn("000001", "系统内部错误")
	// 本地查询绑定关系
	bindRelation, err := mongo.FindBindingRelation(be.MerId, be.BindingId)
	if err != nil {
		//TODO返回什么应答码
		g.Debug("not found any bindRelation (%s)", err)
		return
	}
	// 转换绑定关系、请求
	be.ChanMerId = bindRelation.ChanMerId
	be.ChanBindingId = bindRelation.SysBindingId
	ret = cfca.ProcessBindingEnquiry(be)

	return ret
}

// ProcessBindingPayment 绑定支付
func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 默认返回
	ret = model.NewBindingReturn("000001", "系统内部错误")
	// 本地查询绑定关系
	bindRelation, err := mongo.FindBindingRelation(be.MerId, be.BindingId)
	if err != nil {
		g.Debug("not found any bindRelation (%s)", err)
		return
	}
	// 根据绑定关系得到渠道商户信息
	chanMer := &model.ChanMer{
		ChanCode:  bindRelation.ChanCode,
		ChanMerId: bindRelation.ChanMerId,
	}
	if err = mongo.FindChanMer(chanMer); err != nil {
		g.Debug("not found any chanMer (%s)", err)
		return
	}
	// 记录这笔交易
	trans := &model.Trans{
		Chan:    *chanMer,
		Payment: *be,
	}
	if err = mongo.AddTrans(trans); err != nil {
		g.Debug("add trans fail  (%s)", err)
		return
	}
	be.SettFlag = chanMer.SettFlag
	be.ChanBindingId = bindRelation.SysBindingId
	be.ChanMerId = bindRelation.ChanMerId
	// 支付

	ret = cfca.ProcessBindingPayment(be)

	// 处理结果
	if ret.RespCode == "000000" {
		trans.Flag = 1
		// 不关心是否更新成功
		mongo.ModifyTrans(trans)
	}

	return
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	g.Info("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
