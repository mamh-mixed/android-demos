package scanpay2

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/log"
)

var alipayPublicKey *rsa.PublicKey

// 读证书
func init() {
	certPemFile := goconf.Config.AlipayScanPay.AlipayPubKey
	certPem, err := ioutil.ReadFile(certPemFile)
	if err != nil {
		fmt.Printf("read cfca cert error: %s", err)
		os.Exit(3)
	}

	pemBlock, _ := pem.Decode(certPem)
	if pemBlock == nil {
		log.Fatalf("Could not parse Certificate PEM")
	}
	if pemBlock.Type != "PUBLIC KEY" {
		log.Fatalf("Found wrong key type" + pemBlock.Type)
	}
	publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	alipayPublicKey = publicKey.(*rsa.PublicKey)
}

func LoadPrivateKey(privateKeyPem []byte) *rsa.PrivateKey {
	pemBlock, _ := pem.Decode(privateKeyPem)
	if pemBlock == nil {
		log.Errorf("Could not parse RSA Private Key PEM")
		return nil
	}
	if pemBlock.Type != "RSA PRIVATE KEY" {
		log.Errorf("Found wrong key type: %s", pemBlock.Type)
		return nil
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		log.Errorf("parse PKCS1 private key error: %s", err)
		return nil
	}
	return privateKey
}

// Sha1WithRsa 通过私钥用 SHA1WithRSA 签名，返回 basd64 编码后的签名
func Sha1WithRsa(data []byte, priKey *rsa.PrivateKey) (b64Signed string, err error) {
	log.Debugf("%s", data)

	hashed := sha1.Sum(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA1, hashed[:])
	if err != nil {
		log.Errorf("fail to sign with Sha1WithRsa %s", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// Verify 通过证书用 SHA1WithRSA 验签，如果验签通过，err 值为 nil
func Verify(origin []byte, b64Sign string) (err error) {
	sign, err := base64.StdEncoding.DecodeString(b64Sign)
	if err != nil {
		log.Errorf("base64 decode error %s", err)
		return err
	}

	hashed := sha1.Sum(origin)
	err = rsa.VerifyPKCS1v15(alipayPublicKey, crypto.SHA1, hashed[:], sign)
	if err != nil {
		log.Errorf("signature error %s", err)
	}
	return err
}
