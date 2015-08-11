package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type groupCollection struct {
	name string
}

// GroupColl 集团商户 Collection
var GroupColl = groupCollection{"group"}

// Find 根据集团代码查找
func (col *groupCollection) Find(groupCode string) (g *model.Group, err error) {

	bo := bson.M{
		"groupCode": groupCode,
	}
	g = new(model.Group)
	err = database.C(col.name).Find(bo).One(g)
	if err != nil {
		log.Errorf("Find Group condition is: %+v;error is %s", bo, err)
	}
	return
}

// Add 增加一个集团商户
func (col *groupCollection) Add(g *model.Group) error {
	bo := bson.M{
		"groupCode": g.GroupCode,
	}
	_, err := database.C(col.name).Upsert(bo, g)

	return err
}

// Modify 更新集团商户信息
func (col *groupCollection) Update(g *model.Group) error {
	bo := bson.M{
		"groupCode": g.GroupCode,
	}
	return database.C(col.name).Update(bo, g)
}

// FindByCode 得到某个集团的名称
func (col *groupCollection) FindByCode(groupCode string) ([]*model.Group, error) {
	var cs []*model.Group
	err := database.C(col.name).Find(bson.M{"groupCode": groupCode}).All(&cs)
	return cs, err
}

// FindByCondition 根据条件查找集团商户
func (col *groupCollection) FindByCondition(cond *model.Group) (results []model.Group, err error) {
	results = make([]model.Group, 1)
	err = database.C(col.name).Find(cond).All(&results)
	if err != nil {
		log.Errorf("Find all group error: %s", err)
		return nil, err
	}

	return
}

// Remove 删除代理商
func (col *groupCollection) Remove(groupCode, groupName string) (err error) {
	bo := bson.M{}
	if groupCode != "" {
		bo["groupCode"] = groupCode
	}
	if groupName != "" {
		bo["groupName"] = groupName
	}

	err = database.C(col.name).Remove(bo)

	return err
}
