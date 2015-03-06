package core

import (
	"quickpay/channel"
	"quickpay/model"
)

// CreateBinding 绑卡
func CreateBinding(in *model.BindingCreateIn) (out *model.BindingCreateOut) {

	// 路由
	// 风控
	// 这部分可以提出一个公共模块

	// 判断卡bin ，决定走哪个渠道
	c := channel.GetBindingpayChannel("chinapayment").(*channel.ChinaPayment)

	out = c.CreateBinding(in)

	return out
}
