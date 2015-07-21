package scanpay

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

var openAPIURL = goconf.Config.AlipayScanPay.OpenAPIURL + "?charset=utf-8"

func sendRequest(req BaseReq, body BaseBody, resp BaseResp) error {
	v, err := prepareData(req)
	if err != nil {
		return err
	}

	log.Infof("to alipay message: %s", v.Encode())
	ret, err := send(v)
	if err != nil {
		return err
	}

	err = processResponseBody(ret, body, resp)
	if err != nil {
		return err
	}

	return nil
}

func prepareData(d BaseReq) (v url.Values, err error) {
	v = d.Values()

	// Req to json, then set to biz_content
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	v.Set("biz_content", string(jsonBytes))

	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接 QueryString
	buf := bytes.Buffer{}
	for i, key := range keys {
		if i != 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(v.Get(key))
	}

	// 计算签名
	sign, err := Sha1WithRsa(buf.Bytes(), d.PrivateKey())
	if err != nil {
		log.Errorf("sign error: %s", err)
		return nil, err
	}
	v.Set("sign", sign)

	return v, nil
}

func send(v url.Values) (body []byte, err error) {
	res, err := http.PostForm(openAPIURL, v)

	if err != nil {
		log.Errorf("unable to connect alipay gateway: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	// 处理返回报文
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("unable to read from res %s", err)
		return nil, err
	}
	log.Infof("alipay return message: %s", string(body))

	return body, nil
}

func processResponseBody(data []byte, body BaseBody, resp BaseResp) error {
	// 解析第一层报文，以便验签
	err := json.Unmarshal(data, body)
	if err != nil {
		log.Errorf("json(%s) to struct error: %s", string(data), err)
		return err
	}

	// 验签
	err = Verify(body.GetRaw(), body.GetSign())
	if err != nil {
		log.Errorf("verify sign error, raw=%s, sing=%s: %s", body.GetRaw(), body.GetSign(), err)
		return err
	}

	// 解析第二层报文，得到报文内容
	err = json.Unmarshal(body.GetRaw(), resp)
	if err != nil {
		log.Errorf("unmarshal json(%s) error: %s", body.GetRaw(), err)
		return err
	}
	return nil
}
