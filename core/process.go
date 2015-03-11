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
	be.SettlementFlag = chanMer.SettlementFlag
	ret = cfca.ProcessBindingPayment(be)

	// post process

	return ret
}
