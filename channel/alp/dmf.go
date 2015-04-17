package alp

import (
	"github.com/CardInfoLink/quickpay/model"
)

var Obj alp

// alp 当面付，扫码支付
type alp struct{}

// ProcessBindingEnquiry 无需实现
func (a *alp) ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	return
}

// ProcessBindingCreate 无需实现
func (a *alp) ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn) {
	return
}

// ProcessBindingRemove 无需实现
func (a *alp) ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn) {
	return
}

// ProcessBindingPayment 下单
func (a *alp) ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	return
}

// ProcessPaymentEnquiry 支付查询
func (a *alp) ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	return a.processEnquiry(be)
}

// ProcessBindingRefund 退款
func (a *alp) ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {
	return
}

// ProcessRefundEnquiry 退款查询
func (a *alp) ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	return a.processEnquiry(be)
}

// processEnquiry 查询，包含支付、退款
func (a *alp) processEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	return
}
