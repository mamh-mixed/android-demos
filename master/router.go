package master

import (
	"net/http"
	"time"

	"github.com/CardInfoLink/quickpay/util"
)

// Route 后台管理的请求统一入口
func Route() (mux *MyServeMux) {
	mux = NewMyServeMux()

	mux.HandleFunc("/master/trade/query", tradeQueryHandle)
	mux.HandleFunc("/master/trade/report", tradeReportHandle)
	mux.HandleFunc("/master/trade/stat", tradeQueryStatsHandle)
	mux.HandleFunc("/master/trade/stat/report", tradeQueryStatsReportHandle)
	mux.HandleFunc("/master/merchant/find", merchantFindHandle)
	mux.HandleFunc("/master/merchant/one", merchantFindOneHandle)
	mux.HandleFunc("/master/merchant/save", merchantSaveHandle)
	mux.HandleFunc("/master/merchant/delete", merchantDeleteHandle)
	mux.HandleFunc("/master/router/save", routerSaveHandle)
	mux.HandleFunc("/master/router/find", routerFindHandle)
	mux.HandleFunc("/master/router/one", routerFindOneHandle)
	mux.HandleFunc("/master/router/delete", routerDeleteHandle)
	mux.HandleFunc("/master/channelMerchant/find", channelMerchantFindHandle)
	mux.HandleFunc("/master/channelMerchant/match", channelMerchantMatchHandle)
	mux.HandleFunc("/master/channelMerchant/findByMerIdAndCardBrand", channelFindByMerIdAndCardBrandHandle)
	mux.HandleFunc("/master/channelMerchant/save", channelMerchantSaveHandle)
	mux.HandleFunc("/master/channelMerchant/delete", channelMerchantDeleteHandle)
	mux.HandleFunc("/master/agent/find", agentFindHandle)
	mux.HandleFunc("/master/agent/delete", agentDeleteHandle)
	mux.HandleFunc("/master/agent/save", agentSaveHandle)
	mux.HandleFunc("/master/group/find", groupFindHandle)
	mux.HandleFunc("/master/group/delete", groupDeleteHandle)
	mux.HandleFunc("/master/group/save", groupSaveHandle)
	mux.HandleFunc("/master/qiniu/uptoken", handleUptoken)
	mux.HandleFunc("/master/qiniu/uploaded", handleDownURL)
	mux.HandleFunc("/master/user/find", userFindHandle)
	mux.HandleFunc("/master/user/create", userCreateHandle)

	return mux
}

// MyServeMux 权限拦截器
type MyServeMux struct {
	http.ServeMux
}

// NewMyServeMux allocates and returns a new ServeMux.
func NewMyServeMux() *MyServeMux {
	return &MyServeMux{*http.NewServeMux()}
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *MyServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// c, err := r.Cookie("QUICKMASTERID")
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		// 未登录
	// 		// if uri == "login" { do login}
	// 		// else { error }
	// 		// return
	// 	}
	// 	http.Error(w, err.Error(), http.StatusNotAcceptable)
	// 	return
	// }
	// if c != nil {
	// 	log.Debugf("url=%s, cookie: %s", r.URL, c.String())
	// 	// 验证是否有权限
	// }删除删除

	h, _ := mux.Handler(r)

	// loginCtrl.login 里有如下代码，登录成功后写 Cookie
	if r.URL.Path == "/master/agent/find" {
		http.SetCookie(w, &http.Cookie{
			Name:    "QUICKMASTERID",
			Value:   util.SerialNumber(),
			Path:    "/master/",
			Expires: time.Now().Add(20 * time.Minute),
		})
	}
	h.ServeHTTP(w, r)
}
