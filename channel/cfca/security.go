package cfca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// var chinaPaymentPriKey *rsa.PrivateKey
var chinaPaymentCert *x509.Certificate

// 缓存商户密钥
var keyCache = cache.New(model.Cache_ChanMerRSAPrivKey)

// 读私钥
func initPrivKey(priKeyPem string) *rsa.PrivateKey {

	// 从缓存中查询
	mk, found := keyCache.Get(priKeyPem)

	// 存在 返回
	if found {
		// log.Debug("get key from cache")
		pk := mk.(*rsa.PrivateKey)
		return pk
	}
	// 没有则创建一个
	PEMBlock, _ := pem.Decode([]byte(priKeyPem))
	if PEMBlock == nil {
		log.Fatalf("Could not parse Rsa Private Key PEM")
	}
	if PEMBlock.Type != "RSA PRIVATE KEY" {
		log.Fatalf("Found wrong key type" + PEMBlock.Type)
	}
	chinaPaymentPriKey, err := x509.ParsePKCS1PrivateKey(PEMBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	keyCache.Set(priKeyPem, chinaPaymentPriKey, cache.NoExpiration)
	// keyCache[priKeyPem] = chinaPaymentPriKey

	return chinaPaymentPriKey
}

// 读证书
func init() {
	certPemFile := goconf.GetFile("cfca", "cert")
	certPem, err := ioutil.ReadFile(certPemFile)
	if err != nil {
		fmt.Printf("read cfca cert error: %s", err)
		os.Exit(3)
	}

	PEMBlock, _ := pem.Decode(certPem)
	if PEMBlock == nil {
		log.Fatalf("Could not parse Certificate PEM")
	}
	if PEMBlock.Type != "CERTIFICATE" {
		log.Fatalf("Found wrong key type" + PEMBlock.Type)
	}
	// var err error
	chinaPaymentCert, err = x509.ParseCertificate(PEMBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
}

// SignatureUseSha1WithRsa 通过私钥用 SHA1WithRSA 签名，返回 hex 签名
func signatureUseSha1WithRsa(origin []byte, priKeyPem string) string {
	// gen privatekey
	// TODO 优化，只需要初始化一次
	chinaPaymentPriKey := initPrivKey(priKeyPem)
	hashed := sha1.Sum(origin)

	sign, err := rsa.SignPKCS1v15(rand.Reader, chinaPaymentPriKey, crypto.SHA1, hashed[:])
	if err != nil {
		log.Errorf("fail to sign with Sha1WithRsa %s", err)
	}

	return hex.EncodeToString(sign)
}

// CheckSignatureUseSha1WithRsa 通过证书用 SHA1WithRSA 验签，如果验签通过，err 值为 nil
func checkSignatureUseSha1WithRsa(origin []byte, hexSign string) (err error) {
	sign, err := hex.DecodeString(hexSign)
	if err != nil {
		log.Errorf("hex decode error %s", err)
		return err
	}

	err = chinaPaymentCert.CheckSignature(x509.SHA1WithRSA, origin, sign)
	if err != nil {
		log.Errorf("signature error %s", err)
	}
	return err
}
