package channel

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

var mockClient MockBindingPay

// MockBindingPay 用于测试核心逻辑
type MockBindingPay struct {
}

// ProcessBindingCreate 建立绑定关系
func (m *MockBindingPay) ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingRemove 解除绑定关系
func (m *MockBindingPay) ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingEnquiry 查询绑定关系
func (m *MockBindingPay) ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingPayment 快捷支付
func (m *MockBindingPay) ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessPaymentEnquiry 快捷支付查询
func (m *MockBindingPay) ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingRefund 快捷支付退款
func (m *MockBindingPay) ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessRefundEnquiry 快捷支付退款查询
func (m *MockBindingPay) ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)

	return model.NewBindingReturn("000000", "请求处理成功")
}
