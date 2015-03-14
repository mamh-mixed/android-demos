package bindingpay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quickpay/model"
	"testing"

	"github.com/omigo/g"
)

func TestBindingCreateHandle(t *testing.T) {
	merId := "499999999"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"bindingId":"10000000001003","acctName":"张三","acctNum":"6222022003008481261","identType":"0","identNum":"440583199111031012","phoneNum":"15600009909","acctType":"20","validDate":"1903","cvv2":"232","bankId":"700","sendSmsId":"1000000000009","smsCode":"12353"}`
	doPost("POST", url, body, t)
}

// func TestBindingCreateHandleWhenAcctTypeIs10(t *testing.T) {
// 	merId := "10000001"
// 	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
// 	body := `{"bindingId":"1000000000001","acctName":"张三","acctNum":"6210948000000219","identType":"1","identNum":"36050219880401","phoneNum":"15600009909","acctType":"10","validDate":"","cvv2":"","sendSmsId":"1000000000009","smsCode":"12353"}`
// 	response := bindingCreateRequestHandle("POST", url, body, t)
// 	g.Debug("%+v", response)
// 	if response.RespCode != "000000" {
// 		t.Error("验证 '借记卡' 失败")
// 	} else {
// 		t.Logf("%+v", response)
// 	}
// }

func doPost(method, url, body string, t *testing.T) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		t.Error("创建POST请求失败")
	}

	w := httptest.NewRecorder()
	BindingPay(w, req)
	g.Info("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}

	var out model.BindingReturn
	err = json.Unmarshal([]byte(w.Body.String()), &out)
	if err != nil {
		t.Error("Unmarshal response error")
	}

	if out.RespCode == "" {
		t.Error("测试失败")
	}
}

func TestBindingRemoveHandle(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + merId
	body := `{"bindingId": "1000000001"}`
	doPost("POST", url, body, t)
}

func TestBindingEnquiryHandle(t *testing.T) {
	merId := "001405"
	url := "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + merId
	body := `{"bindingId": "1000000000001"}`
	doPost("POST", url, body, t)
}

func TestBindingPaymentHandle(t *testing.T) {
	merId := "499999999"
	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + merId
	body := `{
		"subMerId": "",
		"merOrderNum": "100000000300091",
		"transAmt": 900,
		"bindingId": "10000000001003",
		"sendSmsId": "",
		"smsCode": "",
		"merId":"499999999"
	}`
	doPost("POST", url, body, t)
}

func TestRefundHandle(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/refund?merId=" + merId
	body := `{
		"merOrderNum": "1000000004",
		"transAmt": -10000,
		"origOrderNum": "1000000003"
	}`
	doPost("POST", url, body, t)
}

func TestNoTrackPaymentHandle(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/noTrackPayment?merId=" + merId
	body := `{
		"subMerId": "",
		"merOrderNum": "1000000003",
		"transAmt": 1,
		"acctName":"张三",
		"acctNum":"6210948000000219",
		"identType":"0",
		"identNum":"36050219880401",
		"phoneNum":"15600009909",
		"acctType":"20",
		"validDate":"1903",
		"cvv2":"232",
		"sendSmsId": "",
		"smsCode": ""
	}`
	doPost("POST", url, body, t)
}
