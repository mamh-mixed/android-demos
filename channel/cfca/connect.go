package cfca

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/CardInfoLink/quickpay/config"
	"github.com/omigo/log"
)

var requestURL = config.GetValue("cfca", "url")

var cli *http.Client

// 初始化中金 HTTPS 客户端
func init() {
	ccaCertFile, err := config.GetFile("cfca", "ccaCert")
	if err != nil {
		fmt.Printf("cfca ev_cca_cert config error: %s", err)
		os.Exit(2)
	}
	cfcaEvCcaCrt, err := ioutil.ReadFile(ccaCertFile)
	if err != nil {
		fmt.Printf("read cfca ev_cca_cert error: %s", err)
		os.Exit(3)
	}

	rootCert, err := config.GetFile("cfca", "rootCert")
	if err != nil {
		fmt.Printf("cfca root_cert config error: %s", err)
		os.Exit(2)
	}
	cfcaEvRootCrt, err := ioutil.ReadFile(rootCert)
	if err != nil {
		fmt.Printf("read cfca ev_root_cert error: %s", err)
		os.Exit(3)
	}

	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM(cfcaEvCcaCrt)
	certs.AppendCertsFromPEM(cfcaEvRootCrt)

	// 发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certs,
			// InsecureSkipVerify: true, // only for testing
		},
	}
	cli = &http.Client{Transport: tr}
}

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
	log.Debugf("to send: %s", xmlBytes)

	// 对 xml 作 base64 编码
	b64Str := base64.StdEncoding.EncodeToString(xmlBytes)
	log.Tracef("base64: %s", b64Str)

	// 对 xml 签名
	hexSign := signatureUseSha1WithRsa(xmlBytes, req.SignCert)
	log.Tracef("signed: %s", hexSign)

	// 准备参数
	v = &url.Values{}
	v.Add("message", b64Str)
	v.Add("signature", hexSign)

	return v
}

func send(v *url.Values) (body []byte) {
	resp, err := cli.PostForm(requestURL, *v)
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
	log.Debugf("resp: [%s]", body)

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
	log.Debugf("received: %s", rxmlBytes)

	// 返回消息验签失败的可能性极小，所以异步验签，提高效率
	go func() {
		rhexSign := strings.TrimSpace(result[1])
		log.Tracef("signed: %s", rhexSign)
		err = checkSignatureUseSha1WithRsa(rxmlBytes, rhexSign)
		if err != nil {
			log.Errorf("check sign failed，xml=%s, sign=%s: %s", string(rxmlBytes), rhexSign, err)
		}
	}()

	// 解编 xml
	resp = new(BindingResponse)
	err = xml.Unmarshal(rxmlBytes, resp)
	if err != nil {
		log.Errorf("unable to unmarshal xml %s", err)
		return nil
	}

	return resp
}
