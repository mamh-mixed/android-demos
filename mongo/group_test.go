package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
)

var (
	groupCode = "TEST001"
)

func TestGroupPaginationFind(t *testing.T) {
	groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName := "", "", "", "", "", ""
	size, page := 10, 1
	results, total, err := GroupColl.PaginationFind(groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName, size, page)
	if err != nil {
		log.Errorf("fail: %s", err)
	}

	t.Logf("total is %d; collections are %#v", total, results)

	t.Logf("current count is %d", len(results))
}

func TestGroupFind(t *testing.T) {
	group, err := GroupColl.Find(groupCode)
	if err != nil {
		t.Error("find group unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("find group success %s", group)
}

func TestGroupAdd(t *testing.T) {
	group := &model.Group{
		GroupCode: "TESTCODE",
	}
	err := GroupColl.Add(group)
	if err != nil {
		t.Error("add group unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("add group success %s", group)
}

func TestGroupUpdate(t *testing.T) {
	group := &model.Group{
		GroupCode: groupCode,
	}
	err := GroupColl.Update(group)
	if err != nil {
		t.Error("update group unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("update group success %s", group)
}

func TestGroupFindByCode(t *testing.T) {
	cs, err := GroupColl.FindByCode(groupCode)
	if err != nil {
		t.Errorf("findAll group unsuccessful %s", err)
		t.FailNow()
	}
	log.Debugf("%+v", cs)
}

func TestGroupFindByCondition(t *testing.T) {
	cms, err := GroupColl.FindByCondition(nil)
	if err != nil {
		t.Error("出错啦")
	}
	t.Logf("result is %+v", cms)
}

func TestGroupRemove(t *testing.T) {
	groupCode := "TESTCODE"
	err := GroupColl.Remove(groupCode)
	if err != nil {
		t.Error("remove group unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("remove group success %s", groupCode)
}
