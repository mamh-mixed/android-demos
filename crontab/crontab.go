package crontab

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
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

	// 检查是否一直没被更新
	go checkDeathLock()

}

func checkDeathLock() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		<-ticker.C
		ts, err := mongo.TaskCol.FindAll()
		if err != nil {
			continue
		}

		for _, t := range ts {
			ut, _ := time.ParseInLocation("2006-01-02 15:04:05", t.UpdateTime, time.Local)
			// 超过时间并且一直在doing
			if mt, ok := tasks[t.Name]; ok {
				if time.Since(ut) > 2*mt.D && t.IsDoing {
					mongo.TaskCol.Finish(&t)
				}
			}
		}
	}
}

func doTick(t *model.Task) {
	ticker := time.NewTicker(t.D)
	var doTask = t.F

	time.Sleep(10 * time.Minute) // 休眠10分钟后开始
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
