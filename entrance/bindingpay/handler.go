package bindingpay

import (
	"encoding/json"
	"errors"
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
		w.Write([]byte("only post method allowed"))
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
	// case "/quickpay/bindingCreate":
	// 	ret, err = bindingCreateHandle(data)
	// case "/quickpay/bindingRemove":
	// 	ret, err = bindingRemoveHandle(data)
	case "/quickpay/bindingEnquiry":
		ret = bindingEnquiryHandle(data)
	// case "/quickpay/bindingPayment":
	// 	ret, err = bindingPaymentHandle(data)
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
func bindingCreateHandle(data []byte) ([]byte, error) {
	// json to obj
	var (
		in  model.BindingCreate
		out = &model.BindingReturn{}
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
	validityCode, validityErr := bindingCreateRequestValidity(in)
	if validityErr == nil {
		// 验证参数OK

		// 业务处理
		// out2 := core.CreateBinding(&in)
		// out.RespCode = out2.RespCode
		// out.RespMsg = out2.RespMsg
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
		in  model.BindingRemove
		out model.BindingReturn
		err error
	)

	err = json.Unmarshal(data, &in)
	if err != nil {
		out.RespCode = "200020"
		out.RespMsg = "解析报文错误"
	} else {
		// 验证请求报文格式
		validityCode, validityErr := bindingRemoveRequestValidity(in)
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
func bindingEnquiryHandle(data []byte) (ret *model.BindingReturn) {
	var be model.BindingEnquiry
	err := json.Unmarshal(data, &be)
	if err != nil {
		ret = &model.BindingReturn{
			RespCode: "200020",
			RespMsg:  "解析报文错误",
		}
		return ret
	}

	// 验证请求报文格式
	validityCode, validityErr := bindingEnquiryRequestValidity(be)
	if validityErr != nil {
		ret = &model.BindingReturn{
			RespCode: validityCode,
			RespMsg:  validityErr.Error(),
		}
		return ret
	}

	ret = core.ProcessBindingEnquiry(&be)

	return ret
}

// 绑定支付关系
func bindingPaymentHandle(data []byte) ([]byte, error) {
	var (
		in  model.BindingPayment
		out model.BindingReturn
		err error
	)

	err = json.Unmarshal(data, &in)
	if err != nil {
		g.Error("Unmarshal request body error msg: ", err.Error())
		out.RespCode = "200020"
		out.RespMsg = "解析报文错误"
	} else {
		// 验证请求报文格式
		validityCode, validityErr := bindingPaymentRequestValidity(in)
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
