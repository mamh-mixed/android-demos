package entrance

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/entrance/applepay"
	"github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/CardInfoLink/quickpay/entrance/scanpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
)

// Scanpay 扫码支付入口
func Scanpay(w http.ResponseWriter, r *http.Request) {
	log.Debugf("url = %s", r.URL.String())

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "read body error", http.StatusNotAcceptable)
		return
	}

	// 请求扫码支付
	retBytes := scanpay.ScanPayHandle(bytes)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(retBytes)
}

// Quickpay 快捷支付统一入口
func Quickpay(w http.ResponseWriter, r *http.Request) {
	log.Infof("url = %s", r.URL.String())

	merId, sign, data, status, err := prepareData(r)
	if err != nil {
		log.Errorf(err.Error())
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}
	log.Infof("from merchant message: %s", data)

	var ret *model.BindingReturn

	result, ret := CheckSignature(data, merId, sign)
	if ret != nil { // 商户不存在等错误
		log.Errorf("merchant error: merId=%s, err=(%s)%s", merId, ret.RespCode, ret.RespMsg)
	} else if !result { // 签名错误
		log.Errorf("check sign error: data=%s, merId=%s, sign=%s", string(data), merId, sign)
		ret = mongo.RespCodeColl.Get("200010")
	} else {
		ret = route(r.URL.Path, data, merId, w)
	}

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	// 签名，并返回
	sign = Signature(rdata, merId)
	w.Header().Set("X-Sign", sign)

	log.Infof("to merchant message: %s", rdata)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(rdata)
}

func route(uri string, data []byte, merId string, w http.ResponseWriter) (ret *model.BindingReturn) {
	switch uri {
	case "/quickpay/bindingCreate":
		ret = bindingpay.BindingCreateHandle(data, merId)
	case "/quickpay/bindingRemove":
		ret = bindingpay.BindingRemoveHandle(data, merId)
	case "/quickpay/bindingEnquiry":
		ret = bindingpay.BindingEnquiryHandle(data, merId)
	case "/quickpay/bindingPayment":
		ret = bindingpay.BindingPaymentHandle(data, merId)
	case "/quickpay/refund":
		ret = bindingpay.BindingRefundHandle(data, merId)
	case "/quickpay/orderEnquiry":
		ret = bindingpay.OrderEnquiryHandle(data, merId)
	case "/quickpay/billingDetails":
		ret = bindingpay.BillingDetailsHandle(data, merId)
	case "/quickpay/billingSummary":
		ret = bindingpay.BillingSummaryHandle(data, merId)
	case "/quickpay/noTrackPayment":
		ret = bindingpay.NoTrackPaymentHandle(data, merId)
	case "/quickpay/applePay":
		ret = applepay.ApplePayHandle(data, merId)
	case "/quickpay/bindingPayWithSms":
		ret = bindingpay.BindingPayWithSMS(data, merId)
	case "/quickpay/sendBindingPaySms":
		ret = bindingpay.SendBindingPaySMS(data, merId)
	default:
		w.WriteHeader(404)
	}
	return ret
}

func prepareData(r *http.Request) (merId, sign string, data []byte, status int, err error) {
	if r.Method != "POST" {
		return "", "", nil, 405, errors.New("only 'POST' method allowed, but actual '" + r.Method + "'")
	}

	v := r.URL.Query()
	merId = v.Get("merId")
	if merId == "" {
		return "", "", nil, 412, errors.New("parameter `merId` required")
	}

	sign = r.Header.Get("X-Sign")
	// 商户可以选择不验签，那么可以不传这个字段
	// if sign == "" {
	// 	return "", "", nil, 412, errors.New("header `X-Sign` required")
	// }

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return "", "", nil, 406, err
	}

	return merId, sign, data, 200, nil
}

// AsyncNotify 处理支付宝、微信异步通知
func AsyncNotify(w http.ResponseWriter, r *http.Request) {
	log.Debugf("url = %s", r.URL.Path)

	var (
		retBytes []byte
		err      error
	)

	switch r.URL.Path {
	case "/qp/back/weixin":
		retBytes, err = weixinNotify(r)
	case "/qp/back/alipay":
		retBytes, err = alipayNotify(r)
	default:
		retBytes = []byte("404")
	}

	if err != nil {
		log.Errorf("read http body error: %s", err)
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	w.Write(retBytes)
}

// WeixinNotify 接受微信异步通知
func weixinNotify(r *http.Request) ([]byte, error) {

	ret := &weixin.WeixinNotifyResp{ReturnCode: "SUCCESS", ReturnMsg: "OK"}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("read http body error: %s", err)
		ret.ReturnCode = "FAIL"
		ret.ReturnMsg = "报文读取错误"
		return xml.Marshal(ret)
	}

	var req weixin.WeixinNotifyReq
	err = xml.Unmarshal(data, &req)
	if err != nil {
		log.Errorf("unmarshal body error: %s, body: %s", err, string(data))
		ret.ReturnCode = "FAIL"
		ret.ReturnMsg = "报文读取错误"
		return xml.Marshal(ret)
	}

	core.ProcessWeixinNotify(&req)

	return xml.Marshal(ret)
}

// AlipayNotify 接受支付宝异步通知
func alipayNotify(r *http.Request) ([]byte, error) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// gbk-utf8
	d := mahonia.NewDecoder("gbk")
	utf8 := d.ConvertString(string(data))

	vs, err := url.ParseQuery(utf8)
	if err != nil {
		return nil, err
	}
	// 处理异步通知
	core.ProcessAlpNotify(vs)

	return []byte("success"), nil
}
