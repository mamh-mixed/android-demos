package security

import (
	"crypto/sha1"
	"fmt"
)

// SHA1 SHA1 哈希
func SHA1(data []byte) string {
	sig := sha1.Sum(data)
	return fmt.Sprintf("%x", sig[:])
}

// SHA1WithKey SHA1 签名，key 附加直接拼接在数据后面
func SHA1WithKey(data, key string) string {
	sig := sha1.Sum([]byte(data + key))
	return fmt.Sprintf("%x", sig[:])
}
