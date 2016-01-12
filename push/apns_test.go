package push

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestAPush(t *testing.T) {
	req := new(model.PushMessageReq)
	req.Device_token = "4951b12f901bb4799c8cb82740cd978ba9fcfa468cb58418ceacf3d38bf89c19"
	req.Title = "test"
	req.Message = "test push"

	err := ApnsPush.APush(req)
	if err != nil {
		t.Error("find agent unsuccessful ", err)
		t.FailNow()
	}
}
