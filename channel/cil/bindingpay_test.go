package cil

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"strconv"
	"testing"
	"time"
)

// ConsumeByApplePay ApplePay消费
// 目前3DSecure数据类型的万事达卡可以跑测试
func TestConsumeByApplePay(t *testing.T) {
	ap := &model.ApplePay{
		MerId:         applePayMerId,
		TransType:     "SALE",
		MerOrderNum:   fmt.Sprintf("%d", time.Now().UnixNano()),
		TransactionId: "49170302b04f74b56b0060f33e11a135134e48e8af80a50cefea6c079353b419",
		ApplePayData: model.ApplePayData{
			ApplicationPrimaryAccountNumber: "5180841200282463",
			ApplicationExpirationDate:       "180531",
			CurrencyCode:                    "156",
			TransactionAmount:               590,
			DeviceManufacturerIdentifier:    "050110030273",
			PaymentDataType:                 "3DSecure",
			PaymentData: model.PaymentData{
				OnlinePaymentCryptogram: "AOZSYAeX7VKTAAKv5hDuAoABFA==",
				EciIndicator:            "5",
			},
		},
		Chcd:       "00000050",
		Mchntid:    "050310058120002",
		TerminalId: "00000001",
		CliSN:      mongo.DaySNColl.GetDaySN("050310058120002", "00000001"),
		SysSN:      mongo.SnColl.GetSysSN(),
	}

	ret := ConsumeByApplePay(ap)

	t.Logf("%+v\n", ret)
}

func TestConsumeNoTrack(t *testing.T) {
	t.Log("=====================")
	p := &model.NoTrackPayment{
		MerId:       applePayMerId,
		SubMerId:    "SM123456",
		MerOrderNum: strconv.FormatInt(time.Now().UnixNano(), 10),
		TransAmt:    10,
		CurrCode:    "156",
		AcctName:    "Peter",
		AcctNum:     testAccountNO,
		IdentType:   "0",
		IdentNum:    testIdentNum,
		PhoneNum:    testPhoneNum,
		AcctType:    "10",
		ValidDate:   testValiDate,
		Cvv2:        testCVV2,
		SendSmsId:   "",
		SmsCode:     "",
		Chcd:        testChcd,
		Mchntid:     testMchntId,
		Terminalid:  testTerminalId,
		CliSN:       mongo.DaySNColl.GetDaySN("050310058120002", "00000001"),
		SysSN:       mongo.SnColl.GetSysSN(),
	}
	Consume(p)
}
