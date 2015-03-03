package tools

import (
	"log"
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
func SignatureUseSha1(data []byte, key string) []byte {
	log.Println("unimplement")

	return data
}

// CheckSignatureUseSha1 使用 SHA1 算法验签， sha1(data + "//" + key).Hex()
func CheckSignatureUseSha1(data []byte, key string) []byte {
	log.Println("unimplement")

	return data
}
