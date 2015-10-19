package crontab

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

var tasks = make(map[string]*model.Task)

func RegisterTask(d time.Duration, taskName string, taskFunc func()) {
	t := &model.Task{
		D: d, Name: taskName, F: taskFunc,
	}
	tasks[taskName] = t
}

func Start() {

	for _, t := range tasks {
		err := mongo.TaskCol.Add(t)
		if err != nil {
			log.Errorf("add task fail: name=%s,error=%s", t.Name, err)
			continue
		}
		go doTick(t)
	}

}

func doTick(t *model.Task) {
	ticker := time.NewTicker(t.D)
	var doTask = t.F

	for {

		// 取任务
		err := mongo.TaskCol.Pop(t)
		if err != nil {
			log.Debugf("get no %s task .. continue", t.Name)
			<-ticker.C
			continue
		}

		// 执行任务
		doTask() // 不可使用t.F()

		// 完成任务
		mongo.TaskCol.Finish(t)

		<-ticker.C
	}
}

// ShutDown 关闭正在做的任务，kill2时使用
func ShutDown() {
	for _, t := range tasks {
		if t.IsDoing {
			mongo.TaskCol.Finish(t)
		}
	}
}
