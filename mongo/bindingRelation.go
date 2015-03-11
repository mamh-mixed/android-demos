package mongo

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
)

type BindingRelation struct {
	CardInfo         model.BindingCreate `json:"cardInfo" bson:"cardInfo,omitempty"`                 //卡片信息
	Router           RouterPolicy        `json:"router" bson:"router,omitempty"`                     //路由信息
	ChannelBindingId string              `json:"channelBindingId" bson:"channelBindingId,omitempty"` //渠道绑定ID
}

const BindingRelationCollectionName = "bindingRelation"

// 插入一条绑定关系到数据库中
func InsertOneBindingRelation(br *BindingRelation) error {
	c := db.C(BindingRelationCollectionName)
	if err := c.Insert(br); err != nil {
		return err
	}
	return nil
}

// 根据源商户号和绑定ID查找一条绑定关系
func FindOneBindingRelationByMerCodeAndBindingId(merCode, bindingId string) (br *BindingRelation, err error) {
	br = new(BindingRelation)
	c := db.C(BindingRelationCollectionName)
	err = c.Find(bson.M{"cardInfo.bindingId": bindingId, "router.origMerCode": merCode}).One(br)
	return br, err
}

// 更新一条绑定关系
func UpdateOneBindingRelation(br *BindingRelation) error {
	if br.CardInfo.BindingId == "" {
		return errors.New("BindingId must required")
	}

	if br.Router.OrigMerCode == "" {
		return errors.New("OrigMerCode must required")
	}

	c, q := db.C(BindingRelationCollectionName), bson.M{"cardInfo.bindingId": br.CardInfo.BindingId, "router.origMerCode": br.Router.OrigMerCode}
	err := c.Update(q, br)
	return err
}
