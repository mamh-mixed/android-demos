package core

import (
	"quickpay/types"
)

// CreateBinding 绑卡
func CreateBinding(in *types.BindingCreateIn) (out *types.BindingCreateOut) {

	// 路由
	// 风控
	// 这部分可以提出一个公共模块

	// 判断卡bin ，决定走哪个渠道
	// channel = newChannel()

	// out = chinapay.CreateBinding(in)

	return
}
