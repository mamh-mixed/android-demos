package alp

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"github.com/omigo/log"
	"net/http"
	"net/url"
	"sort"
)

const (
	requestURL = "https://mapi.alipay.com/gateway.do"
)

// sendRequest 发送请求
func sendRequest(params map[string]string, key string) *AlpResponse {

	toSign := preContent(params)
	toSign += key
	signed := md5.Sum([]byte(toSign))
	params["sign"] = string(signed[:])
	params["sign_type"] = "MD5"

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	// req := values.Encode()
	res, err := http.PostForm(requestURL, values)
	if err != nil {
		log.Errorf("connect %s fail : %s", requestURL, err)
	}
	log.Debugf("response body : %s", res)

	return handleResponseBody(res)
}

// handleResponseBody 处理结果集
func handleResponseBody(body []byte) *AlpResponse {

	alpResp := new(AlpResponse)
	err := xml.Unmarshal(body, alpResp)
	if err != nil {
		log.Errorf("unmarsal body(%s) fail : %s", body, err)
	}

	// TODO 验证签名

	return alpResp
}

// preContent 待签名字符串
func preContent(params map[string]string) string {
	s := make([]string, 0, len(params))
	for k, _ := range params {
		s = append(s, k)
	}
	// 排序
	sort.Strings(s)

	var buf bytes.Buffer
	for _, v := range s {
		param := params[v]
		// 过滤掉空串
		if param == "" {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(v + "=" + param)
		}
	}
	return buf.String()
}
