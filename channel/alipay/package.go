// 包含dmf1.0,dmf2.0版本
package alipay

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1/domestic"
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1/oversea"
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay2"
)

var (
	NotifyUrl = domestic.NotifyPath
	Oversea   = &oversea.DefaultClient
	Domestic  = &domestic.DefaultClient
	Alipay2   = &scanpay2.DefaultClient
)
