package tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io"
	"log"
	"strings"

	"github.com/omigo/g"
)

var privateKey *rsa.PrivateKey

const (
	privKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCvJC9MMGRKmxRBI0KMjDtz2KooIc6XOljHPWhTfAamhV3A5v5y
PiZr4haMDpulU08Y0JxsegwDwfbscQrhG7nvilIqIa+HiI1xkfFxjtNUrMN5hpvO
8HUUfwqzb5EdllQcv/C0xxBkeCECIb86JJry7ty4mNBkN2idbGxldMi90QIDAQAB
AoGATvTIIdfbDss06Vyk/smlb8dohmkfQov6Q/AKHUDXmrCbIIDCiuw70/z73y4i
uviAuxYovrqSugryb4tStUMTogmft4methz1/O/083XHwBNKBPnS2fobYDfBxqkX
tH26woCjrEr/O/wngo6iFp7b5yJlyXapN0x+iOF3CShIhAECQQD2gZ6LLYdxSP8i
aRYAPOh10mF5IHt2dl89eOjNiqVGMlkV5aXNT80jAQr/kWGZfIjscb/xkawSKQKs
ovcn99GRAkEAteL02mBrCLfn2idBwXTdil+yeigReAZmRpqQuAfTRZN4RM+5Dw3q
X0IiCkR3oyiwx89n1eGmz1JTZRxoY1AIQQJAWVbQ5xAxLlWOYiJD3wI0Hb+JpCSp
ml18VwMjHJtLGw3US6NXW/m4Fx+hpM5D2STRWyA+uIZbHpnOZlMJ0Gp4gQJBAK38
66JV5y1Q1r2tHc6UHzQ1tMH7wDIjVQSm6FbSTXxZxAt29Rx8gD0dQvi1ZAg0bV7F
fRtwnqPlqZaoJQcTUMECQQD1Dh+Mu3OMb5AHnrtbk9l1qjM3U81QBKdyF0RY+djo
b3cR9I7+hurpqhJmQ7yuvAWe2xWc+YNTQ48FDJTogXlB
-----END RSA PRIVATE KEY-----`

	asn1Data = `-----BEGIN CERTIFICATE-----
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

func init() {

	pemBlock, _ := pem.Decode([]byte(privKeyPEM))
	if pemBlock == nil {
		g.Error("private key wrong (%s)", pemBlock)
	}
	privateKey, _ = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	// loadCer()
	// privateKey, _ = rsa.GenerateKey(rand.Reader, 1024)
}

// ChinaPaySignature 中金支付渠道签名
// message  采用Base64编码
// signature 采用Sha1WithRsa签名后用Hex编码
func ChinaPaySignature(data string) (message, signature string) {
	// to xml
	xmlBytes := ToXML(data)

	return EncodeBase64(xmlBytes), EncodeHex(SignatureUseSha1WithRsa(xmlBytes))
}

// SignatureUseSha1WithRsa 使用 SHA1WithRSA 私钥签名
func SignatureUseSha1WithRsa(data []byte) []byte {
	// hasded
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)
	// sign
	sgined, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)

	if err != nil {
		g.Error("fail to sign with Sha1WithRsa ", err)
	}
	return sgined
}

// CheckSignatureUseSha1WithRsa 使用 SHA1WithRSA 公钥验签
func CheckSignatureUseSha1WithRsa(signed, signature []byte) error {

	PEMBlock, _ := pem.Decode([]byte(asn1Data))
	if PEMBlock == nil {
		log.Fatal("Could not parse Public Key PEM")
	}
	if PEMBlock.Type != "CERTIFICATE" {
		log.Fatal("Found wrong key type" + PEMBlock.Type)
	}
	cert, err := x509.ParseCertificate(PEMBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	err = cert.CheckSignature(x509.SHA1WithRSA, signed, signature)
	if err != nil {
		g.Error("signature error ", err)
	}
	return err
}

// SignatureUseSha1 使用 SHA1 算法签名， sha1(data + "//" + key).Hex()
func SignatureUseSha1(data, key string) string {
	log.Println("SignatureUseSha1")
	s := sha1.New()
	io.WriteString(s, data+key)
	return hex.EncodeToString(s.Sum(nil))
}

// CheckSignatureUseSha1 使用 SHA1 算法验签， sha1(data + "//" + key).Hex()
func CheckSignatureUseSha1(data, key, signature string) bool {
	log.Println("CheckSignatureUseSha1")
	result := SignatureUseSha1(data, key)
	return strings.EqualFold(result, signature)
}
