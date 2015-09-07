package mongo

import (
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
