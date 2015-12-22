package push

import (
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/anachronistic/apns"
	"github.com/omigo/log"
)

const (
	pemDir  = "%s/pem/%s"
	certPem = "APNS_CloudCashier_Dev_Cert.pem"
	keyFile = "APNS_CloudCashier_Dev_insecure_key.pem"
	gateWay = "gateway.sandbox.push.apple.com:2195"
)

var ApnsPush apnsPush

type apnsPush struct{}

func (*apnsPush) APush(req *model.PushMessageReq) error {
	dict := apns.NewAlertDictionary()
	dict.Title = req.Title
	// dict.Body = req.Message

	payload := apns.NewPayload()
	payload.Alert = dict
	payload.Category = req.Message

	pn := apns.NewPushNotification()
	pn.DeviceToken = req.Device_token
	pn.AddPayload(payload)

	client := apns.NewClient(gateWay, fmt.Sprintf(pemDir, util.WorkDir, certPem), fmt.Sprintf(pemDir, util.WorkDir, keyFile))
	rsp := client.Send(pn)

	alert, err := pn.PayloadString()
	if err != nil {
		log.Errorf("can't push message %s", alert)
		return errors.New("can't push message")
	}

	err = SavePushMessage(req)

	if err != nil {
		log.Errorf("add push message to table fail %s", err)
		return err
	}

	if !rsp.Success {
		log.Errorf("push message fail %s", alert)
		return errors.New("push message fail")
	}

	return nil
}
