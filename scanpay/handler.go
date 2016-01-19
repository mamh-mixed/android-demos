package scanpay

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/CardInfoLink/log"
	// "github.com/CardInfoLink/mahonia"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

// scanpayUnifiedHandle 扫码支付入口
func scanpayUnifiedHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("only `POST` method allowed"))
		return
	}

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
	var unescaped string
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("read http body error: %s", err)
		ret.ReturnCode = "FAIL"
		ret.ReturnMsg = "报文读取错误"
	} else {
		log.Infof("weixin notify: %s", data)
		// 微信有时会返回转义后的 xml，json 解析前必须解转义
		unescaped, err = url.QueryUnescape(string(data))
		if err != nil {
			log.Errorf("unescape xml error: %s", err)
			ret.ReturnCode = "FAIL"
			ret.ReturnMsg = "报文解转义错误"
		} else {
			var req weixin.WeixinNotifyReq
			err = xml.Unmarshal([]byte(unescaped), &req)
			if err != nil {
				log.Errorf("unmarshal body error: %s, body: %s", err, string(data))
				ret.ReturnCode = "FAIL"
				ret.ReturnMsg = "报文读取错误"
			} else {
				if err = weixinNotifyCtrl(&req); err != nil {
					log.Errorf("handle weixin notify error: %s, req: %+v", err, req)
				}
			}
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
	log.Infof("alipay notify: %s", data)

	// log.Debugf("before decoder: %s", string(data))
	// gbk-utf8
	unescape, err := url.QueryUnescape(string(data))
	if err != nil {
		log.Errorf("alp notify: %s, unescape error: %s ", string(data), err)
	}

	// d := mahonia.NewDecoder("gbk")
	// utf8 := d.ConvertString(unescape)

	vs, err := url.ParseQuery(unescape)
	if err != nil {
		log.Infof("handle alipay notify error: %s", err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// 处理异步通知
	if err = alipayNotifyCtrl(vs); err != nil {
		log.Errorf("handle alipay notify error: %s, req: %+v", err, vs)
		http.Error(w, "fail", http.StatusOK)
		return
	}

	http.Error(w, "success", http.StatusOK)
}

// alipay2NotifyHandle 接受支付宝异步通知
func alipay2NotifyHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	log.Infof("alipay2 notify: %s", data)

	// log.Debugf("before decoder: %s", string(data))
	// gbk-utf8
	unescape, err := url.QueryUnescape(string(data))
	if err != nil {
		log.Errorf("alp notify: %s, unescape error: %s ", string(data), err)
	}

	// d := mahonia.NewDecoder("gbk")
	// utf8 := d.ConvertString(unescape)

	vs, err := url.ParseQuery(unescape)
	if err != nil {
		log.Infof("handle alipay notify error: %s", err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// 处理异步通知
	if err = alipay2NotifyCtrl(vs); err != nil {
		log.Errorf("handle alipay notify error: %s, req: %+v", err, vs)
		http.Error(w, "fail", http.StatusOK)
		return
	}

	http.Error(w, "success", http.StatusOK)
}

// testReceiveNotifyHandle 测试接受异步通知
func testReceiveNotifyHandle(w http.ResponseWriter, r *http.Request) {

	// data, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusNotAcceptable)
	// 	return
	// }
	log.Infof("receive notify, data: %s", r.URL.RawQuery)

	// response
	respCode := ""
	if rand.Intn(5) == 0 { //  1/5 的概率失败，要求重发
		respCode = "09"
	} else {
		respCode = "00"
	}
	ret := &model.ScanPayResponse{Respcd: respCode}

	retBytes, _ := json.Marshal(ret)

	log.Infof("return notify: %s", retBytes)
	w.Write(retBytes)
}

// scanFixedMerInfoHandle 扫固定码获取用户信息接口
func scanFixedMerInfoHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	b64MerId := r.FormValue("merchantCode")
	if b64MerId == "" {
		http.Error(w, `{"response":"01","errorDetail":"params should not be null"}`, http.StatusOK)
		return
	}
	// 解b64
	mbytes, err := base64.StdEncoding.DecodeString(b64MerId)
	if err != nil {
		log.Errorf("decode merId=%s fail: %s", b64MerId, err)
		http.Error(w, `{"response":"01","errorDetail":"params decode error"}`, http.StatusOK)
		return
	}
	result := query.GetMerInfo(string(mbytes))
	rbytes, err := json.Marshal(result)
	if err != nil {
		log.Errorf("json marshal error:%s", err)
	}
	w.Write(rbytes)
}

// scanFixedOrderInfoHandle 扫固定码获取用户订单信息
func scanFixedOrderInfoHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	uniqueId := r.FormValue("merchantCode")
	log.Debugf("uniqueId: %s", uniqueId)
	if uniqueId == "" {
		http.Error(w, `{"response":"01","errorDetail":"params should not be null"}`, http.StatusOK)
		return
	}
	result := query.GetOrderInfo(uniqueId)
	rbytes, err := json.Marshal(result)
	if err != nil {
		log.Errorf("json marshal error:%s", err)
	}
	w.Write(rbytes)
}

var webAppUrl = goconf.Config.MobileApp.WebAppUrl

// weChatAuthHandle 重定向到微信获取code
func weChatAuthHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	b64MerId := r.FormValue("merchantCode")
	if b64MerId == "" {
		http.Error(w, `{"response":"01","errorDetail":"params should not be null"}`, http.StatusOK)
		return
	}
	// 解b64
	mbytes, err := base64.StdEncoding.DecodeString(b64MerId)
	if err != nil {
		log.Errorf("decode merId=%s fail: %s", b64MerId, err)
		http.Error(w, `{"response":"01","errorDetail":"params decode error"}`, http.StatusOK)
		return
	}

	// Get chanMer info
	pa, err := query.GetPublicAccount(string(mbytes))
	if err != nil {
		http.Error(w, `{"response":"01","errorDetail":"no appID found"}`, http.StatusOK)
		return
	}

	var redirectUri = webAppUrl + "/pay.html" + url.QueryEscape("?merchantCode="+b64MerId+"&showwxpaytitle=1")
	wxpUri := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect"

	// 告诉页面重定向到微信
	http.Redirect(w, r, fmt.Sprintf(wxpUri, pa.AppID, redirectUri), http.StatusFound)
}
