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

// AesCBCMode 如果key位base64编码过的字符串
// 必须先调用DecodeKey方法，再进行加解密
type AesCBCMode struct {
	Key []byte
	Err error
}

type AesCFBMode struct {
	Key []byte
	Err error
}

// TODO 测试时统一使用 key＝000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f 共 32 字节
// 生产上，每个商户都有自己的 key， 需要把这些 key 放入数据库，类似签名
// var key []byte

// func init() {
// 	key, _ = base64.StdEncoding.DecodeString("AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8=")
// }

// aesCFBEncrypt aes 加密
// 对商户敏感信息加密
func (a *AesCFBMode) Encrypt(pt string) string {

	block, err := aes.NewCipher(a.Key)
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

	block, err := aes.NewCipher(a.Key)
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

	block, err := aes.NewCipher(a.Key)
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
			a.Err = err.(error)
		}
	}()

	ct = strings.TrimSpace(ct)
	ciphertext, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		a.Err = err
		return ct
	}
	block, err := aes.NewCipher(a.Key)
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

// DecodeKey 将base64编码过的44位key
// 转成32位字节数组
func (a *AesCBCMode) DecodeKey(key string) {
	a.Key, a.Err = base64.StdEncoding.DecodeString(key)
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
