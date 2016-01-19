package unionlive

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/log"
)

var encryptKey = goconf.Config.UnionLive.EncryptKey
var signKey = goconf.Config.UnionLive.SignKey

/*
   enc = DES(S1, encryptKey)
   b64Enc = Base64Encode(enc)
   toSign = signKey + b64Enc
   signed = Lowwer(Hex(MD5(toSign)))
   报文：body = signed + b64Enc
*/
func encryptAndSign(src []byte) ([]byte, error) {
	// log.Debugf("src=%s, encryptKey=%s", src, encryptKey)
	// 1. 【加密报文】通过双方约定的密钥 encryptKey 使用 DES 加密报文明文 src,得到密文字节数组 enc;
	enc, err := security.DESEncrypt(src, encryptKey)
	if err != nil {
		return nil, err
	}

	// 2. 【加密报文】取得密文字节数组的 Base64 表示形式 b64Enc;
	b64Enc := base64.StdEncoding.EncodeToString(enc)

	// log.Debugf("signKey=%s, b64Enc=%s", signKey, b64Enc)
	// 3. 【计算签名】将双方约定的密钥 signKey(可以等于 encryptKey)拼接在 b64Enc 的前部,得到 toSign;
	toSign := signKey + b64Enc

	// 4. 【计算签名】计算 toSign 的 MD5,并取得 Hex 字符串的小写形式 signed;
	signed := fmt.Sprintf("%x", md5.Sum([]byte(toSign)))

	// 5. 【拼接结果】将 signed 和 b64Enc 拼接得到最终密文报文: R
	return []byte(signed + b64Enc), nil
}

/*
	msg='aba342.....hoyq8p....'

	signed = msg[:32]
	b64Enc = msg[32:]
	toSign = signKey + b64Enc
	actual = Lowwer(Hex(MD5(toSign)))

	if actual != signed {
	  Fail
	  return
	}

	b64Body = Base64Decode(b64Enc)
	body = DES(b64Body, encryptKey)
*/
func checkSignAndDecrypt(msg []byte) ([]byte, error) {
	// 1. 【拆解密文】取得密文报文字符串 msg 的前 32 个字符,得到签名 signed,取得 32 个字符后的内容, 得到 b64Enc;
	signed, b64Enc := msg[:32], msg[32:]

	// 2. 【验证签名】将双方约定的密钥 signKey 拼接在 b64Enc 的前部,得到 toSign;
	toSign := append([]byte(signKey), b64Enc...)

	// 3. 【验证签名】计算 toSign 的 MD5,并取得 Hex 字符串的小写形式 actual
	actual := fmt.Sprintf("%x", md5.Sum([]byte(toSign)))

	// 4. 【验证签名】判断 M1 是否等于 M2,不相等则验签失败,否则继续如下解密步骤;
	if string(actual[:]) != string(signed) {
		log.Errorf("check sign error: msg=%s, signed=%s, actual=%s", msg, signed, actual)
		return nil, errors.New("验签失败")
	}

	// 5. 【解密报文】取得 b64Enc 的反 Base64 后的密文字节数组 body;
	enc, err := base64.StdEncoding.DecodeString(string(b64Enc))
	if err != nil {
		log.Errorf("base64 decode error: %s", b64Enc)
		return nil, err
	}

	// 6. 【解密报文】通过双方约定的密钥 K2(可以等于 K1)使用 DES 解密密文字节数组 B1,得到明
	// 文报文:R
	return security.DESDecrypt(enc, encryptKey)
}
