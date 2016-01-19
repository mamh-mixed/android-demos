package mongo

import (
	"errors"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type bindingInfoCollection struct {
	name string
}

var BindingInfoColl = bindingInfoCollection{"bindingInfo"}

// Insert 插入一条绑定信息到数据库中
func (c *bindingInfoCollection) Insert(bi *model.BindingInfo) error {
	return database.C(c.name).Insert(bi)
}

// Find 根据商户号和绑定ID查找一条商家绑定信息
func (c *bindingInfoCollection) Find(merId, bindingId string) (bi *model.BindingInfo, err error) {
	bi = new(model.BindingInfo)
	q := bson.M{"bindingId": bindingId, "merId": merId}
	err = database.C(c.name).Find(q).One(bi)
	if err != nil {
		log.Errorf("Error message is: %s\n;'FindBindingInfo' condition: %+v", err.Error(), q)
		return nil, err
	}

	return bi, err
}

// Update 更新一条商家绑定信息
func (c *bindingInfoCollection) Update(bi *model.BindingInfo) error {
	if bi.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if bi.MerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"bindingId": bi.BindingId, "merId": bi.MerId}
	err := database.C(c.name).Update(q, bi)
	return err
}
