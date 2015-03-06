package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"quickpay/core"
	"quickpay/model"
	"quickpay/validity"

	"github.com/omigo/g"
)

// Quickpay 快捷支付入口
func Quickpay(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte // 读取request请求的数据
		out  []byte // 业务处理结束后返回的json字符串
		err  error  //错误信息
	)
	g.Debug("url = %s", r.URL.Path)

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		g.Error("read body error: ", err)
	}

	g.Debug("请求报文的内容： %s", data)
	// 验签，如果失败，立即返回
	// if checkSignature(data, merId)

	// 执行业务逻辑
	switch r.URL.Path {
	case "/quickpay/bindingCreate":
		out, err = bindingCreateHandle(data)
	case "/quickpay/bindingRemove":
		out, err = bindingRemoveHandle(data)
	case "/quickpay/bindingEnquiry":
		out, err = bindingEnquiryHandle(data)
	case "/quickpay/bindingPayment":
		out, err = bindingPaymentHandle(data)
	default:
		w.WriteHeader(204)
	}
	// todo 签名，并返回
	// sign = signature(out, in.merId)
	if err != nil {
		fmt.Fprint(w, "mashal data error")
	} else {
		g.Info("响应的报文: %s", out)
		fmt.Fprintf(w, "%s", out)
	}
}

// 建立绑定关系
func bindingCreateHandle(data []byte) ([]byte, error) {
	// json to obj
	var (
		in  model.BindingCreateIn
		out = &model.BindingCreateOut{}
	)
	err := json.Unmarshal(data, &in)
	if err != nil {
		g.Error("unmashal data error \n", err)
		out.RespCode = "200020"
		out.RespMsg = "解析报文错误"
	} else {
		g.Debug("%+v", in)
		out.BindingId = in.BindingId
	}
	// 验证请求报文是否完整，格式是否正确
	validityCode, validityErr := validity.BindingCreateRequestValidity(in)
	if validityErr == nil {
		// 验证参数OK

		// 业务处理
		out2 := core.CreateBinding(&in)
		out.RespCode = out2.RespCode
		out.RespMsg = out2.RespMsg
	} else {
		// 验证参数失败
		out.RespCode = validityCode
		out.RespMsg = validityErr.Error()
	}

	// obj to json
	body, err := json.Marshal(out)
	if err != nil {
		return nil, errors.New("mashal data error")
	} else {
		return body, nil
	}
}

// 解除绑定关系
func bindingRemoveHandle(data []byte) ([]byte, error) {
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
	// obj to json
	body, err := json.Marshal(out)
	if err != nil {
		return nil, errors.New("mashal data error")
	} else {
		return body, nil
	}
}

// 查询绑定关系
func bindingEnquiryHandle(data []byte) ([]byte, error) {
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
		return nil, errors.New("mashal data error")
	} else {
		return body, nil
	}
}

// 绑定支付关系
func bindingPaymentHandle(data []byte) ([]byte, error) {
	var (
		in  model.BindingPaymentIn
		out model.BindingPaymentOut
		err error
	)

	err = json.Unmarshal(data, &in)
	if err != nil {
		g.Error("Unmarshal request body error msg: ", err.Error())
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
		return nil, errors.New("mashal data error")
	} else {
		return body, nil
	}
}
