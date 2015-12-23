package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type pushMessage struct {
	name string
}

// 推送消息 Collection
var PushMessageColl = pushMessage{"pushMessage"}

func (col *pushMessage) Insert(t *model.PushMessageRsp) error {

	element := new(model.PushInfo)
	element.UserName = t.UserName
	element.Title = t.Title
	element.Message = t.Message
	element.PushTime = t.PushTime

	return database.C(col.name).Insert(element)
}

func (col *pushMessage) FindByUser(t *model.PushMessageRsp) (results []*model.PushInfo, err error) {
	results = make([]*model.PushInfo, 0)

	con := bson.M{}
	if t.LastTime != "" {
		con["pushtime"] = bson.M{"$gt": t.LastTime}
	}

	con["username"] = t.UserName

	err = database.C(col.name).Find(con).All(&results)

	return results, err
}
