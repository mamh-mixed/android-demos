package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type userCollection struct {
	name string
}

// UserColl 用户 Collection
var UserColl = userCollection{"user"}

// FindOneUser 根据userName,mail,phoneNum查找
func (col *userCollection) FindOneUser(userName, mail, phoneNum string) (u *model.User, err error) {

	bo := bson.M{}
	if userName != "" {
		bo["userName"] = userName
	}
	if mail != "" {
		bo["mail"] = mail
	}
	if phoneNum != "" {
		bo["phoneNum"] = phoneNum
	}
	u = new(model.User)
	err = database.C(col.name).Find(bo).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Add 增加一个用户
func (col *userCollection) Add(u *model.User) error {
	err := database.C(col.name).Insert(u)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新用户信息
func (col *userCollection) Update(u *model.User) error {
	bo := bson.M{
		"userName": u.UserName,
	}
	err := database.C(col.name).Update(bo, u)
	if err != nil {
		return err
	}
	return nil
}

// PaginationFind 分页查找user
func (col *userCollection) PaginationFind(user *model.User, size, page int) (results []model.User, total int, err error) {
	results = make([]model.User, 1)

	match := bson.M{}
	if user.UserName != "" {
		match["userName"] = user.UserName
	}
	if user.NickName != "" {
		match["nickName"] = user.NickName
	}
	if user.Mail != "" {
		match["mail"] = user.Mail
	}
	if user.PhoneNum != "" {
		match["phoneNum"] = user.PhoneNum
	}
	if user.UserType != "" {
		match["userType"] = user.UserType
	}
	if user.AgentCode != "" {
		match["agentCode"] = user.AgentCode
	}
	// if user.AgentCode != "" {
	// 	match["subAgentCode"] = user.SubAgentCode
	// }
	if user.GroupCode != "" {
		match["groupCode"] = user.GroupCode
	}
	if user.MerId != "" {
		match["merId"] = user.MerId
	}

	// 计算总数
	total, err = database.C(col.name).Find(match).Count()
	if err != nil {
		log.Errorf("PaginationFind Count err,%s", err)
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	sort := bson.M{"$sort": bson.M{"userName": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, sort, skip, limit)

	err = database.C(col.name).Pipe(cond).All(&results)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Remove 删除用户
func (col *userCollection) Remove(userName string) (err error) {
	bo := bson.M{}
	if userName != "" {
		bo["userName"] = userName
	}
	err = database.C(col.name).Remove(bo)
	return err
}
