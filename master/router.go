package master

import (
	"errors"
	"net/http"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

var agentURLArr = []string{"/master/trade/query", "/master/trade/report", "/master/trade/stat", "/master/trade/stat/report",
	"/master/trade/findOne", "/master/trade/message"}

// Route 后台管理的请求统一入口
func Route() (mux *MyServeMux) {
	mux = NewMyServeMux()

	mux.HandleFunc("/master/trade/query", tradeQueryHandle)
	mux.HandleFunc("/master/trade/findOne", tradeFindOneHandle)
	mux.HandleFunc("/master/trade/report", tradeReportHandle)
	mux.HandleFunc("/master/trade/stat", tradeQueryStatsHandle)
	mux.HandleFunc("/master/trade/stat/report", tradeQueryStatsReportHandle)
	mux.HandleFunc("/master/trade/message", tradeMsgHandle)
	mux.HandleFunc("/master/merchant/find", merchantFindHandle)
	mux.HandleFunc("/master/merchant/one", merchantFindOneHandle)
	mux.HandleFunc("/master/merchant/save", merchantSaveHandle)
	mux.HandleFunc("/master/merchant/update", merchantUpdateHandle)
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
	mux.HandleFunc("/master/qiniu/uptoken", uptokenHandle)
	mux.HandleFunc("/master/qiniu/uploaded", downURLHandle)
	mux.HandleFunc("/master/user/find", userFindHandle)
	mux.HandleFunc("/master/user/create", userCreateHandle)
	mux.HandleFunc("/master/user/update", userUpdateHandle)
	mux.HandleFunc("/master/user/updatePwd", userUpdatePwdHandle)
	mux.HandleFunc("/master/respCode/match", respCodeMatchHandle)

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
	// if r.RequestURI == "*" {
	// 	if r.ProtoAtLeast(1, 1) {
	// 		w.Header().Set("Connection", "close")
	// 	}
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
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
	// 验证session是否过期
	session, err := sessionProcess(w, r)
	if err != nil {
		log.Infof("%s", err)
		return
	}

	// 验证是否有权限
	user := session.User
	err = authProcess(w, r, *user)
	if err != nil {
		log.Debugf("%s,url=%s", err, r.URL.Path)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

// authProcess 权限处理
func authProcess(w http.ResponseWriter, r *http.Request, user model.User) (err error) {
	if user.UserType != "admin" {
		perFlag := false
		for _, url := range agentURLArr {
			if url == r.URL.Path {
				perFlag = true
				break
			}
		}
		if perFlag {
			params := r.URL.Query()
			log.Debugf("agentCode=%s,groupCode=%s,merId=%s", params.Get("agentCode"), params.Get("groupCode"), params.Get("merId"))
			if user.UserType == "agent" {
				agentCode := params.Get("agentCode")
				if agentCode != user.AgentCode {
					return errors.New("permission denied")
				}
			} else if user.UserType == "group" {
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

// sessionProcess 处理session
func sessionProcess(w http.ResponseWriter, r *http.Request) (session *model.Session, err error) {
	// 查看请求中有没有cookie
	c, err := r.Cookie("QUICKMASTERID")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return session, err
		}
	}
	// 查询session是否过期，如果还剩最后5分钟则给此session延期，如果已经过期则返回失败
	session, err = mongo.SessionColl.Find(c.Value)
	if err != nil {
		log.Debugf("查询session(%s)出错:%s", c.Value, err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return session, err
	}
	expire, _ := time.ParseInLocation("2006-01-02 15:04:05", session.Expires, time.Local)
	subTime := expire.Sub(time.Now())
	log.Debugf("subTime:%d", subTime)
	if subTime > 0 && subTime < goconf.Config.App.MinExpires*time.Minute {
		newExpire := expire.Add(goconf.Config.App.Expires * time.Minute)
		session.Expires = newExpire.Format("2006-01-02 15:04:05")
		err = mongo.SessionColl.Add(session)
		if err != nil {
			log.Errorf("update session err,%s", err)
		} else {
			c.Expires = newExpire
			http.SetCookie(w, c)
		}

	} else if subTime < 0 {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return session, errors.New("session过期")
	}
	log.Debugf("url=%s, cookie: %s", r.URL.Path, c.String())
	return session, nil
}
