package mongo

import (
	"errors"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type bindingMapCollection struct {
	name string
}

// BindingMapColl 绑定关系 Collection
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
		log.Errorf("Error message is: %s;'FindBindingMap' condition: %+v", err.Error(), q)
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
	if err != nil {
		log.Errorf("update bindingMap(%+v) fail: %s", bm, err)
	}
	return err
}

func (c *bindingMapCollection) Count(merId, bindingId string) (count int, err error) {
	if merId == "" {
		return 0, errors.New("商户ID为空")
	}

	if bindingId == "" {
		return 0, errors.New("绑定ID为空")
	}

	q := bson.M{"bindingId": bindingId, "merId": merId}
	count, err = database.C(c.name).Find(q).Count()
	return
}
