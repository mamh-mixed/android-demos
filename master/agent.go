package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type agent struct{}

var Agent agent

// Find 根据条件分页查找代理
func (a *agent) FindOne(agentCode string) (result *model.ResultBody) {
	log.Debugf("agentCode=%s", agentCode)

	agent, err := mongo.AgentColl.Find(agentCode)
	if err != nil {
		log.Errorf("select agent (%s)error:%s", agentCode, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    agent,
	}

	return result
}

// Find 根据条件分页查找代理。
func (a *agent) Find(agentCode, agentName string, size, page int) (result *model.ResultBody) {
	log.Debugf("agentCode=%s; agentName=%s", agentCode, agentName)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	agents, total, err := mongo.AgentColl.PaginationFind(agentCode, agentName, size, page)
	if err != nil {
		log.Errorf("查询所有代理出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(agents),
		Data:  agents,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// Save 保存代理信息
func (i *agent) Save(data []byte) (result *model.ResultBody) {
	a := new(model.Agent)
	err := json.Unmarshal(data, a)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}
	if a.AgentCode == "" {
		log.Error("Required fields can not be empty")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}
	isExist := true
	// 查看agentCode是否存在
	_, err = mongo.AgentColl.Find(a.AgentCode)
	if err != nil {
		if err.Error() == "not found" {
			isExist = false
		} else {
			return model.NewResultBody(4, "SELECT_ERROR")
		}

	}
	if isExist {
		return model.NewResultBody(1, "AGENT_CODE_EXIST")
	}

	if a.AgentName == "" {
		log.Error("no AgentName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err = mongo.AgentColl.Insert(a)
	if err != nil {
		log.Errorf("create agent error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "CREATE_AGENT_SUCCESS",
		Data:    a,
	}

	return result
}

// Update 更新代理信息
func (i *agent) Update(data []byte) (result *model.ResultBody) {
	a := new(model.Agent)
	err := json.Unmarshal(data, a)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if a.AgentCode == "" {
		log.Error("no AgentCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if a.AgentName == "" {
		log.Error("no AgentName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err = mongo.AgentColl.Update(a)
	if err != nil {
		log.Errorf("update agent error:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "UPDATE_AGENT_SUCCESS",
		Data:    a,
	}

	return result
}

// Delete 删除代理
func (a *agent) Delete(agentCode string) (result *model.ResultBody) {

	err := mongo.AgentColl.Remove(agentCode)

	if err != nil {
		log.Errorf("delete agent error: %s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "DELETE_AGENT_SUCCESS",
	}

	return result
}
