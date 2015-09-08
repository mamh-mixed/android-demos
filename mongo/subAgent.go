package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type subAgentCollection struct {
	name string
}

// SubAgentColl 代理商 Collection
var SubAgentColl = subAgentCollection{"subAgent"}

// Find 根据代理商代码查找
func (col *subAgentCollection) Find(subAgentCode string) (s *model.SubAgent, err error) {

	bo := bson.M{
		"subAgentCode": subAgentCode,
	}
	s = new(model.SubAgent)
	err = database.C(col.name).Find(bo).One(s)
	if err != nil {
		log.Errorf("Find SubAgent condition is: %+v;error is %s", bo, err)
	}
	return
}

// Add 增加一个代理商
func (col *subAgentCollection) Add(s *model.SubAgent) error {
	bo := bson.M{
		"subAgentCode": s.SubAgentCode,
	}
	_, err := database.C(col.name).Upsert(bo, s)

	return err
}

// Update 更新代理商信息
func (col *subAgentCollection) Update(s *model.SubAgent) error {
	bo := bson.M{
		"subAgentCode": s.SubAgentCode,
	}
	return database.C(col.name).Update(bo, s)
}

// Upsert 有则修改，没有则插入
func (col *subAgentCollection) Upsert(s *model.SubAgent) error {
	bo := bson.M{
		"subAgentCode": s.SubAgentCode,
	}
	_, err := database.C(col.name).Upsert(bo, s)
	return err
}

// FindByCode 得到某个代理商的名称
func (col *subAgentCollection) FindByCode(subAgentCode string) ([]*model.SubAgent, error) {
	var cs []*model.SubAgent
	err := database.C(col.name).Find(bson.M{"subAgentCode": subAgentCode}).All(&cs)
	return cs, err
}

// FindByCondition 根据代理商的条件查找代理商
func (col *subAgentCollection) FindByCondition(cond *model.SubAgent) (results []model.SubAgent, err error) {
	results = make([]model.SubAgent, 1)
	err = database.C(col.name).Find(cond).All(&results)
	if err != nil {
		log.Errorf("Find all subAgent error: %s", err)
		return nil, err
	}

	return results, err
}

// Remove 删除代理商
func (col *subAgentCollection) Remove(subAgentCode string) (err error) {
	bo := bson.M{}
	if subAgentCode != "" {
		bo["subAgentCode"] = subAgentCode
	}
	err = database.C(col.name).Remove(bo)
	return err
}

// PaginationFind 分页查找机构商户
func (c *subAgentCollection) PaginationFind(subAgentCode, subAgentName, agentCode, agentName string, size, page int) (results []model.SubAgent, total int, err error) {
	results = make([]model.SubAgent, 1)

	match := bson.M{}
	if subAgentCode != "" {
		match["subAgentCode"] = subAgentCode
	}
	if subAgentName != "" {
		match["subAgentName"] = subAgentName
	}
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

	sort := bson.M{"$sort": bson.M{"subAgentCode": 1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, sort, skip, limit)

	err = database.C(c.name).Pipe(cond).All(&results)

	return results, total, err
}
