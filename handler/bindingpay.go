package handler

import (
	"fmt"
	"net/http"

	"github.com/omigo/g"
)

// Quickpay 快捷支付入口
func Quickpay(w http.ResponseWriter, r *http.Request) {
	g.Debug("url = %s", r.URL.String())

	// data, err := ioutil.ReadAll(r.Body)
	// 验签，如果失败，立即返回
	// if checkSignature(data, merId)

	// 执行业务逻辑
	switch r.URL.String() {
	case "/quickpay/bindingCreate":
		// bindingCreateHandle(w, r, data)

	default:
		w.WriteHeader(404)
	}
}

func bindingCreateHandle(w http.ResponseWriter, r *http.Request, data []byte) {

	// json to obj
	// err := json.Unmarshal(data, &in)

	// out := core.CreateBinding(in)

	// obj to json
	// body, err := json.Marshell(ret)

	// 签名，并返回
	// sign = signature(out, in.merId)

	fmt.Fprint(w, "unimplement exception")
}
