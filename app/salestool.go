// 销售工具接口
package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"image/jpeg"
	"net/http"
	"time"
)

var tokenMap = make(map[string]*model.User)
var qrImage = "tools/qr/image/%s/%s.jpg"

// CompanyLogin 销售人员-公司级别登录
func CompanyLogin(w http.ResponseWriter, r *http.Request) {

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

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

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	users, err := mongo.AppUserCol.Find(&model.AppUserContiditon{SubAgentCode: agentUser.SubAgentCode})
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
		UserFrom:     model.SalesToolsRegister,
		BelongsTo:    agentUser.UserName,
		SubAgentCode: agentUser.SubAgentCode,
	}
	// 注册

	w.Write(jsonMarshal(User.register(req)))
}

// GetQiniuToken
func GetQiniuToken(w http.ResponseWriter, r *http.Request) {

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

	var agentUser *model.User
	var ok bool

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	appUser, err := mongo.AppUserCol.FindOne(r.FormValue("username"))
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}
	// log.Debug(appUser)

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

	// 默认返回
	result := model.SUCCESS1

	var merchant *model.Merchant
	// 还没申请商户
	if appUser.MerId == "" {
		merchant = &model.Merchant{
			AgentCode:    agentUser.AgentCode,
			SubAgentCode: agentUser.SubAgentCode,
			Permission:   []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf},
			MerStatus:    model.MerStatusNormal,
			Remark:       "agent_register",
			TransCurr:    "156",
			RefundType:   model.CurrentDayRefund, // 只能当天退
			IsNeedSign:   true,
			SignKey:      fmt.Sprintf("%x", randBytes(16)),
			Detail: model.MerDetail{
				MerName:       req.MerName,
				CommodityName: req.MerName,
				Province:      req.Province,
				City:          req.City,
				OpenBankName:  req.BankOpen,
				BankName:      req.BranchBank,
				BankId:        req.BankNo,
				AcctName:      req.Payee,
				AcctNum:       req.PayeeCard,
				ContactTel:    req.PhoneNum,
				TitleOne:      "欢迎光临",
				TitleTwo:      req.MerName,
				Images:        req.Images,
			},
		}

		subAgent, err := mongo.SubAgentColl.Find(agentUser.SubAgentCode)
		if err == nil {
			merchant.AgentName = subAgent.AgentName
			merchant.SubAgentName = subAgent.SubAgentName
		}

		if err := genMerId(merchant, subAgent.AgentCode+"0"); err != nil {
			w.Write(jsonMarshal(model.SYSTEM_ERROR))
			return
		}
		if err := genRouter(merchant); err != nil {
			w.Write(jsonMarshal(model.SYSTEM_ERROR))
			return
		}
	} else {
		result = User.updateSettInfo(req)
		if result.State != model.SUCCESS {
			w.Write(jsonMarshal(result))
			return
		}
		merchant = req.m
	}

	appUser.MerId = merchant.MerId
	appUser.UniqueId = merchant.UniqueId
	appUser.SignKey = merchant.SignKey
	appUser.AgentCode = merchant.AgentCode
	appUser.SubAgentCode = agentUser.SubAgentCode
	mongo.AppUserCol.Upsert(appUser)
	result.User = appUser
	appUser.Password = "" //不显示

	w.Write(jsonMarshal(result))
}

// UserActivate 用户激活
func UserActivate(w http.ResponseWriter, r *http.Request) {
	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	username := r.FormValue("username")
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

	// 异步处理邮件和上传图片
	go func() {

		// 生成支付设计图
		payImage := genImageWithQrCode(payImageTemplate, m.Detail.PayUrl, m.Detail.MerName)

		// 写入缓冲池
		var isEnocdeSuccess = true
		payBuf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(payBuf, payImage, nil)
		if err != nil {
			log.Errorf("jpeg encode fail: %s", err)
			isEnocdeSuccess = false
		}

		// 发email告知用户已开通成功
		email := &email.Email{To: username, Title: open.Title}
		if isEnocdeSuccess {
			jpg64 := base64.StdEncoding.EncodeToString(payBuf.Bytes())
			image := fmt.Sprintf(`<img src="data:image/jpeg;base64,%s"/>`, jpg64)
			email.Body = fmt.Sprintf(open.Body, m.Detail.MerName, username, m.MerId, m.SignKey, image)
		} else {
			email.Body = fmt.Sprintf(open.Body, m.Detail.MerName, username, m.MerId, m.SignKey, "")
		}
		err = email.Send()
		if err != nil {
			log.Errorf("send email fail: %s, To=%s, body=%s", err, username, email.Body)
		}

		// upload payImage
		if isEnocdeSuccess {
			err = qiniu.Put(fmt.Sprintf(qrImage, m.MerId, "pay"), int64(payBuf.Len()), payBuf)
			if err != nil {
				log.Errorf("fail to upload image : %s", err)
			}
		}

		// 生成账单设计图
		billImage := genImageWithQrCode(billImageTemplate, m.Detail.BillUrl, m.Detail.MerName)
		billBuf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(billBuf, billImage, nil)
		if err != nil {
			log.Errorf("jpeg encode fail: %s", err)
			isEnocdeSuccess = false
		} else {
			isEnocdeSuccess = true
		}

		// upload billImage
		if isEnocdeSuccess {
			err = qiniu.Put(fmt.Sprintf(qrImage, m.MerId, "bill"), int64(billBuf.Len()), billBuf)
			if err != nil {
				log.Errorf("fail to upload image : %s", err)
			}
		}

	}()

	w.Write(jsonMarshal(model.SUCCESS1))

}

// GetDownloadUrl 生成下载地址
func GetDownloadUrl(w http.ResponseWriter, r *http.Request) {
	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	merId := r.FormValue("merId")
	imageType := r.FormValue("imageType")

	if merId == "" || imageType == "" {
		w.Write(jsonMarshal(model.PARAMS_EMPTY))
		return
	}

	result := model.SUCCESS1
	result.DownloadUrl = qiniu.MakePrivateUrl(fmt.Sprintf(qrImage, merId, imageType))
	w.Write(jsonMarshal(result))
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
	tokenMap[token] = s.User

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

// NotifySalesman 每天汇总当天用户数据给业务人员
func NotifySalesman() {
	day := time.Now().Format("2006-01-02")
	all, err := mongo.AppUserCol.Find(&model.AppUserContiditon{
		RegisterFrom: model.SalesToolsRegister,
		StartTime:    day + " 00:00:00",
		EndTime:      day + " 23:59:59",
	})
	if err != nil {
		log.Errorf("find appUser error:%s", err)
		return
	}

	// 归属
	c := make(map[string][]*model.AppUser)
	for _, u := range all {
		if users, ok := c[u.BelongsTo]; ok {
			users = append(users, u)
		} else {
			c[u.BelongsTo] = []*model.AppUser{u}
		}
	}

	for k, v := range c {
		for _, u := range v {
			log.Debugf("k=%s c=%s u=%s", k, u, v)
		}
	}

}
