package core

import (
	"github.com/omigo/log"
	"testing"
)

func TestAesCFBEncryptAndDecrypt(t *testing.T) {

	// key := []byte("1234567890123456")
	pt := "中国最好，中国最棒，ye"
	encrypted := aesCFBEncrypt(pt)
	log.Debugf("%s", encrypted)
	decrypted := aesCFBDecrypt(encrypted)
	if string(decrypted) != pt {
		t.Error("decrypt fail")
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}

func TestCBCAesEncryptAndDecrypt(t *testing.T) {

	// key := []byte("1234567890123456")
	pt := "中国最好，中国最棒，ye"
	encrypted := aesCBCEncrypt(pt)
	log.Debugf("%s", encrypted)
	decrypted := aesCBCDecrypt(encrypted)
	if string(decrypted) != pt {
		t.Error("decrypt fail")
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}
