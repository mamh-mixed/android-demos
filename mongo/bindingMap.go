package mongo

import (
	"errors"
	"quickpay/model"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

type bindingMapCollection struct {
	name string
}

var BindingMapColl = bindingMapCollection{"bindingMap"}

// Insert 插入一条绑定映射关系到数据库中
func (c *bindingMapCollection) Insert(bm *model.BindingMap) error {
	return database.C(c.name).Insert(bm)
}

// Find 根据商户号和绑定ID查找一条绑定关系映射
func (c *bindingMapCollection) Find(merId, bindingId string) (bm *model.BindingMap, err error) {
	bm = new(model.BindingMap)
	q := bson.M{"bindingId": bindingId, "merId": merId}
	err = database.C(c.name).Find(q).One(bm)
	if err != nil {
		g.Error("Error message is: %s\n;'FindBindingMap' condition: %+v", err.Error(), q)
		return nil, err
	}

	return bm, err
}

// Update 更新一条绑定关系映射
func (c *bindingMapCollection) Update(bm *model.BindingMap) error {
	if bm.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if bm.MerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"bindingId": bm.BindingId, "merId": bm.MerId}
	err := database.C(c.name).Update(q, bm)
	return err
}
