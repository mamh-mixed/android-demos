package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"net/http"
	"sort"
)

func registerHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
	userName := r.FormValue("username")
	password := r.FormValue("password")
	transtime := r.FormValue("transtime")
	sign := r.FormValue("sign")

	ret := User.register(userName, password, transtime, sign)

	w.Write(jsonMarshal(ret))
}

// loginHandle 登录
func loginHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
	userName := r.FormValue("username")
	password := r.FormValue("password")
	transtime := r.FormValue("transtime")
	sign := r.FormValue("sign")

	ret := User.login(userName, password, transtime, sign)

	w.Write(jsonMarshal(ret))
}

// reqActivateHandle 请求发送激活邮件
func reqActivateHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
	userName := r.FormValue("username")
	password := r.FormValue("password")
	transtime := r.FormValue("transtime")
	sign := r.FormValue("sign")

	ret := User.reqActivate(userName, password, transtime, sign)

	w.Write(jsonMarshal(ret))
}

// activateHandle 激活
func activateHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
	userName := r.FormValue("username")
	code := r.FormValue("code")

	ret := User.activate(userName, code)

	w.Write(jsonMarshal(ret))
}

// improveInfoHandle 补充清算信息
func improveInfoHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// getOrderHandle 获得单个订单信息
func getOrderHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// billHandle 获取账单信息
func billHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// getTotalHandle 获取某天总交易金额
func getTotalHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// getRefdHandle 获得某笔交易已退款金额
func getRefdHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// passwordHandle 密码修改
func passwordHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// promoteLimitHandle 提升限额
func promoteLimitHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// updateHandle 修改清算帐号信息
func updateHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

// getInfoHandle 获取清算帐号信息
func getInfoHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}
}

func checkSign(r *http.Request) bool {
	var keys []string
	for k, _ := range r.Form {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	var sign string
	for _, v := range keys {
		// sign不参与签名
		if v == "sign" {
			sign = v
			continue
		}
		value := r.FormValue(v)
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(v + "=" + value)
	}
	content := buf.String()
	valid := fmt.Sprintf("%x", sha1.Sum([]byte(content)))
	if sign != valid {
		log.Warnf("check sign error, expect %s ,get %s", valid, sign)
		return false
	}
	return true
}

func jsonMarshal(result *model.AppResult) []byte {
	data, err := json.Marshal(result)
	if err != nil {
		log.Error("json marshal error: %s", err)
		return []byte(model.JSON_ERROR)
	}
	log.Debugf("response message: %s", string(data))
	return data
}
