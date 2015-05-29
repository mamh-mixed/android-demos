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
	"github.com/CardInfoLink/quickpay/config"
	"github.com/omigo/log"
	"io"
	"os"
	"strings"
)

var sysKey []byte

func init() {
	firstPart := config.GetValue("app", "encryptKey")
	whole := firstPart + "TEZMUboYmBLVfjnduURAk4="
	bytes, err := base64.StdEncoding.DecodeString(whole)
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}
	sysKey = bytes
}

// AesCBCMode 如果key位base64编码过的字符串
type AesCBCMode struct {
	Key    []byte
	Err    error
	sysAes *AesCBCMode
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
		sysAes: &AesCBCMode{
			Key: sysKey,
		},
	}
}

// DcyAndUseSysKeyEcy 解密商户字段后用系统的key进行加密
// decrypted 解密后的明文 encrypted 使用新key后的密文
func (a *AesCBCMode) DcyAndUseSysKeyEcy(ct string) (decrypted, encrypted string) {

	if a.sysAes == nil {
		a.sysAes = &AesCBCMode{Key: sysKey}
	}

	// decrypt
	decrypted = a.Decrypt(ct)

	if a.Err != nil {
		return decrypted, decrypted
	}
	// encrypt
	encrypted = a.sysAes.Encrypt(decrypted)

	if a.sysAes.Err != nil {
		// 将错误传递到a
		a.Err = a.sysAes.Err
	}
	return
}

// UseSysKeyDcyAndMerEcy 使用系统的key解密再用商户的key加密
func (a *AesCBCMode) UseSysKeyDcyAndMerEcy(ct string) string {

	if a.sysAes == nil {
		a.sysAes = &AesCBCMode{Key: sysKey}
	}

	decrypted := a.sysAes.Decrypt(ct)

	// log.Debugf("orig: %s, decrypted: %s", ct, decrypted)

	if a.sysAes.Err != nil {
		// 将错误传递到a
		a.Err = a.sysAes.Err
	}
	encrypted := a.Encrypt(decrypted)
	return encrypted
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
