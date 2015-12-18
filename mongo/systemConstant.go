package mongo

import (
	"github.com/CardInfoLink/quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

type systemConstantCollection struct {
	name string
}

const SysConstId = "QUICKPAY_SYS_CONST"

var SysConstColl = systemConstantCollection{"systemConstant"}

// Find 查找系统常量
func (c *systemConstantCollection) Find() (result *model.SystemConstant, err error) {
	result = new(model.SystemConstant)

	err = database.C(c.name).Find(bson.M{"id": SysConstId}).One(result)

	return result, err
}

// Find 查找系统常量
func (c *systemConstantCollection) Upsert(cond *model.SystemConstant) (err error) {
	cond.ID = SysConstId
	_, err = database.C(c.name).Upsert(bson.M{"id": SysConstId}, cond)
	return err
}
