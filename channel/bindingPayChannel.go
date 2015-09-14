package channel

import (
	"github.com/CardInfoLink/quickpay/channel/cfca"
	"github.com/CardInfoLink/quickpay/model"
)

// BindingPayChan 渠道对象统一接口
type BindingPayChan interface {
	// ProcessBindingEnquiry 查询绑定关系
	ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn)

	// ProcessBindingCreate 建立绑定关系
	ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn)

	// ProcessBindingRemove 解除绑定关系
	ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn)

	// ProcessBindingPayment 快捷支付
	ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn)

	// ProcessPaymentEnquiry 快捷支付查询
	ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn)

	// ProcessBindingRefund 快捷支付退款
	ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn)

	// ProcessRefundEnquiry 快捷支付退款查询
	ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn)

	// ProcessSendBindingPaySMS 快捷支付发送短信验证码
	ProcessSendBindingPaySMS(be *model.BindingPayment) (ret *model.BindingReturn)

	// ProcessPaymentWithSMS 快捷支付短信验证支付
	ProcessPaymentWithSMS(be *model.BindingPayment) (ret *model.BindingReturn)

	// ProcessPaySettlement 支付结算
	ProcessPaySettlement(be *model.PaySettlement) (ret *model.BindingReturn)

	// ProcessTransChecking 交易对账，清算
	// ProcessTransChecking(chanMerId, settDate, signCert string) (resp *BindingResponse)
}

// GetChan 根据chanCode获得渠道对象
func GetChan(chanCode string) BindingPayChan {

	switch chanCode {
	case "Mock":
		return &mockClient
	case "CFCA":
		return &cfca.DefaultClient
	}
	return nil
}
