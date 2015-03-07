package channel

import (
	"quickpay/model"
)

// BindingPaymentChannel 绑定支付接口
type BindingPaymentChannel interface {

	//查询绑定关系
	ProcessBindingEnquiry(be *model.BindingEnquiry) *model.BindingReturn
}

// GetBindingPaymentChannel 根据渠道 Id 取渠道对象
func GetBindingPaymentChannel(c string) BindingPaymentChannel {
	return &ChinaPayment{}
}
