package goconf

import (
	"encoding/base64"
	"testing"

	"github.com/CardInfoLink/quickpay/security"
)

func TestEncryptMongoURL(t *testing.T) {
	testEncryptField(Config.Mongo.URL, t)
}
func TestEncryptUnionLiveEncryptKey(t *testing.T) {
	testEncryptField(Config.UnionLive.EncryptKey, t)
}
func TestEncryptUnionLiveSignKey(t *testing.T) {
	testEncryptField(Config.UnionLive.SignKey, t)
}

func TestDecryptMongoEncryptURL(t *testing.T) {
	testDecryptField(Config.Mongo.EncryptURL, t)
}
func TestDecryptUnionLiveEncryptEncryptKey(t *testing.T) {
	testDecryptField(Config.UnionLive.EncryptEncryptKey, t)
}

func TestDecryptUnionLiveEncryptSignKey(t *testing.T) {
	testDecryptField(Config.UnionLive.EncryptSignKey, t)
}

func testEncryptField(field string, t *testing.T) {
	t.Logf("origin: %s", field)

	ciphertext, err := security.RSAEncrypt([]byte(field), publicKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	b64Cipher := base64.StdEncoding.EncodeToString(ciphertext)
	t.Logf("encrypt: %s", b64Cipher)
}

func testDecryptField(field string, t *testing.T) {
	t.Logf("encrypt: %s", field)

	origData, err := security.RSADecryptBase64(field, privateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("origin: %s", origData)
}
