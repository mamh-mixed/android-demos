package model

// Menu 菜单
type Menu struct {
	NameCN      string `json:"nameCN" bson:"nameCN"`
	NameEN      string `json:"nameEN" bson:"nameEN"`
	Icon        string `json:"icon" bson:"icon"`
	URL         string `json:"url" bson:"url"`
	Route       string `json:"route" bson:"route"`
	Children    []Menu `json:"children" bson:"children"`
	Level       int    `json:"level" bson:"level"` // 菜单的级别 1：一级菜单，2：二级菜单
	ParentRoute string `json:"parentRoute" bson:"parentRoute"`
}

// Role 角色
type Role struct {
	RoleID string `json:"roleID" bson:"roleID"`
	Name   string `json:"name" bson:"name"`
	Desc   string `json:"desc" bson:"desc"` // 角色描述
	Menus  []Menu `json:"menus" bson:"menus"`
}

// User 用户表
type User struct {
	UserName string `json:"userName" bson:"userName"`
	NickName string `json:"nickName" bson:"nickName"`
	Password string `json:"password" bson:"password"`
	Role     Role   `json:"role" bson:"role"`
}
