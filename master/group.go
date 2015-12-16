package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type group struct{}

var Group group

// Find 根据条件分页查找集团商户
func (g *group) FindOne(groupCode string) (result *model.ResultBody) {
	log.Debugf("groupCode=%s", groupCode)

	group, err := mongo.GroupColl.Find(groupCode)

	if err != nil {
		log.Errorf("查询集团(%s)出错: %s", groupCode, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    group,
	}

	return result
}

// Find 根据条件分页查找商户。
func (g *group) Find(groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName string, size, page int) (result *model.ResultBody) {
	log.Debugf("groupCode=%s; groupName=%s;agentCode=%s;agentName=%s;subAgentCode=%s;subAgentName=%s,", groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	groups, total, err := mongo.GroupColl.PaginationFind(groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName, size, page)
	if err != nil {
		log.Errorf("查询所有集团出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(groups),
		Data:  groups,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// Save 保存集团商户信息
func (i *group) Save(data []byte) (result *model.ResultBody) {
	g := new(model.Group)
	err := json.Unmarshal(data, g)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if g.GroupCode == "" {
		log.Error("no GroupCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	isExist := true

	_, err = mongo.GroupColl.Find(g.GroupCode)
	if err != nil {
		if err.Error() == "not found" {
			isExist = false
		} else {
			return model.NewResultBody(4, "SELECT_ERROR")
		}
	}
	if isExist {
		return model.NewResultBody(1, "GROUP_CODE_EXIST")
	}

	if g.GroupName == "" {
		log.Error("no GroupName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err = mongo.GroupColl.Insert(g)
	if err != nil {
		log.Errorf("create group error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "CREATE_GROUP_SUCCESS",
		Data:    g,
	}

	return result
}

// Delete 删除集团商户
func (g *group) Delete(groupCode string) (result *model.ResultBody) {

	err := mongo.GroupColl.Remove(groupCode)

	if err != nil {
		log.Errorf("delete group error: %s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "DELETE_GROUP_SUCCESS",
	}

	return result
}

// Update 更新集团商户信息
func (i *group) Update(data []byte) (result *model.ResultBody) {
	g := new(model.Group)
	err := json.Unmarshal(data, g)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_ERROR")
	}

	if g.GroupCode == "" {
		log.Error("NO GroupCode")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	if g.GroupName == "" {
		log.Error("NO GroupName")
		return model.NewResultBody(3, "REQUIRED_FILED_NOT_BE_EMPTY")
	}

	err = mongo.GroupColl.Update(g)
	if err != nil {
		log.Errorf("update group error:%s", err)
		return model.NewResultBody(1, "ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "UPDATE_GROUP_SUCCESS",
		Data:    g,
	}

	return result
}
