package weixin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/mahonia"
)

type WeixinRequest interface {
	copyData(scanPayReq *model.ScanPay)
}

type WeixinResponse interface {
	convertToScanPayResp() *model.ScanPayResponse
}

type WeixinPay struct {
	WeixinRequestFactory
	WeixinResponseFactory
}

func (c *WeixinPay) requestWeixin(m WeixinRequest, url string) WeixinResponse {
	rep := sendRequestWithXMLBody(m, url)
	fmt.Println("http response from weixin:", rep)
	body := rep.Body
	defer rep.Body.Close()

	d := xml.NewDecoder(body)

	d.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		dec := mahonia.NewDecoder(s)
		if dec == nil {
			return nil, fmt.Errorf("not support %s", s)
		}
		return dec.NewReader(r), nil
	}

	weixinResp := c.createResponseData(m)
	err := d.Decode(weixinResp)
	if err != nil {
		// log.Errorf("unmarsal body fail : %s", err)
		log.Printf("unmarsal body fail : %s", err)
	}
	fmt.Println("weixinResp:", weixinResp)
	return weixinResp
}

func sendRequestWithXMLBody(m WeixinRequest, url string) *http.Response {
	buf, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(buf))

	bytebuf := bytes.NewBuffer(buf)
	// the http response from weixin
	rep, err := http.Post(url, "text/xml", bytebuf)

	if err != nil {
		log.Println("http request errors", err)
	}

	return rep
}
