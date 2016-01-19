package master

import (
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"net/http"
)

const (
	SessionKey = "QUICKMASTERID"
)

type sessionService struct{}

var Session sessionService
var expiredTime = time.Duration(goconf.Config.App.SessionExpiredTime)

func init() {
	go func() {
		timingClearSession()
	}()
}
func timingClearSession() {
	// 每段时间清理一次 Session
	refetchTime := expiredTime / 5
	for {
		select {
		case <-time.After(refetchTime):
			num, err := mongo.SessionColl.RemoveByTime()
			if err != nil {
				log.Errorf("clear session err,%s", err)
			}
			log.Infof("clear %d sessions", num)
		}
	}
}

// Get
func (s *sessionService) Get(req *http.Request) (*model.Session, error) {
	c, err := req.Cookie(SessionKey)
	if err != nil {
		return nil, err
	}
	return mongo.SessionColl.Find(c.Value)
}

// 新建会话
func (s *sessionService) Save(session *model.Session) (ret *model.ResultBody) {
	session.UserType = session.User.UserType
	session.NickName = session.User.NickName
	err := mongo.SessionColl.Add(session)
	if err != nil {
		log.Errorf("创建session失败,%s", err)
		return model.NewResultBody(1, "创建session失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "创建session成功",
		Data:    session,
	}
	return ret
}

// FindOne 根据sessionID查找session
func (s *sessionService) FindOne(sessionID string) (ret *model.ResultBody) {
	log.Debugf("sessionID=%s", sessionID)

	session, err := mongo.SessionColl.Find(sessionID)
	if err != nil {
		log.Errorf("find session(%s) err: %s", sessionID, err)
		return model.NewResultBody(1, "查询失败")
	}
	user := session.User
	// user.UserName = ""
	user.PhoneNum = ""
	user.Mail = ""
	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    user,
	}

	return ret
}

// Delete 根据sessionID删除session
func (s *sessionService) Delete(sessionID string) (ret *model.ResultBody) {
	log.Debugf("sessionID=%s", sessionID)

	err := mongo.SessionColl.Remove(sessionID)
	if err != nil {
		log.Errorf("删除session(%s)出错:%s", sessionID, err)
		return model.NewResultBody(1, "删除失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "删除成功",
	}

	return ret
}
