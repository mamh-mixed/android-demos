package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var TaskCol = taskCollection{"task"}

type taskCollection struct {
	name string
}

type task struct {
	Name            string `bson:"name"`
	IsSingleProcess bool   `bson:"isSignleProcess"`
	// IsProcessSuccess bool   `bson:"isProcessSuccess"`
	IsDoing    bool   `bson:"isDoing"`
	UpdateTime string `bson:"updateTime"`
}

// Add 增加一个任务
func (t *taskCollection) Add(taskName string, isSingleProcess bool) error {

	if isSingleProcess {
		_, err := database.C(t.name).Upsert(bson.M{"name": taskName}, bson.M{
			"$set": bson.M{"name": taskName},
		})
		return err
	}

	return database.C(t.name).Insert(&task{
		Name: taskName,
	})
}

// Pop 推出一个任务
func (t *taskCollection) Pop(taskName string) error {
	query := bson.M{
		"name":    taskName,
		"isDoing": bson.M{"$ne": true},
	}

	change := mgo.Change{}
	change.Update = bson.M{
		"$set": bson.M{"isDoing": true,
			"updateTime": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	_, err := database.C(t.name).Find(query).Apply(change, &task{})

	return err
}

// Finish 完成一个任务
func (t *taskCollection) Finish(taskName string) error {
	query := bson.M{
		"name":    taskName,
		"isDoing": true,
	}

	change := mgo.Change{}
	change.Update = bson.M{
		"$set": bson.M{"isDoing": false},
	}

	_, err := database.C(t.name).Find(query).Apply(change, &task{})

	return err
}
