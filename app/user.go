package app

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type user struct{}

var User user

// register 注册
func (u *user) register(userName, password, transtime, sign string) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s,sign=%s", userName, password, transtime, sign)
	// 参数不能为空
	if userName == "" || password == "" || transtime == "" || sign == "" {
		return model.PARAMS_EMPTY
	}

	user := &model.AppUser{
		UserName: userName,
		Password: password,
		Activate: "false",
	}
	// 用户是否存在
	num, err := mongo.AppUserCol.FindCountByUserName(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}
	if num != 0 {
		return model.USERNAME_EXIST
	}
	// 保存用户信息
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}

	result = model.SUCCESS1
	return result
}

// login 登录
func (u *user) login(userName, password, transtime, sign string) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s,sign=%s", userName, password, transtime, sign)
	// 参数不能为空
	if userName == "" || password == "" || transtime == "" || sign == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}
	// 密码是否正确
	if user.Password != password {
		return model.USERNAME_PASSWORD_ERROR
	}
	// 用户是否激活
	if user.Activate == "false" {
		return model.USER_NO_ACTIVATE
	}

	// 查找uniqueId
	if user.MerId != "" {
		merchant, err := mongo.MerchantColl.Find(user.MerId)
		if err != nil {
			log.Errorf("find database err,%s", err)
			return model.SYSTEM_ERROR
		}
		user.UniqueId = merchant.UniqueId
	}

	result = &model.AppResult{
		State: model.SUCCESS,
		Error: "",
		User:  user,
	}

	return result
}

// reqActivate 请求发送激活链接
func (u *user) reqActivate(userName, password, transtime, sign string) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s,sign=%s", userName, password, transtime, sign)
	// 参数不能为空
	if userName == "" || password == "" || transtime == "" || sign == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}

	// 密码是否正确
	if user.Password != password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 如果用户已激活，直接返回success
	if user.Activate == "true" {
		return model.SUCCESS1
	}

	// 发送激活链接到注册时提供的邮箱
	code := fmt.Sprintf("%d", rand.Int31())
	hostAddress := goconf.Config.App.NotifyURL
	activateUrl := fmt.Sprintf("%s/app/activate?username=%s&code=%s", hostAddress, userName, code)

	email := &email.Email{
		To:    userName,
		Title: activation.Title,
		Body:  fmt.Sprintf(activation.Body, activateUrl, "点我激活"),
	}
	err = email.Send()

	e := &model.Email{
		UserName:  userName,
		Code:      code,
		Success:   true,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err != nil {
		e.Success = false
	}

	// 保存email信息
	err = mongo.EmailCol.Upsert(e)
	if err != nil {
		log.Errorf("save email err")
	}

	return model.SUCCESS1
}

// activate 激活
func (u *user) activate(userName, code string) (result *model.AppResult) {
	log.Debugf("userName=%s,code=%s", userName, code)
	// 参数不能为空
	if userName == "" || code == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}

	// 判断code是否正确
	e, err := mongo.EmailCol.FindOne(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}
	if code != e.Code {

	}

	// 更新activate为已激活
	user.Activate = "true"
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}

	return model.SUCCESS1
}

// promoteLimit 提升限额
func (u *user) promoteLimit(req *reqParams) (result *model.AppResult) {
	// 用户名不为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 发送邮件通知Andy.Li
	return
}

// getSettInfo 获得清算信息
func (u *user) getSettInfo(req *reqParams) (result *model.AppResult) {
	// 用户名不为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 返回
	result = model.NewAppResult(model.SUCCESS, "")
	result.Payee = user.Payee
	result.BankOpen = user.BankOpen
	result.PayeeCard = user.PayeeCard
	result.PhoneNum = user.PhoneNum
	return
}

// updateSettInfo 更新清算信息
func (u *user) updateSettInfo(req *reqParams) (result *model.AppResult) {
	// 用户名不为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 修改
	user.Payee = req.Payee
	user.BankOpen = req.BankOpen
	user.PayeeCard = req.PayeeCard
	user.PhoneNum = req.PhoneNum

	if err = mongo.AppUserCol.Upsert(user); err != nil {
		return model.SYSTEM_ERROR
	}

	return model.SUCCESS1
}
