package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var NotifyRecColl = notifyRecCollection{"notifyRecord"}

type notifyRecCollection struct {
	Name string
}

// Add 增加一条消息
func (n *notifyRecCollection) Add(r *model.NotifyRecord) error {

	r.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	return database.C(n.Name).Insert(r)
}

// Count 计算某个订单是否已经接受过异步消息消息通知
func (n *notifyRecCollection) Count(merId, orderNum string) (int, error) {
	return database.C(n.Name).Find(bson.M{"merId": merId, "orderNum": orderNum}).Count()
}

// Update 更新一条记录
func (n *notifyRecCollection) Update(r *model.NotifyRecord) error {
	r.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	return database.C(n.Name).Update(bson.M{"merId": r.MerId, "orderNum": r.OrderNum}, r)
}

// FindOne
func (n *notifyRecCollection) FindOne(merId, orderNum string) (*model.NotifyRecord, error) {
	result := new(model.NotifyRecord)
	err := database.C(n.Name).Find(bson.M{"merId": merId, "orderNum": orderNum}).One(result)
	return result, err
}
