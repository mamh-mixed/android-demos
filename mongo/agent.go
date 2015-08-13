package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type agentCollection struct {
	name string
}

// AgentColl 代理商 Collection
var AgentColl = agentCollection{"agent"}

// Find 根据代理商代码查找
func (col *agentCollection) Find(agentCode string) (a *model.Agent, err error) {

	bo := bson.M{
		"agentCode": agentCode,
	}
	a = new(model.Agent)
	err = database.C(col.name).Find(bo).One(a)
	if err != nil {
		log.Errorf("Find Agent condition is: %+v;error is %s", bo, err)
	}
	return
}

// Add 增加一个代理商
func (col *agentCollection) Add(a *model.Agent) error {
	bo := bson.M{
		"agentCode": a.AgentCode,
	}
	_, err := database.C(col.name).Upsert(bo, a)

	return err
}

// Modify 更新代理商信息
func (col *agentCollection) Update(a *model.Agent) error {
	bo := bson.M{
		"agentCode": a.AgentCode,
	}
	return database.C(col.name).Update(bo, a)
}

// FindByCode 得到某个代理商的名称
func (col *agentCollection) FindByCode(agentCode string) ([]*model.Agent, error) {
	var cs []*model.Agent
	err := database.C(col.name).Find(bson.M{"agentCode": agentCode}).All(&cs)
	return cs, err
}

// FindByCondition 根据代理商的条件查找代理商
func (col *agentCollection) FindByCondition(cond *model.Agent) (results []model.Agent, err error) {
	results = make([]model.Agent, 1)
	err = database.C(col.name).Find(cond).All(&results)
	if err != nil {
		log.Errorf("Find all agent error: %s", err)
		return nil, err
	}

	return results, err
}

// Remove 删除代理商
func (col *agentCollection) Remove(agentCode string) (err error) {
	bo := bson.M{}
	if agentCode != "" {
		bo["agentCode"] = agentCode
	}
	err = database.C(col.name).Remove(bo)
	return err
}

// PaginationFind 分页查找机构商户
func (c *agentCollection) PaginationFind(agentCode, agentName string, size, page int) (results []model.Agent, total int, err error) {
	results = make([]model.Agent, 1)

	match := bson.M{}
	if agentCode != "" {
		match["agentCode"] = agentCode
	}
	if agentName != "" {
		match["agentName"] = agentName
	}

	// 计算总数
	total, err = database.C(c.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, skip, limit)

	err = database.C(c.name).Pipe(cond).All(&results)

	return results, total, err
}
