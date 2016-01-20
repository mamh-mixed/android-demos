package channel

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"time"
)

var mockClient MockBindingPay
var d = 500 * time.Millisecond

// MockBindingPay 用于测试核心逻辑
type MockBindingPay struct {
}

// ProcessPaySettlement 支付结算
func (m *MockBindingPay) ProcessPaySettlement(be *model.PaySettlement) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingCreate 建立绑定关系
func (m *MockBindingPay) ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingRemove 解除绑定关系
func (m *MockBindingPay) ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingEnquiry 查询绑定关系
func (m *MockBindingPay) ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingPayment 快捷支付
func (m *MockBindingPay) ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessPaymentEnquiry 快捷支付查询
func (m *MockBindingPay) ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessBindingRefund 快捷支付退款
func (m *MockBindingPay) ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessRefundEnquiry 快捷支付退款查询
func (m *MockBindingPay) ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessSendBindingPaySMS 快捷支付发送短信验证码
func (m *MockBindingPay) ProcessSendBindingPaySMS(be *model.BindingPayment) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ProcessPaymentWithSMS 快捷支付短信验证支付
func (m *MockBindingPay) ProcessPaymentWithSMS(be *model.BindingPayment) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// Consume 模拟一个消费(无卡直接支付)的处理。
func (m *MockBindingPay) Consume(be *model.NoTrackPayment) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}

// ConsumeByApplePay 模拟一个Apple pay消费的处理。
func (m *MockBindingPay) ConsumeByApplePay(be *model.ApplePay) (ret *model.BindingReturn) {
	// TODO validate feilds
	log.Infof("MOCK - %#v", be)
	time.Sleep(d)

	return model.NewBindingReturn("000000", "请求处理成功")
}
