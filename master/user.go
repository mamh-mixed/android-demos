package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type user struct{}

var User user

// Login 登陆
// func (u *user) Login(userName, password string) (ret *model.ResultBody) {
// 	log.Debugf("userName:%s,password:%s", userName, password)
// 	if userName == "" || password == "" {
// 		log.Errorf("用户名或密码不能为空")
// 		return model.NewResultBody(1, "用户名或密码不用为空")
// 	}
// 	user, err := mongo.UserColl.FindByUserName(userName)
// 	if err != nil {
// 		log.Errorf("查询用户(%s)出错:%s", userName, err)
// 		return model.NewResultBody(2, "用户名错误")
// 	}
// 	if password != user.Password {
// 		log.Errorf("密码错误")
// 		return model.NewResultBody(3, "密码错误")
// 	}
// 	// 先每次登陆从数据中查找1级菜单
// 	menusTemp, err := mongo.MenuColl.FindByLevel(1)
// 	if err != nil {
// 		log.Errorf("查找所有1级菜单失败")
// 		return model.NewResultBody(4, "查找所有1级菜单失败")
// 	}
// 	// 传给前端的菜单
// 	menus := make([]model.Menu, 1)
// 	for _, menu1 := range menusTemp {
// 		menu1.Children = make([]model.Menu, 1)
// 		for _, menu2 := range user.Role.Menus {
// 			if menu2.ParentRoute == menu1.Route {
// 				menu1.Children = append(menu1.Children, menu2)
// 			}
// 		}
// 		if len(menu1.Children) > 0 {
// 			menus = append(menus, menu1)
// 		}
// 	}
//
// 	ret = &model.ResultBody{
// 		Status:  0,
// 		Message: "登陆成功",
// 		Data:    menus,
// 	}
// 	return ret
// }

// 新建用户
func (u *user) CreateUser(data []byte) (ret *model.ResultBody) {
	log.Debugf("create user:%s", string(data))
	user := &model.User{}
	err := json.Unmarshal(data, user)
	if err != nil {
		log.Errorf("json unmsrshal err,%s", err)
		return model.NewResultBody(1, "json失败")
	}
	// 设置默认密码
	user.Password = "12345678"

	// 判断必填项是否为空
	if user.UserName == "" || user.Mail == "" || user.PhoneNum == "" {
		log.Errorf("必填项不能为空")
		return model.NewResultBody(2, "必填项不能为空")
	}
	if user.UserType == "agent" {
		if user.AgentCode == "" || user.AgentName == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
	} else if user.UserType == "group" {
		if user.GroupCode == "" || user.GroupName == "" {
			log.Errorf("必填项不能为空")
			return model.NewResultBody(2, "必填项不能为空")
		}
	} else if user.UserType == "merchant" {
		if user.MerId == "" || user.MerName == "" {
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
	_, err = mongo.UserColl.FindOneUser("", user.Mail, "")
	if err == nil {
		log.Errorf("邮箱已存在,userName=%s", user.UserName)
		return model.NewResultBody(4, "邮箱已存在")
	}
	// 手机号码不能重复
	_, err = mongo.UserColl.FindOneUser("", "", user.PhoneNum)
	if err == nil {
		log.Errorf("用手机号码已存在,userName=%s", user.UserName)
		return model.NewResultBody(5, "手机号码已存在")
	}

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
// func (u *user) UpdateUser(data []byte) (ret *model.ResultBody) {
// 	log.Debugf("update user:%s", string(data))
// 	user := &model.User{}
// 	err := json.Unmarshal(data, user)
// 	if err != nil {
// 		log.Errorf("json unmsrshal err,%s", err)
// 		return model.NewResultBody(1, "json失败")
// 	}
// 	if user.UserName == "" || user.Password == "" {
// 		log.Errorf("用户名和密码不能为空")
// 		return model.NewResultBody(2, "用户名和密码不能为空")
// 	}
// 	err = mongo.UserColl.Update(user)
// 	if err != nil {
// 		log.Errorf("更新用户失败,%s", err)
// 		return model.NewResultBody(5, "更新用户失败")
// 	}
// 	ret = &model.ResultBody{
// 		Status:  0,
// 		Message: "更新用户成功",
// 	}
// 	return ret
// }

// 分页查找用户
func (u *user) Find(user *model.User, size, page int) (ret *model.ResultBody) {

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
func (u *user) RemoveUser(userName string) (ret *model.ResultBody) {
	log.Debugf("removeUser,userName=%s", userName)
	if userName == "" {
		log.Errorf("用户名不能为空")
		return model.NewResultBody(1, "用户名不能为空")
	}
	err := mongo.UserColl.Remove(userName)
	if err != nil {
		log.Errorf("删除用户失败，%s", err)
		return model.NewResultBody(2, "删除用户失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "删除用户成功",
	}
	return ret
}
