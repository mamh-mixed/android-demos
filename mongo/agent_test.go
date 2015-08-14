package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

var (
	agentCode = "TESTCODE"
	agentName = "TEST001"
)

func TestAgentPaginationFind(t *testing.T) {
	agentCode, agentName := "", ""
	size, page := 10, 1
	results, total, err := AgentColl.PaginationFind(agentCode, agentName, size, page)
	if err != nil {
		log.Errorf("fail: %s", err)
	}

	t.Logf("total is %d; collections are %#v", total, results)

	t.Logf("current count is %d", len(results))
}

func TestAgentFind(t *testing.T) {
	agent, err := AgentColl.Find(agentCode)
	if err != nil {
		t.Error("find agent unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("find agent success %s", agent)
}

func TestAgentAdd(t *testing.T) {
	agent := &model.Agent{
		AgentCode: "TESTCODE",
		AgentName: "TESTNAME",
	}
	err := AgentColl.Add(agent)
	if err != nil {
		t.Error("add agent unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("add agent success %s", agent)
}

func TestAgentUpdate(t *testing.T) {
	agent := &model.Agent{
		AgentCode: agentCode,
		AgentName: "TESTNAME01",
	}
	err := AgentColl.Update(agent)
	if err != nil {
		t.Error("update agent unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("update agent success %s", agent)
}

func TestAgentFindByCode(t *testing.T) {
	cs, err := AgentColl.FindByCode(agentCode)
	if err != nil {
		t.Errorf("findAll agent unsuccessful %s", err)
		t.FailNow()
	}
	log.Debugf("%+v", cs)
}

func TestAgentFindByCondition(t *testing.T) {
	cms, err := AgentColl.FindByCondition(nil)
	if err != nil {
		t.Error("出错啦")
	}
	t.Logf("result is %+v", cms)
}

func TestAgentRemove(t *testing.T) {
	agentCode := "TESTCODE"
	err := AgentColl.Remove(agentCode)
	if err != nil {
		t.Error("remove agent unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("remove agent success %s", agentCode)
}
