package security

import (
	"encoding/base64"
	"testing"
)

var b64cipher = "hbW0++NBuvYCXBvk/pZW//F1YVin1ZvLNozWKORC6jz0WXbTSyW9ct52owIaBM5lpC+9b0GnZnOPojUKfKM3ZXLrR5bVzF03nElUQYd3Vb2xk7OCxzj1AqvrxhevetzdhA7lU2DbbldhJPbaDwQR0ryFdINOT0Sr+J6lp5O/UMQ="
var origData = "shangxuejin@gmail.com"

func TestRSAEncrypt(t *testing.T) {
	ciphertext, err := RSAEncrypt([]byte(origData), publicKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 相同的源文，每次加密后的密文都不相同，不能根据加密后的数据是否相等来判断加密算法是否正确
	// 要把加密后得数据再解码，如果和源文一致，说明加密算法正确
	actual, err := RSADecrypt(ciphertext, privateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(actual) != origData {
		t.Errorf("decrypt failed:  expected=%s, actual=%s", origData, actual)
	}
}

func TestRSADecrypt(t *testing.T) {
	ciphertext, err := base64.StdEncoding.DecodeString(b64cipher)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	actual, err := RSADecrypt(ciphertext, privateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(actual) != origData {
		t.Errorf("decrypt failed: ciphertext=%s, expected=%s, actual=%s", ciphertext, origData, actual)
	}
}

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)
