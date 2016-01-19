package domestic

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/logs"
	"github.com/CardInfoLink/log"
	// "github.com/CardInfoLink/mahonia"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
)

var requestURL = goconf.Config.AlipayScanPay.URL

// sendRequest 发送请求
func sendRequest(alpReq *alpRequest) (*alpResponse, error) {

	req := alpReq.SpReq
	if req == nil {
		return nil, fmt.Errorf("%s", "no params spReq found")
	}

	// 记录日志
	logs.SpLogs <- req.GetChanReqLogs(alpReq)

	params := make(map[string]string)
	params = toMap(alpReq)

	toSign := preContent(params)

	toSign += req.SignKey
	fmt.Printf("the toSign:%s\n", toSign)
	signed := md5.Sum([]byte(toSign))
	params["sign"] = hex.EncodeToString(signed[:])
	params["sign_type"] = "MD5"

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	log.Infof("post alipay data: %s", values.Encode())

	var res *http.Response
	var err error
	var count = 0
	// 重试
	for {
		count++
		res, err = http.PostForm(requestURL, values)
		if err != nil {
			log.Errorf("connect %s fail : %s, retry ... %d", requestURL, err, count)
			if count == 3 {
				return nil, err
			}
			continue
		}
		break
	}

	defer res.Body.Close()

	alpResp, err := handleResponseBody(res.Body)
	if err != nil {
		return nil, err
	}

	// 记录日志
	logs.SpLogs <- req.GetChanRetLogs(alpResp)

	return alpResp, nil
}

// handleResponseBody 处理结果集
func handleResponseBody(reader io.Reader) (*alpResponse, error) {

	alpResp := new(alpResponse)

	// 重写CharsetReader，使Decoder能解析gbk
	// d := xml.NewDecoder(reader)
	// d.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
	// 	dec := mahonia.NewDecoder(s)
	// 	if dec == nil {
	// 		return nil, fmt.Errorf("not support %s", s)
	// 	}
	// 	return dec.NewReader(r), nil
	// }
	// err := d.Decode(alpResp)
	// if err != nil {
	// 	log.Errorf("unmarsal body fail : %s", err)
	// 	return nil, err
	// }
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(bs, alpResp)
	if err != nil {
		return nil, err
	}
	log.Infof("alp response body: \n %+v \n", alpResp)
	// TODO:验证签名

	return alpResp, nil
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

//账单请求
func sendSettleRequest(alpReq *alpRequest) (*alpSettleResponse, error) {

	req := alpReq.SpReq
	if req == nil {
		return nil, fmt.Errorf("%s", "no params spReq found")
	}

	// 记录日志
	// logs.SpLogs <- req.GetChanReqLogs(alpReq)

	params := make(map[string]string)
	params = toSettleMap(alpReq)

	toSign := preContent(params)

	toSign += req.SignKey
	signed := md5.Sum([]byte(toSign))
	params["sign"] = hex.EncodeToString(signed[:])
	params["sign_type"] = "MD5"

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	log.Infof("post alipay data: %s", values.Encode())

	var res *http.Response
	var err error
	var count = 0
	// 重试
	for {
		count++
		res, err = http.PostForm(requestURL, values)
		if err != nil {
			log.Errorf("connect %s fail : %s, retry ... %d", requestURL, err, count)
			if count == 3 {
				return nil, err
			}
			continue
		}
		break
	}

	defer res.Body.Close()

	alpResp, err := handleSettleResponseBody(res.Body)
	if err != nil {
		return nil, err
	}

	// 记录日志
	// logs.SpLogs <- req.GetChanRetLogs(alpResp)

	return alpResp, nil

}

//查询账单
func handleSettleResponseBody(reader io.Reader) (*alpSettleResponse, error) {

	alpResp := new(alpSettleResponse)

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	log.Debugf("alp response body: %s", string(bs))
	err = xml.Unmarshal(bs, alpResp)
	if err != nil {
		return nil, err
	}
	// TODO:验证签名

	return alpResp, nil
}
