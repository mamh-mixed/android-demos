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

	match := bson.M{}
	match["username"] = t.UserName
	match["pushtime"] = bson.M{"$gte": t.LastTime}

	cond := []bson.M{
		{"$match": match},
	}

	if t.Index == "" {
		t.Index = "0"
	}

	if t.Size == "" {
		t.Size = "10"
	}

	sort := bson.M{"$sort": bson.M{"pushtime": -1}}

	skip := bson.M{"$skip": t.Index}

	limit := bson.M{"$limit": t.Size}

	cond = append(cond, sort, skip, limit)

	err = database.C(col.name).Pipe(cond).All(&results)

	return results, err
}
