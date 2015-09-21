package mongo

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/omigo/log"
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

func TestInsertURLValues(t *testing.T) {
	c := database.C("test.urlValues")

	v := &url.Values{}
	v.Set("key", "Values2")
	v.Set("timestamp", time.Now().String())

	err := c.Insert(v)
	if err != nil {
		log.Errorf("insert error %s", err)
	}

}
