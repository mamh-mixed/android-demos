package master

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/qiniu/log"
)

type session struct{}

var Session session

// 新建用户
func (s *session) Save(session *model.Session) (ret *model.ResultBody) {

	err := mongo.SessionColl.Add(session)
	if err != nil {
		log.Errorf("创建session失败,%s", err)
		return model.NewResultBody(1, "创建session失败")
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "创建session成功",
		Data:    session.SessionID,
	}
	return ret
}

// FindOne 根据sessionID查找session
func (s *session) FindOne(sessionID string) (ret *model.ResultBody) {
	log.Debugf("sessionID=%s", sessionID)

	session, err := mongo.SessionColl.Find(sessionID)
	if err != nil {
		log.Errorf("查询session(%s)出错:%s", sessionID, err)
		return model.NewResultBody(1, "查询失败")
	}
	user := &model.User{
		NickName: session.User.NickName,
		UserType: session.User.UserType,
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    user,
	}

	return ret
}

// Delete 根据sessionID删除session
func (s *session) Delete(sessionID string) (ret *model.ResultBody) {
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
