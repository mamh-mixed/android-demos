package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type emailCollection struct {
	name string
}

var EmailCol = emailCollection{"email"}

func (col *emailCollection) Upsert(e *model.Email) (err error) {

	bo := bson.M{
		"username": e.UserName,
	}
	_, err = database.C(col.name).Upsert(bo, e)
	if err != nil {
		log.Errorf("upsert email err,%s", err)
		return err
	}
	return nil
}

func (col *emailCollection) FindOne(userName string) (e *model.Email, err error) {
	bo := bson.M{
		"username": userName,
	}
	e = new(model.Email)
	err = database.C(col.name).Find(bo).One(e)
	if err != nil {
		log.Errorf("find email by userName err,userName=%s", err)
		return nil, err
	}
	return e, nil
}

// FindOneByCode 根据code查找对应的发送邮箱纪录
func (col *emailCollection) FindOneByCode(code string) (e *model.Email, err error) {
	bo := bson.M{
		"code": code,
	}
	e = new(model.Email)
	err = database.C(col.name).Find(bo).Sort("-timestamp").One(e)
	if err != nil {
		log.Errorf("find email by userName err,userName=%s", err)
		return nil, err
	}
	return e, nil
}
