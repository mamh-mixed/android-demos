package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
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

// Update 修改
func (col *sessionCollection) Update(s *model.Session) error {
	return database.C(col.name).Update(bson.M{"sessionId": s.SessionID}, s)
}

// Find 根据sessionID查找
func (col *sessionCollection) Find(sessionID string) (s *model.Session, err error) {

	bo := bson.M{
		"sessionId": sessionID,
	}
	s = new(model.Session)
	err = database.C(col.name).Find(bo).One(s)
	return s, err
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
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}
