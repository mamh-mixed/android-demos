package core

import (
	"quickpay/channel/cfca"
	"quickpay/model"
)

// ProcessBindingEnquiry 绑定关系查询
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {

	// 本地查询绑定关系
	// merId = ass.merId

	ret = cfca.ProcessBindingEnquiry(be)

	// post process

	return ret
}
