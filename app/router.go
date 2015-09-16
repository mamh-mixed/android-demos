package app

import "net/http"

// Route app请求统一入口
func Route() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/app/register", registerHandle)

	return mux
}
