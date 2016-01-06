package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type pushMessage struct {
	name string
}

// 推送消息 Collection
var PushMessageColl = pushMessage{"pushMessage"}

// Insert 插入一条消息
func (col *pushMessage) Insert(t *model.PushMessage) error {
	return database.C(col.name).Insert(t)
}

// UpdateByID 更新某条信息状态
func (col *pushMessage) UpdateStatusByID(msgId string, status int) error {

	update := bson.M{
		"$set": bson.M{
			"updateTime": time.Now().Format("20060102150405"),
			"status":     status,
		},
	}
	return database.C(col.name).Update(bson.M{"msgId": msgId}, update)
}

// Find 查找
func (col *pushMessage) Find(t *model.PushMessage) (result []model.PushMessage, err error) {

	con := bson.M{}
	if t.LastTime != "" {
		con["pushtime"] = bson.M{"$gt": t.LastTime}
	}
	if t.MaxTime != "" {
		con["pushtime"] = bson.M{"$lt": t.MaxTime}
	}

	con["username"] = t.UserName
	con["status"] = 0
	q := database.C(col.name).Find(con).Sort("-pushtime")
	if t.Size != 0 {
		err = q.Limit(t.Size).All(&result)
	} else {
		err = q.All(&result)
	}

	return result, err
}
