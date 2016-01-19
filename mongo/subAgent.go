package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type subAgentCollection struct {
	name string
}

// SubAgentColl 二级代理 Collection
var SubAgentColl = subAgentCollection{"subAgent"}

// Find 根据二级代理代码查找
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

// Add 增加一个二级代理
func (col *subAgentCollection) Add(s *model.SubAgent) error {
	bo := bson.M{
		"subAgentCode": s.SubAgentCode,
	}
	_, err := database.C(col.name).Upsert(bo, s)

	return err
}

// Update 更新二级代理信息
func (col *subAgentCollection) Update(s *model.SubAgent) error {
	s.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
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

// FindByCode 得到某个二级代理的名称
func (col *subAgentCollection) FindByCode(subAgentCode string) ([]*model.SubAgent, error) {
	var cs []*model.SubAgent
	err := database.C(col.name).Find(bson.M{"subAgentCode": subAgentCode}).All(&cs)
	return cs, err
}

// FindByCondition 根据条件查找二级代理
func (col *subAgentCollection) FindByCondition(cond *model.SubAgent) (results []model.SubAgent, err error) {
	results = make([]model.SubAgent, 1)
	err = database.C(col.name).Find(cond).All(&results)
	if err != nil {
		log.Errorf("Find all subAgent error: %s", err)
		return nil, err
	}

	return results, err
}

// Remove 删除二级代理
func (col *subAgentCollection) Remove(subAgentCode string) (err error) {
	bo := bson.M{}
	if subAgentCode != "" {
		bo["subAgentCode"] = subAgentCode
	}
	err = database.C(col.name).Remove(bo)
	return err
}

// PaginationFind 分页查找二级代理
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
func (col *subAgentCollection) Insert(s *model.SubAgent) error {
	s.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	s.UpdateTime = s.CreateTime
	return database.C(col.name).Insert(s)
}
