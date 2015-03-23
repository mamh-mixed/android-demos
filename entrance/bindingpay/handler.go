package bindingpay

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/g"
)

// BindingPay 绑定支付入口
func BindingPay(w http.ResponseWriter, r *http.Request) {
	g.Debug("url = %s", r.URL.Path)

	merId, sign, data, status, err := prepareData(r)
	if err != nil {
		g.Error(err.Error())
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}
	g.Debug("商户报文: %s", data)

	var ret *model.BindingReturn

	result, ret := CheckSignature(data, merId, sign)
	if ret == nil && !result {
		g.Error("check sign error ", err)
		ret = model.NewBindingReturn("200010", "验证签名失败")
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
		case "/quickpay/noTrackPayment":
			ret = noTrackPaymentHandle(data, merId)
		default:
			w.WriteHeader(404)
		}
	}
	g.Debug("处理后报文: %+v", ret)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	// 签名，并返回
	sign = Signature(rdata, merId)
	w.Header().Set("X-Sign", sign)

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
		return model.NewBindingReturn("200002", "解析报文错误")
	}
	bc.MerId = merId

	// 验证请求报文是否完整，格式是否正确
	ret = validateBindingCreate(bc)
	if ret != nil {
		return ret
	}

	//todo 业务处理
	ret = core.ProcessBindingCreate(bc)

	return ret
}

// 解除绑定关系
func bindingRemoveHandle(data []byte, merId string) (ret *model.BindingReturn) {
	br := new(model.BindingRemove)
	err := json.Unmarshal(data, br)
	if err != nil {
		return model.NewBindingReturn("200002", "解析报文错误")
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
		return model.NewBindingReturn("200002", "解析报文错误")
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
		return model.NewBindingReturn("200002", "解析报文错误")
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
		return model.NewBindingReturn("200002", "解析报文错误")
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
		return model.NewBindingReturn("200002", "解析报文错误")
	}
	b.MerId = merId

	// 验证请求报文格式
	//	ret = validateBillingSummary(&b)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	// mock return
	ret = model.NewBindingReturn("000000", "虚拟数据")
	return ret
}

// 交易对账明细
func billingDetailsHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BillingDetails)
	err := json.Unmarshal(data, b)
	if err != nil {
		return model.NewBindingReturn("200002", "解析报文错误")
	}
	b.MerId = merId

	// 验证请求报文格式
	//	ret = validateBillingDetails(&b)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	// mock return
	ret = model.NewBindingReturn("000000", "虚拟数据")
	return ret
}

// 查询订单状态
func orderEnquiryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.OrderEnquiry)
	err := json.Unmarshal(data, b)
	if err != nil {
		g.Error("解析报文错误 :%s", err)
		return model.NewBindingReturn("200002", "解析报文错误")
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
		g.Error("解析报文错误 :%s", err)
		return model.NewBindingReturn("200002", "解析报文错误")
	}
	b.MerId = merId

	ret = validateNoTrackPayment(b)
	if ret != nil {
		return ret
	}

	//  todo 无卡支付暂不开放；业务处理
	ret = model.NewBindingReturn("100030", "不支持此类交易")
	return ret
}
