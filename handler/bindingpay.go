package handler

import (
	"encoding/json"
	"fmt"
	"github.com/omigo/g"
	"io/ioutil"
	"net/http"
	"quickpay/domain"
	"quickpay/validity"
)

// Quickpay 快捷支付入口
func Quickpay(w http.ResponseWriter, r *http.Request) {
	g.Debug("url = %s", r.URL.Path)

	data, err := ioutil.ReadAll(r.Body)
	g.Error("read body error: %s", err)
	// 验签，如果失败，立即返回
	// if checkSignature(data, merId)

	// 执行业务逻辑
	switch r.URL.Path {
	case "/quickpay/bindingCreate":
		bindingCreateHandle(w, r, data)

	default:
		w.WriteHeader(204)
	}
}

func bindingCreateHandle(w http.ResponseWriter, r *http.Request, data []byte) {

	// json to obj
	var request domain.BindingCreateRequest
	var response domain.BindingCreateResponse
	err := json.Unmarshal(data, &request)
	if err != nil {
		g.Error("unmashal data error \n", err)
		fmt.Fprint(w, "unmashal data error")
	} else {
		g.Debug("%+v", request)
		response.MerBindingId = request.MerBindingId
	}
	// todo 验证参数
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
