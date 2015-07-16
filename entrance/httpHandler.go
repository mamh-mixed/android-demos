package entrance

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/entrance/applepay"
	"github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/CardInfoLink/quickpay/entrance/scanpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Scanpay 扫码支付入口(测试页面)
func Scanpay(w http.ResponseWriter, r *http.Request) {
	log.Debugf("url = %s", r.URL.String())

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "读取数据出错", http.StatusNotAcceptable)
		return
	}

	switch r.URL.Path {
	case "/scanpay/":
		// 请求扫码支付
		retBytes := scanpay.ScanPayHandle(bytes)
		w.Write(retBytes)
	case "/scanpay/query":
		q := &model.QueryCondition{}
		err = json.Unmarshal(bytes, q)
		if err != nil {
			http.Error(w, "数据格式错误: "+err.Error(), http.StatusNotAcceptable)
			return
		}
		ret := core.TransQuery(q)
		retBytes, err := json.Marshal(ret)
		if err != nil {
			http.Error(w, "数据格式错误: "+err.Error(), http.StatusNotAcceptable)
			return
		}
		w.Write(retBytes)
	}
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
	if ret == nil && !result {
		log.Errorf("check sign error %s", err)
		ret = mongo.RespCodeColl.Get("200010")
	}

	// 验签通过，执行业务逻辑
	if result {
		switch r.URL.Path {
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
		default:
			w.WriteHeader(404)
		}
	}

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	// 签名，并返回
	sign = Signature(rdata, merId)
	w.Header().Set("X-Sign", sign)

	log.Infof("to merchant message: %s", rdata)
	w.Write(rdata)
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

	ret := &model.WeixinNotifyResp{ReturnCode: "SUCCESS", ReturnMsg: "OK"}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("read http body error: %s", err)
		ret.ReturnCode = "FAIL"
		ret.ReturnMsg = "报文读取错误"
		return xml.Marshal(ret)
	}

	var req model.WeixinNotifyReq
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
