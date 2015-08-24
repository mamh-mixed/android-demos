package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type role struct{}

// Role角色
var Role role = role{}

// CreateRole 创建角色
func (r *role) CreateRole(data []byte) (ret *model.ResultBody) {
	log.Debugf("CreateRole,%s", string(data))
	var role *model.Role
	err := json.Unmarshal(data, role)
	if err != nil {
		log.Errorf("json unmarshal err,%s", err)
		return model.NewResultBody(1, "json解析失败")
	}
	if role.RoleID == "" || role.Name == "" {
		log.Errorf("roleID和name为空")
		return model.NewResultBody(2, "roleID和name不能为空")
	}
	roleOld, err := mongo.RoleColl.FindByRoleID(role.RoleID)
	if err != nil {
		log.Errorf("FindByRoleID失败，%s", err)
		return model.NewResultBody(3, "操作数据库失败")
	}
	if roleOld != nil {
		log.Errorf("roleID已经存在")
		return model.NewResultBody(4, "roleID已经存在")
	}
	err = mongo.RoleColl.Add(role)
	if err != nil {
		log.Errorf("Add Role err,%s", err)
		return model.NewResultBody(5, "操作数据库失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "添加角色成功",
	}
	log.Infof("添加角色成功,routeId:%s", role.RoleID)
	return ret
}

// CreateRole 创建角色
func (r *role) UpdateRole(data []byte) (ret *model.ResultBody) {
	log.Debugf("UpdateRole,%s", string(data))
	var role *model.Role
	err := json.Unmarshal(data, role)
	if err != nil {
		log.Errorf("json unmarshal err,%s", err)
		return model.NewResultBody(1, "json解析失败")
	}
	if role.RoleID == "" || role.Name == "" {
		log.Errorf("roleID和name为空")
		return model.NewResultBody(2, "roleID和name不能为空")
	}
	err = mongo.RoleColl.Update(role)
	if err != nil {
		log.Errorf("Add Role err,%s", err)
		return model.NewResultBody(3, "操作数据库失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "更新角色成功",
	}
	log.Infof("更新角色成功,routeId:%s", role.RoleID)
	return ret
}
func (r *role) UpdateMenu(data []byte) (ret *model.ResultBody) {
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

// 删除角色
func (r *role) RemoveRole(roleID string) (ret *model.ResultBody) {
	log.Debugf("RemoveRole,roleID=%s", roleID)
	if roleID == "" {
		log.Errorf("roleID为空")
		return model.NewResultBody(1, "roleID不能为空")
	}
	err := mongo.RoleColl.Remove(roleID)
	if err != nil {
		log.Errorf("Remove err,%s", err)
		return model.NewResultBody(2, "操作数据库失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "删除角色成功",
	}
	return ret
}

// 分页查找菜单
func (r *role) Find(roleID, name string, size, page int) (ret *model.ResultBody) {
	log.Debugf("roleID is %s; name is %s", roleID, name)

	if page <= 0 {
		return model.NewResultBody(1, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	roles, total, err := mongo.RoleColl.PaginationFind(roleID, name, size, page)
	if err != nil {
		log.Errorf("查询所有角色出错:%s", err)
		return model.NewResultBody(2, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(roles),
		Data:  roles,
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return ret
}
