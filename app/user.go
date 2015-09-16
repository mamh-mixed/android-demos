package app

import (
	"crypto/sha1"
	"fmt"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type user struct{}

var User user

func (u *user) register(userName, password, transtime, sign string) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s,sign=%s", userName, password, transtime, sign)
	// 参数不能为空
	if userName == "" || password == "" || transtime == "" {
		log.Errorf("params are empty")
		return model.NewAppResult(model.FAIL, model.PARAMS_EMPTY)
	}

	//验签
	str := fmt.Sprintf("username=%s&password=%s&transtime=%s%s", userName, password, transtime, model.KEY)
	value := sha1.Sum([]byte(str))
	valueStr := fmt.Sprintf("%x", value)
	log.Debugf("sign(%s)=%s", str, valueStr)
	if sign != valueStr {
		log.Errorf("check signature err")
		return model.NewAppResult(model.FAIL, model.SIGN_FAIL)
	}

	user := &model.AppUser{
		UserName: userName,
		Password: password,
	}
	// 用户是否存在
	num, err := mongo.AppUserCol.FindCountByUserName(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.NewAppResult(model.FAIL, model.SYSTEM_ERROR)
	}
	if num != 0 {
		log.Errorf("userName is exist,userName=%s", userName)
		return model.NewAppResult(model.FAIL, model.USERNAME_EXIST)
	}
	// 保存用户信息
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.NewAppResult(model.FAIL, model.SYSTEM_ERROR)
	}

	result = model.NewAppResult(model.SUCCESS, "")
	return result
}

// login 登录
func (u *user) login(userName, password, transtime, sign string) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s,sign=%s", userName, password, transtime, sign)
	// 参数不能为空
	if userName == "" || password == "" || transtime == "" {
		log.Errorf("params are empty")
		return model.NewAppResult(model.FAIL, model.PARAMS_EMPTY)
	}

	//验签
	str := fmt.Sprintf("username=%s&password=%s&transtime=%s%s", userName, password, transtime, model.KEY)
	value := sha1.Sum([]byte(str))
	valueStr := fmt.Sprintf("%x", value)
	log.Debugf("sign(%s)=%s", str, valueStr)
	if sign != valueStr {
		log.Errorf("check signature err")
		return model.NewAppResult(model.FAIL, model.SIGN_FAIL)
	}

	// 用户是否存在
	num, err := mongo.AppUserCol.FindCountByUserName(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.NewAppResult(model.FAIL, model.SYSTEM_ERROR)
	}
	if num == 0 {
		log.Errorf("userName is exist,userName=%s", userName)
		return model.NewAppResult(model.FAIL, model.USERNAME_NO_EXIST)
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(userName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.NewAppResult(model.FAIL, model.SYSTEM_ERROR)
	}
	if user.Password != password {
		log.Errorf("password is invalid")
		return model.NewAppResult(model.FAIL, model.USERNAME_PASSWORD_ERROR)
	}

	return result
}
