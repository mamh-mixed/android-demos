package push

import (
	"errors"
	"fmt"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/anachronistic/apns"
	"github.com/CardInfoLink/log"
)

var (
	pemDir  = "%s/push/pem/%s"
	certPem = goconf.Config.Push.APNSCertPath
	keyFile = goconf.Config.Push.APNSKeyPath
	gateWay = goconf.Config.Push.APNSGw
)

var client *apns.Client

func init() {
	certDir := fmt.Sprintf(pemDir, util.WorkDir, certPem)
	keyDir := fmt.Sprintf(pemDir, util.WorkDir, keyFile)
	log.Debugf("push client load pem, cert=%s, key=%s", certDir, keyDir)
	client = apns.NewClient(gateWay, certDir, keyDir)
}

var ApnsPush apnsPush

type apnsPush struct{}

func (*apnsPush) APush(req *model.PushMessageReq) error {

	payload := apns.NewPayload()
	payload.Alert = req.Title
	payload.Category = req.Message
	payload.Sound = "ding.mp3"

	pn := apns.NewPushNotification()
	pn.DeviceToken = req.DeviceToken
	pn.AddPayload(payload)

	rsp := client.Send(pn)
	alert, err := pn.PayloadString()
	if err != nil {
		log.Errorf("can't push message %s", alert)
		return errors.New("can't push message")
	}

	if !rsp.Success {
		log.Errorf("push message fail %s, error: %s", alert, rsp.Error)
		return errors.New("push message fail")
	}

	return nil
}
