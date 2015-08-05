package master

import (
	"net/http"

	"github.com/omigo/log"
)

// MasterRoute 后台管理的请求统一入口
func MasterRoute(w http.ResponseWriter, r *http.Request) {
	log.Infof("url = %s", r.URL.String())

	switch r.URL.Path {
	case "/master/trade/report":
		tradeReport(w, r)
		return
	default:
		w.WriteHeader(404)
	}

}
