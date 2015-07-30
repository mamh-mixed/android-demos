package enterprisepay

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
	"io/ioutil"
	"net/http"
	"os"
)

var url = goconf.Config.WeixinScanPay.URL

var cli *http.Client

// 初始化微信 HTTPS 客户端
func init() {
	cliCrt, err := tls.LoadX509KeyPair(goconf.Config.WeixinScanPay.ClientCert,
		goconf.Config.WeixinScanPay.ClientKey)
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		os.Exit(4)
	}

	// 发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// InsecureSkipVerify: true, // only for testing
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	cli = &http.Client{Transport: tr}
}

const (
	enterprisePayURI   = "/mmpaymkttransfers/promotion/transfers"
	enterpriseQueryURI = "/mmpaymkttransfers/gettransferinfo"
)

func getUri(req BaseReq) string {
	switch v := req.(type) {

	case *EnterprisePayReq:
		return enterprisePayURI
	default:
		log.Errorf("unknown BaseReq type: %#v", v)
		return "/404"
	}
}

func sendRequest(req BaseReq, resp BaseResp) error {
	xmlBytes, err := prepareData(req)
	if err != nil {
		return err
	}

	log.Infof(">>> send to weixin: %s", string(xmlBytes))
	ret, err := send(getUri(req), xmlBytes)
	if err != nil {
		return err
	}
	log.Infof("<<< return from weixin: %s", string(ret))

	return processResponseBody(ret, resp)
}

func prepareData(d BaseReq) (xmlBytes []byte, err error) {
	d.GenSign()

	xmlBytes, err = xml.Marshal(d)
	if err != nil {
		log.Errorf("struct(%#v) to xml error: %s", d, err)
		return nil, err
	}

	return xmlBytes, nil
}

func send(uri string, body []byte) (ret []byte, err error) {
	var resp *http.Response

	resp, err = cli.Post(url+uri, "text/xml", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("unable to connect WeixinEnterprisePay gateway: %s", err)
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

func processResponseBody(body []byte, respData interface{}) error {
	err := xml.Unmarshal(body, respData)
	if err != nil {
		log.Errorf("xml(%s) to struct error: %s", string(body), err)
		return err
	}

	return nil
}
