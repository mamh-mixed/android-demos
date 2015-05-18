package channel

import (
	"github.com/CardInfoLink/quickpay/channel/cil"
	"github.com/CardInfoLink/quickpay/model"
)

// DirectPayChan 表示无卡直接支付、apple pay等统一的渠道接口
type DirectPayChan interface {
	// Consume 直接消费（订购消费）
	Consume(be *model.NoTrackPayment) (ret *model.BindingReturn)

	// ConsumeByApplePay ApplePay 消费
	ConsumeByApplePay(ap *model.ApplePay) (ret *model.BindingReturn)
}

// GetChan 根据chanCode获得渠道对象
func GetDirectPayChan(chanCode string) DirectPayChan {

	switch chanCode {
	case "Mock":
		return &mockClient
	case "CIL":
		return &cil.DefaultCILPayClient
	}
	return nil
}
