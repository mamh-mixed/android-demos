package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var RoleSettCol = roleSettCollection{name: "roleSett"}

// roleSettCollection 根绝清算角色清算
type roleSettCollection struct {
	name string
}

// FindOne 查找一条记录
func (r *roleSettCollection) FindOne(role, date string) (*model.RoleSett, error) {
	result := &model.RoleSett{}
	err := database.C(r.name).Find(bson.M{"settRole": role, "settDate": date}).One(result)
	return result, err
}

// Upsert
func (r *roleSettCollection) Upsert(rs *model.RoleSett) error {
	rs.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	if rs.CreateTime == "" {
		rs.CreateTime = rs.UpdateTime
	}
	_, err := database.C(r.name).Upsert(bson.M{"settRole": rs.SettRole, "settDate": rs.SettDate}, rs)
	return err
}
