package master

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type group struct{}

var Group group

//// Find 根据条件查找集团商户。
func (g *group) Find(groupCode, groupName string) (result *model.ResultBody) {
	log.Debugf("groupCode is %s; groupName is %s", groupCode, groupName)

	cond := new(model.Group)

	if groupCode != "" {
		cond.GroupCode = groupCode
	}

	if groupName != "" {
		cond.GroupName = groupName
	}

	groups, err := mongo.GroupColl.FindByCondition(cond)

	if err != nil {
		log.Errorf("查询所有集团商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    groups,
	}

	return result
}

// Save 保存集团商户信息，能同时用于新增或者修改的时候
func (i *group) Save(data []byte) (result *model.ResultBody) {
	g := new(model.Group)
	err := json.Unmarshal(data, g)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if g.GroupCode == "" {
		log.Error("没有GroupCode")
		return model.NewResultBody(3, "缺失必要元素GroupCode")
	}

	if g.GroupName == "" {
		log.Error("没有GroupName")
		return model.NewResultBody(3, "缺失必要元素GroupName")
	}

	err = mongo.GroupColl.Add(g)
	if err != nil {
		log.Errorf("新增集团商户失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    g,
	}

	return result
}

// Delete 删除集团商户，参数是 groupCode, groupName
func (g *group) Delete(groupCode, groupName string) (result *model.ResultBody) {

	err := mongo.GroupColl.Remove(groupCode, groupName)

	if err != nil {
		log.Errorf("删除集团商户失败: %s", err)
		return model.NewResultBody(1, "删除集团商户失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return result
}
