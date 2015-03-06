package channel

import (
	"quickpay/model"
)

type Bindingpay interface {

	//建立绑定关系
	CreateBinding() *model.ChannelRes

	//查询绑定关系
	QueryBinding() *model.ChannelRes

	//快捷支付
	QuickPay() *model.ChannelRes

	//快捷支付查询
	QuickPayQuery() *model.ChannelRes

	//快捷支付退款
	QuickPayRefund() *model.ChannelRes

	//快捷支付退款查询
	QuickPayRefundQuery() *model.ChannelRes

	//交易对账单
	TradePayments() *model.ChannelRes
}
