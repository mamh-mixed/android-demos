package entrance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/entrance/applepay"
	"github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/omigo/log"
)

// BindingPay 绑定支付入口
func BindingPay(w http.ResponseWriter, r *http.Request) {
	log.Debugf("url = %s", r.URL.Path)

	merId, sign, data, status, err := prepareData(r)
	if err != nil {
		log.Errorf(err.Error())
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}
	log.Debugf("from mer msg: %s", data)

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

	log.Debugf("to mer msg: %s", rdata)
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
	if sign == "" {
		return "", "", nil, 412, errors.New("parameter `X-Sign` required")
	}

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return "", "", nil, 406, err
	}

	return merId, sign, data, 200, nil
}
