package bindingpay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"quickpay/core"
	"quickpay/model"

	"github.com/omigo/g"
)

// BindingPay 绑定支付入口
func BindingPay(w http.ResponseWriter, r *http.Request) {
	g.Debug("url = %s", r.URL.Path)

	if r.Method != "POST" {
		g.Error("methond(%s) not allowed", r.Method)
		w.WriteHeader(405)
		w.Write([]byte("only 'POST' method allowed"))
		return
	}

	// merId 可以放到 json 里
	v := r.URL.Query()
	merId := v.Get("merId")
	if merId == "" {
		w.WriteHeader(412)
		w.Write([]byte("parameter merId must required"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		g.Error("read body error: ", err)
		w.WriteHeader(406)
		w.Write([]byte("can not read request body"))
		return
	}

	data := body
	g.Debug("商户报文: %s", data)

	// 验签，如果失败，立即返回
	// if checkSignature(data, merId)

	// 执行业务逻辑
	var ret *model.BindingReturn
	switch r.URL.Path {
	case "/quickpay/bindingCreate":
		ret = bindingCreateHandle(data)
	case "/quickpay/bindingRemove":
		ret = bindingRemoveHandle(data)
	case "/quickpay/bindingEnquiry":
		ret = bindingEnquiryHandle(data)
	case "/quickpay/bindingPayment":
		ret = bindingPaymentHandle(data)
	case "/quickpay/refund":
		ret = refundHandle(data)
	default:
		w.WriteHeader(404)
	}

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	// todo 签名，并返回
	// sign = signature(out, merId)
	_ = merId

	rbody := rdata
	w.Write(rbody)
}

// 建立绑定关系
func bindingCreateHandle(data []byte) (ret *model.BindingReturn) {
	var bc model.BindingCreate
	err := json.Unmarshal(data, &bc)
	if ret = checkUnmarshalError(err); ret != nil {
		return ret
	}

	// 验证请求报文是否完整，格式是否正确
	ret = bindingCreateRequestValidity(bc)
	if ret != nil {
		return ret
	}

	//todo 业务处理
	// mock return
	ret = &model.BindingReturn{
		RespCode: "000000",
		RespMsg:  "虚拟数据",
	}
	return ret
}

// 解除绑定关系
func bindingRemoveHandle(data []byte) (ret *model.BindingReturn) {
	var br model.BindingRemove
	err := json.Unmarshal(data, &br)
	if ret = checkUnmarshalError(err); ret != nil {
		return ret
	}

	ret = bindingRemoveRequestValidity(br)
	if ret != nil {
		return ret
	}
	// todo 业务处理
	// mock return
	ret = &model.BindingReturn{
		RespCode: "000000",
		RespMsg:  "虚拟数据",
	}
	return ret
}

// 查询绑定关系
func bindingEnquiryHandle(data []byte) (ret *model.BindingReturn) {
	var be model.BindingEnquiry
	err := json.Unmarshal(data, &be)
	if ret = checkUnmarshalError(err); ret != nil {
		return ret
	}

	// 验证请求报文格式
	ret = bindingEnquiryRequestValidity(be)
	if ret != nil {
		return ret
	}

	ret = core.ProcessBindingEnquiry(&be)

	return ret
}

// 绑定支付关系
func bindingPaymentHandle(data []byte) (ret *model.BindingReturn) {
	var in model.BindingPayment
	err := json.Unmarshal(data, &in)
	if ret = checkUnmarshalError(err); ret != nil {
		return ret
	}

	// 验证请求报文格式
	ret = bindingPaymentRequestValidity(in)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	// mock return
	ret = &model.BindingReturn{
		RespCode: "000000",
		RespMsg:  "虚拟数据",
	}
	return ret
}

// 退款处理
func refundHandle(data []byte) (ret *model.BindingReturn) {
	var in model.BindingRefund

	err := json.Unmarshal(data, &in)
	if ret = checkUnmarshalError(err); ret != nil {
		return ret
	}

	// 验证请求报文格式
	ret = bindingRefundRequestValidity(&in)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	// mock return
	ret = &model.BindingReturn{
		RespCode: "000000",
		RespMsg:  "虚拟数据",
	}
	return ret
}

func checkUnmarshalError(err error) (ret *model.BindingReturn) {
	if err != nil {
		ret = &model.BindingReturn{
			RespCode: "200020",
			RespMsg:  "解析报文错误",
		}
		return ret
	}
	return nil
}
