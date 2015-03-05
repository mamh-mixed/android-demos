package tools

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"strings"
)

// SignatureUseSha1WithRsa 使用 SHA1WithRSA 私钥签名
func SignatureUseSha1WithRsa(data []byte, privateKey string) []byte {
	log.Println("unimplement")

	return data
}

// CheckSignatureUseSha1WithRsa 使用 SHA1WithRSA 公钥验签
func CheckSignatureUseSha1WithRsa(data []byte, publicKey string) []byte {
	log.Println("unimplement")

	return data
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
