package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/omigo/log"
	"io"
)

// 16位
var key = []byte("1234567890123456")

// aesCFBEncrypt aes 加密
// 对商户敏感信息加密
func aesCFBEncrypt(pt string) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}
	plaintext := []byte(pt)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Panicln(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext)
}

// aesCFBDecrypt aes 解密
func aesCFBDecrypt(ct string) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}
	ciphertext, _ := hex.DecodeString(ct)
	if len(ciphertext) < aes.BlockSize {
		log.Panicln("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}

// aesCBCDecrypt cbc mode
func aesCBCEncrypt(pt string) string {

	plaintext := PKCS5Padding([]byte(pt), aes.BlockSize)

	if len(plaintext)%aes.BlockSize != 0 {
		log.Panicln("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Panicln(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext)
}

// aesCBCDecrypt cbc mode
func aesCBCDecrypt(ct string) string {

	ciphertext, _ := hex.DecodeString(ct)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}

	if len(ciphertext) < aes.BlockSize {
		log.Panicln("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		log.Panicln("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)
	return string(ciphertext)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
