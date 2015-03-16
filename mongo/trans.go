package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"quickpay/model"
	"time"
)

// Add 添加一笔交易
func AddTrans(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.TransTime = time.Now().Unix()
	t.TransFlag = 0
	return db.trans.Insert(t)
}

// Modify 通过Add时生成的Id来修改
func ModifyTrans(t *model.Trans) error {
	return db.trans.Update(bson.M{"_id": t.Id}, t)
}
