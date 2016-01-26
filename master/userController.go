package master

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"unicode"

	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/log"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

const (
	LOCKTIME      = 30000 //锁定三小时
	LOGINDIFFTIME = 10000 //1小时以内
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
	user, err := mongo.UserColl.FindOne(userName)
	var isAppUser = false
	if err != nil {

		// 兼容appUser用户登录
		appUser, err := mongo.AppUserCol.FindOne(userName)
		if err != nil {
			log.Errorf("find user(%s) error: %s", userName, err)
			return model.NewResultBody(2, "USERNAME_PASSWORD_ERROR")
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
				MerId:     appUser.MerId,
				LockTime:  appUser.LockTime,
				LoginTime: appUser.LoginTime,
			}
			isAppUser = true
		}

	} else {
		encryptPass = fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD + "{" + userName + "}" + password))))
	}

	//判断锁定情况
	if user.LockTime != "" {
		localTime := time.Now().Format("2006-01-02 15:04:05")
		stemp := ""
		for _, v := range localTime {
			if (v >= '0') && (v <= '9') {
				stemp += string(v)
			}
		}
		localTime = stemp
		date1 := user.LockTime[:8]
		date2 := localTime[:8]
		time1 := user.LockTime[8:]
		time2 := localTime[8:]
		time1Int, err := strconv.ParseInt(string(time1), 10, 32)
		if err != nil {
			log.Errorf("convert string to int error value:%s, error:%s", time1, err)
			return model.NewResultBody(3, model.SYSTEM_ERROR.Error)
		}
		time2Int, err := strconv.ParseInt(string(time2), 10, 32)
		if err != nil {
			log.Errorf("convert string to int error value:%s, error:%s", time2, err)
			return model.NewResultBody(3, model.SYSTEM_ERROR.Error)
		}
		if string(date1) != string(date2) {
			time2Int += 240000
		}

		if (time2Int - time1Int) > LOCKTIME {
			if isAppUser { //解锁
				mongo.AppUserCol.UpdateLoginTime(userName, "", "")
			} else {
				mongo.UserColl.UpdateLoginTime(userName, "", "")
			}
		} else {
			return model.NewResultBody(3, model.USER_LOCK.Error)
		}
	}

	// 校验
	if encryptPass != user.Password {
		localTime := time.Now().Format("2006-01-02 15:04:05")
		stemp := ""
		for _, v := range localTime {
			if (v >= '0') && (v <= '9') {
				stemp += string(v)
			}
		}
		localTime = stemp
		if user.LoginTime != "" {
			timeArray := strings.Split(user.LoginTime, ",")
			var count = 0
			var loginTime = ""
			date2 := localTime[:8]
			time2 := localTime[8:]
			time2Int, err := strconv.ParseInt(string(time2), 10, 32)
			if err != nil {
				log.Errorf("convert string to int is value:%s, error:%s", time2, err)
				return model.NewResultBody(3, model.SYSTEM_ERROR.Error)
			}
			for _, timeElement := range timeArray {
				date1 := timeElement[:8]
				time1 := timeElement[8:]

				time1Int, err := strconv.ParseInt(string(time1), 10, 32)
				if err != nil {
					log.Errorf("convert string to int is value:%s, error:%s", time1, err)
					continue
				}

				if string(date1) != string(date2) { //日期不同，加24小时
					time2Int += 240000
				}

				if (time2Int - time1Int) > LOGINDIFFTIME { //超过一小时，舍弃
					continue
				} else {
					count++
					//记录
					if loginTime == "" {
						loginTime = timeElement
					} else {
						loginTime += ","
						loginTime += timeElement
					}
				}
			}
			//判断count是否达到10次
			if count == 9 {
				loginTime += ","
				loginTime += localTime
				if isAppUser {
					mongo.AppUserCol.UpdateLoginTime(userName, loginTime, localTime)
				} else {
					mongo.UserColl.UpdateLoginTime(userName, loginTime, localTime)
				}
				return model.NewResultBody(3, model.USER_LOCK.Error) //锁定
			} else {
				var ret model.ResultBody
				ret.Status = 3
				if count == 6 {
					ret.Message = model.USER_THREE_TIMES.Error
				} else if count == 7 {
					ret.Message = model.USER_TWO_TIMES.Error
				} else if count == 8 {
					ret.Message = model.USER_ONE_TIMES.Error
				} else {
					ret.Message = model.USERNAME_PASSWORD_ERROR.Error
				}

				if loginTime == "" {
					loginTime = localTime
				} else {
					loginTime += ","
					loginTime += localTime
				}
				if isAppUser {
					mongo.AppUserCol.UpdateLoginTime(userName, loginTime, "")
				} else {
					mongo.UserColl.UpdateLoginTime(userName, loginTime, "")
				}

				return &ret
			}
		} else {
			if isAppUser {
				mongo.AppUserCol.UpdateLoginTime(userName, localTime, "")
			} else {
				mongo.UserColl.UpdateLoginTime(userName, localTime, "")
			}
		}
		log.Errorf("wrong password, expect %s but get %s", user.Password, encryptPass)
		return model.NewResultBody(3, "USERNAME_PASSWORD_ERROR")
	}

	//密码正确，清空登陆记录
	if isAppUser {
		mongo.AppUserCol.UpdateLoginTime(userName, "", "")
	} else {
		mongo.UserColl.UpdateLoginTime(userName, "", "")
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

	// 用户名不能重复
	num, err := mongo.UserColl.FindCountByUserName(user.UserName)
	if err != nil {
		log.Errorf("find database err,username=%s,%s", user.UserName, err)
		return model.NewResultBody(6, "系统错误，请重试")
	}
	if num != 0 {
		return model.NewResultBody(3, "用户名已存在")
	}

	// 设置默认密码
	passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + model.DEFAULT_PWD)
	user.Password = fmt.Sprintf("%x", sha1.Sum(passData))

	// 判断必填项是否为空
	if user.UserName == "" {
		return model.NewResultBody(2, "必填项不能为空")
	}

	// 验证用户类型是否合法
	if errResult := validUserType(user); errResult != nil {
		return errResult
	}

	user.LoginTime = ""
	user.LockTime = ""

	err = mongo.UserColl.Add(user)
	if err != nil {
		log.Errorf("user save error, username=%s, err: %s", user.UserName, err)
		return model.NewResultBody(6, "系统错误，请重试")
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
		return model.NewResultBody(1, "系统错误，请重试")
	}
	if user.UserName == "" || user.Password == "" {
		return model.NewResultBody(2, "用户名或密码为空")
	}

	// 验证用户类型是否合法
	if errResult := validUserType(user); errResult != nil {
		return errResult
	}

	err = mongo.UserColl.Update(user)
	if err != nil {
		log.Errorf("user update error, username=%s, err: %s", user.UserName, err)
		return model.NewResultBody(5, "系统错误，请重试")
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
		log.Debugf("remove user(%s) fail: %s", userName, err)
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

	user, err := mongo.UserColl.FindOne(userPwd.UserName)
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

	user, err := mongo.UserColl.FindOne(userPwd.UserName)
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

	if len(newPwd) >= 20 {
		log.Errorf("new password's length must less than 21; now is %d", len(newPwd))
		return model.NewResultBody(6, "PASSWORD_LENGTH_TWO_LONG")
	}

	// 密码复杂度校验
	if !isPasswordOk(newPwd) {
		log.Errorf("new password is not conplicated enough: %s", newPwd)
		return model.NewResultBody(7, "PASSWORD_NOT_COMPLICATED_ENOUGH")
	}

	// passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + model.DEFAULT_PWD)
	passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + newPwd)
	user.Password = fmt.Sprintf("%x", sha1.Sum(passData))
	user.LoginTime = ""
	user.LockTime = ""

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

// PasswordReset 密码重置整理
func (u *userController) PasswordReset(data []byte) (ret *model.ResultBody) {
	log.Debugf("data:%s", string(data))
	resetUser := &model.ResetUser{}
	err := json.Unmarshal(data, resetUser)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(6, "JSON_ERROR")
	}

	if resetUser.PassWord == "" || resetUser.CheckCode == "" {
		return model.NewResultBody(1, "MISS_REQUIRED_PARAMETER")
	}

	// 查找发送邮件的纪录
	emailInfo, err := mongo.EmailCol.FindOneByCode(resetUser.CheckCode)
	if err != nil {
		log.Errorf("find email sending record error, checkcode:%s, error:%s", resetUser.CheckCode, err)
		return model.NewResultBody(2, "INVALID_CHECK_CODE")
	}
	resetUser.UserName = emailInfo.UserName

	// 是否已经激活过
	if emailInfo.IsOperated {
		return model.NewResultBody(6, "ALREADY_OPERATED")
	}

	// 是否操作过时
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", emailInfo.Timestamp, time.Local)
	if time.Now().Sub(timestamp) > 12*time.Hour {
		return model.NewResultBody(5, "OPERATION_OUT_OF_DATE")
	}

	appUser, err := mongo.AppUserCol.FindOne(resetUser.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.NewResultBody(3, "USERNAME_NOT_MATCH")
		}
		log.Errorf("select user by userName err,userName=%s,%s", resetUser.UserName, err)
		return model.NewResultBody(2, "SYSTEM_ERROR")
	}

	appUser.Password = resetUser.PassWord
	appUser.LoginTime = ""
	appUser.LockTime = ""
	err = mongo.AppUserCol.Update(appUser)
	if err != nil {
		log.Errorf("reset password err,userName=%s,%s", resetUser.UserName, err)
		return model.NewResultBody(2, "SYSTEM_ERROR")
	}

	// 更新到EMAIL集合
	e := &model.Email{
		UserName:   resetUser.UserName,
		Code:       resetUser.CheckCode,
		Success:    true,
		IsOperated: true,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
	}
	mongo.EmailCol.Upsert(e)

	ret = &model.ResultBody{
		Status:  0,
		Message: "SUCCESS",
	}
	return ret
}

// validUserType 校验用户类型是否合法
func validUserType(user *model.User) (ret *model.ResultBody) {
	if user.UserType == model.UserTypeAgent {
		if user.AgentCode == "" {
			return model.NewResultBody(2, "必填项不能为空")
		}
		agent, err := mongo.AgentColl.Find(user.AgentCode)
		if err != nil {
			return model.NewResultBody(3, "无此代理代码")
		}
		user.AgentName = agent.AgentName

	} else if user.UserType == model.UserTypeCompany {
		if user.SubAgentCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		subAgent, err := mongo.SubAgentColl.Find(user.SubAgentCode)
		if err != nil {
			return model.NewResultBody(3, "无此公司代码")
		}
		user.AgentCode = subAgent.AgentCode
		user.AgentName = subAgent.AgentName
		user.SubAgentName = subAgent.SubAgentName

	} else if user.UserType == model.UserTypeMerchant {
		if user.GroupCode == "" {
			return model.NewResultBody(2, "必填项不能为空")
		}
		group, err := mongo.GroupColl.Find(user.GroupCode)
		if err != nil {
			return model.NewResultBody(3, "无此商户代码")
		}
		user.AgentCode = group.AgentCode
		user.AgentName = group.AgentName
		user.SubAgentCode = group.SubAgentCode
		user.SubAgentName = group.SubAgentName
		user.GroupName = group.GroupName

	} else if user.UserType == model.UserTypeShop {
		if user.MerId == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
		merchant, err := mongo.MerchantColl.FindNotInCache(user.MerId)
		if err != nil {
			return model.NewResultBody(3, "无此门店代码")
		}
		user.AgentCode = merchant.AgentCode
		user.AgentName = merchant.AgentName
		user.SubAgentCode = merchant.SubAgentCode
		user.SubAgentName = merchant.SubAgentName
		user.GroupCode = merchant.GroupCode
		user.GroupName = merchant.GroupName
		user.MerName = merchant.Detail.MerName
	}
	return nil
}
