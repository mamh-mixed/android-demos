package master

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

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
	user, err := mongo.UserColl.FindOneUser(userName, "", "")
	if err != nil {
		log.Errorf("find user(%s) error: %s", userName, err)
		return model.NewResultBody(2, "无此用户名")
	}
	passSha1 := fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD + "{" + userName + "}" + password))))
	if passSha1 != user.Password {
		log.Errorf("wrong password")
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
	} else if user.UserType == model.UserTypeCompany {
		if user.AgentCode == "" || user.SubAgentCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
	} else if user.UserType == model.UserTypeMerchant {
		if user.AgentCode == "" || user.SubAgentCode == "" || user.GroupCode == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
	} else if user.UserType == model.UserTypeShop {
		if user.AgentCode == "" || user.SubAgentCode == "" || user.GroupCode == "" || user.MerId == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
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
	user, err := mongo.UserColl.FindOneUser(userPwd.UserName, "", "")
	if err != nil {
		return model.NewResultBody(2, "查询数据库失败")
	}
	if user.Password != fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD+"{"+userPwd.UserName+"}"+userPwd.Password)))) {
		return model.NewResultBody(3, "原密码错误")
	}
	user.Password = fmt.Sprintf("%x", sha1.Sum([]byte(model.RAND_PWD+"{"+userPwd.UserName+"}"+userPwd.NewPwd)))
	err = mongo.UserColl.Update(user)
	if err != nil {
		log.Infof("修改密码失败,%s", err)
		return model.NewResultBody(4, "修改密码失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "修改密码成功",
	}
	return ret
}

func (u *userController) ResetPwd(userName string) (ret *model.ResultBody) {
	user, err := mongo.UserColl.FindOneUser(userName, "", "")
	if err != nil {
		if err.Error() == "not found" {
			return model.NewResultBody(1, "无此用户名")
		}
		log.Errorf("select user by userName err,userName=%s,%s", userName, err)
		return model.NewResultBody(2, "查询数据库失败")
	}
	passData := []byte(model.RAND_PWD + "{" + user.UserName + "}" + model.DEFAULT_PWD)
	user.Password = fmt.Sprintf("%x", sha1.Sum(passData))
	err = mongo.UserColl.Update(user)
	if err != nil {
		log.Errorf("reset password err,userName=%s,%s", userName, err)
		return model.NewResultBody(3, "重置密码失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "重置密码成功",
	}
	return ret
}
