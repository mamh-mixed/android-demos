package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
	// "time"
)

// func TestInsert(t *testing.T) {
// 	rsp := new(model.PushMessageRsp)
// 	rsp.UserName = "842712881@qq.Com"
// 	rsp.Title = "test"
// 	rsp.Message = "success"
// 	rsp.PushTime = time.Now().Format("2006-01-02 15:04:05")

// 	PushMessageColl.Insert(rsp)
// }

func TestPushMessageFind(t *testing.T) {
	push := new(model.PushMessage)
	push.UserName = "842712881@qq.Com"
	push.Size = 50
	push.LastTime = "2015-12-21 13:00:00"
	message, err := PushMessageColl.Find(push)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v", message)
}
