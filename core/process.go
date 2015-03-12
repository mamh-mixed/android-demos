package core

import (
	"quickpay/channel/cfca"
	"quickpay/model"
	"quickpay/mongo"
)

// ProcessBindingEnquiry 绑定关系查询
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {

	// 本地查询绑定关系
	// merId = ass.merId

	ret = cfca.ProcessBindingEnquiry(be)

	// post process

	return ret
}

func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 默认返回
	ret = &model.BindingReturn{
		RespCode: "000001",
		RespMsg:  "系统错误",
	}
	// 本地查询绑定关系
	bindRelation := mongo.FindOneBindingRelationByMerCodeAndBindingId(be.InstitutionId, be.BindingId)
	// 根据绑定关系得到渠道商户信息
	chanMer := mongo.ChanMer{
		ChanCode:  bindRelation.Router.ChannelCode,
		ChanMerId: bindRelation.Router.ChannelMerCode,
	}
	err := chanMer.Init()
	if err != nil {
		//not found
		return
	}
	// 记录这笔交易
	trans := mongo.Trans{
		Chan:    chanMer,
		Payment: be,
	}
	err = trans.Add()
	if err != nil {
		// 添加操作发生错误
		return
	}
	be.SettlementFlag = chanMer.SettlementFlag
	be.BindingId = bindRelation.ChannelBindingId
	be.InstitutionId = bindRelation.Router.ChannelMerCode
	ret = cfca.ProcessBindingPayment(be)

	// 处理结果
	if ret.RespCode == "000000" {
		trans.Flag = 1
		// 不关心是否更新成功
		trans.Modify()
	}

	return
}
