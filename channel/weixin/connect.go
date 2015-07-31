package weixin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"github.com/omigo/validator"
)

var url = goconf.Config.WeixinScanPay.URL

// Execute 发送报文执行微信支付
func Execute(req BaseReq, resp BaseResp) error {
	if err := validator.Validate(req); err != nil {
		log.Errorf("validate error, %s", err)
		return err
	}

	xmlBytes, err := prepareData(req)
	if err != nil {
		return err
	}

	log.Infof(">>> send to weixin: %s", string(xmlBytes))
	ret, err := send(req.GetHTTPClient(), req.GetURI(), xmlBytes)
	if err != nil {
		return err
	}
	log.Infof("<<< return from weixin: %s", string(ret))

	return processRespBody(ret, req.GetSignKey(), resp)
}

func prepareData(d BaseReq) (xmlBytes []byte, err error) {
	buf, err := util.Query(d)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	buf.WriteString("&key=" + d.GetSignKey())
	// log.Debugf("%s", buf.String())

	sign := md5.Sum(buf.Bytes())
	d.SetSign(strings.ToUpper(hex.EncodeToString(sign[:])))

	xmlBytes, err = xml.Marshal(d)
	if err != nil {
		log.Errorf("struct(%#v) to xml error: %s", d, err)
		return nil, err
	}

	return xmlBytes, nil
}

func send(cli *http.Client, uri string, body []byte) (ret []byte, err error) {
	var resp *http.Response

	// 如果连接失败，重试 3 次，休眠 3s、6s
	for i := 1; i <= 3; i++ {
		resp, err = cli.Post(url+uri, "text/xml", bytes.NewBuffer(body))
		if err == nil {
			break
		}
		log.Warnf("unable to connect WeixinScanPay gateway: %s", err)
		sleepTime := time.Duration(i*3) * time.Second
		log.Warnf("after %s retry...", sleepTime)
		time.Sleep(sleepTime)
	}

	// 如果有错，说明重试了 3 次，都失败，那么直接返回
	if err != nil {
		log.Errorf("unable to connect WeixinScanPay gateway: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 处理返回报文
	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("unable to read from resp %s", err)
		return nil, err
	}

	return ret, nil
}

func processRespBody(body []byte, signKey string, resp BaseResp) error {
	err := xml.Unmarshal(body, resp)
	if err != nil {
		log.Errorf("xml(%s) to struct error: %s", string(body), err)
		return err
	}

	// 验签
	buf, err := util.Query(resp)
	if err != nil {
		log.Error(err)
		return err
	}
	buf.WriteString("&key=" + signKey)
	log.Debugf("%s", buf.String())

	sign := md5.Sum(buf.Bytes())
	actual := strings.ToUpper(hex.EncodeToString(sign[:]))

	if actual != resp.GetSign() {
		log.Errorf("check sign error: query={%s}, expected=%s, actual=%s", buf.String(), resp.GetSign(), actual)
		return errors.New("check sign error")
	}

	return nil
}
