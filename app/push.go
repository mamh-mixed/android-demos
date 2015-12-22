package app

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/push"
	"github.com/omigo/log"
	"strings"
)

// 推送消息类型
const (
	MsgType_A = "A"
	MsgType_B = "B"
	MsgType_C = "C"
)

// 支持1W的并发量
var MsgChan = make(chan *model.PushMessageReq, 1e4) //缓存大小，channel没有内容时，阻塞，直到有内容写入

func StartPush() {
	go pushMsg()
}

func pushMsg() {
	for {
		req := <-MsgChan
		if req == nil {
			// log.Errorf("the element from PushChan is nil")
			continue
		}

		// 单独的goruntine处理
		// 防止发推送时卡住
		go func() {
			var appUser *model.AppUser
			var err error
			// 判断消息类型
			switch strings.ToLower(req.From) {
			case "wap":
				// 此类为A类消息，向所有用户推
				appUsers, err := mongo.AppUserCol.FindByMerId(req.MerID)
				if err != nil {
					log.Errorf("find appUser error: %s ", err)
					return
				}

				// 逐个推
				for _, user := range appUsers {
					req.Device_token = user.Device_token
					req.To = user.Device_type
					if req.Device_token != "" {
						push.Do(req)
					}
				}

			case "ios", "android":
				// 此类为B类消息
				appUser, err = mongo.AppUserCol.FindOne(req.UserName)
				if err != nil {
					log.Errorf("prepare to push, but appUser(%s) not found, err:%s", req.UserName, err)
					return
				}
				// 组装报文推送
				req.Device_token = appUser.Device_token
				req.To = appUser.Device_type
				if req.Device_token != "" {
					push.Do(req)
				}
			case "notify":
				// 此类为C类消息
				// TODO
			default:
				log.Errorf("prepare to push, but unknown from=%s", req.From)
			}
		}()
	}
}
