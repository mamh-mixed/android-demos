package mongo

import (
	"errors"
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"

	"github.com/omigo/g"
)

// BindingRelation 绑定关系
type BindingRelation struct {
	CardInfo      model.BindingCreate `json:"cardInfo" bson:"cardInfo,omitempty"`                 //卡片信息
	Router        model.RouterPolicy  `json:"router" bson:"router,omitempty"`                     //路由信息
	ChanBindingId string              `json:"channelBindingId" bson:"channelBindingId,omitempty"` //渠道绑定ID
}

// InsertBindingRelation 插入一条绑定关系到数据库中
func InsertBindingRelation(br *BindingRelation) error {
	if err := db.bindingRelation.Insert(br); err != nil {
		return err
	}
	return nil
}

// FindBindingRelation 根据源商户号和绑定ID查找一条绑定关系
func FindBindingRelation(merCode, bindingId string) (br *BindingRelation, err error) {
	br = new(BindingRelation)
	q := bson.M{"cardInfo.bindingId": bindingId, "router.origMerId": merCode}
	g.Debug("'FindBindingRelation' condition: %+v", q)
	err = db.bindingRelation.Find(q).One(br)
	return br, err
}

// UpdateBindingRelation 更新一条绑定关系
func UpdateBindingRelation(br *BindingRelation) error {
	if br.CardInfo.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if br.Router.OrigMerId == "" {
		return errors.New("OrigMerId must required")
	}

	q := bson.M{"cardInfo.bindingId": br.CardInfo.BindingId, "router.origMerId": br.Router.OrigMerId}
	err := db.bindingRelation.Update(q, br)
	return err
}
