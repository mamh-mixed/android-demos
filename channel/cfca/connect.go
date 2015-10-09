package cfca

import (
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

var requestURL = goconf.Config.CFCA.URL

// sendRequest 对中金接口访问的统一处理
func sendRequest(req *BindingRequest) *BindingResponse {
	values := prepareRequestData(req)
	if values == nil {
		return nil
	}

	body := send(values)
	if body == nil {
		return nil
	}

	return processResponseBody(body)
}

func prepareRequestData(req *BindingRequest) (v *url.Values) {
	// xml 编组
	xmlBytes, err := xml.Marshal(req)
	if err != nil {
		log.Errorf("unable to marshal xml: %s", err)
		return nil
	}
	log.Infof("to send cfca: %s", xmlBytes)

	// 对 xml 作 base64 编码
	b64Str := base64.StdEncoding.EncodeToString(xmlBytes)
	log.Tracef("base64: %s", b64Str)

	// 对 xml 签名
	hexSign, err := signatureUseSha1WithRsa(xmlBytes, req.PrivateKey)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Tracef("signed: %s", hexSign)

	// 准备参数
	v = &url.Values{}
	v.Add("message", b64Str)
	v.Add("signature", hexSign)

	return v
}

func send(v *url.Values) (body []byte) {
	resp, err := http.PostForm(requestURL, *v)
	if err != nil {
		log.Errorf("unable to connect Cfca gratway %s", err)
		return nil
	}
	defer resp.Body.Close()

	// 处理返回报文
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("unable to read from resp %s", err)
		return nil
	}
	log.Debugf("resp from cfca: [%s]", body)

	return body
}

func processResponseBody(body []byte) (resp *BindingResponse) {
	// 得到报文和签名
	result := strings.Split(string(body), ",")
	rb64Str := strings.TrimSpace(result[0])
	// 数据 base64 解码
	rxmlBytes, err := base64.StdEncoding.DecodeString(rb64Str)
	if err != nil {
		log.Errorf("unable to decode base64 content %s", err)
	}
	log.Infof("received from cfca: %s", rxmlBytes)

	// 返回消息验签失败的可能性极小，所以异步验签，提高效率
	// go func() {
	rhexSign := strings.TrimSpace(result[1])
	log.Tracef("signed: %s", rhexSign)
	err = checkSignatureUseSha1WithRsa(rxmlBytes, rhexSign)
	if err != nil {
		log.Errorf("check sign failed，xml=%s, sign=%s: %s", string(rxmlBytes), rhexSign, err)
		return nil
	}
	// }()

	// 解编 xml
	resp = new(BindingResponse)
	err = xml.Unmarshal(rxmlBytes, resp)
	if err != nil {
		log.Errorf("unable to unmarshal xml %s", err)
		return nil
	}

	return resp
}
