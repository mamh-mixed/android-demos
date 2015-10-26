package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
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
