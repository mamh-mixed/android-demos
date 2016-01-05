package push

import (
	"regexp"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
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
	formatTime := time.Now().Format("2006-01-02 15:04:05")
	reg, err := regexp.Compile("[^0-9]")
	if err != nil {
		log.Errorf("reg error in function SavePushMessage in file pushCommon")
		return err
	}
	formatTime = reg.ReplaceAllString(formatTime, "")
	rsp.PushTime = formatTime
	rsp.MsgId = util.SerialNumber()
	return mongo.PushMessageColl.Insert(rsp)
}
