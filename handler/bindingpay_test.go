package handler

import (
	"bytes"
	"encoding/json"
	"github.com/omigo/g"
	"net/http"
	"net/http/httptest"
	"quickpay/domain"
	"quickpay/model"
	"testing"
)

func bindingCreateRequestHandle(method, url, body string, t *testing.T) (response domain.BindingCreateResponse) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		g.Fatal("", err)
		panic(err)
	}

	// sign := signature(merId, []byte(body))
	// req.Header.Set("X-Signature", sign)

	w := httptest.NewRecorder()
	Quickpay(w, req)

	g.Info("%d - %s", w.Code, w.Body.String())

	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Error("Unmarshal response error")
	}
	return response
}

func TestBindingCreateHandle(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"bindingId":"1000000000001","acctName":"张三","acctNum":"6210948000000219","identType":"0","identNum":"36050219880401","phoneNum":"15600009909","acctType":"20","validDate":"1903","cvv2":"232","sendSmsId":"1000000000009","smsCode":"12353"}`

	response := bindingCreateRequestHandle("POST", url, body, t)
	g.Debug("%+v", response)
}

func TestBindingCreateHandleWhenIdentTypeWrong(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"bindingId":"1000000000001","acctName":"张三","acctNum":"6210948000000219","identType":"12","identNum":"36050219880401","phoneNum":"15600009909","acctType":"20","validDate":"1903","cvv2":"232","sendSmsId":"1000000000009","smsCode":"12353"}`

	response := bindingCreateRequestHandle("POST", url, body, t)
	g.Debug("%+v", response)
	if response.RespCode != "200120" {
		t.Error("验证 '证件类型有误' 失败")
	} else {
		t.Logf("%+v", response)
	}
}
func TestBindingCreateHandleWhenPhoneNumWrong(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"bindingId":"1000000000001","acctName":"张三","acctNum":"6210948000000219","identType":"0","identNum":"36050219880401","phoneNum":"059586832309","acctType":"20","validDate":"1903","cvv2":"232","sendSmsId":"1000000000009","smsCode":"12353"}`

	response := bindingCreateRequestHandle("POST", url, body, t)
	g.Debug("%+v", response)
	if response.RespCode != "200130" {
		t.Error("验证 '手机好嘛格式错误' 失败")
	} else {
		t.Logf("%+v", response)
	}
}

func TestBindingCreateHandleWhenAcctTypeIs10(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"bindingId":"1000000000001","acctName":"张三","acctNum":"6210948000000219","identType":"1","identNum":"36050219880401","phoneNum":"15600009909","acctType":"10","validDate":"","cvv2":"","sendSmsId":"1000000000009","smsCode":"12353"}`

	response := bindingCreateRequestHandle("POST", url, body, t)
	g.Debug("%+v", response)
	if response.RespCode != "000000" {
		t.Error("验证 '借记卡' 失败")
	} else {
		t.Logf("%+v", response)
	}
}

func TestBindingRemoveHandle(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + merId
	body := `{"bindingId": "1000000001"}`
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		t.Error("创建POST请求失败")
	} else {
		w := httptest.NewRecorder()
		Quickpay(w, req)
		g.Info("%d - %s", w.Code, w.Body.String())
		if w.Code != 200 {
			t.Errorf("response error with status %d", w.Code)
		}
		var out model.BindingRemoveOut
		err = json.Unmarshal([]byte(w.Body.String()), &out)
		if err != nil {
			t.Error("Unmarshal response error")
		} else {
			if out.RespCode != "000000" {
				t.Error("测试失败")
			} else {
				t.Log("测试OK")
			}
		}
	}
}
