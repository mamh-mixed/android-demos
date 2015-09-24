package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var TaskCol = taskCollection{"task"}

type taskCollection struct {
	name string
}

// Add 往队列里增加一个任务
func (t *taskCollection) Add(taskName string, isUnique bool) error {

	selector := bson.M{
		"name": taskName,
	}
	update := bson.M{}
	if isUnique {
		update["$addToSet"] = bson.M{
			"queue": 1,
		}
	} else {
		update["$push"] = bson.M{
			"queue": 1,
		}
	}
	return database.C(t.name).Update(selector, update)
}

// Pop 推出一个任务
func (t *taskCollection) Pop(taskName string) (int, error) {
	query := bson.M{
		"name": taskName,
	}

	change := mgo.Change{}
	change.Update = bson.M{
		"$pop": bson.M{"queue": 1},
	}

	var result = &struct {
		Name  string `bson:"name"`
		Queue []int  `bson:"queue"`
	}{}

	_, err := database.C(t.name).Find(query).Apply(change, result)

	return len(result.Queue), err
}
