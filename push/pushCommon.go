package push

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"strings"
	"time"
)

func Do(req *model.PushMessageReq) {

	if err := SavePushMessage(req); err != nil {
		log.Errorf("save push message error: %s", err)
	}

	switch strings.ToLower(req.To) {
	case model.IOS:
		ApnsPush.APush(req)
	case model.Android:
		UmengPush.UPush(req)
	default:
		log.Errorf("prepare to push,but unknown to=%s", req.To)
	}
}

func SavePushMessage(req *model.PushMessageReq) error {
	rsp := new(model.PushMessage)
	rsp.UserName = req.UserName
	rsp.Title = req.Title
	rsp.DeviceToken = req.DeviceToken
	rsp.Message = req.Message
	rsp.PushTime = time.Now().Format("2006-01-02 15:04:05")
	rsp.MsgId = util.SerialNumber()
	return mongo.PushMessageColl.Insert(rsp)
}
