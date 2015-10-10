package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type sessionCollection struct {
	name string
}

// SessionColl  session Collection
var SessionColl = sessionCollection{"session"}

// Add 增加一个session
func (col *sessionCollection) Add(s *model.Session) error {
	bo := bson.M{
		"sessionId": s.SessionID,
	}
	_, err := database.C(col.name).Upsert(bo, s)

	return err
}

// Find 根据sessionID查找
func (col *sessionCollection) Find(sessionID string) (s *model.Session, err error) {

	bo := bson.M{
		"sessionId": sessionID,
	}
	s = new(model.Session)
	err = database.C(col.name).Find(bo).One(s)
	if err != nil {
		log.Errorf("Find Session condition is: %+v;error is %s", bo, err)
	}
	return
}

// Remove 删除session
func (col *sessionCollection) Remove(sessionID string) (err error) {
	bo := bson.M{}
	if sessionID != "" {
		bo["sessionId"] = sessionID
	}
	err = database.C(col.name).Remove(bo)
	return err
}

// RemoveByTime 删除已经过期的session
func (col *sessionCollection) RemoveByTime() (num int, err error) {
	bo := bson.M{"expires": bson.M{"$lte": time.Now().Format("2006-01-02 15:04:05")}}
	info, err := database.C(col.name).RemoveAll(bo)
	return info.Removed, err
}
