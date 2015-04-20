package channel

import (
	// "github.com/CardInfoLink/quickpay/channel/alp"
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
	case "CIL":
		return nil
	case "ALP":
		return nil
	case "WXP":
		return nil
	}
	return nil
}
