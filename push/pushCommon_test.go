package push

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestPushInfos(t *testing.T) {
	req := new(model.PushMessageRsp)
	req.UserName = "842712881@qq.Com"
	req.Password = "96e79218965eb72c92a549dd5a330112"
	req.Size = "10"
	req.Index = "1"
	req.LastTime = "2015-12-22 08:19:00"

	rsp := PushInfos(req)
	fmt.Println("info:", rsp.Error, rsp.Count)
}
