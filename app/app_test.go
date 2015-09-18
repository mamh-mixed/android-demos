package app

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type testFunc func(w http.ResponseWriter, r *http.Request)

func post(values url.Values, f testFunc) (result *model.AppResult, err error) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(w, r) }))

	// 签名
	values.Add("sign", fmt.Sprintf("%x", sha1.Sum([]byte(signContent(values)+sha1Key))))

	// post
	resp, err := http.PostForm(ts.URL, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析结果
	result = new(model.AppResult)
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func TestRegisterHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))

	result, err := post(values, registerHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestGetOrderHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("orderNum", "1442537558990")
	result, err := post(values, getOrderHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestBillHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("month", "20150918")
	values.Add("status", "success")
	values.Add("index", "3")
	result, err := post(values, billHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestGetTotalHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("date", "20150918")
	result, err := post(values, getTotalHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestGetRefdHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("orderNum", "1442560835714")
	result, err := post(values, getRefdHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestPasswordHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("oldpassword", "awdwadsdasdawdwad")
	values.Add("newpassword", "dwainczncjawhduha")
	values.Add("transtime", time.Now().Format("20060102150405"))
	result, err := post(values, passwordHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestPromoteLimitHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("payee", "锐哥")
	values.Add("phone_num", "15618103236")
	values.Add("email", "379630413@qq.com")
	result, err := post(values, promoteLimitHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestUpdateSettInfoHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("bank_open", "中国工商银行")
	values.Add("phone_num", "15618103236")
	values.Add("payee", "陈芝锐")
	values.Add("payee_card", "6222022003008481261")
	result, err := post(values, updateSettInfoHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestGetSettInfoHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630413@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))
	result, err := post(values, getSettInfoHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}
