package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type menuCollection struct {
	name string
}

// AgentColl 菜单 Collection
var MenuColl = menuCollection{"menu"}

// Find 根据Route查找
func (col *menuCollection) FindByRoute(route string) (m *model.Menu, err error) {
	bo := bson.M{
		"route": route,
	}
	m = new(model.Menu)
	err = database.C(col.name).Find(bo).One(m)
	if err != nil {
		log.Errorf("Find Menu condition is: %+v;error is %s", bo, err)
		return nil, err
	}
	return m, nil

}

// Find 根据level查找
func (col *menuCollection) FindByLevel(level int) (results []model.Menu, err error) {
	results = make([]model.Menu, 1)
	bo := bson.M{
		"level": level,
	}
	err = database.C(col.name).Find(bo).All(results)
	if err != nil {
		log.Errorf("Find Menu condition is: %+v;error is %s", bo, err)
		return nil, err
	}
	return results, nil
}

// Add 增加一个菜单
func (col *menuCollection) Add(m *model.Menu) error {
	err := database.C(col.name).Insert(m)
	if err != nil {
		log.Debugf("Add menu err,%s", err)
		return err
	}
	return nil
}

// Update 更新菜单信息
func (col *menuCollection) Update(m *model.Menu) error {
	bo := bson.M{
		"route": m.Route,
	}
	err := database.C(col.name).Update(bo, m)
	if err != nil {
		log.Debugf("update menu err,%s", err)
		return err
	}
	return nil
}

// PaginationFind 分页查找menu
func (col *menuCollection) PaginationFind(nameCN, route string, size, page int) (results []model.Menu, total int, err error) {
	results = make([]model.Menu, 1)

	match := bson.M{}
	if nameCN != "" {
		match["nameCN"] = nameCN
	}
	if route != "" {
		match["route"] = route
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
		log.Errorf("find menue paging err,%s", err)
		return nil, 0, err
	}

	return results, total, nil
}

// Remove 删除菜单
func (col *menuCollection) Remove(route string) (err error) {
	bo := bson.M{}
	if route != "" {
		bo["route"] = route
	}
	err = database.C(col.name).Remove(bo)
	return err
}
