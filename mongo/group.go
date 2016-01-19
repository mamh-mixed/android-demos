package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
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

// Update 更新集团商户信息
func (col *groupCollection) Update(g *model.Group) error {
	g.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	bo := bson.M{
		"groupCode": g.GroupCode,
	}
	return database.C(col.name).Update(bo, g)
}

// Upsert 有则更新无则插入
func (col *groupCollection) Upsert(g *model.Group) error {
	bo := bson.M{
		"groupCode": g.GroupCode,
	}
	_, err := database.C(col.name).Upsert(bo, g)
	return err
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

	return results, err
}

// Remove 删除代理商
func (col *groupCollection) Remove(groupCode string) (err error) {
	bo := bson.M{}
	if groupCode != "" {
		bo["groupCode"] = groupCode
	}
	err = database.C(col.name).Remove(bo)
	return err
}

// PaginationFind 分页查找机构商户
func (c *groupCollection) PaginationFind(groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName string, size, page int) (results []model.Group, total int, err error) {
	results = make([]model.Group, 1)

	match := bson.M{}
	if groupCode != "" {
		match["groupCode"] = groupCode
	}
	if groupName != "" {
		match["groupName"] = groupName
	}
	if agentCode != "" {
		match["agentCode"] = agentCode
	}
	if agentName != "" {
		match["agentName"] = agentName
	}

	if subAgentCode != "" {
		match["subAgentCode"] = subAgentCode
	}
	if subAgentName != "" {
		match["subAgentName"] = subAgentName
	}

	// 计算总数
	total, err = database.C(c.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	sort := bson.M{"$sort": bson.M{"groupCode": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, sort, skip, limit)

	err = database.C(c.name).Pipe(cond).All(&results)

	return results, total, err
}

// Insert 更新集团商户信息
func (col *groupCollection) Insert(g *model.Group) error {
	g.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	g.UpdateTime = g.CreateTime
	return database.C(col.name).Insert(g)
}
