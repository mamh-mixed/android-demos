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

func RegisterTask(d time.Duration, taskName string, IsSingleProcess bool, taskFunc func()) {
	t := task{
		D: d, Name: taskName, F: taskFunc,
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

	var index int
	for {
		<-ticker.C
		// 取任务
		if index == 0 {
			v, err := mongo.TaskCol.Pop(t.Name)
			if err != nil {
				log.Errorf("task pop error: %s", err)
				continue
			}
			// 拿不到任务
			if v == 0 {
				log.Warn("no task getted. quit!")
				ticker.Stop()
				return
			}
		}
		// 执行任务
		t.F()
		index++
	}
}
