package scanpay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

var url = goconf.Config.WeixinScanPay.URL

const (
	ScanPayURI      = "/pay/micropay"
	ScanPayQueryURI = "/pay/orderquery"
)

func sendRequest(uri string, d BaseData, respData interface{}) error {
	xmlBytes, err := prepareData(d)
	if err != nil {
		return err
	}

	ret, err := send(uri, xmlBytes)
	if err == nil {
		return err
	}

	return processResponseBody(ret, respData)
}

func prepareData(d BaseData) (xmlBytes []byte, err error) {
	d.GenSign()

	xmlBytes, err = xml.Marshal(d)
	if err != nil {
		log.Errorf("struct(%#v) to xml error: %s", d, err)
		return nil, err
	}

	return xmlBytes, nil
}

func send(uri string, body []byte) (ret []byte, err error) {
	resp, err := http.Post(url+uri, "text/xml", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("unable to connect WeixinScanPay gratway %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 处理返回报文
	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("unable to read from resp %s", err)
		return nil, err
	}
	log.Debugf("resp: [%s]", ret)

	return ret, nil
}

func processResponseBody(body []byte, respData interface{}) error {
	err := xml.Unmarshal(body, respData)
	if err != nil {
		log.Errorf("xml(%s) to struct error: %s", string(body), err)
		return err
	}
	return nil
}
