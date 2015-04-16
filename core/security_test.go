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
	encrypted := AesCBCEncrypt(pt)
	log.Debugf("%s", encrypted)
	decrypted := AesCBCDecrypt(encrypted)
	if string(decrypted) != pt {
		t.Error("decrypt fail")
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}
func TestEncrypt(t *testing.T) {
	accnum := AesCBCEncrypt("6222020302062061901")
	accname := AesCBCEncrypt("陈芝锐")
	// cvv2 := aesCBCEncrypt("")
	identnum := AesCBCEncrypt("440583199111031012")
	validdate := AesCBCEncrypt("09/18")
	log.Debugf("%s,%s,%s,%s", accnum, accname, identnum, validdate)
}

func TestDecrypt(t *testing.T) {

	s := "60202215176842555995459843018306154894ce1da849aa0af4699d5334fa5b"
	decrypted := AesCBCDecrypt(s)
	if decrypted != "张三" {
		t.Errorf("expect 张三 , but get %s", decrypted)
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}
