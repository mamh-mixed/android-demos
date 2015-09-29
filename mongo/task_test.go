package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestAddTask(t *testing.T) {

	err := TaskCol.Add("test", false)
	if err != nil {
		t.Error(err)
	}
}

func TestPopTask(t *testing.T) {

	err := TaskCol.Pop("test")
	if err != nil {
		t.Error(err)
	}
}

func TestFindLogs(t *testing.T) {

	ls, err := SpChanLogsCol.Find(&model.QueryCondition{
		MerId:    "100000000010001",
		OrderNum: "1443517963485",
		ReqIds:   []string{"6eac314baaeb47b96b4e8798d9070673", "29aa667500d94e5752140fc18b8894d5"},
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("%+v", ls)
}
