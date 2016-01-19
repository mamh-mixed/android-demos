package weixin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/logs"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"github.com/omigo/validator"
)

// DownloadBill 下载对账单
func DownloadBill(req BaseReq) ([]byte, error) {
	xmlBytes, err := prepareData(req)
	if err != nil {
		return nil, err
	}

	ret, err := send(req.GetHTTPClient(), req.GetURI(), xmlBytes)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Execute 发送报文执行微信支付
func Execute(req BaseReq, resp BaseResp) error {

	m := req.GetSpReq()
	if m == nil {
		return fmt.Errorf("%s", "no params spReq found")
	}

	// 记录请求渠道日志
	logs.SpLogs <- m.GetChanReqLogs(req)

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

	err = processRespBody(ret, req.GetSignKey(), resp)

	// 记录渠道返回日志
	logs.SpLogs <- m.GetChanRetLogs(resp)

	return err
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
		// start := time.Now()
		resp, err = cli.Post(goconf.Config.WeixinScanPay.URL+uri, "text/xml", bytes.NewBuffer(body))
		// end := time.Now()
		// log.Infof("=== %s === %s%s", end.Sub(start), goconf.Config.WeixinScanPay.URL, uri)
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

	// log.Debugf("content Type : %s", resp.Header.Get("Content-Type")) // application/octet-stream

	return ret, nil
}

func processRespBody(body []byte, signKey string, resp BaseResp) error {
	err := xml.Unmarshal(body, resp)
	if err != nil {
		log.Errorf("xml(%s) to struct error: %s", string(body), err)
		return err
	}

	// 没有返回签名时，不验签，只打印日志，不中止逻辑
	if resp.GetSign() == "" {
		log.Error("sign is blank, skip check sign")
		return nil
	}

	// 验签, 如果验签失败，只打印日志，不中止逻辑
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
		return nil
	}

	return nil
}

//结算请求
func SettleExecute(req BaseReq, resp BaseResp) (string, error) {
	m := req.GetSpReq()
	if m == nil {
		return "", fmt.Errorf("%s", "no params spReq found")
	}

	// 记录请求渠道日志
	// logs.SpLogs <- m.GetChanReqLogs(req)

	if err := validator.Validate(req); err != nil {
		log.Errorf("validate error, %s", err)
		return "", err
	}

	xmlBytes, err := prepareData(req)
	if err != nil {
		return "", err
	}

	log.Infof(">>> send to weixin: %s", string(xmlBytes))
	ret, err := send(req.GetHTTPClient(), req.GetURI(), xmlBytes)
	if err != nil {
		return "", err
	}

	return string(ret), err
}
