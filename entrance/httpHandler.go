package entrance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/entrance/applepay"
	"github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/CardInfoLink/quickpay/entrance/scanpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/omigo/log"
)

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

// QuickPayBack 网关接受支付宝、微信等异步通知
func QuickPayBack(w http.ResponseWriter, r *http.Request) {
	log.Infof("url = %s", r.URL.String())

	content, values := "", r.URL.Query()
	switch r.URL.Path {
	case "/quickpay/back/alp":
		// TODO check sign
		values.Add("scanpay_chcd", "ALP")
		content = "success"
	case "/quickpay/back/wxp":
		// TODO check sign
		values.Add("scanpay_chcd", "WXP")
		content = "success"
	default:
		http.Error(w, "invalid request!", http.StatusNotFound)
		return
	}

	// 处理异步通知
	scanpay.AsyncNotifyRouter(values)
	w.Write([]byte(content))
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
