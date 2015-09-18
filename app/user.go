package app

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
func (u *user) register(req *reqParams) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	user := &model.AppUser{
		UserName: req.UserName,
		Password: req.Password,
		Activate: "false",
	}
	// 用户是否存在
	num, err := mongo.AppUserCol.FindCountByUserName(req.UserName)
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

	return model.SUCCESS1
}

// login 登录
func (u *user) login(req *reqParams) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}
	// 密码是否正确
	if user.Password != req.Password {
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
func (u *user) reqActivate(req *reqParams) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}

	// 密码是否正确
	if user.Password != req.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 如果用户已激活，直接返回success
	if user.Activate == "true" {
		return model.SUCCESS1
	}

	// 发送激活链接到注册时提供的邮箱
	code := fmt.Sprintf("%d", rand.Int31())
	hostAddress := goconf.Config.App.NotifyURL
	activateUrl := fmt.Sprintf("%s/app/activate?username=%s&code=%s", hostAddress, req.UserName, code)

	email := &email.Email{
		To:    req.UserName,
		Title: activation.Title,
		Body:  fmt.Sprintf(activation.Body, activateUrl, "点我激活"),
	}
	err = email.Send()

	e := &model.Email{
		UserName:  req.UserName,
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
		return model.SYSTEM_ERROR
	}

	if e.Success {
		return model.SUCCESS1
	} else {
		return model.SYSTEM_ERROR
	}
}

// activate 激活
func (u *user) activate(req *reqParams) (result *model.AppResult) {
	log.Debugf("userName=%s,code=%s", req.UserName, req.Code)
	// 参数不能为空
	if req.UserName == "" || req.Code == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}
	// 如果用户已激活，则返回成功
	if user.Activate == "true" {
		return model.SUCCESS1
	}

	// 判断code是否正确
	e, err := mongo.EmailCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}
	if req.Code != e.Code {
		return model.CODE_ERROR
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

// improveInfo 信息完善
func (u *user) improveInfo(req *reqParams) (result *model.AppResult) {
	if req.UserName == "" || req.Password == "" || req.BankOpen == "" || req.Payee == "" || req.PayeeCard == "" ||
		req.PhoneNum == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}

	// 密码是否正确
	if user.Password != req.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 创建商户
	uniqueId := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31())
	randStr := fmt.Sprintf("%d", rand.Int31())
	merchant := &model.Merchant{
		MerStatus:  model.MerStatusNormal,
		TransCurr:  "156",
		UniqueId:   uniqueId,
		IsNeedSign: true,
		SignKey:    fmt.Sprintf("%x", base64.StdEncoding.EncodeToString([]byte(randStr))),
		Detail: model.MerDetail{
			MerName:       "云收银",
			CommodityName: "讯联云收银在线注册商户",
		},
	}
	for {
		// 设置merId
		maxMerId, err := mongo.MerchantColl.FindMaxMerId()
		if err != nil {
			log.Errorf("find database  err,%s", err)
			return model.SYSTEM_ERROR
		}
		maxMerIdNum, err := strconv.Atoi(maxMerId)
		if err != nil {
			log.Errorf("format maxMerId(%s) err", maxMerId)
			return model.SYSTEM_ERROR
		}
		merchant.MerId = fmt.Sprintf("%d", maxMerIdNum+1)

		err = mongo.MerchantColl.Insert2(merchant)
		if err != nil {
			isDuplicateMerId := strings.Contains(err.Error(), "E11000 duplicate key error index")
			if !isDuplicateMerId {
				log.Errorf("create merchant err,%s", err)
				return model.SYSTEM_ERROR
			}

		} else {
			break
		}
	}

	// 创建路由,支付宝，微信
	alpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "ALP",
		ChanCode:  "ALP",
		ChanMerId: "2088811767473826",
	}
	err = mongo.RouterPolicyColl.Insert(alpRoute)
	if err != nil {
		log.Errorf("create routePolicy err,%s", err)
		return model.SYSTEM_ERROR
	}

	wxpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "WXP",
		ChanCode:  "WXP",
		ChanMerId: "1239305502",
	}
	err = mongo.RouterPolicyColl.Insert(wxpRoute)
	if err != nil {
		log.Errorf("create routePolicy err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 更新用户信息
	user.BankOpen = req.BankOpen
	user.Payee = req.Payee
	user.PayeeCard = req.PayeeCard
	user.PhoneNum = req.PhoneNum
	user.MerId = merchant.MerId
	user.Limit = "true"
	user.SignKey = merchant.SignKey
	user.AgentCode = merchant.AgentCode
	user.UniqueId = merchant.UniqueId

	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}
	result = &model.AppResult{
		State: model.SUCCESS,
		Error: "",
		User:  user,
	}
	return result
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
