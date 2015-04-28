package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	applePayMerId  = "APPTEST"  // apple pay 测试用商户号
	testTerminalId = "00000001" // 测试用渠道商户的终端号

	// Apple Pay测试数据
	testAPPCard       = "5457210001000019"
	testAPPExpireDate = "141231"
)

func applePayPost(method, url, body string, t *testing.T) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		t.Error("创建POST请求失败")
	}
	req.Header.Set("X-Sign", entrance.SignatureUseSha1([]byte(body), "0123456789")) // TODO

	w := httptest.NewRecorder()
	entrance.Quickpay(w, req)
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
		t.Error("测试失败")
	}
}

func TestApplePayHandle(t *testing.T) {
	subMerId := fmt.Sprintf("%05d", time.Now().UnixNano())
	merOrderNum := fmt.Sprintf("%012d", time.Now().UnixNano())
	transactionId := fmt.Sprintf("%020d", time.Now().UnixNano())
	url := "https://api.xxxx.com/quickpay/applePay?merId=" + applePayMerId
	b := `{
		"transType":"SALE",
		"subMerId":"` + subMerId + `",
		"terminalId":"` + testTerminalId + `",
		"merOrderNum":"` + merOrderNum + `",
		"transactionId":"` + transactionId + `",
		"applePayData": {
		    "applicationPrimaryAccountNumber": "` + testAPPCard + `",
		    "applicationExpirationDate": "` + testAPPExpireDate + `",
		    "currencyCode": "156",
		    "transactionAmount": 120,
		    "deviceManufacturerIdentifier": "040010030273",
		    "paymentDataType": "3DSecure",
		    "paymentData": {
		        "onlinePaymentCryptogram": "AOZSYAeX7VKTAAKv5hDuAoABFA==",
		        "eciIndicator": "5"
		    }
		}
	}`
	applePayPost("POST", url, b, t)
}
