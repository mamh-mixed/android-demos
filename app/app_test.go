package app

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
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
	values.Add("username", "379630414@qq.com")
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

func TestLoginHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630414@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))

	result, err := post(values, loginHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestReqActivateHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630414@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("transtime", time.Now().Format("20060102150405"))

	result, err := post(values, reqActivateHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestImproveInfoHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630414@qq.com")
	values.Add("password", "awdwadsdasdawdwad")
	values.Add("bank_open", "招商银行")
	values.Add("payee", "测试")
	values.Add("payee_card", "123456789")
	values.Add("phone_num", "13148143570")
	values.Add("transtime", time.Now().Format("20060102150405"))

	result, err := post(values, improveInfoHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestMaxMerIdHandle(t *testing.T) {
	maxMerId, err := mongo.MerchantColl.FindMaxMerId()
	if err != nil {
		if err.Error() == "not found" {
			t.Logf(" set mix merId is 999118880000001")
			maxMerId = "999118880000001"
		} else {
			return
		}
	}
	maxMerIdNum, err := strconv.Atoi(maxMerId)
	if err != nil {
		log.Errorf("format maxMerId(%s) err", maxMerId)
		return
	}
	maxMerId = fmt.Sprintf("%d", maxMerIdNum+1)
	t.Logf("%s", maxMerId)
}
