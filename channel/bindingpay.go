package channel

import (
	"quickpay/model"
)

var channelFactory map[string]Channel

func init() {

	// ChinaPay
	chinaPay := Chinapay{"https://test.china-clearing.com/Gateway/InterfaceII"}

	channelFactory["ChinaPay"] = &chinaPay

}

// Channel 绑定支付接口
type Channel interface {

	//建立绑定关系
	CreateBinding(data *model.BindingCreateIn) *model.BindingCreateOut

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

// GetChannel 根据渠道 Id 取渠道对象
func GetChannel(c string) Channel {
	return channelFactory[c]
}
