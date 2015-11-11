package scanpay2

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

var openAPIURL = goconf.Config.AlipayScanPay.OpenAPIURL + "?charset=" + CharsetUTF8

func sendRequest(req BaseReq, resp BaseResp) error {
	v, err := prepareData(req)
	if err != nil {
		return err
	}
	log.Infof(">>> to alipay message: %s", v.Encode())

	body, err := send(v)
	if err != nil {
		return err
	}
	log.Infof("<<< alipay return message: %s", string(body))

	err = parseBody(body, resp)
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
	sign, err := Sha1WithRsa(buf.Bytes(), d.GetPrivateKey())
	if err != nil {
		log.Errorf("sign error: %s", err)
		return nil, err
	}
	v.Set("sign", sign)

	return v, nil
}

func send(v url.Values, gw string) (body []byte, err error) {
	res, err := http.PostForm(gw, v)

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

	return body, nil
}

func parseBody(data []byte, resp BaseResp) error {
	// 解析第一层报文，以便验签
	err := json.Unmarshal(data, resp)
	if err != nil {
		log.Errorf("json(%s) to struct error: %s", string(data), err)
		return err
	}

	// 验签
	err = Verify(resp.GetRaw(), resp.GetSign())
	if err != nil {
		log.Errorf("verify sign error, raw=%s, sing=%s: %s", resp.GetRaw(), resp.GetSign(), err)
		return err
	}

	// 解析第二层报文，得到报文内容
	err = json.Unmarshal(resp.GetRaw(), resp)
	if err != nil {
		log.Errorf("unmarshal json(%s) error: %s", resp.GetRaw(), err)
		return err
	}
	return nil
}
