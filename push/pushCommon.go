package push

import (
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
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
	rsp.OrderNum = req.OrderNum
	rsp.Title = req.Title
	rsp.DeviceToken = req.DeviceToken
	rsp.Message = req.Message
	rsp.PushTime = time.Now().Format("20060102150405")
	rsp.MsgId = util.SerialNumber()
	rsp.MsgType = req.MsgType
	return mongo.PushMessageColl.Insert(rsp)
}
