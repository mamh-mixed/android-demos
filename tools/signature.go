package tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"github.com/omigo/g"
	"log"
)

var privateKey *rsa.PrivateKey

var opts rsa.PSSOptions

func init() {

	// pemBlock, _ := pem.Decode([]byte(key))
	// pk, _ := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	// loadCer()
	privateKey, _ = rsa.GenerateKey(rand.Reader, 1024)
	opts.SaltLength = rsa.PSSSaltLengthAuto
}

// ChinaPaySignature 中金支付渠道签名
// message  采用Base64编码
// signature 采用Sha1WithRsa签名后用Hex编码
func ChinaPaySignature(data string) (message, signature string) {
	// to xml
	xmlBytes := ToXML(data)

	return EncodeBase64(xmlBytes), EncodeHex(SignatureUseSha1WithRsa(xmlBytes))
}

// CheckChinaPaySignature 中金支付渠道验签
func CheckChinaPaySignature(data string, signature string) bool {
	// encode base64
	message := DecodeBase64(data)
	// ecode hex
	sign := DecodeHex(signature)
	// verify
	err := CheckSignatureUseSha1WithRsa(message, sign)

	return err == nil
}

// SignatureUseSha1WithRsa 使用 SHA1WithRSA 私钥签名
func SignatureUseSha1WithRsa(data []byte) []byte {
	// hasded
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)
	// sign
	sgined, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA1, hashed, &opts)

	if err != nil {
		g.Error("fail to sign with Sha1WithRsa (%s)", err)
	}
	return sgined
}

// CheckSignatureUseSha1WithRsa 使用 SHA1WithRSA 公钥验签
func CheckSignatureUseSha1WithRsa(data []byte, signature []byte) error {

	// hashed
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)

	return rsa.VerifyPSS(&privateKey.PublicKey, crypto.SHA1, hashed, signature, &opts)

}

// SignatureUseSha1 使用 SHA1 算法签名， sha1(data + "//" + key).Hex()
func SignatureUseSha1(data []byte, key string) []byte {
	log.Println("unimplement")

	return data
}

// CheckSignatureUseSha1 使用 SHA1 算法验签， sha1(data + "//" + key).Hex()
func CheckSignatureUseSha1(data []byte, key string) []byte {
	log.Println("unimplement")

	return data
}
