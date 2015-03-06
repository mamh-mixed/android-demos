package handler

import (
	"encoding/json"
	"fmt"
	"github.com/omigo/g"
	"io/ioutil"
	"net/http"
	"quickpay/domain"
	"quickpay/model"
	"quickpay/validity"
)

// Quickpay 快捷支付入口
func Quickpay(w http.ResponseWriter, r *http.Request) {
	g.Debug("url = %s", r.URL.Path)

	data, err := ioutil.ReadAll(r.Body)
	g.Error("read body error: ", err)
	// 验签，如果失败，立即返回
	// if checkSignature(data, merId)

	// 执行业务逻辑
	switch r.URL.Path {
	case "/quickpay/bindingCreate":
		bindingCreateHandle(w, r, data)
	case "/quickpay/bindingRemove":
		bindingRemoveHandle(w, r, data)
	case "/quickpay/bindingEnquiry":
		bindingEnquiryHandle(w, r, data)
	case "/quickpay/bindingPayment":
		bindingPaymentHandle(w, r, data)
	default:
		w.WriteHeader(204)
	}
}

// 建立绑定关系
func bindingCreateHandle(w http.ResponseWriter, r *http.Request, data []byte) {

	// json to obj
	var request domain.BindingCreateRequest
	var response domain.BindingCreateResponse
	err := json.Unmarshal(data, &request)
	if err != nil {
		g.Error("unmashal data error \n", err)
		response.RespCode = "200020"
		response.RespMsg = "解析报文错误"
	} else {
		g.Debug("%+v", request)
		response.BindingId = request.BindingId
	}
	// 验证请求报文是否完整，格式是否正确
	validityCode, validityErr := validity.BindingCreateRequestValidity(request)
	if validityErr == nil {
		// 验证参数OK

		// 业务处理
		// out := core.CreateBinding(in)

		// 虚拟数据，假设成功
		response.RespCode = "000000"
		response.RespMsg = "Success"
	} else {
		// 验证参数失败
		response.RespCode = validityCode
		response.RespMsg = validityErr.Error()
	}

	// 签名，并返回
	// sign = signature(out, in.merId)

	// obj to json
	body, err := json.Marshal(response)
	if err != nil {
		fmt.Fprint(w, "mashal data error")
	} else {
		fmt.Fprintf(w, "%s", body)
	}
}

// 解除绑定关系
func bindingRemoveHandle(w http.ResponseWriter, r *http.Request, data []byte) {
	var (
		in  model.BindingRemoveIn
		out model.BindingRemoveOut
		err error
	)

	err = json.Unmarshal(data, &in)
	if err != nil {
		out.RespCode = "200020"
		out.RespMsg = "解析报文错误"
	} else {
		// 验证请求报文格式
		validityCode, validityErr := validity.BindingRemoveRequestValidity(in)
		if validityErr != nil {
			out.RespCode = validityCode
			out.RespMsg = validityErr.Error()
		} else {
			// todo 业务处理，这里先返回OK响应码
			out.RespCode = "000000"
			out.RespMsg = "Success"
		}
	}
	//  todo 签名并返回
	// obj to json
	body, err := json.Marshal(out)
	if err != nil {
		fmt.Fprint(w, "mashal data error")
	} else {
		fmt.Fprintf(w, "%s", body)
	}
}

// 查询绑定关系
func bindingEnquiryHandle(w http.ResponseWriter, r *http.Request, data []byte) {
	var (
		in  model.BindingEnquiryIn
		out model.BindingEnquiryOut
		err error
	)

	err = json.Unmarshal(data, &in)
	if err != nil {
		out.RespCode = "200020"
		out.RespMsg = "解析报文错误"
	} else {
		// 验证请求报文格式
		validityCode, validityErr := validity.BindingEnquiryRequestValidity(in)
		if validityErr != nil {
			out.RespCode = validityCode
			out.RespMsg = validityErr.Error()
		} else {
			// todo 业务处理，这里先返回OK响应码
			out.RespCode = "000000"
			out.RespMsg = "Success"
		}
	}
	//  todo 签名并返回
	// obj to json
	body, err := json.Marshal(out)
	if err != nil {
		fmt.Fprint(w, "mashal data error")
	} else {
		fmt.Fprintf(w, "%s", body)
	}
}

// 绑定支付关系
func bindingPaymentHandle(w http.ResponseWriter, r *http.Request, data []byte) {
	var (
		in  model.BindingPaymentIn
		out model.BindingPaymentOut
		err error
	)

	err = json.Unmarshal(data, &in)
	if err != nil {
		g.Error("Unmarshal request body error msg: %s", err.Error())
		out.RespCode = "200020"
		out.RespMsg = "解析报文错误"
	} else {
		// 验证请求报文格式
		validityCode, validityErr := validity.BindingPaymentRequestValidity(in)
		if validityErr != nil {
			out.RespCode = validityCode
			out.RespMsg = validityErr.Error()
		} else {
			// todo 业务处理，这里先返回OK响应码
			out.RespCode = "000000"
			out.RespMsg = "Success"
		}
	}
	//  todo 签名并返回
	// obj to json
	body, err := json.Marshal(out)
	if err != nil {
		fmt.Fprint(w, "mashal data error")
	} else {
		fmt.Fprintf(w, "%s", body)
	}
}
