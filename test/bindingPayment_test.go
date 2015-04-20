package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardInfoLink/quickpay/conf"
	. "github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

const (
	testMerId = "012345678901234"
	testSign  = "0123456789"
)

var (
	bindingId string
	orderNum  string
)

func init() {
	conf.Initialize()
}

func doPost(url string, m interface{}, t *testing.T) {
	j, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Sign", SignatureUseSha1(j, testSign))

	w := httptest.NewRecorder()
	BindingPay(w, req)
	log.Infof("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}

	var out model.BindingReturn
	err = json.Unmarshal(w.Body.Bytes(), &out)
	if err != nil {
		t.Errorf("Unmarshal response error (%s)", err)
	}

	if out.RespCode != "000000" {
		t.Error(out)
	}
}

func TestBindingCreate(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId

	bindingId = tools.Millisecond()

	b := model.BindingCreate{
		MerId:     testMerId,
		BindingId: bindingId,
		AcctName:  "测试账号",
		AcctNum:   "6222022003008481261",
		IdentType: "0",
		IdentNum:  "440583199111031012",
		PhoneNum:  "18205960039",
		AcctType:  "20",
		ValidDate: "0612",
		Cvv2:      "793",
		SendSmsId: "",
		SmsCode:   "",
		BankId:    "102",
	}

	var aes = tools.AesCBCMode{}

	b.AcctName = aes.Encrypt(b.AcctName)
	b.AcctNum = aes.Encrypt(b.AcctNum)
	b.IdentNum = aes.Encrypt(b.IdentNum)
	b.PhoneNum = aes.Encrypt(b.PhoneNum)
	b.ValidDate = aes.Encrypt(b.ValidDate)
	b.Cvv2 = aes.Encrypt(b.Cvv2)
	if aes.Err != nil {
		panic(aes.Err)
	}

	doPost(url, b, t)
}

func TestBindingEnquiryHandle(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + testMerId
	b := model.BindingEnquiry{BindingId: bindingId}
	doPost(url, b, t)
}

func TestBindingPaymentHandle(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + testMerId
	orderNum = tools.Millisecond()
	b := model.BindingPayment{
		MerOrderNum: orderNum,
		TransAmt:    1000,
		BindingId:   bindingId,
		MerId:       testMerId,
	}
	doPost(url, b, t)
}

func TestOrderEnquiry(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/orderEnquiry?merId=" + testMerId
	b := model.OrderEnquiry{
		OrigOrderNum: orderNum,
	}
	doPost(url, b, t)
}

func TestBindingRefundHandle(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/refund?merId=" + testMerId
	b := model.BindingRefund{
		OrigOrderNum: orderNum,
		MerOrderNum:  tools.Millisecond(),
		TransAmt:     1000,
	}
	doPost(url, b, t)
}

func TestNoTrackPaymentHandle(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/noTrackPayment?merId=" + testMerId

	b := model.NoTrackPayment{
		MerId:       testMerId,
		MerOrderNum: tools.Millisecond(),
		AcctName:    "测试账号",
		AcctNum:     "6222022003008481261",
		IdentType:   "0",
		IdentNum:    "440583199111031012",
		PhoneNum:    "18205960039",
		AcctType:    "20",
		ValidDate:   "0612",
		Cvv2:        "793",
		TransAmt:    1000,
	}

	var aes = tools.AesCBCMode{}

	b.AcctName = aes.Encrypt(b.AcctName)
	b.AcctNum = aes.Encrypt(b.AcctNum)
	b.IdentNum = aes.Encrypt(b.IdentNum)
	b.PhoneNum = aes.Encrypt(b.PhoneNum)
	b.ValidDate = aes.Encrypt(b.ValidDate)
	b.Cvv2 = aes.Encrypt(b.Cvv2)
	if aes.Err != nil {
		panic(aes.Err)
	}
	doPost(url, b, t)
}

func TestBindingRemoveHandle(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + testMerId
	b := model.BindingEnquiry{BindingId: bindingId}
	doPost(url, b, t)
}
