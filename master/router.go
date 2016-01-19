package master

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

// 路径中包含以下关键字，则记录到数据库
var logKeysArr = []string{
	"create",
	"save",
	"update",
	"delete",
	"reset",
	"login",
	"logout",
}

// Route 后台管理的请求统一入口
func Route() (mux *MyServeMux) {
	mux = NewMyServeMux()

	mux.HandleFunc("/master/trade/query", tradeQueryHandle)
	mux.HandleFunc("/master/trade/findOne", tradeFindOneHandle)
	mux.HandleFunc("/master/trade/report", tradeReportHandle)
	mux.HandleFunc("/master/trade/stat", tradeQueryStatsHandle)
	mux.HandleFunc("/master/trade/stat/report", tradeQueryStatsReportHandle)
	mux.HandleFunc("/master/trade/transfer/query", tradeTransferQueryHandle)
	mux.HandleFunc("/master/trade/transfer/report", tradeTransferReportHandle)
	mux.HandleFunc("/master/trade/message", tradeMsgHandle)
	mux.HandleFunc("/master/trade/settle/journal", tradeSettleJournalHandle)
	mux.HandleFunc("/master/trade/settle/report", tradeSettleReportHandle)
	mux.HandleFunc("/master/trade/settle/refresh", tradeSettleRefreshHandle)
	mux.HandleFunc("/master/merchant/find", merchantFindHandle)
	mux.HandleFunc("/master/merchant/one", merchantFindOneHandle)
	mux.HandleFunc("/master/merchant/save", merchantSaveHandle)
	mux.HandleFunc("/master/merchant/update", merchantUpdateHandle)
	mux.HandleFunc("/master/merchant/import", importMerchant)
	mux.HandleFunc("/master/merchant/delete", merchantDeleteHandle)
	mux.HandleFunc("/master/merchant/export", merchantExportHandle)
	mux.HandleFunc("/master/router/save", routerSaveHandle)
	mux.HandleFunc("/master/router/update", routerUpdateHandle)
	mux.HandleFunc("/master/router/find", routerFindHandle)
	mux.HandleFunc("/master/router/one", routerFindOneHandle)
	mux.HandleFunc("/master/router/delete", routerDeleteHandle)
	mux.HandleFunc("/master/channelMerchant/find", channelMerchantFindHandle)
	mux.HandleFunc("/master/channelMerchant/match", channelMerchantMatchHandle)
	mux.HandleFunc("/master/channelMerchant/findByMerIdAndCardBrand", channelFindByMerIdAndCardBrandHandle)
	mux.HandleFunc("/master/channelMerchant/save", channelMerchantSaveHandle)
	mux.HandleFunc("/master/channelMerchant/update", channelMerchantUpdateHandle)
	mux.HandleFunc("/master/channelMerchant/delete", channelMerchantDeleteHandle)
	mux.HandleFunc("/master/agent/find", agentFindHandle)
	mux.HandleFunc("/master/agent/delete", agentDeleteHandle)
	mux.HandleFunc("/master/agent/save", agentSaveHandle)
	mux.HandleFunc("/master/agent/update", agentUpdateHandle)
	mux.HandleFunc("/master/subAgent/find", subAgentFindHandle)
	mux.HandleFunc("/master/subAgent/delete", subAgentDeleteHandle)
	mux.HandleFunc("/master/subAgent/save", subAgentSaveHandle)
	mux.HandleFunc("/master/subAgent/update", subAgentUpdateHandle)
	mux.HandleFunc("/master/group/find", groupFindHandle)
	mux.HandleFunc("/master/group/delete", groupDeleteHandle)
	mux.HandleFunc("/master/group/save", groupSaveHandle)
	mux.HandleFunc("/master/group/update", groupUpdateHandle)
	mux.HandleFunc("/master/qiniu/uptoken", uptokenHandle)
	mux.HandleFunc("/master/qiniu/uploaded", downURLHandle)
	mux.HandleFunc("/master/qiniu/download", qiniuDownloadHandle)
	mux.HandleFunc("/master/respCode/match", respCodeMatchHandle)
	mux.HandleFunc("/master/user/find", userFindHandle)
	mux.HandleFunc("/master/user/create", userCreateHandle)
	mux.HandleFunc("/master/user/update", userUpdateHandle)
	mux.HandleFunc("/master/user/updatePwd", userUpdatePwdHandle)
	mux.HandleFunc("/master/user/delete", userDeleteHandle)
	mux.HandleFunc("/master/user/resetPwd", userResetPwdHandle)
	mux.HandleFunc("/master/app/locale", appLocaleHandle)
	mux.HandleFunc("/master/app/resetPwd", appResetPwdHandle)
	mux.HandleFunc("/master/list", kvListHandle)
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
	if r.URL.Path == "/master/logout" {
		sessionDeleteHandle(w, r)
		return
	}
	// 验证 session 是否过期
	session, err := sessionProcess(w, r)
	if err != nil {
		log.Infof("%s", err)

		// 将QUICKMASTERID设成失效
		http.SetCookie(w, &http.Cookie{
			Name:     "QUICKMASTERID",
			Value:    "",
			HttpOnly: true,
			Path:     "/master",
			MaxAge:   -1,
		})

		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// 验证是否有权限访问这个 URL
	user := session.User
	err = authProcess(user, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 记录平台操作日志
	HandleMasterLog(w, r, user)

	fillUserTypeParam(r, user)
	// log.Debugf("query: %#v", r.URL.Query())

	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

func fillUserTypeParam(r *http.Request, user *model.User) {
	log.Debugf("user: %#v", user)

	query := r.URL.Query()
	query.Set("userType", user.UserType)

	switch user.UserType {
	case model.UserTypeCIL:
	case model.UserTypeGenAdmin:
	case model.UserTypeShop:
		query.Set("merId", user.MerId)
	case model.UserTypeMerchant:
		query.Set("groupCode", user.GroupCode)
	case model.UserTypeCompany:
		query.Set("subAgentCode", user.SubAgentCode)
	case model.UserTypeAgent:
		query.Set("agentCode", user.AgentCode)
	default:
		log.Errorf("user type error: %s", user.UserType)
	}

	r.URL.RawQuery = query.Encode()

}

// authProcess 权限处理
func authProcess(user *model.User, url string) (err error) {
	has := false
	switch user.UserType {
	case model.UserTypeCIL:
		has = true
	case model.UserTypeGenAdmin:
		has = util.StringInSlice(url, genAdminURLArr)
	case model.UserTypeAgent:
		has = util.StringInSlice(url, agentURLArr)
	case model.UserTypeCompany:
		has = util.StringInSlice(url, commonURLArr)
	case model.UserTypeMerchant:
		has = util.StringInSlice(url, commonURLArr)
	case model.UserTypeShop:
		has = util.StringInSlice(url, commonURLArr)
	default:
		log.Errorf("user type error: %s", user.UserType)
		return fmt.Errorf("用户类型（%s）配置错误", user.UserType)
	}

	if !has {
		log.Errorf("permission deney: username=%s, url=%s", user.UserName, url)
		return errors.New("Unauthorized")
	}

	return nil
}

var refreshTime = expiredTime / 2

// sessionProcess 处理session
func sessionProcess(w http.ResponseWriter, r *http.Request) (session *model.Session, err error) {
	// 查看请求中有没有cookie
	c, err := r.Cookie(SessionKey)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, err
		}
	}
	log.Debugf("==================sessionId=%s", c.Value)
	// 查询 session 是否过期，如果接近失效则给此 session 延期，如果已经过期则返回失败
	session, err = mongo.SessionColl.Find(c.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "QUICKMASTERID",
			Value:    "",
			HttpOnly: true,
			Path:     "*",
			MaxAge:   -1,
		})
		log.Debugf("session(%s) not exist: %s", c.Value, err)
		return nil, err
	}

	// 计算现在到失效时间还有多久
	now := time.Now()
	subTime := session.Expires.Sub(now)
	log.Debugf("session time remain: %s", subTime)

	// 会话已过期
	if subTime < 0 {
		http.SetCookie(w, &http.Cookie{
			Name:     "QUICKMASTERID",
			Value:    "",
			HttpOnly: true,
			Path:     "*",
			MaxAge:   -1,
		})
		log.Infof("session(%s) expired", c.Value)
		return nil, errors.New("会话已过期")
	}

	// 会话接近失效，延长会话失效时间
	if subTime < refreshTime {
		session.Expires = session.Expires.Add(refreshTime)
		session.UpdateTime = now
		err = mongo.SessionColl.Add(session)
		if err != nil {
			log.Errorf("update session err,%s", err)
		} else {
			c.Expires = session.Expires
			c.HttpOnly = true
			c.Path = "/master"
			http.SetCookie(w, c)
		}
		log.Infof("prolong session(%s) to %s", c.Value, session.Expires)
	}

	return session, nil
}

// HandleMasterLog 记录平台操作日志
func HandleMasterLog(w http.ResponseWriter, r *http.Request, user *model.User) {
	path := r.URL.Path
	// 增删改操作记录到数据库
	isLog := false
	for _, key := range logKeysArr {
		if strings.Contains(path, key) {
			isLog = true
			break
		}
	}
	if !isLog {
		return
	}
	var body []byte
	var err error
	if r.Method == "POST" {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("read body err,%s", err)
			return
		}
		r.Body.Close()

		// r.Body 只能被读取一次，读完之后再写入
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	// 如果是修改密码操作，则不需要记录body中数据
	if path == "/master/user/updatePwd" {
		body = []byte("")
	}
	InsertMasterLog(r, user, body)
}

// InsertMasterLog 操作日志入库
func InsertMasterLog(r *http.Request, user *model.User, body []byte) {
	// 取客户端 IP，优先取 X-Forwarded-For 第一个 IP，
	// 如果没有，再取 X-Real-IP，最后是 RemoteAddr
	clientIP := r.RemoteAddr
	forwordedFor := r.Header.Get("X-Forwarded-For")
	if forwordedFor != "" {
		clientIP = strings.Split(forwordedFor, ", ")[0]
	} else {
		realIP := r.Header.Get("X-Real-IP")
		if realIP != "" {
			clientIP = realIP
		}
	}

	masterLog := &model.MasterLog{
		UserName: user.UserName,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Path:     r.URL.Path,
		Method:   r.Method,
		Query:    r.URL.RawQuery,
		Body:     string(body),
		IP:       clientIP,
	}
	mongo.MasterLogColl.Insert(masterLog)
}
