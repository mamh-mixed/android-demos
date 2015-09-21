package master

import (
	"errors"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
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
	mux.HandleFunc("/master/merchant/import", importMerchant)
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
	mux.HandleFunc("/master/subAgent/find", subAgentFindHandle)
	mux.HandleFunc("/master/subAgent/delete", subAgentDeleteHandle)
	mux.HandleFunc("/master/subAgent/save", subAgentSaveHandle)
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
	// 处理登陆
	if r.URL.Path == "/master/login" {
		loginHandle(w, r)
		return
	}
	// 查找session
	if r.URL.Path == "/master/session/find" {
		findSessionHandle(w, r)
		return
	}
	// 删除session
	if r.URL.Path == "/master/session/delete" {
		sessionDeleteHandle(w, r)
		return
	}

	// 查看请求中有没有cookie
	c, err := r.Cookie("QUICKMASTERID")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
	}
	// 验证是否有权限
	if c != nil {
		log.Infof("url=%s, cookie: %s", r.URL.Path, c.String())
		session, err := mongo.SessionColl.Find(c.Value)
		if err != nil {
			http.Error(w, "查找session失败", http.StatusNotAcceptable)
			return
		}
		user := session.User
		err = authProcess(w, r, *user)
		if err != nil {
			log.Errorf("%s,url=%s", err, r.URL.Path)
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

	}

	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

// authProcess 权限处理
func authProcess(w http.ResponseWriter, r *http.Request, user model.User) (err error) {
	if user.UserType != "admin" {
		if r.URL.Path == "/master/trade/query" || r.URL.Path == "/master/trade/report" ||
			r.URL.Path == "/master/trade/stat" || r.URL.Path == "/master/trade/stat/report" {
			log.Infof("agentCode=%s,groupCode=%s,merId=%s", user.AgentCode, user.GroupCode, user.MerId)
			params := r.URL.Query()
			log.Infof("agentCode1=%s,groupCode1=%s,merId1=%s", params.Get("agentCode"), params.Get("groupCode"), params.Get("merId"))
			if user.UserType == "agent" && r.URL.Path != "/master/trade/report" {
				agentCode := params.Get("agentCode")
				if agentCode != user.AgentCode {
					return errors.New("permission denied")
				}
			} else if user.UserType == "group" && r.URL.Path == "/master/trade/query" {
				groupCode := params.Get("groupCode")
				if groupCode != user.GroupCode {
					return errors.New("permission denied")
				}
			} else if user.UserType == "merchant" {
				merId := params.Get("merId")
				if merId != user.MerId {
					return errors.New("permission denied")
				}
			}
		} else {
			return errors.New("permission denied")
		}

	}
	return nil
}
