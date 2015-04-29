package entrance

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// SignatureUseSha1 使用 SHA1 算法签名， sha1(data + key).Hex()
func SignatureUseSha1(data []byte, key string) string {
	s := sha1.New()
	s.Write(data)
	io.WriteString(s, key)
	return hex.EncodeToString(s.Sum(nil))
}

// CheckSignatureUseSha1 使用 SHA1 算法验签， sha1(data + key).Hex()
func CheckSignatureUseSha1(data []byte, key, signature string) bool {
	result := SignatureUseSha1(data, key)
	return strings.EqualFold(result, signature)

	// return true // TODO only for testing
}

// CheckSignature 根据商户ID到数据库查找签名密钥，然后进行验签
func CheckSignature(data []byte, merId, signature string) (result bool, ret *model.BindingReturn) {
	m, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		if err.Error() == "not found" {
			ret = mongo.RespCodeColl.Get("200063")
			return false, ret
		}
		return false, mongo.RespCodeColl.Get("000001")
	}
	result = CheckSignatureUseSha1(data, m.SignKey, signature)

	// only for test
	result = true

	return result, nil
}

// Signature 根据商户ID到数据库查找签名密钥，然后拼接到数据后面，签名
func Signature(data []byte, merId string) string {
	m, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		log.Errorf("Signature find Merchant error")
		return ""
	}
	return SignatureUseSha1(data, m.SignKey)
}
