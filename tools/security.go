package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

type AesCBCMode struct {
	Key string
	Err error
}

type AesCFBMode struct {
	Key string
	Err error
}

var key []byte

func init() {
	key, _ = base64.StdEncoding.DecodeString("AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8=")
}

// aesCFBEncrypt aes 加密
// 对商户敏感信息加密
func (a *AesCFBMode) Encrypt(pt string) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		a.Err = err
	}
	plaintext := []byte(pt)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		a.Err = err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext)
}

// aesCFBDecrypt aes 解密
func (a *AesCFBMode) Decrypt(ct string) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		a.Err = err
	}
	ciphertext, _ := hex.DecodeString(ct)
	if len(ciphertext) < aes.BlockSize {
		a.Err = err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}

// aesCBCEncrypt cbc mode
func (a *AesCBCMode) Encrypt(pt string) string {

	if a.Err != nil {
		return pt
	}
	plaintext := PKCS5Padding([]byte(pt), aes.BlockSize)

	if len(plaintext)%aes.BlockSize != 0 {
		a.Err = errors.New(fmt.Sprintf("%s : plaintext is not a multiple of the block size", pt))
		return pt
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		a.Err = err
		return pt
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	// 随机生成16个字节数组
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		a.Err = err
		return pt
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	// mode.CryptBlocks(ciphertext, plaintext)
	// return hex.EncodeToString(ciphertext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

// aesCBCDecrypt cbc mode
func (a *AesCBCMode) Decrypt(ct string) string {

	if a.Err != nil {
		return ct
	}
	defer func() {
		if err := recover(); err != nil {
			a.Err, _ = err.(error)
		}
	}()

	ct = strings.TrimSpace(ct)
	ciphertext, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		a.Err = err
		return ct
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		a.Err = err
		return ct
	}

	if len(ciphertext) < aes.BlockSize {
		a.Err = errors.New(fmt.Sprintf("%s : ciphertext too short", ct))
		return ct
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		a.Err = errors.New(fmt.Sprintf("%s : ciphertext is not a multiple of the block size", ct))
		return ct
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

func Xxx(ct string) string {

	ct = strings.TrimSpace(ct)
	ciphertext, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		panic(fmt.Sprintf("%s : ciphertext is not a multiple of the block size", ct))
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return string(ciphertext)
}
