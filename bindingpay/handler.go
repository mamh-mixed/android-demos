package bindingpay

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/log"
)

// BindingpayHandle 绑定支付统一入口
func BindingpayHandle(w http.ResponseWriter, r *http.Request) {
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
	case "/bindingpay/bindingCreate":
		ret = BindingCreateHandle(data, merId)
	case "/bindingpay/bindingRemove":
		ret = BindingRemoveHandle(data, merId)
	case "/bindingpay/bindingEnquiry":
		ret = BindingEnquiryHandle(data, merId)
	case "/bindingpay/bindingPayment":
		ret = BindingPaymentHandle(data, merId)
	case "/bindingpay/refund":
		ret = BindingRefundHandle(data, merId)
	case "/bindingpay/orderEnquiry":
		ret = OrderEnquiryHandle(data, merId)
	case "/bindingpay/billingDetails":
		ret = BillingDetailsHandle(data, merId)
	case "/bindingpay/billingSummary":
		ret = BillingSummaryHandle(data, merId)
	case "/bindingpay/bindingPayWithSms":
		ret = BindingPayWithSMS(data, merId)
	case "/bindingpay/sendBindingPaySms":
		ret = SendBindingPaySMS(data, merId)
	case "/bindingpay/bindingPaymentSettlement":
		ret = BindingPaymentSettlementHandle(data, merId)
	case "/bindingpay/getCardInfo":
		ret = GetCardInfoHandle(data, merId)
	case "/bindingpay/noTrackPayment":
		ret = NoTrackPaymentHandle(data, merId)
	case "/bindingpay/applePay":
		ret = ApplePayHandle(data, merId)
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

// CheckSignature 根据商户ID到数据库查找签名密钥，然后进行验签
func CheckSignature(data []byte, merId, expected string) (result bool, ret *model.BindingReturn) {
	m, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		if err.Error() == "not found" {
			ret = mongo.RespCodeColl.Get("200063")
			return false, ret
		}
		return false, mongo.RespCodeColl.Get("000001")
	}

	// 如果商户无需验签，验签结果直接返回 true
	if !m.IsNeedSign {
		return true, nil
	}

	if expected == "" {
		return false, nil
	}

	actual := security.SHA1WithKey(string(data), m.SignKey)

	return actual == expected, nil
}

// Signature 根据商户ID到数据库查找签名密钥，然后拼接到数据后面，签名
func Signature(data []byte, merId string) string {
	m, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		log.Errorf("Signature find Merchant error")
		return ""
	}
	return security.SHA1WithKey(string(data), m.SignKey)
}
