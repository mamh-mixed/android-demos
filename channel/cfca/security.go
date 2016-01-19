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
	"github.com/CardInfoLink/log"
)

// var chinaPaymentPriKey *rsa.PrivateKey
var chinaPaymentCert *x509.Certificate

// 缓存商户密钥
var keyCache = cache.New(model.Cache_CfcaMerRSAPrivKey)

// 读私钥
func initPrivKey(priKeyPem string) (*rsa.PrivateKey, error) {
	// 从缓存中查询
	mk, found := keyCache.Get(priKeyPem)

	// 存在 返回
	if found {
		// log.Debug("get key from cache")
		pk := mk.(*rsa.PrivateKey)
		return pk, nil
	}
	// 没有则创建一个
	PEMBlock, _ := pem.Decode([]byte(priKeyPem))
	if PEMBlock == nil {
		return nil, fmt.Errorf("for input privateKeyPem:%s, %s", priKeyPem,
			"Could not parse Rsa Private Key PEM")
	}
	if PEMBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Found wrong key type %s", PEMBlock.Type)
	}
	insPriKey, err := x509.ParsePKCS1PrivateKey(PEMBlock.Bytes)
	if err != nil {
		return nil, err
	}
	keyCache.Set(priKeyPem, insPriKey, cache.NoExpiration)

	return insPriKey, nil
}

// 读证书
func init() {
	certPemFile := goconf.Config.CFCA.CheckSignPublicKey
	certPem, err := ioutil.ReadFile(certPemFile)
	if err != nil {
		fmt.Printf("read cfca cert error: %s\n", err)
		os.Exit(3)
	}

	PEMBlock, _ := pem.Decode(certPem)
	if PEMBlock == nil {
		fmt.Println("Could not parse Certificate PEM")
		os.Exit(3)
	}
	if PEMBlock.Type != "CERTIFICATE" {
		fmt.Printf("Found wrong key type: %s\n", PEMBlock.Type)
		os.Exit(3)
	}
	// var err error
	chinaPaymentCert, err = x509.ParseCertificate(PEMBlock.Bytes)
	if err != nil {
		fmt.Printf("parse certificate error: %s", err)
		os.Exit(3)
	}
}

// SignatureUseSha1WithRsa 通过私钥用 SHA1WithRSA 签名，返回 hex 签名
func signatureUseSha1WithRsa(origin []byte, priKeyPem string) (string, error) {
	// gen privatekey
	chinaPaymentPriKey, err := initPrivKey(priKeyPem)
	if err != nil {
		log.Error(err)
		return "", err
	}
	hashed := sha1.Sum(origin)

	sign, err := rsa.SignPKCS1v15(rand.Reader, chinaPaymentPriKey, crypto.SHA1, hashed[:])
	if err != nil {
		log.Errorf("fail to sign with Sha1WithRsa %s", err)
	}

	return hex.EncodeToString(sign), nil
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
