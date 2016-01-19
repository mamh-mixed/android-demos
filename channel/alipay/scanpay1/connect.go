// dmf1.0
package scanpay1

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/CardInfoLink/quickpay/logs"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Execute 发送报文执行微信支付
func Execute(req BaseReq, resp BaseResp) error {

	m := req.GetSpReq()
	if m == nil {
		return fmt.Errorf("%s", "no params spReq found")
	}

	// 记录请求渠道日志
	logs.SpLogs <- m.GetChanReqLogs(req)

	params, err := prepareData(req)
	if err != nil {
		return err
	}
	log.Infof(">>> send to alipay: %+v", params)

	ret, err := send(req.GetURI(), params)
	if err != nil {
		return err
	}
	log.Infof("<<< return from alipay: %s", string(ret))

	err = processRespBody(ret, resp)

	// 记录渠道返回日志
	logs.SpLogs <- m.GetChanRetLogs(resp)

	return err
}

func prepareData(req BaseReq) (url.Values, error) {

	buf, err := util.Query(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debugf("sign content: %s", buf.String())

	// gbk encoding
	content := buf.String()
	signed := md5.Sum([]byte(content + req.GetSignKey()))
	params, err := url.ParseQuery(content)
	if err != nil {
		return nil, err
	}

	params.Add("sign", hex.EncodeToString(signed[:]))
	params.Add("sign_type", "MD5")

	return params, nil
}

func send(gw string, params url.Values) ([]byte, error) {
	var res *http.Response
	var err error
	var count = 0
	// 重试
	for {
		count++
		res, err = http.PostForm(gw, params)
		if err != nil {
			log.Errorf("connect %s fail : %s, retry ... %d", gw, err, count)
			if count == 3 {
				return nil, err
			}
			continue
		}
		break
	}

	defer res.Body.Close()

	ret, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// processRespBody 处理结果集
func processRespBody(ret []byte, resp BaseResp) error {

	// 重写CharsetReader，使Decoder能解析gbk
	// d := xml.NewDecoder(bytes.NewReader(ret))
	// d.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
	// 	dec := mahonia.NewDecoder(s)
	// 	if dec == nil {
	// 		return nil, fmt.Errorf("not support %s", s)
	// 	}
	// 	return dec.NewReader(r), nil
	// }
	// err := d.Decode(resp)
	// if err != nil {
	// 	log.Errorf("unmarsal body fail : %s", err)
	// 	return err
	// }
	err := xml.Unmarshal(ret, resp)

	// TODO:验证签名
	return err
}
