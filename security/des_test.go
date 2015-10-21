package security

import (
	"encoding/base64"
	"testing"
)

const encryptKey = "12345678"

var src = "1234567890"
var enc = "PXWVqYv/gJ00/ZV9q+rrCg=="

func TestDESEncrypt(t *testing.T) {

	encrypted, err := DESEncrypt([]byte(src), encryptKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	actual := base64.StdEncoding.EncodeToString(encrypted)

	if enc != string(actual) {
		t.Errorf("encryted error: expected=%s, actual=%s", enc, actual)
	}
}

func TestDESDecrypt(t *testing.T) {
	encrypted, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	decrypted, err := DESDecrypt(encrypted, encryptKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(decrypted) != src {
		t.Errorf("encryted error: expected=%s, actual=%s", src, decrypted)
	}
}
