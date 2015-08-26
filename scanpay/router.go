package scanpay

import (
	"net/http"

	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/channel/weixin"
)

// Route 后台管理的请求统一入口
func Route() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/scanpay/unified", scanpayUnifiedHandle)
	mux.HandleFunc(weixin.NotifyURL, weixinNotifyHandle)
	mux.HandleFunc(alipay.NotifyURL, alipayNotifyHandle)

	return mux
}
