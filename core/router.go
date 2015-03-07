package core

import (
	"quickpay/channel"
)

// GetChannelByCardBin 根据卡 Bin 决定交易路由
func GetChannelByCardBin(card string) (c channel.BindingPaymentChannel) {

	c = channel.GetBindingPaymentChannel("chinapayment")
	return c
}

// GetChannelByBindingId 根据绑定关系决定交易路由
func GetChannelByBindingId(bid string) (c channel.BindingPaymentChannel) {

	c = channel.GetBindingPaymentChannel("chinapayment")
	return c
}
