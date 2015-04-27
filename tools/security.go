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

	"github.com/omigo/log"
)

// AesCBCMode 如果key位base64编码过的字符串
// 必须先调用DecodeKey方法，再进行加解密
type AesCBCMode struct {
	Key []byte
	Err error
}

// NewAESCBCEncrypt 创建一个 AES 加密对象，使用 CBC 模式
func NewAESCBCEncrypt(b64Key string) *AesCBCMode {
	bytesKey, err := base64.StdEncoding.DecodeString(b64Key)

	if err != nil {
		log.Errorf("AES key(%s) base64 decode error: %s", b64Key, err)
	}

	return &AesCBCMode{
		Key: bytesKey,
		Err: err,
	}
}

// Encrypt cbc mode
func (a *AesCBCMode) Encrypt(pt string) string {

	if a.Err != nil {
		return pt
	}
	plaintext := PKCS7Padding([]byte(pt), aes.BlockSize)

	if len(plaintext)%aes.BlockSize != 0 {
		a.Err = fmt.Errorf("%s : plaintext is not a multiple of the block size", pt)
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

// Decrypt cbc mode
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
	ciphertext = PKCS7UnPadding(ciphertext)
	return string(ciphertext)
}

// DecodeKey 将base64编码过的44位key
// 转成32位字节数组
func (a *AesCBCMode) DecodeKey(key string) {
	a.Key, a.Err = base64.StdEncoding.DecodeString(key)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type AesCFBMode struct {
	Key []byte
	Err error
}

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
