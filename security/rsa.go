package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RSAEncrypt RSA 加密
func RSAEncrypt(origData, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
}

// RSADecrypt RSA 解密
func RSADecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(ciphertext))
}

// RSADecryptBase64 RSA 解密 Base64 密文
func RSADecryptBase64(b64Cipher string, privateKey []byte) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(b64Cipher)
	if err != nil {
		return nil, err
	}

	return RSADecrypt(cipherText, privateKey)
}
