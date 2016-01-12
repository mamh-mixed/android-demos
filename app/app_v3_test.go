package app

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
)

type testV3Func func(w http.ResponseWriter, r *http.Request)

func postV3(values url.Values, f testV3Func) (result *model.AppResult, err error) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(w, r) }))

	// 签名
	values.Add("sign", fmt.Sprintf("%x", sha256.Sum256([]byte(signContent(values)+sha1Key))))

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

func TestCouponsHandler(t *testing.T) {
	values := url.Values{}
	values.Add("username", "453481716@qq.com")
	values.Add("password", "670b14728ad9902aecba32e22fa4f6bd")
	values.Add("transtime", time.Now().Format("20060102150405"))
	// values.Add("month", "201512")
	values.Add("index", "0")
	values.Add("clientId", "99911888")
	values.Add("size", "5")

	result, err := postV3(values, couponsHandler)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestTotalSummaryHandler(t *testing.T) {
	values := url.Values{}
	values.Add("username", "453481716@qq.com")
	values.Add("password", "670b14728ad9902aecba32e22fa4f6bd")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("day", "20160107")
	values.Add("reportType", "2")

	result, err := postV3(values, totalSummaryHandler)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}
