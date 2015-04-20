package bindingpay

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"

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
			ret = bindingCreateHandle(data, merId)
		case "/quickpay/bindingRemove":
			ret = bindingRemoveHandle(data, merId)
		case "/quickpay/bindingEnquiry":
			ret = bindingEnquiryHandle(data, merId)
		case "/quickpay/bindingPayment":
			ret = bindingPaymentHandle(data, merId)
		case "/quickpay/refund":
			ret = bindingRefundHandle(data, merId)
		case "/quickpay/orderEnquiry":
			ret = orderEnquiryHandle(data, merId)
		case "/quickpay/billingDetails":
			ret = billingDetailsHandle(data, merId)
		case "/quickpay/billingSummary":
			ret = billingSummaryHandle(data, merId)
		case "/quickpay/noTrackPayment":
			ret = noTrackPaymentHandle(data, merId)
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

// 建立绑定关系
func bindingCreateHandle(data []byte, merId string) (ret *model.BindingReturn) {
	bc := new(model.BindingCreate)
	err := json.Unmarshal(data, bc)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	bc.MerId = merId

	// 解密特定字段
	aes := new(tools.AesCBCMode)
	bc.AcctNumDecrypt = aes.Decrypt(bc.AcctNum)
	bc.AcctNameDecrypt = aes.Decrypt(bc.AcctName)
	bc.IdentNumDecrypt = aes.Decrypt(bc.IdentNum)
	bc.PhoneNumDecrypt = aes.Decrypt(bc.PhoneNum)
	if bc.AcctType == "20" {
		bc.ValidDateDecrypt = aes.Decrypt(bc.ValidDate)
		bc.Cvv2Decrypt = aes.Decrypt(bc.Cvv2)
	}
	// 报文解密错误
	if aes.Err != nil {
		log.Errorf("decrypt fail : merId=%s, request=%+v, err=%s", merId, bc, aes.Err)
		return mongo.RespCodeColl.Get("200021")
	}
	log.Debugf("after decrypt field : acctNum=%s, acctName=%s, phoneNum=%s, identNum=%s, validDate=%s, cvv2=%s",
		bc.AcctNumDecrypt, bc.AcctNameDecrypt, bc.PhoneNumDecrypt, bc.IdentNumDecrypt, bc.ValidDateDecrypt, bc.Cvv2Decrypt)

	// 验证请求报文是否完整，格式是否正确
	ret = validateBindingCreate(bc)
	if ret != nil {
		return ret
	}

	// 业务处理
	ret = core.ProcessBindingCreate(bc)

	return ret
}

// 解除绑定关系
func bindingRemoveHandle(data []byte, merId string) (ret *model.BindingReturn) {
	br := new(model.BindingRemove)
	err := json.Unmarshal(data, br)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	br.MerId = merId

	ret = validateBindingRemove(br)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBindingReomve(br)
	return ret
}

// 查询绑定关系
func bindingEnquiryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	be := new(model.BindingEnquiry)
	err := json.Unmarshal(data, be)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	be.MerId = merId

	// 验证请求报文格式
	ret = validateBindingEnquiry(be)
	if ret != nil {
		return ret
	}

	ret = core.ProcessBindingEnquiry(be)

	return ret
}

// 绑定支付关系
func bindingPaymentHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBindingPayment(b)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	ret = core.ProcessBindingPayment(b)

	return ret
}

// 退款处理
func bindingRefundHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingRefund)
	err := json.Unmarshal(data, b)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBindingRefund(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBindingRefund(b)

	return ret
}

// 交易对账汇总
func billingSummaryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BillingSummary)
	err := json.Unmarshal(data, b)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}

	b.MerId = merId
	// 验证请求报文格式
	ret = validateBillingSummary(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBillingSummary(b)
	// mock return
	return ret
}

// 交易对账明细
func billingDetailsHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BillingDetails)
	err := json.Unmarshal(data, b)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBillingDetails(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBillingDetails(b)
	// mock return
	return ret
}

// 查询订单状态
func orderEnquiryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.OrderEnquiry)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("解析报文错误 :%s", err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateOrderEnquiry(b)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	ret = core.ProcessOrderEnquiry(b)

	return ret
}

// 无卡直接支付的处理
func noTrackPaymentHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.NoTrackPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("解析报文错误 :%s", err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	ret = validateNoTrackPayment(b)
	if ret != nil {
		return ret
	}

	//  todo 无卡支付暂不开放；业务处理
	ret = mongo.RespCodeColl.Get("100030")
	return ret
}
