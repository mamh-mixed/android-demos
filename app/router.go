package app

import (
	"net/http"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
)

// Route app请求统一入口
func Route() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/app/register", registerHandle)
	mux.HandleFunc("/app/login", loginHandle)
	mux.HandleFunc("/app/request_activate", reqActivateHandle)
	mux.HandleFunc("/app/activate", activateHandle)
	mux.HandleFunc("/app/improveinfo", improveInfoHandle)
	mux.HandleFunc("/app/getOrder", getOrderHandle)
	mux.HandleFunc("/app/bill", billHandle)
	mux.HandleFunc("/app/getTotal", getTotalHandle)
	mux.HandleFunc("/app/getrefd", getRefdHandle)
	mux.HandleFunc("/app/updatepassword", passwordHandle)
	mux.HandleFunc("/app/limitincrease", promoteLimitHandle)
	mux.HandleFunc("/app/updateinfo", updateSettInfoHandle)
	mux.HandleFunc("/app/getinfo", getSettInfoHandle)
	mux.HandleFunc("/app/ticket", ticketHandle)
	mux.HandleFunc("/app/findOrder", findOrderHandle)
	mux.HandleFunc("/app/updateMessage", updateMessageHandle)
	mux.HandleFunc("/app/forgetpassword", forgetPasswordHandle)
	mux.HandleFunc("/app/getQiniuToken", getQiniuTokenHandle)
	mux.HandleFunc("/app/improveCertInfo", improveCertInfoHandle)
	mux.HandleFunc("/app/pullinfo", pullInfoHandle)

	// 地推工具api
	mux.HandleFunc("/app/tools/login", CompanyLogin)
	mux.HandleFunc("/app/tools/users", UserList)
	mux.HandleFunc("/app/tools/register", UserRegister)
	mux.HandleFunc("/app/tools/uploadToken", GetQiniuToken)
	mux.HandleFunc("/app/tools/update", UpdateUserInfo)
	mux.HandleFunc("/app/tools/activate", UserActivate)
	mux.HandleFunc("/app/tools/download", GetDownloadUrl)

	return mux
}

// RouteV3 APPv3请求统一入口
func RouteV3() (mux *AppV3ServeMux) {
	mux = NewAppV3ServeMux()

	// app3.0接口
	mux.HandleFunc("/app/v3/bill", billV3Handle)

	return mux
}

// AppV3ServeMux APPv3拦截器
type AppV3ServeMux struct {
	http.ServeMux
}

// NewAppV3ServeMux allocates and returns a new ServeMux.
func NewAppV3ServeMux() *AppV3ServeMux {
	return &AppV3ServeMux{*http.NewServeMux()}
}

func (mux *AppV3ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("**********%s************", "<<<<<<>>>>>>")
	log.Debugf("request url is %s; sign is %s", r.URL.Path, r.FormValue("sign"))
	log.Debugf("username is %s", r.FormValue("username"))

	// 验签
	if !checkSignSha256(r) {
		// 验签失败
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
