package mongo

import (
	"errors"
	"quickpay/model"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

// InsertBindingInfo 插入一条绑定信息到数据库中
func InsertBindingInfo(bi *model.BindingInfo) error {
	return db.bindingInfo.Insert(bi)
}

// InsertBindingMap 插入一条绑定映射关系到数据库中
func InsertBindingMap(bm *model.BindingMap) error {
	return db.bindingMap.Insert(bm)
}

// FindBindingInfo 根据商户号和绑定ID查找一条商家绑定信息
func FindBindingInfo(merId, bindingId string) (bi *model.BindingInfo, err error) {
	bi = new(model.BindingInfo)
	q := bson.M{"bindingId": bindingId, "merId": merId}
	err = db.bindingInfo.Find(q).One(bi)
	if err != nil {
		g.Error("Error message is: %s\n;'FindBindingInfo' condition: %+v", err.Error(), q)
		return nil, err
	}

	return bi, err
}

// FindBindingMap 根据商户号和绑定ID查找一条绑定关系映射
func FindBindingMap(merId, bindingId string) (bm *model.BindingMap, err error) {
	bm = new(model.BindingMap)
	q := bson.M{"bindingId": bindingId, "merId": merId}
	err = db.bindingMap.Find(q).One(bm)
	if err != nil {
		g.Error("Error message is: %s\n;'FindBindingMap' condition: %+v", err.Error(), q)
		return nil, err
	}

	return bm, err
}

// UpdateBindingInfo 更新一条商家绑定信息
func UpdateBindingInfo(bi *model.BindingInfo) error {
	if bi.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if bi.MerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"bindingId": bi.BindingId, "merId": bi.MerId}
	err := db.bindingInfo.Update(q, bi)
	return err
}

// UpdateBindingRelation 更新一条绑定关系映射
func UpdateBindingMap(bm *model.BindingMap) error {
	if bm.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if bm.MerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"bindingId": bm.BindingId, "merId": bm.MerId}
	err := db.bindingMap.Update(q, bm)
	return err
}
