package scanpay

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
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

// scanpayBillsHandle 清算对账
func scanpayBillsHandle(w http.ResponseWriter, r *http.Request) {

	var errorResp = `{"respcd":"%s","errorDetail":"%s"}`

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "read body error", http.StatusNotAcceptable)
		return
	}

	var billsReq = &struct {
		MerId        string `json:"mchntid"`
		Busicd       string `json:"busicd"`
		SettDate     string `json:"settDate"`
		NextOrderNum string `json:"nextOrderNum"`
		Sign         string `json:"sign"`
	}{}

	err = json.Unmarshal(bytes, billsReq)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(errorResp, "30", "报文错误")))
		return
	}

	// TODO:验签

	// 获取对账单
	result := query.GetBills(&model.QueryCondition{
		MerId:        billsReq.MerId,
		Busicd:       billsReq.Busicd,
		StartTime:    billsReq.SettDate + " 00:00:00",
		EndTime:      billsReq.SettDate + " 23:59:59",
		NextOrderNum: billsReq.NextOrderNum,
	})

	retBytes, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(errorResp, "01", "系统错误，请重试。")))
		return
	}

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
		log.Infof("weixin notify: %s", data)

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

	log.Infof("return weixin: %s", retBytes)
	w.Write(retBytes)
}

// alipayNotifyHandle 接受支付宝异步通知
func alipayNotifyHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	log.Infof("alipay notify(GBK): %s", data)

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
		log.Infof("return alipay: %s", err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// 处理异步通知
	err = alipayNotifyCtrl(vs)
	if err != nil {
		log.Info("return alipay: fail")
		http.Error(w, "fail", http.StatusOK)
		return
	}

	log.Info("return alipay: success")
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
