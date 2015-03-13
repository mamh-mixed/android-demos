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
	// ret = validateSmsCode(bc.SendSmsId, bc.SmsCode)
	// if ret != nil {
	// 	return ret
	// }

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
	g.Info("'BindingRelation' is: %+v", br)
	// 绑定关系入库
	if err := mongo.InsertBindingRelation(br); err != nil {
		// todo 插入绑定关系失败的错误码
		return model.NewBindingReturn("-100000", err.Error())
	}

	// bc(BindingCreate)用来向渠道发送请求，增加一些渠道要求的数据。
	bc.ChanMerId = rp.ChanMerId
	bc.ChanBindingId = br.SysBindingId
	g.Info("'BindingCreate' is: %+v", bc)
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

// ProcessBindingEnquiry 绑定关系查询。
// 先到本地库去查找，如果本地库查找的结果是正在处理中，就到渠道查找；查找完更新到数据库中。
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 默认返回
	ret = model.NewBindingReturn("000001", "系统内部错误")

	// 本地查询绑定关系
	bindRelation, err := mongo.FindBindingRelation(be.MerId, be.BindingId)
	if err != nil {
		//TODO返回什么应答码
		g.Debug("not found any bindRelation (%s)", err)
		return ret
	}

	// 非处理中，直接返回结果
	if bindRelation.BindingStatus != "000009" {
		return mongo.GetRespCode(bindRelation.BindingStatus)
	}

	// 正在处理中，到渠道那边查找
	// 转换绑定关系、请求
	be.ChanMerId = bindRelation.ChanMerId
	be.ChanBindingId = bindRelation.SysBindingId
	// todo 查找该商户配置的渠道，这里为了简单，到中金查找。
	ret = cfca.ProcessBindingEnquiry(be)

	// 更新绑定关系的状态
	bindRelation.BindingStatus = ret.RespCode
	if err = mongo.UpdateBindingRelation(bindRelation); err != nil {
		// todo 更新数据库错误码
		return model.NewBindingReturn("-100000", err.Error())
	}

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
		return ret
	}

	// 根据绑定关系得到渠道商户信息
	chanMer := &model.ChanMer{
		ChanCode:  bindRelation.ChanCode,
		ChanMerId: bindRelation.ChanMerId,
	}
	if err = mongo.FindChanMer(chanMer); err != nil {
		g.Debug("not found any chanMer (%s)", err)
		return ret
	}

	// 记录这笔交易

	trans := &model.Trans{Payment: *be}
	if err = mongo.AddTrans(trans); err != nil {
		g.Debug("add trans fail  (%s)", err)
		return ret
	}

	// 支付
	be.SettFlag = chanMer.SettFlag
	be.ChanBindingId = bindRelation.SysBindingId
	be.ChanMerId = bindRelation.ChanMerId
	ret = cfca.ProcessBindingPayment(be)

	// 处理结果
	if ret.RespCode == "000000" {
		trans.TransFlag = 1
		if err = mongo.ModifyTrans(trans); err != nil {
			g.Error("update trans status fail ", err)
		}
	}

	return ret
}

// todo 校验短信验证码，短信验证通过就返回nil
func validateSmsCode(sendSmsId, smsCode string) (ret *model.BindingReturn) {
	g.Info("SendSmsId is: %s;SmsCode is: %s", sendSmsId, smsCode)
	return ret
}
