package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type agent struct{}

var Agent agent

//// Find 根据条件查找代理商。
func (a *agent) Find(agentCode, agentName string) (result *model.ResultBody) {
	log.Debugf("agentCode is %s; agentName is %s", agentCode, agentName)

	cond := new(model.Agent)

	if agentCode != "" {
		cond.AgentCode = agentCode
	}

	if agentName != "" {
		cond.AgentName = agentName
	}

	agents, err := mongo.AgentColl.FindByCondition(cond)

	if err != nil {
		log.Errorf("查询所有代理商出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    agents,
	}

	return result
}

// Save 保存代理商信息，能同时用于新增或者修改的时候
func (i *agent) Save(data []byte) (result *model.ResultBody) {
	a := new(model.Agent)
	err := json.Unmarshal(data, a)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if a.AgentCode == "" {
		log.Error("没有AgentCode")
		return model.NewResultBody(3, "缺失必要元素AgentCode")
	}

	if a.AgentName == "" {
		log.Error("没有AgentName")
		return model.NewResultBody(3, "缺失必要元素AgentName")
	}

	err = mongo.AgentColl.Add(a)
	if err != nil {
		log.Errorf("新增代理商失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    a,
	}

	return result
}

// Delete 删除代理商，参数是 agentCode, agentName
func (a *agent) Delete(agentCode, agentName string) (result *model.ResultBody) {

	err := mongo.AgentColl.Remove(agentCode, agentName)

	if err != nil {
		log.Errorf("删除代理商失败: %s", err)
		return model.NewResultBody(1, "删除代理商失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return result
}
