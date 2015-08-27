package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type menu struct{}

var Menu menu

// CreateMenu 创建菜单
func (m *menu) CreateMenu(data []byte) (ret *model.ResultBody) {
	var menu *model.Menu
	err := json.Unmarshal(data, menu)
	if err != nil {
		log.Errorf("json unmarshal err,%s", err)
		return model.NewResultBody(1, "json 解析失败")
	}
	if menu.Route == "" || menu.Level == 0 {
		log.Errorf("route和level为空")
		return model.NewResultBody(2, "route和level不能为空")
	}
	if menu.Level == 2 {
		if menu.ParentRoute == "" {
			log.Errorf("parentRoute为空")
			return model.NewResultBody(3, "parentRoute不能为空")
		}
	}
	menuOld, err := mongo.MenuColl.FindByRoute(menu.Route)
	if err != nil {
		log.Errorf("FindByRoute err,%s", err)
		return model.NewResultBody(4, "操作数据库失败")
	}
	if menuOld != nil {
		log.Errorf("route已经存在")
		return model.NewResultBody(5, "route已经存在")
	}
	err = mongo.MenuColl.Add(menu)
	if err != nil {
		log.Errorf("add menu err,%s", err)
		return model.NewResultBody(6, "操作数据库失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "创建菜单成功",
	}
	return ret
}
func (m *menu) UpdateMenu(data []byte) (ret *model.ResultBody) {
	var menu *model.Menu
	err := json.Unmarshal(data, menu)
	if err != nil {
		log.Errorf("json unmarshal err,%s", err)
		return model.NewResultBody(1, "json 解析失败")
	}
	if menu.Route == "" || menu.Level == 0 {
		log.Errorf("route和level为空")
		return model.NewResultBody(2, "route和level不能为空")
	}
	if menu.Level == 2 {
		if menu.ParentRoute == "" {
			log.Errorf("parentRoute为空")
			return model.NewResultBody(3, "parentRoute不能为空")
		}
	}
	err = mongo.MenuColl.Update(menu)
	if err != nil {
		log.Errorf("add menu err,%s", err)
		return model.NewResultBody(4, "操作数据库失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "创建菜单成功",
	}
	return ret
}

// 删除菜单
func (m *menu) RemoveMenu(route string) (ret *model.ResultBody) {
	log.Debugf("RemoveMenu,route=%s", route)
	if route == "" {
		log.Errorf("route为空")
		return model.NewResultBody(1, "route不能为空")
	}
	err := mongo.MenuColl.Remove(route)
	if err != nil {
		log.Errorf("Remove err,%s", err)
		return model.NewResultBody(2, "操作数据库失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "删除菜单成功",
	}
	return ret
}

// 分页查找菜单
func (m *menu) Find(nameCN, route string, size, page int) (ret *model.ResultBody) {
	log.Debugf("nameCN is %s; route is %s", nameCN, route)

	if page <= 0 {
		return model.NewResultBody(1, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	menus, total, err := mongo.MenuColl.PaginationFind(nameCN, route, size, page)
	if err != nil {
		log.Errorf("查询所有菜单出错:%s", err)
		return model.NewResultBody(2, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(menus),
		Data:  menus,
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return ret
}
