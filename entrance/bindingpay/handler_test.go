package bindingpay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"

	"github.com/omigo/g"
)

func post(req *http.Request, t *testing.T) {
	w := httptest.NewRecorder()
	BindingPay(w, req)
	g.Info("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}

	var out model.BindingReturn
	err := json.Unmarshal(w.Body.Bytes(), &out)
	if err != nil {
		t.Errorf("Unmarshal response error (%s)", err)
	}

	if out.RespCode == "" {
		t.Error("测试失败")
	}
}

func TestBindingCreateWithSignHandle(t *testing.T) {
	b := model.BindingCreate{
		MerId:         "99001405",            // 商户ID
		BindingId:     tools.Millisecond(),   // 银行卡绑定ID
		AcctName:      "张三",                  // 账户名称
		AcctNum:       "6222020302062061908", // 账户号码
		IdentType:     "0",                   // 证件类型
		IdentNum:      "350583199009153732",  // 证件号码
		PhoneNum:      "18205960039",         // 手机号
		AcctType:      "20",                  // 账户类型
		ValidDate:     "1903",                // 信用卡有限期
		Cvv2:          "232",                 // CVV2
		SendSmsId:     "",                    // 发送短信验证码的交易流水
		SmsCode:       "",                    // 短信验证码
		BankId:        "102",                 // 银行ID
		ChanBindingId: "",                    // 渠道绑定ID
		ChanMerId:     "",                    // 渠道商户ID
	}
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + b.MerId

	body, _ := json.Marshal(b)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Error("创建POST请求失败")
	}
	sign := SignatureUseSha1(body, "0123456789")
	req.Header.Set("X-Sign", sign)

	post(req, t)
}

func TestBindingCreateHandle(t *testing.T) {
	// todo 生成一个随机的商户号
	rdMerId := tools.Millisecond()
	// todo 在路由策略里面插入新商户号码的路由策略
	rp := &model.RouterPolicy{
		MerId:     rdMerId,
		CardBrand: "CUP",
		ChanCode:  "CFCA",
		ChanMerId: "001405",
	}

	if err := mongo.RouterPolicyColl.Insert(rp); err != nil {
		t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
	}
	// 生成一个随机的绑定ID
	rdBindingId := tools.Millisecond()
	merId := rdMerId
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"bindingId":"` + rdBindingId + `","acctName":"张三","acctNum":"6222022003008481261","identType":"0","identNum":"440583199111031012","phoneNum":"18205960039","acctType":"10","BankId":"102","sendSmsId":"1000000000009","smsCode":"12353"}`

	doPost("POST", url, body, t)
}

func doPost(method, url, body string, t *testing.T) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		t.Error("创建POST请求失败")
	}
	req.Header.Set("X-Sign", SignatureUseSha1([]byte(body), "0123456789")) // TODO

	w := httptest.NewRecorder()
	BindingPay(w, req)
	g.Info("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}

	var out model.BindingReturn
	err = json.Unmarshal(w.Body.Bytes(), &out)
	if err != nil {
		t.Errorf("Unmarshal response error (%s)", err)
	}

	if out.RespCode == "" {
		t.Error("测试失败")
	}
}

func TestBindingRemoveHandle(t *testing.T) {
	merId := "1426583281344"
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + merId
	body := `{"bindingId": "1426583281402"}`
	doPost("POST", url, body, t)
}

func TestBindingEnquiryHandle(t *testing.T) {
	merId := "99001405"
	url := "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + merId
	body := `{"bindingId": "1426689409857"}`
	doPost("POST", url, body, t)
}

func TestBindingPaymentHandle(t *testing.T) {
	merId := "1426818873711"
	merOrderNum := tools.Millisecond()
	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + merId
	body := `{
		"subMerId": "",
		"merOrderNum": "` + merOrderNum + `",
		"transAmt": 2000,
		"bindingId": "1426818873771",
		"sendSmsId": "",
		"smsCode": "",
		"merId":"1426818873711"
	}`
	doPost("POST", url, body, t)
}

func TestRefundHandle(t *testing.T) {
	merId := "1426818873711"
	merOrderNum := tools.Millisecond()
	url := "https://api.xxxx.com/quickpay/refund?merId=" + merId
	body := `{
		"merOrderNum": "` + merOrderNum + `",
		"transAmt": 2000,
		"origOrderNum": "1426830050043"
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

func TestOrderEnquiry(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/orderEnquiry?merId=1426818873711"
	body := `{
		"origOrderNum":"1426830241979",
		"merId":"1426818873711",
		"showOrigInfo":"1"
		}`
	//"showOrigInfo":1,
	doPost("POST", url, body, t)
}
