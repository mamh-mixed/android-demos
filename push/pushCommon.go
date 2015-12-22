package push

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"strings"
	"time"
)

const (
	IOS         = "ios"
	ANDROID     = "android"
	buffer_size = 1024
)

func Do(req *model.PushMessageReq) {
	switch strings.ToLower(req.To) {
	case IOS:
		ApnsPush.APush(req)
	case ANDROID:
		UmengPush.UPush(req)
	default:
		log.Errorf("prepare to push,but unknown to=%s", req.To)
	}
}

func SavePushMessage(req *model.PushMessageReq) error {
	rsp := new(model.PushMessageRsp)
	rsp.UserName = req.UserName
	rsp.Title = req.Title
	rsp.Message = req.Message
	rsp.PushTime = time.Now().Format("2006-01-02 15:04:05")

	return mongo.PushMessageColl.Insert(rsp)
}

func PushInfos(req *model.PushMessageRsp) (rsp *PushInfoRsp) {
	rsp = new(PushInfoRsp)
	rsp.Error = "true"
	if req.UserName == "" || req.Password == "" {
		rsp.Error = model.PARAMS_EMPTY.Error
		return rsp
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			rsp.Error = model.USERNAME_PASSWORD_ERROR.Error
			return rsp
		}
		log.Errorf("find database err,%s", err)
		rsp.Error = model.SYSTEM_ERROR.Error
		return rsp
	}

	// 密码不对
	if req.Password != user.Password {
		rsp.Error = model.USERNAME_PASSWORD_ERROR.Error
		return rsp
	}

	infos, err := mongo.PushMessageColl.FindByUser(req)

	if err != nil {
		infos = make([]*model.PushInfo, 0)
	}

	rsp.Count = len(infos)
	rsp.Message = infos

	return rsp
}
