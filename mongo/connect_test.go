package mongo

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/CardInfoLink/log"
)

type TestModel struct {
	Timestamp time.Time
	Value     string
}

func TestConnect(t *testing.T) {
	for {
		c := database.C("testConnect")

		m := &TestModel{time.Now(), "test..."}
		m2 := &TestModel{time.Now().Add(time.Millisecond), "test2..."}

		err := c.Insert(m, m2)
		if err != nil {
			log.Errorf("insert error %s", err)
		}

		cond := bson.M{"timestamp": m.Timestamp}
		change := bson.M{"$set": bson.M{"updateTime": time.Now()}}
		err = c.Update(cond, change)
		if err != nil {
			log.Errorf("update error %s", err)
		}

		result := TestModel{}
		err = c.Find(cond).One(&result)
		if err != nil {
			log.Errorf("find error %s", err)
		}

		fmt.Println("timestamp:", result.Timestamp)

		err = c.Remove(cond)

		time.Sleep(time.Second)
	}
}
