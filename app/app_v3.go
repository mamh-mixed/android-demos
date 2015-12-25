package app

import (
	"github.com/CardInfoLink/quickpay/model"
	"net/http"
)

// billV3Handle 获取账单信息
func billV3Handle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getUserBill(&reqParams{
		UserName:    r.FormValue("username"),
		Password:    r.FormValue("password"),
		Month:       r.FormValue("month"),
		Date:        r.FormValue("day"),
		Status:      r.FormValue("status"),
		Transtime:   r.FormValue("transtime"),
		Index:       r.FormValue("index"),
		OrderDetail: r.FormValue("order_detail"),
		Size:        r.FormValue("size"),
		TransType:   model.PayTrans,
	})

	w.Write(jsonMarshal(result))
}
