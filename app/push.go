package app

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/push"
	"github.com/CardInfoLink/log"
	"strings"
)

// 推送消息类型
const (
	MsgType_A = "MSG_TYPE_A"
	MsgType_B = "MSG_TYPE_B"
	MsgType_C = "MSG_TYPE_C"
)

func StartPush() {
	go pushMsg()
}

func pushMsg() {
	for {
		t := <-core.MsgQueue
		if t == nil {
			// log.Errorf("the element from PushChan is nil")
			continue
		}

		req := pushMsgReq(t)

		// 单独的goruntine处理
		// 防止发推送时卡住
		go func() {
			var appUser *model.AppUser
			var err error
			// 判断消息类型
			switch req.MsgType {
			case MsgType_A:
				// 此类为A类消息，向所有用户推
				appUsers, err := mongo.AppUserCol.FindByMerId(req.MerID)
				if err != nil {
					log.Errorf("find appUser error: %s ", err)
					return
				}

				// 逐个推
				for _, u := range appUsers {
					req.DeviceToken = u.DeviceToken
					req.To = u.DeviceType
					req.UserName = u.UserName
					if req.DeviceToken != "" {
						log.Debugf("push to user=%s,token=%s", u.UserName, u.DeviceToken)
						push.Do(req)
					}
				}

			case MsgType_B:
				// 此类为B类消息
				appUser, err = mongo.AppUserCol.FindOne(req.UserName)
				if err != nil {
					log.Errorf("prepare to push, but appUser(%s) not found, err:%s", req.UserName, err)
					return
				}
				// 组装报文推送
				req.DeviceToken = appUser.DeviceToken
				req.To = appUser.DeviceType
				if req.DeviceToken != "" {
					log.Debugf("push to user=%s,token=%s", appUser.UserName, appUser.DeviceToken)
					push.Do(req)
				}
			case MsgType_C:
				// 此类为C类消息
				// TODO
			default:
				log.Errorf("prepare to push, but unknown from=%s", req.MsgType)
			}
		}()
	}
}

// A类：校验码[5789],金额 67.89元，交易成功，来自云收银收款码。
// B类：您于14：03分的交易，金额为78.98元，交易成功，点我查看详情。
// C类：全部内容均来源于外部填写。后面会规划平台功能去实现。
func pushMsgReq(t *model.Trans) *model.PushMessageReq {
	req := &model.PushMessageReq{}
	switch strings.ToLower(t.TradeFrom) {
	case model.IOS, model.Android:
		req.MsgType = MsgType_B
		req.Message = MsgType_B
		transTime := t.CreateTime[11:16]
		req.Title = fmt.Sprintf("您于 %s 分的交易，金额为 %0.2f 元，交易成功，点我查看详情。", transTime, float64(t.TransAmt)/100)
	case model.Wap:
		req.MsgType = MsgType_A
		req.Message = MsgType_A
		req.Title = fmt.Sprintf("校验码[%s]，金额 %0.2f 元，交易成功，来自云收银收款码。", t.VeriCode, float64(t.TransAmt)/100)
	default:
		req.MsgType = MsgType_C
		req.Message = MsgType_C
	}
	req.OrderNum = t.OrderNum
	req.MerID = t.MerId
	req.UserName = t.Terminalid
	return req
}
