package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
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
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if s.SubAgentCode == "" {
		log.Error("no SubAgentCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}
	isExist := true
	_, err = mongo.SubAgentColl.Find(s.SubAgentCode)
	if err != nil {
		if err.Error() == "not found" {
			isExist = false
		} else {
			return model.NewResultBody(1, "SELECT_ERROR")
		}

	}
	if isExist {
		return model.NewResultBody(1, "SUB_AGENT_CODE_EXIST")
	}
	if s.SubAgentName == "" {
		log.Error("no SubAgentName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if s.AgentCode == "" {
		log.Error("no AgentCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if s.AgentName == "" {
		log.Error("no AgentName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err = mongo.SubAgentColl.Insert(s)
	if err != nil {
		log.Errorf("create subAgent error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "CREATE_SUB_AGENT_SUCCESS",
		Data:    s,
	}

	return result
}

// Delete 删除二级代理
func (s *subAgent) Delete(subAgentCode string) (result *model.ResultBody) {

	err := mongo.SubAgentColl.Remove(subAgentCode)

	if err != nil {
		log.Errorf("delete subAgent error: %s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "DELETE_SUB_AGENT_SUCCESS",
	}

	return result
}

// Update
func (i *subAgent) Update(data []byte) (result *model.ResultBody) {
	s := new(model.SubAgent)
	err := json.Unmarshal(data, s)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if s.SubAgentCode == "" {
		log.Error("no SubAgentCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if s.SubAgentName == "" {
		log.Error("no SubAgentName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if s.AgentCode == "" {
		log.Error("no AgentCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if s.AgentName == "" {
		log.Error("没有AgentName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err = mongo.SubAgentColl.Update(s)
	if err != nil {
		log.Errorf("update subAgent error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "UPDATE_SUB_AGENT_SUCCESS",
		Data:    s,
	}

	return result
}
