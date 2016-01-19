package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/CardInfoLink/log"
)

// AESCBCMode 如果key位base64编码过的字符串
type AESCBCMode struct {
	Key    []byte
	Err    error
	sysAES *AESCBCMode
}

// NewAESCBCEncrypt 创建一个 AES 加密对象，使用 CBC 模式
func NewAESCBCEncrypt(b64Key, sysKey string) *AESCBCMode {
	bytesKey, err := base64.StdEncoding.DecodeString(b64Key)

	if err != nil {
		log.Errorf("AES key(%s) base64 decode error: %s", b64Key, err)
	}

	return &AESCBCMode{
		Key: bytesKey,
		Err: err,
		sysAES: &AESCBCMode{
			Key: []byte(sysKey),
		},
	}
}

// DcyAndUseSysKeyEcy 解密商户字段后用系统的key进行加密
// decrypted 解密后的明文 encrypted 使用新key后的密文
func (a *AESCBCMode) DcyAndUseSysKeyEcy(ct, sysKey string) (decrypted, encrypted string) {
	if a.sysAES == nil {
		a.sysAES = &AESCBCMode{Key: []byte(sysKey)}
	}

	// decrypt
	decrypted = a.Decrypt(ct)

	if a.Err != nil {
		return decrypted, decrypted
	}
	// encrypt
	encrypted = a.sysAES.Encrypt(decrypted)

	if a.sysAES.Err != nil {
		// 将错误传递到a
		a.Err = a.sysAES.Err
	}
	return
}

// UseSysKeyDcyAndMerEcy 使用系统的key解密再用商户的key加密
func (a *AESCBCMode) UseSysKeyDcyAndMerEcy(ct string, sysKey []byte) string {

	if a.sysAES == nil {
		a.sysAES = &AESCBCMode{Key: []byte(sysKey)}
	}

	decrypted := a.sysAES.Decrypt(ct)

	// log.Debugf("orig: %s, decrypted: %s", ct, decrypted)

	if a.sysAES.Err != nil {
		// 将错误传递到a
		a.Err = a.sysAES.Err
	}
	encrypted := a.Encrypt(decrypted)
	return encrypted
}

// Encrypt cbc mode
func (a *AESCBCMode) Encrypt(pt string) string {

	if a.Err != nil {
		return pt
	}
	plaintext := pkcs7Padding([]byte(pt), aes.BlockSize)

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
func (a *AESCBCMode) Decrypt(ct string) string {

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
		a.Err = fmt.Errorf("%s : ciphertext too short", ct)
		return ct
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		a.Err = fmt.Errorf("%s : ciphertext is not a multiple of the block size", ct)
		return ct
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = pkcs7UnPadding(ciphertext)
	return string(ciphertext)
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AESCFBMode CFB 模式的 AES 加密
type AESCFBMode struct {
	Key []byte
	Err error
}

// Encrypt aesCFBEncrypt aes 加密  对商户敏感信息加密
func (a *AESCFBMode) Encrypt(pt string) string {

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

// Decrypt aes 解密
func (a *AESCFBMode) Decrypt(ct string) string {

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
