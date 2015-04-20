package core

import (
	"testing"

	"github.com/omigo/log"
)

var aesCBF, aesCBC = AesCFBMode{}, AesCBCMode{}

func TestAesCFBEncryptAndDecrypt(t *testing.T) {

	// key := []byte("1234567890123456")
	aesCBF := new(AesCFBMode)
	pt := "中国最好，中国最棒，ye"
	encrypted := aesCBF.Encrypt(pt)
	log.Debugf("%s", encrypted)
	decrypted := aesCBF.Decrypt(encrypted)
	if string(decrypted) != pt {
		t.Error("decrypt fail")
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}

func TestCBCAesEncryptAndDecrypt(t *testing.T) {

	// key := []byte("1234567890123456")
	pt := "中国最好，中国最棒，ye"
	encrypted := aesCBC.Encrypt(pt)
	log.Debugf("%s", encrypted)
	decrypted := aesCBC.Decrypt(encrypted)
	if string(decrypted) != pt {
		t.Error("decrypt fail")
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}
func TestEncrypt(t *testing.T) {
	accnum := aesCBC.Encrypt("6222020302062061901")
	accname := aesCBC.Encrypt("陈芝锐")
	// cvv2 := aesCBCEncrypt("")
	identnum := aesCBC.Encrypt("440583199111031012")
	validdate := aesCBC.Encrypt("09/18")
	log.Debugf("%s,%s,%s,%s", accnum, accname, identnum, validdate)
}

func TestDecrypt2(t *testing.T) {

	s := "8XSOZyOvovSrpsmPyz/8CAUS6lXdQqG9gyRTBubsRZg="
	decrypted := aesCBC.Decrypt(s)
	if aesCBC.Err != nil {
		t.Error(aesCBC.Err)
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}

func TestDecrypt(t *testing.T) {

	s := "44906806872556215819411164477969f321fdb9d00279ac6755565d0348274ea456823ee5210e7bb0eedb3bbd8035a3"
	decrypted := aesCBC.Decrypt(s)
	if aesCBC.Err != nil {
		t.Error(aesCBC.Err)
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}
