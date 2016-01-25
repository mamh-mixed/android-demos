package app

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/push"
	"strings"
	"testing"
)

func TestGuangBo(t *testing.T) {
	apps, err := mongo.AppUserCol.Find(&model.AppUserContiditon{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("find apps=%d", len(apps))
	// return

	var ss int
	for _, u := range apps {
		if u.DeviceToken != "" && u.MerId != "" {
			ss++
			var to = strings.ToLower(u.DeviceType)
			push.Do(&model.PushMessageReq{
				MerID:       u.MerId,
				UserName:    u.UserName,
				Title:       "您好，经与微信确认及我们持续观察，目前微信支付已可正常使用。给您带来的不便敬请谅解。",
				Message:     "微信交易恢复通知",
				DeviceToken: u.DeviceToken,
				MsgType:     MsgType_C,
				To:          to,
				// OrderNum:    "",
			})
		}
	}
	t.Logf("send %d", ss)
}
