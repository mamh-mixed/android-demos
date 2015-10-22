package app

import "net/http"

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

	// 地推工具api
	mux.HandleFunc("/app/tools/login", CompanyLogin)
	mux.HandleFunc("/app/tools/users", UserList)
	mux.HandleFunc("/app/tools/register", UserRegister)
	mux.HandleFunc("/app/tools/uploadToken", GetQiniuToken)
	mux.HandleFunc("/app/tools/update", UpdateUserInfo)

	return mux
}
