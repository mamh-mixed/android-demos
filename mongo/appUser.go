package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	UPDATE_LOCK_TIME        = 0
	UPDATE_DEVICE_INFO      = 1
	UPDATE_DEVICE_LOCK_INFO = 2
)

type appUserCollection struct {
	name string
}

var AppUserCol = appUserCollection{"appUser"}

func (col *appUserCollection) Upsert(user *model.AppUser) (err error) {

	bo := bson.M{
		"username": user.UserName,
	}
	_, err = database.C(col.name).Upsert(bo, user)
	if err != nil {
		log.Errorf("upsert user err,%s", err)
		return err
	}
	return nil
}

// Update 修改密码
func (col *appUserCollection) Update(user *model.AppUser) (err error) {

	bo := bson.M{
		"username": user.UserName,
	}
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err = database.C(col.name).Update(bo, user)
	return err
}

// BatchAdd 批量增加
func (col *appUserCollection) BatchAdd(users []model.AppUser) (err error) {
	var temp []interface{}
	for _, user := range users {
		temp = append(temp, user)
	}
	return database.C(col.name).Insert(temp...)
}

// FindByMerId 根据关联的merId查找
func (col *appUserCollection) FindByMerId(merId string) ([]*model.AppUser, error) {
	var appUsers []*model.AppUser
	err := database.C(col.name).Find(bson.M{"merId": merId}).All(&appUsers)
	return appUsers, err
}

func (col *appUserCollection) FindOne(userName string) (user *model.AppUser, err error) {
	bo := bson.M{
		"username": userName,
	}
	user = new(model.AppUser)
	err = database.C(col.name).Find(bo).One(user)
	if err != nil {
		log.Errorf("find user by userName err,userName=%s", err)
		return nil, err
	}
	return user, nil
}

func (col *appUserCollection) Count(userName string) (int, error) {
	return database.C(col.name).Find(bson.M{
		"username": userName,
	}).Count()
}

// Find 根据条件查找
func (col *appUserCollection) Find(q *model.AppUserContiditon) ([]*model.AppUser, error) {

	query := bson.M{}
	if q.SubAgentCode != "" {
		query["subAgentCode"] = q.SubAgentCode
	}
	if q.RegisterFrom != 0 {
		query["registerFrom"] = q.RegisterFrom
	}
	if q.StartTime != "" && q.EndTime != "" {
		query["createTime"] = bson.M{"$gte": q.StartTime, "$lte": q.EndTime}
	}

	var users []*model.AppUser
	err := database.C(col.name).Find(query).All(&users)
	return users, err
}

// FindPromoteLimit 查找申请提升限额的用户
func (col *appUserCollection) FindPromoteLimit(date string) ([]*model.AppUser, error) {

	q := bson.M{
		"promoteLimitRecord": bson.M{
			"$elemMatch": bson.M{"$regex": bson.RegEx{date, "."}},
		},
	}
	var users []*model.AppUser
	err := database.C(col.name).Find(q).All(&users)
	return users, err
}

func (col *appUserCollection) FindCountByUserName(userName string) (num int, err error) {
	bo := bson.M{
		"username": userName,
	}
	num, err = database.C(col.name).Find(bo).Count()
	if err != nil {
		log.Errorf("find count by userName err,userName=%s", err)
		return 0, err
	}
	return num, nil
}

func (col *appUserCollection) UpdateLoginTime(userName, loginTime, lockTime string) error {
	bo := bson.M{
		"username": userName,
	}

	update := bson.M{
		"$set": bson.M{"loginTime": loginTime,
			"lockTime": lockTime},
	}
	err := database.C(col.name).Update(bo, update)
	return err
}

func (col *appUserCollection) UpdateAppUser(user *model.AppUser, oprType int) error {

	bo := bson.M{
		"username": user.UserName,
	}
	update := bson.M{}

	switch oprType {
	case UPDATE_LOCK_TIME:
		update = bson.M{
			"$set": bson.M{
				"loginTime": user.LoginTime,
				"lockTime":  user.LockTime},
		}
	case UPDATE_DEVICE_INFO:
		update = bson.M{
			"$set": bson.M{
				"deviceType":  user.DeviceType,
				"deviceToken": user.DeviceToken},
		}
	case UPDATE_DEVICE_LOCK_INFO:
		update = bson.M{
			"$set": bson.M{
				"loginTime":   user.LockTime,
				"lockTime":    user.LockTime,
				"deviceType":  user.DeviceType,
				"deviceToken": user.DeviceToken},
		}
	}

	err := database.C(col.name).Update(bo, update)
	return err
}
