package push

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	//"github.com/CardInfoLink/quickpay/mongo"
	"bytes"
	"fmt"
	"github.com/CardInfoLink/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	method            = "POST"
	appkey            = "55ebd98367e58e0f13002f2a"
	app_master_secret = "s75ash0xyhvl2m3vtzip6lth6muwai7i"
	umeng_url         = "http://msg.umeng.com/api/send"

	message_type = "unicast"
	display_type = "notification"
	after_open   = "go_app"
)

var UmengPush umeng

type umeng struct{}

func (*umeng) UPush(req *model.PushMessageReq) error {
	umengBody := UmengBody{
		Ticker:     req.Title,
		Title:      req.Title,
		Text:       req.Message,
		After_open: after_open,
	}
	params := UmengMessage{
		Appkey:        appkey,
		Timestamp:     time.Now().Unix(),
		Device_tokens: req.DeviceToken,
		Type:          message_type,
		Payload: UmengPayload{
			Body:         umengBody,
			Display_type: display_type,
		},
	}

	post_body, err := json.Marshal(params)
	if err != nil {
		log.Errorf("convert UmengMessage to json error,%s", err)
		return err
	}

	signStr := method + umeng_url + string(post_body) + app_master_secret
	stringMd5 := md5.Sum([]byte(signStr))
	byteArray := stringMd5[:]
	signKey := strings.ToLower(hex.EncodeToString(byteArray))

	stempUrl := umeng_url + "?sign=%s"
	var res *http.Response
	var count = 0
	// 重试
	for {
		count++
		res, err = http.Post(fmt.Sprintf(stempUrl, signKey), "application/json", bytes.NewReader(post_body))
		if err != nil {
			log.Errorf("connect %s fail : %s, retry ... %d", stempUrl, err, count)
			if count == 3 {
				return err
			}
			continue
		}
		break
	}

	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Debugf("response body: %s", string(bs))

	return nil
}
