package push

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestUPush(t *testing.T) {

	req := new(model.PushMessageReq)
	req.Device_token = "Ai1W8_QftPI6IBRtofzfKnpG3dqeh-j5OYaaQ4Ek3kPz"
	req.Title = "test"
	req.Message = "test push"

	err := UmengPush.UPush(req)
	if err != nil {
		t.Error("find agent unsuccessful ", err)
		t.FailNow()
	}
}
