package applepay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

func doPost(method, url, body string, t *testing.T) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		t.Error("创建POST请求失败")
	}
	req.Header.Set("X-Sign", entrance.SignatureUseSha1([]byte(body), "0123456789")) // TODO

	w := httptest.NewRecorder()
	entrance.BindingPay(w, req)
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
