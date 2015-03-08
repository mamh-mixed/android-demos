package bindingpay

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
)

// SignatureUseSha1 使用 SHA1 算法签名， sha1(data + key).Hex()
func SignatureUseSha1(data, key string) string {
	s := sha1.New()
	io.WriteString(s, data+key)
	return hex.EncodeToString(s.Sum(nil))
}

// CheckSignatureUseSha1 使用 SHA1 算法验签， sha1(data + key).Hex()
func CheckSignatureUseSha1(data, key, signature string) bool {
	result := SignatureUseSha1(data, key)
	return strings.EqualFold(result, signature)
}
