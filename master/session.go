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
	user := &model.User{
		NickName: session.User.NickName,
		UserType: session.User.UserType,
	}
	ret = &model.ResultBody{
		Status:  0,
		Message: "创建session成功",
		Data:    user,
	}
	return ret
}
