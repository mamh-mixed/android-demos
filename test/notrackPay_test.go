package test

import (
	"strconv"
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
)

func xTestNoTrackPaymentHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/noTrackPayment?merId=" + testMerID

	b := model.NoTrackPayment{
		MerId:       testMerID,
		TransType:   "SALE",
		SubMerId:    "SM123456",
		MerOrderNum: strconv.FormatInt(time.Now().UnixNano(), 10),
		TransAmt:    120,
		CurrCode:    "156",
		AcctName:    "Peter",
		AcctNum:     testMSCCard,
		IdentType:   "0",
		IdentNum:    testCUPIdentNum,
		PhoneNum:    testCUPPhone,
		AcctType:    "10",
		ValidDate:   testMSCValidDate,
		Cvv2:        testMSCCVV2,
	}

	var aes = tools.NewAESCBCEncrypt(testEncryptKey)
	b.AcctName = aes.Encrypt(b.AcctName)
	b.AcctNum = aes.Encrypt(b.AcctNum)
	b.IdentNum = aes.Encrypt(b.IdentNum)
	b.PhoneNum = aes.Encrypt(b.PhoneNum)
	b.ValidDate = aes.Encrypt(b.ValidDate)
	b.Cvv2 = aes.Encrypt(b.Cvv2)
	if aes.Err != nil {
		panic(aes.Err)
	}
	doPost(url, b, t)
}
