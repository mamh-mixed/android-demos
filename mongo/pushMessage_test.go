package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	rsp := new(model.PushMessageRsp)
	rsp.UserName = "842712881@qq.Com"
	rsp.Title = "test"
	rsp.Message = "success"
	rsp.PushTime = time.Now().Format("2006-01-02 15:04:05")

	PushMessageColl.Insert(rsp)
}
