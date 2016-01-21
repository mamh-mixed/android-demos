package master

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/CardInfoLink/log"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

type appUserController struct{}

var AppUser appUserController

func (u *appUserController) ResetPwd(data []byte, user *model.User) (ret *model.ResultBody) {
	log.Debugf("data:%s", string(data))
	userPwd := &model.UserPwd{}
	err := json.Unmarshal(data, userPwd)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(1, "json失败")
	}

	// app用户名是否存在
	appUser, err := mongo.AppUserCol.FindOne(userPwd.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.NewResultBody(3, "USERNAME_NOT_FOUND")
		}
		return model.NewResultBody(2, "查询数据库失败")
	}

	if appUser.MerId == "" {
		return model.NewResultBody(4, "NO_PERMISSION")
	}

	// 该用户是否有权限
	if user.UserType == model.UserTypeAgent {
		merchant, err := mongo.MerchantColl.FindNotInCache(appUser.MerId)

		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(5, "SHOP_NOT_FOUND")
			}
			log.Errorf("查询一个商户(%s)出错: %s", appUser.MerId, err)
			return model.NewResultBody(2, "查询数据库失败")
		}
		if merchant.AgentCode != user.AgentCode {
			return model.NewResultBody(4, "NO_PERMISSION")
		}
	}

	// 新密码解密
	newPwd, err := rsaDecryptFromBrowser(userPwd.NewPwd)
	if err != nil {
		log.Errorf("escrypt password error %s", err)
		return model.NewResultBody(5, "DECRYPT_ERROR")
	}

	if len(newPwd) < 8 {
		log.Errorf("new password's length must greater than 8; now is %d", len(newPwd))
		return model.NewResultBody(6, "PASSWORD_LENGTH_NOT_ENOUGH")
	}

	// 密码复杂度校验
	if !isPasswordOk(newPwd) {
		log.Errorf("new password is not conplicated enough: %s", newPwd)
		return model.NewResultBody(7, "PASSWORD_NOT_COMPLICATED_ENOUGH")
	}

	if appUser != nil {
		pb := md5.Sum([]byte(newPwd))
		appUser.Password = fmt.Sprintf("%x", string(pb[:]))
		appUser.LoginTime = ""
		appUser.LockTime = ""
		if err = mongo.AppUserCol.Update(appUser); err != nil {
			return model.NewResultBody(2, "修改密码失败")
		}
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "重置密码成功",
	}
	return ret
}
