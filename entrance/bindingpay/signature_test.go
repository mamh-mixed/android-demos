package bindingpay

import (
	"strings"
	"testing"
)

func TestSignatureUseSha1(t *testing.T) {
	data, key := `{"merBindingId":"1000000001"}`, "C380BEC2BFD727A4B6845133519F3AD6"
	result, sign := "1886dfb54758972b38c7d045e967573d1eff093b", SignatureUseSha1(data, key)
	t.Logf("sign result: %s\n", sign)
	if strings.EqualFold(sign, result) {
		t.Log("match,successfully")
	} else {
		t.Error("not equal")
	}
	dataCn := `{"respCode":"00","respMsg":"你好，世界"}`
	t.Logf("中文签名: %s\n", SignatureUseSha1(dataCn, key))
}

func TestCheckSignatureUseSha1(t *testing.T) {
	data, key, sign := `{"merBindingId":"1000000001"}`, "C380BEC2BFD727A4B6845133519F3AD6", "1886dfb54758972b38c7d045e967573d1eff093b"
	if CheckSignatureUseSha1(data, key, sign) {
		t.Log("Successfully")
	} else {
		t.Error("Fail")
	}

	dataCn := `{"respCode":"00","respMsg":"你好，世界"}`
	if CheckSignatureUseSha1(dataCn, key, sign) {
		t.Error("Fail")
	} else {
		t.Log("Successfully")
	}
}
