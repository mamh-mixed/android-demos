package cfca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	"github.com/omigo/log"
)

const (
	// 	priKeyPem = `-----BEGIN RSA PRIVATE KEY-----
	// MIICXQIBAAKBgQCvJC9MMGRKmxRBI0KMjDtz2KooIc6XOljHPWhTfAamhV3A5v5y
	// PiZr4haMDpulU08Y0JxsegwDwfbscQrhG7nvilIqIa+HiI1xkfFxjtNUrMN5hpvO
	// 8HUUfwqzb5EdllQcv/C0xxBkeCECIb86JJry7ty4mNBkN2idbGxldMi90QIDAQAB
	// AoGATvTIIdfbDss06Vyk/smlb8dohmkfQov6Q/AKHUDXmrCbIIDCiuw70/z73y4i
	// uviAuxYovrqSugryb4tStUMTogmft4methz1/O/083XHwBNKBPnS2fobYDfBxqkX
	// tH26woCjrEr/O/wngo6iFp7b5yJlyXapN0x+iOF3CShIhAECQQD2gZ6LLYdxSP8i
	// aRYAPOh10mF5IHt2dl89eOjNiqVGMlkV5aXNT80jAQr/kWGZfIjscb/xkawSKQKs
	// ovcn99GRAkEAteL02mBrCLfn2idBwXTdil+yeigReAZmRpqQuAfTRZN4RM+5Dw3q
	// X0IiCkR3oyiwx89n1eGmz1JTZRxoY1AIQQJAWVbQ5xAxLlWOYiJD3wI0Hb+JpCSp
	// ml18VwMjHJtLGw3US6NXW/m4Fx+hpM5D2STRWyA+uIZbHpnOZlMJ0Gp4gQJBAK38
	// 66JV5y1Q1r2tHc6UHzQ1tMH7wDIjVQSm6FbSTXxZxAt29Rx8gD0dQvi1ZAg0bV7F
	// fRtwnqPlqZaoJQcTUMECQQD1Dh+Mu3OMb5AHnrtbk9l1qjM3U81QBKdyF0RY+djo
	// b3cR9I7+hurpqhJmQ7yuvAWe2xWc+YNTQ48FDJTogXlB
	// -----END RSA PRIVATE KEY-----`

	certPem = `-----BEGIN CERTIFICATE-----
MIIDrTCCAxagAwIBAgIQKYs1sciDjU/yBDKECiqedDANBgkqhkiG9w0BAQUFADAk
MQswCQYDVQQGEwJDTjEVMBMGA1UEChMMQ0ZDQSBURVNUIENBMB4XDTEyMDgyODAy
NTc1N1oXDTE0MDYyODAyNTc1N1owczELMAkGA1UEBhMCQ04xFTATBgNVBAoTDENG
Q0EgVEVTVCBDQTERMA8GA1UECxMITG9jYWwgUkExFDASBgNVBAsTC0VudGVycHJp
c2VzMSQwIgYDVQQDFBswNDFAWjIwMTEwODIzQHRlc3RAMDAwMDAwMDEwgZ8wDQYJ
KoZIhvcNAQEBBQADgY0AMIGJAoGBALluXyP1nHglJUTijVciTCSX3T6YxfJTeXqv
PYDI2bQdLdP+M/pQhqSnyICyCjlVewE4s2n/2ssCekuV1+xFotpMYad7rHLds0FG
Mja+eCLqUzpQwXFDJc4y+CIb/zcj1q+6HXdYg7Qr9qkpupdms/fI7dElJcOhHwY0
ikBS/ivHAgMBAAGjggGPMIIBizAfBgNVHSMEGDAWgBRGctwlcp8CTlWDtYD5C9vp
k7P0RTAdBgNVHQ4EFgQUUA/8Hd7EYZgDDwCYt+XmO3gl1lAwCwYDVR0PBAQDAgTw
MAwGA1UdEwQFMAMBAQAwOwYDVR0lBDQwMgYIKwYBBQUHAwEGCCsGAQUFBwMCBggr
BgEFBQcDAwYIKwYBBQUHAwQGCCsGAQUFBwMIMIHwBgNVHR8EgegwgeUwT6BNoEuk
STBHMQswCQYDVQQGEwJDTjEVMBMGA1UEChMMQ0ZDQSBURVNUIENBMQwwCgYDVQQL
EwNDUkwxEzARBgNVBAMTCmNybDEyN18yMzMwgZGggY6ggYuGgYhsZGFwOi8vdGVz
dGxkYXAuY2ZjYS5jb20uY246Mzg5L0NOPWNybDEyN18yMzMsT1U9Q1JMLE89Q0ZD
QSBURVNUIENBLEM9Q04/Y2VydGlmaWNhdGVSZXZvY2F0aW9uTGlzdD9iYXNlP29i
amVjdGNsYXNzPWNSTERpc3RyaWJ1dGlvblBvaW50MA0GCSqGSIb3DQEBBQUAA4GB
ANhD7dsg+uQMBuAcewdtbViOXCZCqXeFw0ZicZq0zkVA+NdjrejEWgcS2S1lNqYY
VDnyTIghECm6UxGO4UEF8/nwYsYpQJKtpdjHGbiDVvja/xcNaGCaH+ER+n08uAdB
ikahaQLV1atGk63K701Jtj061/jqkF2/Drv6FY+Uy+Rn
-----END CERTIFICATE-----`
)

// var chinaPaymentPriKey *rsa.PrivateKey
var chinaPaymentCert *x509.Certificate

// 缓存商户密钥
var keyCache map[string]*rsa.PrivateKey

// 读私钥
func initPrivKey(priKeyPem string) *rsa.PrivateKey {

	// 从缓存中查询
	mk := keyCache[priKeyPem]

	// 存在 返回
	if mk != nil {
		return mk
	}
	// 没有则创建一个
	PEMBlock, _ := pem.Decode([]byte(priKeyPem))
	if PEMBlock == nil {
		log.Fatalf("Could not parse Rsa Private Key PEM")
	}
	if PEMBlock.Type != "RSA PRIVATE KEY" {
		log.Fatalf("Found wrong key type" + PEMBlock.Type)
	}
	chinaPaymentPriKey, err := x509.ParsePKCS1PrivateKey(PEMBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	keyCache[priKeyPem] = chinaPaymentPriKey

	return chinaPaymentPriKey
}

// 读证书
func init() {
	PEMBlock, _ := pem.Decode([]byte(certPem))
	if PEMBlock == nil {
		log.Fatalf("Could not parse Certificate PEM")
	}
	if PEMBlock.Type != "CERTIFICATE" {
		log.Fatalf("Found wrong key type" + PEMBlock.Type)
	}
	var err error
	chinaPaymentCert, err = x509.ParseCertificate(PEMBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	// init cache
	keyCache = make(map[string]*rsa.PrivateKey)
}

// SignatureUseSha1WithRsa 通过私钥用 SHA1WithRSA 签名，返回 hex 签名
func signatureUseSha1WithRsa(origin []byte, priKeyPem string) string {
	// gen privatekey
	// TODO 优化，只需要初始化一次
	chinaPaymentPriKey := initPrivKey(priKeyPem)
	hashed := sha1.Sum(origin)

	sign, err := rsa.SignPKCS1v15(rand.Reader, chinaPaymentPriKey, crypto.SHA1, hashed[:])
	if err != nil {
		log.Errorf("fail to sign with Sha1WithRsa %s", err)
	}

	return hex.EncodeToString(sign)
}

// CheckSignatureUseSha1WithRsa 通过证书用 SHA1WithRSA 验签，如果验签通过，err 值为 nil
func checkSignatureUseSha1WithRsa(origin []byte, hexSign string) (err error) {
	sign, err := hex.DecodeString(hexSign)
	if err != nil {
		log.Errorf("hex decode error %s", err)
		return err
	}

	err = chinaPaymentCert.CheckSignature(x509.SHA1WithRSA, origin, sign)
	if err != nil {
		log.Errorf("signature error %s", err)
	}
	return err
}
