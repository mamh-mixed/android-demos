package weixin

import (
	"crypto/hmac"
	"crypto/sha256"
	"testing"
)

func TestHmacSha256(t *testing.T) {
	plaintText, key := []byte("Hello world"), []byte("120943629")

	mac := hmac.New(sha256.New, key)

	mac.Write(plaintText)

	expectedMac := mac.Sum(nil)

	t.Logf("expected mac is %x", expectedMac)
}
