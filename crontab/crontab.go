package crontab

import (
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

var tasks = make(map[string]task)

type task struct {
	D               time.Duration
	Name            string
	F               func()
	IsSingleProcess bool
}

func RegisterTask(d time.Duration, taskName string, flag bool, taskFunc func()) {
	t := task{
		D: d, Name: taskName, F: taskFunc, IsSingleProcess: flag,
	}
	tasks[taskName] = t
}

func Start() {

	for _, t := range tasks {
		err := mongo.TaskCol.Add(t.Name, t.IsSingleProcess)
		if err != nil {
			log.Errorf("add task fail: name=%s,error=%s", t.Name, err)
			continue
		}
		go doTick(t)
	}

}

func doTick(t task) {

	ticker := time.NewTicker(t.D)

	for {
		<-ticker.C

		// 取任务
		err := mongo.TaskCol.Pop(t.Name)
		if err != nil {
			log.Debug("get no task .. continue")
			continue
		}

		// 执行任务
		t.F()

		// 完成任务
		mongo.TaskCol.Finish(t.Name)
	}
}
