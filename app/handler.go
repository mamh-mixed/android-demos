package app

import (
	"encoding/json"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

func registerHandle(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	password := r.FormValue("password")
	transtime := r.FormValue("transtime")
	sign := r.FormValue("sign")

	ret := User.register(userName, password, transtime, sign)

	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("json marshal err,%s", err)
		w.Write([]byte(model.JSON_ERROR))
		log.Debugf("response message: %s", model.JSON_ERROR)
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// loginHandle 登录
func loginHandle(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	password := r.FormValue("password")
	transtime := r.FormValue("transtime")
	sign := r.FormValue("sign")

	ret := User.login(userName, password, transtime, sign)

	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("json marshal err,%s", err)
		w.Write([]byte(model.JSON_ERROR))
		log.Debugf("response message: %s", model.JSON_ERROR)
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// reqActivateHandle 请求发送激活邮件
func reqActivateHandle(w http.ResponseWriter, r *http.Request) {}

// activateHandle 激活
func activateHandle(w http.ResponseWriter, r *http.Request) {}

// improveInfoHandle 补充清算信息
func improveInfoHandle(w http.ResponseWriter, r *http.Request) {}

// getOrderHandle 获得单个订单信息
func getOrderHandle(w http.ResponseWriter, r *http.Request) {}

// billHandle 获取账单信息
func billHandle(w http.ResponseWriter, r *http.Request) {}

// getTotalHandle 获取某天总交易金额
func getTotalHandle(w http.ResponseWriter, r *http.Request) {}

// getRefdHandle 获得某笔交易已退款金额
func getRefdHandle(w http.ResponseWriter, r *http.Request) {}

// passwordHandle 密码修改
func passwordHandle(w http.ResponseWriter, r *http.Request) {}

// promoteLimitHandle 提升限额
func promoteLimitHandle(w http.ResponseWriter, r *http.Request) {}

// updateHandle 修改清算帐号信息
func updateHandle(w http.ResponseWriter, r *http.Request) {}

// getInfoHandle 获取清算帐号信息
func getInfoHandle(w http.ResponseWriter, r *http.Request) {}
