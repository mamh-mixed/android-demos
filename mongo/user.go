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

// FindByUserName 根据userName查找
func (col *userCollection) FindByUserName(userName string) (u *model.User, err error) {

	bo := bson.M{
		"userName": userName,
	}
	u = new(model.User)
	err = database.C(col.name).Find(bo).One(u)
	if err != nil {
		log.Errorf("Find User condition is: %+v;error is %s", bo, err)
		return nil, err
	}
	return u, nil
}

// Add 增加一个用户
func (col *userCollection) Add(u *model.User) error {
	err := database.C(col.name).Insert(u)
	if err != nil {
		log.Debugf("Add user err,%s", err)
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
		log.Debugf("update user err,%s", err)
		return err
	}
	return nil
}

// PaginationFind 分页查找user
func (col *userCollection) PaginationFind(userName, nickName, roleName string, size, page int) (results []model.User, total int, err error) {
	results = make([]model.User, 1)

	match := bson.M{}
	if userName != "" {
		match["userName"] = userName
	}
	if nickName != "" {
		match["nickName"] = nickName
	}
	if roleName != "" {
		match["role.name"] = roleName
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

	// sort := bson.M{"$sort": bson.M{"name": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, skip, limit)

	err = database.C(col.name).Pipe(cond).All(&results)
	if err != nil {
		log.Errorf("find user paging err,%s", err)
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
