package mongo

import (
	"errors"
	"quickpay/model"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

// InsertBindingRelation 插入一条绑定关系到数据库中
func InsertBindingRelation(br *model.BindingRelation) error {
	if err := db.bindingRelation.Insert(br); err != nil {
		return err
	}
	return nil
}

// FindBindingRelation 根据源商户号和绑定ID查找一条绑定关系
func FindBindingRelation(merCode, bindingId string) (br *model.BindingRelation, err error) {
	br = new(model.BindingRelation)
	q := bson.M{"bindingId": bindingId, "origMerId": merCode}
	g.Debug("'FindBindingRelation' condition: %+v", q)
	err = db.bindingRelation.Find(q).One(br)
	return br, err
}

// UpdateBindingRelation 更新一条绑定关系
func UpdateBindingRelation(br *model.BindingRelation) error {
	if br.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if br.OrigMerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"bindingId": br.BindingId, "origMerId": br.OrigMerId}
	err := db.bindingRelation.Update(q, br)
	return err
}

// 删除一条绑定关系
func DeleteBindingRelation(br *model.BindingRelation) error {
	if br.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if br.OrigMerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"bindingId": br.BindingId, "origMerId": br.OrigMerId}
	g.Debug("'DeleteBindingRelation' condition: %+v", q)

	return db.bindingRelation.Remove(q)
}
