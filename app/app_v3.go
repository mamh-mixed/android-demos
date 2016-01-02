package app

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/qiniu"

	"github.com/omigo/log"
)

// checkSignSha256 APP V3 版本使用SHA256算法签名
func checkSignSha256(r *http.Request) bool {
	sign, content := r.FormValue("sign"), signContent(r.Form)
	log.Debugf("sign content: %s", content)

	valid := fmt.Sprintf("%x", sha256.Sum256([]byte(content+sha1Key)))
	if sign != valid {
		log.Warnf("check sign error, expect %s ,get %s", valid, sign)
		return false
	}

	return true
}

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

// qiniuTokenHandler 获取七牛的上传token
func qiniuTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("username is %s; password is %s", r.FormValue("username"), r.FormValue("password"))
	result := User.getQiniuToken(&reqParams{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
	})

	if result.State == "success" {
		result.UploadToken = qiniu.GetUploadtoken()
	}

	w.Write(jsonMarshal(result))
}
