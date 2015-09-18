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
