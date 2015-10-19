package master

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

var agentURLArr = []string{
	"/master/trade/query",
	"/master/trade/report",
	"/master/trade/stat",
	"/master/trade/stat/report",
	"/master/trade/findOne",
	"/master/trade/message",
}

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
	mux.HandleFunc("/master/respCode/match", respCodeMatchHandle)
	mux.HandleFunc("/master/user/find", userFindHandle)
	mux.HandleFunc("/master/user/create", userCreateHandle)
	mux.HandleFunc("/master/user/update", userUpdateHandle)
	mux.HandleFunc("/master/user/updatePwd", userUpdatePwdHandle)
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

	// 验证 session 是否过期
	session, err := sessionProcess(w, r)
	if err != nil {
		log.Infof("%s", err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// 验证是否有权限访问这个 URL
	user := session.User
	err = authProcess(user, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	fillUserTypeParam(r.URL.Query(), user)
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

func fillUserTypeParam(qurey url.Values, user *model.User) {
	switch user.UserType {
	case model.UserTypeCIL:
	case model.UserTypeAgent:
		qurey.Set("agentCode", user.AgentCode)
	case model.UserTypeCompany:
		qurey.Set("subAgentCode", user.SubAgentCode)
	case model.UserTypeMerchant:
		qurey.Set("groupCode", user.GroupCode)
	case model.UserTypeShop:
		qurey.Set("merId", user.MerId)
	default:
		log.Errorf("user type error: %s", user.UserType)
	}
}

// authProcess 权限处理
func authProcess(user *model.User, url string) (err error) {
	has := false
	switch user.UserType {
	case model.UserTypeCIL:
		has = true
	case model.UserTypeAgent:
		has = util.StringInSlice(url, agentURLArr)
	case model.UserTypeCompany:
		has = util.StringInSlice(url, agentURLArr)
	case model.UserTypeMerchant:
		has = util.StringInSlice(url, agentURLArr) // 暂时用代理的权限
	case model.UserTypeShop:
		has = util.StringInSlice(url, agentURLArr) // 暂时用代理的权限
	default:
		log.Errorf("user type error: %s", user.UserType)
		return fmt.Errorf("用户类型（%s）配置错误", user.UserType)
	}

	if !has {
		log.Errorf("permission deney: username=%s, url=%s", user.UserName, url)
		return fmt.Errorf("用户没有权限访问 `%s`", url)
	}

	return nil
}

// sessionProcess 处理session
func sessionProcess(w http.ResponseWriter, r *http.Request) (session *model.Session, err error) {
	// 查看请求中有没有cookie
	c, err := r.Cookie(SessionKey)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return nil, err
		}
	}
	// 查询 session 是否过期，如果接近失效则给此 session 延期，如果已经过期则返回失败

	session, err = mongo.SessionColl.Find(c.Value)
	if err != nil {
		log.Debugf("session(%s) not exist: %s", c.Value, err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return nil, err
	}

	// 计算现在到失效时间还有多久
	expire, _ := time.ParseInLocation("2006-01-02 15:04:05", session.Expires, time.Local)
	subTime := expire.Sub(time.Now())
	log.Debugf("session time remain: %s", subTime)

	// 会话已过期
	if subTime < 0 {
		log.Infof("session(%s) expired", c.Value)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return nil, errors.New("会话已过期")
	}

	// 会话接近失效，延长会话失效时间
	if subTime < expiredTime/5 {
		newExpire := expire.Add(expiredTime)
		session.Expires = newExpire.Format("2006-01-02 15:04:05")
		err = mongo.SessionColl.Add(session)
		if err != nil {
			log.Errorf("update session err,%s", err)
		} else {
			c.Expires = newExpire
			http.SetCookie(w, c)
		}
	}

	return session, nil
}
