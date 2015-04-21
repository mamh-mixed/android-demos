package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

var VersionColl = versionCollection{"sysVersion"}

type versionCollection struct {
	name string
}

// Find 查找某个类型的版本号
func (v *versionCollection) FindOne(t string) (version *model.Version, err error) {
	version = new(model.Version)
	err = database.C(v.name).Find(bson.M{"vnType": t}).One(version)
	return
}

// Add 增加一条
func (v *versionCollection) Add(version *model.Version) error {
	return database.C(v.name).Insert(version)
}
