package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type roleCollection struct {
	name string
}

// RoleColl 角色 Collection
var RoleColl = roleCollection{"role"}

// FindByRoleID 根据RoleID查找角色
func (col *roleCollection) FindByRoleID(roleID string) (r *model.Role, err error) {

	bo := bson.M{
		"roleID": roleID,
	}
	r = new(model.Role)
	err = database.C(col.name).Find(bo).One(r)
	if err != nil {
		log.Errorf("Find Role condition is: %+v;error is %s", bo, err)
		return nil, err
	}
	return r, nil
}

// Add 增加一个角色
func (col *roleCollection) Add(r *model.Role) error {
	err := database.C(col.name).Insert(r)
	if err != nil {
		log.Debugf("Add role err,%s", err)
		return err
	}
	return nil
}

// Update 更新角色信息
func (col *roleCollection) Update(r *model.Role) error {
	bo := bson.M{
		"roleID": r.RoleID,
	}
	err := database.C(col.name).Update(bo, r)
	if err != nil {
		log.Debugf("update role err,%s", err)
		return err
	}
	return nil
}

// PaginationFind 分页查找role
func (col *roleCollection) PaginationFind(roleID, name string, size, page int) (results []model.Role, total int, err error) {
	results = make([]model.Role, 1)
	match := bson.M{}
	if roleID != "" {
		match["roleID"] = roleID
	}
	if name != "" {
		match["name"] = name
	}
	// 计算总数
	total, err = database.C(col.name).Find(match).Count()
	if err != nil {
		log.Errorf("PaginationFind Count err,%s", err)
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	// sort := bson.M{"$sort": bson.M{"name": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, skip, limit)

	err = database.C(col.name).Pipe(cond).All(&results)
	if err != nil {
		log.Errorf("find role paging err,%s", err)
		return nil, 0, err
	}
	return results, total, nil
}

// Remove 删除角色
func (col *roleCollection) Remove(roleID string) (err error) {
	bo := bson.M{}
	if roleID != "" {
		bo["roleID"] = roleID
	}
	err = database.C(col.name).Remove(bo)
	return err
}
