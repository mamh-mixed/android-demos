package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

type subAgent struct{}

var SubAgent subAgent

// Find 根据条件分页查找二级代理
func (s *subAgent) FindOne(subAgentCode string) (result *model.ResultBody) {
	log.Debugf("subAgentCode=%s", subAgentCode)

	subAgent, err := mongo.SubAgentColl.Find(subAgentCode)
	if err != nil {
		log.Errorf("查询二级代理(%s)出错:%s", subAgentCode, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    subAgent,
	}

	return result
}

// Find 根据条件分页查找二级代理。
func (s *subAgent) Find(subAgentCode, subAgentName, agentCode, agentName string, size, page int) (result *model.ResultBody) {
	log.Debugf("subAgentCode=%s; subAgentName=%s,agentCode=%s,agentName=%s", subAgentCode, subAgentName, agentCode, agentName)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	subAgents, total, err := mongo.SubAgentColl.PaginationFind(subAgentCode, subAgentName, agentCode, agentName, size, page)
	if err != nil {
		log.Errorf("查询所有二级代理出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(subAgents),
		Data:  subAgents,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// Save 保存二级代理信息，能同时用于新增或者修改的时候
func (i *subAgent) Save(data []byte) (result *model.ResultBody) {
	s := new(model.SubAgent)
	err := json.Unmarshal(data, s)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if s.SubAgentCode == "" {
		log.Error("没有SubAgentCode")
		return model.NewResultBody(3, "缺失必要元素SubAgentCode")
	}
	isExist := true
	_, err = mongo.SubAgentColl.Find(s.SubAgentCode)
	if err != nil {
		if err.Error() == "not found" {
			isExist = false
		} else {
			return model.NewResultBody(1, "查询数据库失败")
		}

	}
	if isExist {
		return model.NewResultBody(1, "公司代码已存在")
	}
	if s.SubAgentName == "" {
		log.Error("没有SubAgentName")
		return model.NewResultBody(3, "缺失必要元素SubAgentName")
	}

	if s.AgentCode == "" {
		log.Error("没有AgentCode")
		return model.NewResultBody(3, "缺失必要元素AgentCode")
	}

	if s.AgentName == "" {
		log.Error("没有AgentName")
		return model.NewResultBody(3, "缺失必要元素AgentName")
	}

	err = mongo.SubAgentColl.Insert(s)
	if err != nil {
		log.Errorf("新增二级代理失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    s,
	}

	return result
}

// Delete 删除二级代理
func (s *subAgent) Delete(subAgentCode string) (result *model.ResultBody) {

	err := mongo.SubAgentColl.Remove(subAgentCode)

	if err != nil {
		log.Errorf("删除二级代理失败: %s", err)
		return model.NewResultBody(1, "删除二级代理失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return result
}

// Update
func (i *subAgent) Update(data []byte) (result *model.ResultBody) {
	s := new(model.SubAgent)
	err := json.Unmarshal(data, s)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if s.SubAgentCode == "" {
		log.Error("没有SubAgentCode")
		return model.NewResultBody(3, "缺失必要元素SubAgentCode")
	}

	if s.SubAgentName == "" {
		log.Error("没有SubAgentName")
		return model.NewResultBody(3, "缺失必要元素SubAgentName")
	}

	if s.AgentCode == "" {
		log.Error("没有AgentCode")
		return model.NewResultBody(3, "缺失必要元素AgentCode")
	}

	if s.AgentName == "" {
		log.Error("没有AgentName")
		return model.NewResultBody(3, "缺失必要元素AgentName")
	}

	err = mongo.SubAgentColl.Update(s)
	if err != nil {
		log.Errorf("新增二级代理失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    s,
	}

	return result
}
