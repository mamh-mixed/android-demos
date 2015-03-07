package core

import (
	"quickpay/channel/cfca"
	"quickpay/model"
)

// CreateBinding 绑卡
func CreateBinding(in *model.BindingCreateIn) (out *model.BindingCreateOut) {

	// 建立绑定关系

	// 路由
	// 风控
	// 这部分可以提出一个公共模块

	// 判断卡bin ，决定走哪个渠道

	return out
}

// ProcessBindingEnquiry 绑定关系查询
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {

	// 本地查询绑定关系
	// merId = ass.merId

	ret = cfca.ProcessBindingEnquiry(be)

	// post process

	return ret
}
