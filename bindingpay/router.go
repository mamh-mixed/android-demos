package bindingpay

import "net/http"

// Route 后台管理的请求统一入口
func Route() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/bindingpay/", BindingpayHandle)

	return mux
}
