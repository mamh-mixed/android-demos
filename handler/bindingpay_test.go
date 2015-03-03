package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omigo/g"
)

func TestBindingCreateHandle(t *testing.T) {
	merId := "10000001"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	body := `{"mer_id":"10000001","mer_binding_id":"1000000000001","acct_name":"张三","acct_num":"6210948000000219","card_brand":"CUP","ident_type":"0","ident_num":"36050219880401","phone_num":"15600009909","acct_type":"20","valid_date":"1903","cvv2":"232","send_sms_id":"1000000000009","sms_code":"12353","sign":"FBE7381C0D4FDA71"}`

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		g.Fatal("", err)
	}

	// sign := signature(merId, []byte(body))
	// req.Header.Set("X-Signature", sign)

	w := httptest.NewRecorder()
	Quickpay(w, req)

	g.Info("%d - %s", w.Code, w.Body.String())

	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}
}
