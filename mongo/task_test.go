package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestAddPublicAccount(t *testing.T) {

	a := &model.PublicAccount{
		ChanMerId: "1276970101",
		AppID:     "wxaa785395d3d09403",
		AppSecret: "ff84e8b52e2ca0c439590070c7001c60",
	}
	err := PulicAccountCol.Add(a)
	if err != nil {
		t.Log(err)
	}
}

func TestAddTask(t *testing.T) {

	var task = &model.Task{
		Name:    "test",
		IsDoing: false,
	}

	err := TaskCol.Add(task)
	if err != nil {
		t.Error(err)
	}
}

func TestPopTask(t *testing.T) {

	var task = &model.Task{
		Name: "test",
	}

	err := TaskCol.Pop(task)
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
