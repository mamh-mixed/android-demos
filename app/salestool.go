// 销售工具接口
package app

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"io/ioutil"
	"net/http"
	"time"
)

// tokenMap TODO 将token存放到数据库
var tokenMap = make(map[string]*model.User)

// CompanyLogin 销售人员-公司级别登录
func CompanyLogin(w http.ResponseWriter, r *http.Request) {

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	debugReqParams(r)

	username := r.FormValue("username")
	user, err := mongo.UserColl.FindOneUser(username, "", "")
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	// 不是公司这一级的
	if user.UserType != model.UserTypeCompany {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	password := r.FormValue("password")
	ps := fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD + "{" + username + "}" + password))))
	if user.Password != ps {
		w.Write(jsonMarshal(model.USERNAME_PASSWORD_ERROR))
		return
	}

	if user.SubAgentCode == "" {
		log.Errorf("userType is company,but can not find subAgentCode, username=%s", username)
		w.Write(jsonMarshal(model.USER_DATA_ERROR))
		return
	}

	result := model.SUCCESS1
	result.AccessToken = genAccessToken(user)

	w.Write(jsonMarshal(result))
}

// UserList 用户列表
func UserList(w http.ResponseWriter, r *http.Request) {

	var agentUser *model.User
	var ok bool

	debugReqParams(r)

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	users, err := mongo.AppUserCol.FindBySubAgentCode(agentUser.SubAgentCode)
	if err != nil {
		w.Write(jsonMarshal(model.SYSTEM_ERROR))
		return
	}

	var merIds []string
	var userMap = make(map[string]*model.AppUser)
	for _, user := range users {
		merIds = append(merIds, user.MerId)
		userMap[user.MerId] = user
	}

	// 关联商户信息
	mers, err := mongo.MerchantColl.FuzzyFind(&model.QueryCondition{
		MerIds: merIds,
	})
	if err != nil {
		w.Write(jsonMarshal(model.SYSTEM_ERROR))
		return
	}

	for _, m := range mers {
		if user, ok := userMap[m.MerId]; ok {
			user.BankOpen = m.Detail.OpenBankName
			user.Payee = m.Detail.AcctName
			user.PayeeCard = m.Detail.AcctNum
			user.PhoneNum = m.Detail.ContactTel
			user.SignKey = m.SignKey
			user.AgentCode = m.AgentCode
			user.UniqueId = m.UniqueId
			user.MerName = m.Detail.MerName
			user.Images = m.Detail.Images
			user.Password = "" // 不显示
		}
	}

	// 成功返回
	result := model.SUCCESS1
	result.Users = users

	w.Write(jsonMarshal(result))
}

// UserRegister 用户注册
func UserRegister(w http.ResponseWriter, r *http.Request) {

	var agentUser *model.User
	var ok bool

	debugReqParams(r)

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	req := &reqParams{
		UserName:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		Transtime:    time.Now().Format("20060102150405"),
		Remark:       "company_register",
		SubAgentCode: agentUser.SubAgentCode,
	}
	// 注册
	result := User.register(req)

	// 注册成功，创建商户
	if result.State == model.SUCCESS {
		merchant := &model.Merchant{
			AgentCode:    agentUser.AgentCode,
			SubAgentCode: agentUser.SubAgentCode,
			Permission:   []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf},
			MerStatus:    model.MerStatusNormal,
			Remark:       "agent_register",
			TransCurr:    "156",
			RefundType:   model.CurrentDayRefund, // 只能当天退
			IsNeedSign:   true,
			SignKey:      fmt.Sprintf("%x", randBytes(16)),
		}

		subAgent, err := mongo.SubAgentColl.Find(agentUser.SubAgentCode)
		if err == nil {
			merchant.AgentName = subAgent.AgentName
			merchant.SubAgentName = subAgent.SubAgentName
		}

		if err := genMerId(merchant, subAgent.AgentCode+"0"); err != nil {
			w.Write(jsonMarshal(err))
			return
		}
		if err := genRouter(merchant); err != nil {
			w.Write(jsonMarshal(err))
			return
		}
		user := req.AppUser
		if user != nil {
			user.MerId = merchant.MerId
			user.UniqueId = merchant.UniqueId
			user.SignKey = merchant.SignKey
			user.AgentCode = merchant.AgentCode
			mongo.AppUserCol.Upsert(user)
			result.User = user
			user.Password = "" //不显示
		}
	}

	w.Write(jsonMarshal(result))
}

// GetQiniuToken
func GetQiniuToken(w http.ResponseWriter, r *http.Request) {

	debugReqParams(r)

	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	result := model.SUCCESS1
	result.UploadToken = qiniu.GetUploadtoken()

	w.Write(jsonMarshal(result))
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {

	debugReqParams(r)
	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	appUser, err := mongo.AppUserCol.FindOne(r.FormValue("username"))
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
	}

	req := &reqParams{
		BankOpen:   r.FormValue("bank_open"),
		Payee:      r.FormValue("payee"),
		PayeeCard:  r.FormValue("payee_card"),
		PhoneNum:   r.FormValue("phone_num"),
		Transtime:  r.FormValue("transtime"),
		Province:   r.FormValue("province"),
		City:       r.FormValue("city"),
		BranchBank: r.FormValue("branch_bank"),
		BankNo:     r.FormValue("bankNo"),
		MerName:    r.FormValue("merName"),
		Images:     r.Form["image"],
		AppUser:    appUser, // 带上user
	}

	w.Write(jsonMarshal(User.updateSettInfo(req)))
}

// UserActivate 用户激活
func UserActivate(w http.ResponseWriter, r *http.Request) {
	debugReqParams(r)
	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	username := r.FormValue("username")
	key := r.FormValue("imageUrl")
	if username == "" || key == "" {
		w.Write(jsonMarshal(model.PARAMS_EMPTY))
		return
	}

	appUser, err := mongo.AppUserCol.FindOne(username)
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	m, err := mongo.MerchantColl.Find(appUser.MerId)
	if err != nil {
		w.Write(jsonMarshal(model.MERID_NO_EXIST))
		return
	}

	// 访问七牛获取二维码设计图
	var pngBytes []byte
	if key != "" {
		resp, err := http.Get(qiniu.MakePrivateUrl(key))
		if err == nil {
			pngBytes, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("fail to read resp: %s", err)
			}
		} else {
			log.Errorf("fail to load image: %s, key=%s", err, key)
		}
	}

	// 发email告知用户已开通成功
	email := &email.Email{To: username, Title: open.Title}

	if err != nil && key != "" {
		png64 := base64.StdEncoding.EncodeToString(pngBytes)
		image := fmt.Sprintf(`<img src="data:image/png;base64,%s"/>`, png64)
		email.Body = fmt.Sprintf(open.Body, m.Detail.MerName, username, m.MerId, m.SignKey, image)
	} else {
		email.Body = fmt.Sprintf(open.Body, m.Detail.MerName, username, m.MerId, m.SignKey, "")
	}

	// 异步发送邮件
	go func() {
		err := email.Send()
		if err != nil {
			log.Errorf("send email fail: %s, To=%s, body=%s", err, username, email.Body)
		}
	}()

	w.Write(jsonMarshal(model.SUCCESS1))

}

// checkAccessToken
func checkAccessToken(token string) (*model.User, bool) {

	if user, ok := tokenMap[token]; ok {
		return user, true
	}
	// 向数据库里查找
	s, err := mongo.SessionColl.Find(token)
	if err != nil {
		return nil, false
	}

	// 放到map里
	token[token] = s.User

	return s.User, true
}

func genAccessToken(user *model.User) string {
	s := &model.Session{
		SessionID:  util.SerialNumber(),
		User:       user,
		CreateTime: time.Now(),
	}
	s.UpdateTime = s.CreateTime
	s.Expires = s.CreateTime.Add(24 * time.Hour)
	mongo.SessionColl.Add(s)
	tokenMap[s.SessionID] = user
	return s.SessionID
}

func debugReqParams(r *http.Request) {
	log.Debugf("app tools req: %v", r.Form)
}
