package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

type cilMsgCollection struct {
	name string
}

var CilMsgColl = cilMsgCollection{"cilMsg"}

// Upsert 如果已存在一个CIL报文，就更新，不存在就插入
func (col *cilMsgCollection) Upsert(m *model.CilMsg) (err error) {
	cond := bson.M{
		"uuid": m.UUID,
	}

	_, err = database.C(col.name).Upsert(cond, m)

	if err != nil {
		return err
	}

	return
}
