package master

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"unicode"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type userController struct{}

var User userController

// Login 登陆
func (u *userController) Login(userName, password string) (ret *model.ResultBody) {
	if userName == "" || password == "" {
		log.Errorf("username or passwod must not blank")
		return model.NewResultBody(1, "用户名或密码不用为空")
	}

	var encryptPass string
	user, err := mongo.UserColl.FindOneUser(userName, "", "")
	if err != nil {

		// 兼容appUser用户登录
		appUser, err := mongo.AppUserCol.FindOne(userName)
		if err != nil {
			log.Errorf("find user(%s) error: %s", userName, err)
			return model.NewResultBody(2, "无此用户名")
		} else {
			pb := md5.Sum([]byte(password))
			encryptPass = fmt.Sprintf("%x", pb[:])
			user = &model.User{
				UserName: appUser.UserName,
				NickName: appUser.UserName,
				Password: appUser.Password,
				Mail:     "",
				PhoneNum: "",
				UserType: model.UserTypeShop,
				// AgentCode:    "",
				// SubAgentCode: "",
				// GroupCode:    "",
				MerId: appUser.MerId,
			}
		}

	} else {
		encryptPass = fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD + "{" + userName + "}" + password))))
	}

	// 校验
	if encryptPass != user.Password {
		log.Errorf("wrong password, expect %s but get %s", user.Password, encryptPass)
		return model.NewResultBody(3, "密码错误")
	}

	// 隐藏密码
	user.Password = ""

	ret = &model.ResultBody{
		Status:  0,
		Message: "登陆成功",
		Data:    user,
	}
	return ret
}

// 新建用户
func (u *userController) CreateUser(data []byte) (ret *model.ResultBody) {
	user := &model.User{}
	err := json.Unmarshal(data, user)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(1, "json失败")
	}
	// 设置默认密码
	passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + model.DEFAULT_PWD)
	user.Password = fmt.Sprintf("%x", sha1.Sum(passData))

	// 判断必填项是否为空
	if user.UserName == "" {
		log.Errorf("必填项不能为空")
		return model.NewResultBody(2, "必填项不能为空")
	}
	if user.UserType == model.UserTypeAgent {
		if user.AgentCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		_, err := mongo.AgentColl.Find(user.AgentCode)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此代理代码")
			}
			log.Errorf("查询代理(%s)出错:%s", user.AgentCode, err)
			return model.NewResultBody(1, "查询失败")
		}
	} else if user.UserType == model.UserTypeCompany {
		if user.SubAgentCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		subAgent, err := mongo.SubAgentColl.Find(user.SubAgentCode)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此公司代码")
			}
			log.Errorf("查询二级代理(%s)出错:%s", user.SubAgentCode, err)
			return model.NewResultBody(1, "查询失败")
		}
		user.AgentCode = subAgent.AgentCode

	} else if user.UserType == model.UserTypeMerchant {
		if user.GroupCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		group, err := mongo.GroupColl.Find(user.GroupCode)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此商户代码")
			}
			log.Errorf("查询集团(%s)出错: %s", user.GroupCode, err)
			return model.NewResultBody(1, "查询失败")
		}
		user.AgentCode = group.AgentCode
		user.SubAgentCode = group.SubAgentCode
	} else if user.UserType == model.UserTypeShop {
		if user.MerId == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		merchant, err := mongo.MerchantColl.FindNotInCache(user.MerId)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此门店代码")
			}
			log.Errorf("查询一个商户(%s)出错: %s", user.MerId, err)
			return model.NewResultBody(1, "查询失败")
		}

		user.AgentCode = merchant.AgentCode
		user.SubAgentCode = merchant.SubAgentCode
		user.GroupCode = merchant.GroupCode

	}
	// 用户名不能重复
	_, err = mongo.UserColl.FindOneUser(user.UserName, "", "")
	if err == nil {
		log.Errorf("用户名已存在,userName=%s", user.UserName)
		return model.NewResultBody(3, "用户名已存在")
	}
	// 邮箱不能重复
	// _, err = mongo.UserColl.FindOneUser("", user.Mail, "")
	// if err == nil {
	// 	log.Errorf("邮箱已存在,userName=%s", user.UserName)
	// 	return model.NewResultBody(4, "邮箱已存在")
	// }
	// 手机号码不能重复
	// _, err = mongo.UserColl.FindOneUser("", "", user.PhoneNum)
	// if err == nil {
	// 	log.Errorf("用手机号码已存在,userName=%s", user.UserName)
	// 	return model.NewResultBody(5, "手机号码已存在")
	// }

	err = mongo.UserColl.Add(user)
	if err != nil {
		log.Errorf("创建用户失败,%s", err)
		return model.NewResultBody(6, "创建用户失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "创建用户成功",
	}
	return ret
}

// 修改用户信息
func (u *userController) UpdateUser(data []byte) (ret *model.ResultBody) {
	log.Debugf("update user:%s", string(data))
	user := &model.User{}
	err := json.Unmarshal(data, user)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(1, "json失败")
	}
	if user.UserName == "" || user.Password == "" {
		log.Errorf("用户名和密码不能为空")
		return model.NewResultBody(2, "用户名和密码不能为空")
	}
	if user.UserName == "" {
		log.Errorf("必填项不能为空")
		return model.NewResultBody(2, "必填项不能为空")
	}
	if user.UserType == model.UserTypeAgent {
		if user.AgentCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		_, err := mongo.AgentColl.Find(user.AgentCode)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此代理代码")
			}
			log.Errorf("查询代理(%s)出错:%s", user.AgentCode, err)
			return model.NewResultBody(1, "查询失败")
		}
	} else if user.UserType == model.UserTypeCompany {
		if user.SubAgentCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		subAgent, err := mongo.SubAgentColl.Find(user.SubAgentCode)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此公司代码")
			}
			log.Errorf("查询二级代理(%s)出错:%s", user.SubAgentCode, err)
			return model.NewResultBody(1, "查询失败")
		}
		user.AgentCode = subAgent.AgentCode

	} else if user.UserType == model.UserTypeMerchant {
		if user.GroupCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		group, err := mongo.GroupColl.Find(user.GroupCode)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此商户代码")
			}
			log.Errorf("查询集团(%s)出错: %s", user.GroupCode, err)
			return model.NewResultBody(1, "查询失败")
		}
		user.AgentCode = group.AgentCode
		user.SubAgentCode = group.SubAgentCode
	} else if user.UserType == model.UserTypeShop {
		if user.MerId == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		merchant, err := mongo.MerchantColl.FindNotInCache(user.MerId)
		if err != nil {
			if err.Error() == "not found" {
				return model.NewResultBody(3, "无此门店代码")
			}
			log.Errorf("查询一个商户(%s)出错: %s", user.MerId, err)
			return model.NewResultBody(1, "查询失败")
		}

		user.AgentCode = merchant.AgentCode
		user.SubAgentCode = merchant.SubAgentCode
		user.GroupCode = merchant.GroupCode

	}

	err = mongo.UserColl.Update(user)
	if err != nil {
		log.Errorf("更新用户失败,%s", err)
		return model.NewResultBody(5, "更新用户失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "更新用户成功",
	}
	return ret
}

// 分页查找用户
func (u *userController) Find(user *model.User, size, page int) (ret *model.ResultBody) {

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	users, total, err := mongo.UserColl.PaginationFind(user, size, page)
	if err != nil {
		log.Errorf("查询所有用户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(users),
		Data:  users,
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return ret
}

// 删除用户
func (u *userController) RemoveUser(userName string) (ret *model.ResultBody) {
	if userName == "" {
		log.Debugf("用户名不能为空")
		return model.NewResultBody(1, "用户名不能为空")
	}
	err := mongo.UserColl.Remove(userName)
	if err != nil {
		log.Debugf("删除用户失败，%s", err)
		return model.NewResultBody(2, "删除用户失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "删除用户成功",
	}
	return ret
}

func (u *userController) UpdatePwd(data []byte) (ret *model.ResultBody) {
	log.Debugf("data:%s", string(data))
	userPwd := &model.UserPwd{}
	err := json.Unmarshal(data, userPwd)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(1, "json失败")
	}

	// 原密码解密
	pwd, err := rsaDecryptFromBrowser(userPwd.Password)
	if err != nil {
		log.Errorf("escrypt password error %s", err)
		return model.NewResultBody(2, "DECRYPT_ERROR")
	}

	var appUser *model.AppUser
	oldPwdEncrypt, oldPwd := "", ""

	user, err := mongo.UserColl.FindOneUser(userPwd.UserName, "", "")
	if err != nil {
		appUser, err = mongo.AppUserCol.FindOne(userPwd.UserName)
		if err != nil {
			return model.NewResultBody(2, "查询数据库失败")
		} else {
			pb := md5.Sum([]byte(pwd))
			oldPwd, oldPwdEncrypt = appUser.Password, fmt.Sprintf("%x", string(pb[:]))
		}

	} else {
		oldPwd = user.Password
		oldPwdEncrypt = fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD + "{" + userPwd.UserName + "}" + pwd))))
	}

	// 校验
	if oldPwd != oldPwdEncrypt {
		return model.NewResultBody(3, "OLD_PASSWORD_NOT_MATCH")
	}

	// 新密码解密
	newPwd, err := rsaDecryptFromBrowser(userPwd.NewPwd)
	if err != nil {
		log.Errorf("escrypt password error %s", err)
		return model.NewResultBody(2, "DECRYPT_ERROR")
	}

	// 不能和旧密码一样
	if newPwd == pwd {
		log.Error("new password is equal with old password")
		return model.NewResultBody(8, "NEW_PASSWORD_IS_EQUAL_WITH_OLD_PASSWORD")
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
		if err = mongo.AppUserCol.Update(appUser); err != nil {
			return model.NewResultBody(4, "修改密码失败")
		}
	} else {
		user.Password = fmt.Sprintf("%x", sha1.Sum([]byte(model.RAND_PWD+"{"+userPwd.UserName+"}"+newPwd)))
		if err = mongo.UserColl.Update(user); err != nil {
			return model.NewResultBody(4, "修改密码失败")
		}
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "修改密码成功",
	}
	return ret
}

var mustHaveInPwd = []func(rune) bool{
	unicode.IsUpper,
	unicode.IsLower,
	unicode.IsDigit,
}

// isPasswordOk 判断密码是否足够复杂
func isPasswordOk(p string) bool {
	for _, testRune := range mustHaveInPwd {
		found := false
		for _, r := range p {
			if testRune(r) {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (u *userController) ResetPwd(data []byte, curUser *model.User) (ret *model.ResultBody) {
	log.Debugf("data:%s", string(data))
	userPwd := &model.UserPwd{}
	err := json.Unmarshal(data, userPwd)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(1, "json失败")
	}

	user, err := mongo.UserColl.FindOneUser(userPwd.UserName, "", "")
	if err != nil {
		if err.Error() == "not found" {
			return model.NewResultBody(1, "USERNAME_NOT_FOUND")
		}
		log.Errorf("select user by userName err,userName=%s,%s", userPwd.UserName, err)
		return model.NewResultBody(2, "查询数据库失败")
	}

	// 该用户是否有权限
	if curUser.UserType == model.UserTypeAgent {
		if user.AgentCode != curUser.AgentCode {
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

	// passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + model.DEFAULT_PWD)
	passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + newPwd)
	user.Password = fmt.Sprintf("%x", sha1.Sum(passData))

	err = mongo.UserColl.Update(user)
	if err != nil {
		log.Errorf("reset password err,userName=%s,%s", user.UserName, err)
		return model.NewResultBody(8, "重置密码失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "重置密码成功",
	}
	return ret
}
