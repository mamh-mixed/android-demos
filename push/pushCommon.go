package push

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

const (
	ios         = "iOS"
	android     = "Android"
	buffer_size = 1024
)

var PushChan = make(chan *model.PushMessageReq, buffer_size) //缓存大小，channel没有内容时，阻塞，直到有内容写入

func PushMessage() {
	for true {
		req := <-PushChan
		if req == nil {
			log.Errorf("the element from PushChan is nil")
			continue
		}
		userModel, err := mongo.AppUserCol.FindOne(req.UserName)
		if err != nil { //查到该app用户的信息
			if (userModel.Device_type != "") && (userModel.Device_token != "") {
				req.Device_token = userModel.Device_token
				if userModel.Device_type == ios {
					ApnsPush.APush(req)
				} else if userModel.Device_type == android {
					UmengPush.UPush(req)
				} else {
					log.Errorf("the Device_type is diff, Device_type:%s", userModel.Device_type)
				}
			} else {
				log.Errorf("the app Device_type or Device_token is nil")
			}
		} else {
			log.Errorf("find app user error, err:%s", err)
		}

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
