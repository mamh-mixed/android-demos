package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/CardInfoLink/quickpay/entrance/bindingpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

const (
	testMerId = "012345678901234"
	testSign  = "0123456789"
)

func TestBindingPay(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId

	b := model.BindingCreate{
		MerId:     testMerId,
		BindingId: tools.Millisecond(),
		AcctName:  "张三",
		AcctNum:   "6222022003008481261",
		IdentType: "0",
		IdentNum:  "440583199111031012",
		PhoneNum:  "18205960039",
		AcctType:  "10",
		ValidDate: "0612",
		Cvv2:      "793",
		SendSmsId: "",
		SmsCode:   "",
		BankId:    "102",
	}

	doPost("POST", url, b, t)
}

func doPost(method, url string, m interface{}, t *testing.T) {
	j, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(j))
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

	if out.RespCode == "" {
		t.Error(out)
	}
}
