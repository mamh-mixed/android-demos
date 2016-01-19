package mongo

import (
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

// 模糊搜索配置项，keyValueCondition
type KVListCondition struct {
	ColName   string                 `json:"colName,omitempty"`   // 集合名
	CodeField string                 `json:"codeField,omitempty"` // 键域，模糊搜索目标
	NameField string                 `json:"nameField,omitempty"` // 值域，模糊搜索目标
	FilterMap map[string]interface{} `json:"filterMap,omitempty"` // 过滤条件
	Limit     int                    `json:"limit,omitempty"`     // 单次查询最多数量
}

// 模糊搜索结果
type KVItem struct {
	Code string `json:"code,omitempty"` // key
	Name string `json:"name,omitempty"` // name
}

const (
	KV_LIST_DEFAULT_LIMIT = 10
)

// Find 模糊查找
func (c *KVListCondition) Find(keyWord string) (results []KVItem, err error) {
	if c.Limit == 0 {
		c.Limit = KV_LIST_DEFAULT_LIMIT
	}

	cond, match := bson.M{}, []bson.M{}

	if c.CodeField != "" {
		match = append(match, bson.M{c.CodeField: bson.RegEx{keyWord, "."}})
	}
	if c.NameField != "" {
		match = append(match, bson.M{c.NameField: bson.RegEx{keyWord, "."}})
	}
	cond["$or"] = match

	for key, value := range c.FilterMap {
		cond[key] = value
	}

	log.Debugf("cond is %+v", cond)

	results = make([]KVItem, 0)
	err = database.C(c.ColName).Pipe([]bson.M{
		{"$match": cond},
		{"$limit": c.Limit},
		{
			"$project": bson.M{
				"code": "$" + c.CodeField,
				"name": "$" + c.NameField,
			},
		},
	}).All(&results)

	return results, err
}
