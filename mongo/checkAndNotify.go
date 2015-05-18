package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type notifyCollection struct {
	name string
}

var NotifyColl = notifyCollection{"checkAndNotify"}

// GetAll 加载所有需要缓存的类型信息
func (n *notifyCollection) GetAll() ([]*model.CheckAndNotify, error) {

	var cans []*model.CheckAndNotify

	err := database.C(n.name).Find(nil).All(&cans)

	return cans, err
}

// Update 更新
func (n *notifyCollection) Update(update *model.CheckAndNotify) error {

	// selector
	q := bson.M{
		"bizType": update.BizType,
	}
	return database.C(n.name).Update(q, update)
}
