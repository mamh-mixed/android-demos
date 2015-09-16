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
		"username": user.UserName,
	}
	user = new(model.AppUser)
	err = database.C(col.name).Find(bo).One(user)
	if err != nil {
		log.Errorf("find user by userName err,userName=%s", err)
		return nil, err
	}
	return user, nil
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
