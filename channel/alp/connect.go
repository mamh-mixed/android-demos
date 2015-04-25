package alp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"github.com/omigo/log"
	"github.com/qiniu/iconv"
	"io"
	"net/http"
	"net/url"
	"sort"
)

const (
	requestURL = "https://mapi.alipay.com/gateway.do"
)

// sendRequest 发送请求
func sendRequest(params map[string]string, key string) *alpResponse {

	toSign := preContent(params)

	toSign += key
	signed := md5.Sum([]byte(toSign))
	params["sign"] = hex.EncodeToString(signed[:])
	params["sign_type"] = "MD5"

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	log.Debugf("%s", values.Encode())

	res, err := http.PostForm(requestURL, values)
	if err != nil {
		log.Errorf("connect %s fail : %s", requestURL, err)
	}

	return handleResponseBody(res.Body)
}

// handleResponseBody 处理结果集
func handleResponseBody(reader io.Reader) *alpResponse {

	alpResp := new(alpResponse)

	// 重写CharsetReader，使Decoder能解析gbk
	d := xml.NewDecoder(reader)
	d.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		cd, err := iconv.Open("utf-8", s)
		defer cd.Close()
		return iconv.NewReader(cd, r, iconv.DefaultBufSize), err
	}
	err := d.Decode(alpResp)
	if err != nil {
		log.Errorf("unmarsal body fail : %s", err)
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
			delete(params, v)
			continue
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(v + "=" + param)
	}
	return buf.String()
}
