package mongo

import (
	"fmt"
	"testing"
	"time"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	NameName   string
	PhonePhone string
}

func TestConnect(t *testing.T) {
	c := db.people
	err := c.Insert(&Person{"Ale", "+55 53 8116 9639"}, &Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		g.Fatal("insert error", err)
	}

	cond := bson.M{"name": "Cla"}
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
	err = c.Update(cond, change)
	if err != nil {
		g.Fatal("update error", err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		g.Fatal("find error", err)
	}

	fmt.Println("Phone:", result.PhonePhone)

	err = c.Remove(bson.M{"name": "Ale"})
}
