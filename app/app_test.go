package app

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/log"
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
func TestActivateHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "379630414@qq.com")
	values.Add("code", "awdwadsdasdawdwa")

	result, err := post(values, activateHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}
func TestBillHandle(t *testing.T) {
	values := url.Values{}
	values.Add("username", "330961193@qq.com")
	values.Add("password", "670b14728ad9902aecba32e22fa4f6bd")
	values.Add("transtime", time.Now().Format("20060102150405"))
	values.Add("month", "201509")
	values.Add("day", "30")
	values.Add("status", "all")
	values.Add("index", "1")
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

func TestMaxMerIdHandle(t *testing.T) {
	prefix := "9902041"
	var length = fmt.Sprintf("%d", 15-len(prefix))
	maxMerId, err := mongo.MerchantColl.FindMaxMerId(prefix)
	if err != nil {
		if err.Error() == "not found" {
			t.Logf(" set mix merId is 999118880000001")
			maxMerId = prefix + fmt.Sprintf("%0"+length+"d", 1)
		} else {
			t.Log(err)
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

var token = "2902f8a9caee45a46fb2a8408ec4f401"

func TestAppToolsLogin(t *testing.T) {
	values := url.Values{}
	values.Add("username", "toolstest")
	values.Add("password", "Yun#1016")
	result, err := post(values, CompanyLogin)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestAppToolsUserList(t *testing.T) {
	values := url.Values{}
	values.Add("accessToken", token)
	result, err := post(values, UserList)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestAppToolsRegister(t *testing.T) {
	values := url.Values{}
	values.Add("accessToken", token)
	values.Add("username", "379630413@qq.com")
	values.Add("password", "12345678")

	result, err := post(values, UserRegister)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestAppToolsUpdateUserInfo(t *testing.T) {
	values := url.Values{}
	values.Add("accessToken", token)
	values.Add("username", "379630413@qq.com")
	values.Add("bank_open", "中国工商")
	values.Add("payee", "陈芝锐")
	values.Add("payee_card", "6222022003008481261")
	values.Add("phone_num", "15618103236")
	values.Add("province", "广东省")
	values.Add("city", "汕头市")
	values.Add("branch_bank", "中国工商")
	values.Add("bankNo", "123312312313|123213131323")
	values.Add("merName", "汕头牛肉丸1")
	result, err := post(values, UpdateUserInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestActivateUser(t *testing.T) {
	values := url.Values{}
	values.Add("accessToken", token)
	values.Add("username", "379630413@qq.com")

	result, err := post(values, UserActivate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(30 * time.Second)

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestGetDownloadUrl(t *testing.T) {
	values := url.Values{}
	values.Add("accessToken", token)
	values.Add("merId", "199005050000019")
	values.Add("imageType", "pay")

	result, err := post(values, GetDownloadUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}

func TestDownload(t *testing.T) {

	dlUrl := qiniu.MakePrivateUrl(fmt.Sprintf(qrImage, "199005050000019", "pay"))
	resp, err := http.Get(dlUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	jpg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	f, _ := os.Create("/Users/zhiruichen/Desktop/d.jpg")
	f.Write(jpg)
	f.Close()
}

func TestSendEmail(t *testing.T) {
	NotifySalesman()
}

func TestRandBytes(t *testing.T) {
	bs := randBytes(32)
	t.Logf("%x", bs)
}

func TestTickHandle(t *testing.T) {
	values := url.Values{}
	values.Add("receiptnum", "199005050000019")
	values.Add("ordernum", "15111716332617542")
	values.Add("username", "cherripe.Chen@cardinfolink.com")
	values.Add("password", "96e79218965eb72c92a549dd5a330112")

	result, err := post(values, ticketHandle)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bs, _ := json.Marshal(result)
	t.Logf("%s", string(bs))
}
