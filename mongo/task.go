package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	// "github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var TaskCol = taskCollection{"task"}

type taskCollection struct {
	name string
}

// Add 增加一个任务
func (t *taskCollection) Add(task *model.Task) error {

	count, err := database.C(t.name).Find(bson.M{"name": task.Name}).Count()
	if err != nil {
		return err
	}

	if count == 0 {
		task.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		task.UpdateTime = task.CreateTime
		err = database.C(t.name).Insert(task)
	}

	return err
}

// Pop 推出一个任务
func (t *taskCollection) Pop(task *model.Task) error {
	query := bson.M{
		"name":    task.Name,
		"isDoing": bson.M{"$ne": true},
	}

	change := mgo.Change{}
	change.Update = bson.M{
		"$set": bson.M{"isDoing": true,
			"updateTime": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	change.ReturnNew = true

	_, err := database.C(t.name).Find(query).Apply(change, task)

	return err
}

// Finish 完成一个任务
func (t *taskCollection) Finish(task *model.Task) error {
	query := bson.M{
		"name":    task.Name,
		"isDoing": true,
	}

	change := mgo.Change{}
	change.Update = bson.M{
		"$set": bson.M{"isDoing": false,
			"updateTime": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	change.ReturnNew = true
	_, err := database.C(t.name).Find(query).Apply(change, task)

	return err
}

func (t *taskCollection) FindAll() ([]model.Task, error) {
	var ts []model.Task
	err := database.C(t.name).Find(nil).All(&ts)
	return ts, err
}
