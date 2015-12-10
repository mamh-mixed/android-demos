package mongo

import (
	// "github.com/CardInfoLink/quickpay/model"

	"testing"
)

func TestKVListConditionFind(t *testing.T) {
	// var fuzzy = FuzzyCollection{
	// 	ColName:     "agent",
	// 	MatchFields: []string{"agentCode", "agentName"},
	// 	CodeField:   "agentCode",
	// 	NameField:   "agentName",
	// 	// Type:        "agent",
	// }
	var fuzzy = KVListCondition{
		ColName:   "subAgent",
		CodeField: "subAgentCode",
		NameField: "subAgentName",
	}

	filterMap := make(map[string]interface{})

	filterMap["agentCode"] = "19992900"
	fuzzy.FilterMap = filterMap

	keyWord := ""

	results, err := fuzzy.Find(keyWord)
	if err != nil {
		t.Logf("error is %s", err)
		t.Fail()
	}
	t.Logf("results is %+v; length is %d", results, len(results))
}
