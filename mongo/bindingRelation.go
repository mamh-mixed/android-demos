package mongo

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"

	"github.com/omigo/g"
)

type BindingRelation struct {
	CardInfo         model.BindingCreate `json:"cardInfo" bson:"cardInfo,omitempty"`                 //卡片信息
	Router           RouterPolicy        `json:"router" bson:"router,omitempty"`                     //路由信息
	ChannelBindingId string              `json:"channelBindingId" bson:"channelBindingId,omitempty"` //渠道绑定ID
}

// 插入一条绑定关系到数据库中
func InsertOneBindingRelation(br *BindingRelation) error {
	if err := db.bindingRelation.Insert(br); err != nil {
		return err
	}
	return nil
}

// 根据源商户号和绑定ID查找一条绑定关系
func FindOneBindingRelation(merCode, bindingId string) (br *BindingRelation, err error) {
	br = new(BindingRelation)
	q := bson.M{"cardInfo.bindingId": bindingId, "router.origMerId": merCode}
	g.Debug("'FindOneBindingRelation' condition: %+v", q)
	err = db.bindingRelation.Find(q).One(br)
	return br, err
}

// 更新一条绑定关系
func UpdateOneBindingRelation(br *BindingRelation) error {
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
