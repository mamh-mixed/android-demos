package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
	"strconv"
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

	if t.Index == "" {
		t.Index = "0"
	}

	if t.Size == "" {
		t.Size = "10"
	}

	index, err := strconv.Atoi(t.Index)
	if err != nil {
		index = 0
	}

	size, err := strconv.Atoi(t.Size)
	if err != nil {
		size = 10
	}

	err = database.C(col.name).Find(match).Sort("-pushtime").Skip(index).Limit(size).All(&results)

	return results, err
}
