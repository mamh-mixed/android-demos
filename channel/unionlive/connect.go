package unionlive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/CardInfoLink/log"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/logs"
	"github.com/CardInfoLink/validator"
)

var requestURL = goconf.Config.UnionLive.URL
var channelId = goconf.Config.UnionLive.ChannelId

// Execute 对优方接口访问的统一处理
func Execute(req BaseReq, resp BaseResp) error {

	m := req.GetSpReq()
	if m == nil {
		return fmt.Errorf("%s", "no params spReq found")
	}

	// 记录请求渠道日志
	logs.SpLogs <- m.GetChanReqLogs(req)

	if ok, errs := validator.Validate(req); !ok {
		log.Errorf("validate error, %v", errs)
		return errors.New("validate error")
	}

	message, err := prepareData(req)
	if err != nil {
		return err
	}

	ret, err := send(req.GetT(), message)
	if err != nil {
		return err
	}

	err = processRespBody(ret, resp)

	// 记录渠道返回日志
	logs.SpLogs <- m.GetChanRetLogs(resp)

	return err
}

func prepareData(req BaseReq) (message []byte, err error) {
	// json 编组
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		log.Errorf("unable to marshal json: %s", err)
		return nil, err
	}
	log.Infof("send to unionlive: %s", jsonBytes)

	return encryptAndSign(jsonBytes)
}

func send(t string, body []byte) (ret []byte, err error) {
	var resp *http.Response

	url := fmt.Sprintf("%s?t=%s&channelId=%s", requestURL, t, channelId)

	// 如果连接失败，重试 3 次，休眠 3s、6s
	for i := 1; i <= 3; i++ {
		start := time.Now()
		resp, err = http.Post(url, "text/plain", bytes.NewBuffer(body))
		end := time.Now()
		log.Infof("=== %s === %s", end.Sub(start), url)
		if err == nil {
			break
		}
		log.Warnf("unable to connect UnionLive gateway: %s", err)
		sleepTime := time.Duration(i*3) * time.Second
		log.Warnf("after %s retry...", sleepTime)
		time.Sleep(sleepTime)
	}

	// 如果有错，说明重试了 3 次，都失败，那么直接返回
	if err != nil {
		log.Errorf("unable to connect UnionLive gateway: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 处理返回报文
	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("unable to read from resp %s", err)
		return nil, err
	}

	// log.Debugf("content Type : %s", resp.Header.Get("Content-Type")) // application/octet-stream

	return ret, nil
}

func processRespBody(message []byte, resp BaseResp) error {
	if len(message) < 33 {
		log.Errorf("get error message: %s", string(message))
		return fmt.Errorf("message length error, expected > 32, actual=%d", len(message))
	}

	rbody, err := checkSignAndDecrypt(message)
	if err != nil {
		return err
	}

	log.Debugf("return from unionlive: %s", rbody)

	// 解编 json
	err = json.Unmarshal(rbody, resp)
	if err != nil {
		log.Errorf("unable to unmarshal json(%s): %s", rbody, err)
		return err
	}

	return nil
}
