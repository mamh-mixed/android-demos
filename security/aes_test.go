package security

import (
	"testing"

	"github.com/CardInfoLink/log"
)

var aesCBF, aesCBC = &AESCFBMode{}, &AESCBCMode{}

func init() {
	aesCBC = NewAESCBCEncrypt("AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8=")
}

func TestAESCFBEncryptAndDecrypt(t *testing.T) {

	key := []byte("1234567890123456")
	aesCBF.Key = key
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

func TestCBCAESEncryptAndDecrypt(t *testing.T) {

	pt := "中国最好，中国最棒，ye"
	encrypted := aesCBC.Encrypt(pt)
	log.Debugf("%s", encrypted)

	decrypted, encrypted1 := aesCBC.DcyAndUseSysKeyEcy(encrypted)
	log.Debugf("%s %s", decrypted, encrypted1)

	encrypted2 := aesCBC.UseSysKeyDcyAndMerEcy(encrypted1)

	if aesCBC.Err != nil {
		t.Error(aesCBC.Err)
		t.FailNow()
	}

	decrypted1 := aesCBC.Decrypt(encrypted2)

	if decrypted1 != decrypted {
		t.Errorf("%s not equal %s", decrypted1, decrypted)
		t.FailNow()
	}

}
func TestEncrypt(t *testing.T) {
	accnum := aesCBC.Encrypt("6222020302062061901")
	accname := aesCBC.Encrypt("陈芝锐")
	// cvv2 := aesCBCEncrypt("")
	identnum := aesCBC.Encrypt("440583199111031012")
	validdate := aesCBC.Encrypt("09/18")
	log.Debugf("%s,%s,%s,%s", accnum, accname, identnum, validdate)
}

func TestDecrypt(t *testing.T) {

	s := "zAnIKYwZqy+LUREI0thLomWTMxmRaD1NsHd1pNMsGRRuiBuMK+t6shWrJyIHxggm"
	decrypted := aesCBC.Decrypt(s)
	if aesCBC.Err != nil {
		t.Error(aesCBC.Err)
		t.FailNow()
	}
	log.Debugf("%s", decrypted)
}
