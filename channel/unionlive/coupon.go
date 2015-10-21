package coupon

import "github.com/CardInfoLink/quickpay/security"

var encryptKey = "11111"
var signKey = "xxxxx"

/*
   B1 = DES(S1, K1)
   S1 = Base64Encode(B1)
   S3 = K2 + S1
   M1 = Lowwer(Hex(MD5(S3)))
   报文：R = M1 + S2
*/
func encryptAndSign(src []byte) ([]byte, err) {
	enc, err := security.DESEncrypt(src, encryptKey)
	if err != nil {
		return nil, err
	}

}
