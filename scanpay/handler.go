package scanpay

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
)

// scanpayUnifiedHandle 扫码支付入口
func scanpayUnifiedHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Debugf("url = %s", r.URL.String())

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "read body error", http.StatusNotAcceptable)
		return
	}

	// 请求扫码支付
	retBytes := ScanPayHandle(bytes, false)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(retBytes)
}

// weixinNotifyHandle 接受微信异步通知
func weixinNotifyHandle(w http.ResponseWriter, r *http.Request) {
	ret := &weixin.WeixinNotifyResp{ReturnCode: "SUCCESS", ReturnMsg: "OK"}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("read http body error: %s", err)
		ret.ReturnCode = "FAIL"
		ret.ReturnMsg = "报文读取错误"
	} else {
		var req weixin.WeixinNotifyReq
		err = xml.Unmarshal(data, &req)
		if err != nil {
			log.Errorf("unmarshal body error: %s, body: %s", err, string(data))
			ret.ReturnCode = "FAIL"
			ret.ReturnMsg = "报文读取错误"
		} else {
			err = weixinNotifyCtrl(&req)
		}
	}
	if err != nil {
		ret.ReturnCode = "FAIL"
		ret.ReturnMsg = "SYSTEM_ERROR"
	}

	retBytes, err := xml.Marshal(ret)

	if err != nil {
		log.Errorf("read http body error: %s", err)
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	w.Write(retBytes)
}

// alipayNotifyHandle 接受支付宝异步通知
func alipayNotifyHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// log.Debugf("before decoder: %s", string(data))
	// gbk-utf8
	unescape, err := url.QueryUnescape(string(data))
	if err != nil {
		log.Errorf("alp notify: %s, unescape error: %s ", string(data), err)
	}

	d := mahonia.NewDecoder("gbk")
	utf8 := d.ConvertString(unescape)

	vs, err := url.ParseQuery(utf8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// 处理异步通知
	err = alipayNotifyCtrl(vs)
	if err != nil {
		http.Error(w, "fail", http.StatusOK)
		return
	}

	http.Error(w, "success", http.StatusOK)
}

// testReceiveNotifyHandle 测试接受异步通知
func testReceiveNotifyHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	log.Infof("receive notify: %s", string(data))

	// response
	respCode := ""
	if rand.Intn(5) == 0 { //  1/5 的概率失败，要求重发
		respCode = "09"
	} else {
		respCode = "00"
	}
	ret := &model.ScanPayResponse{Respcd: respCode}

	retBytes, _ := json.Marshal(ret)
	log.Debug("return notify: %s", retBytes)
	w.Write(retBytes)
}
