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
		TransactionId: fmt.Sprintf("%020d", time.Now().UnixNano()),
		ApplePayData: model.ApplePayData{
			ApplicationPrimaryAccountNumber: testMSCCard,
			ApplicationExpirationDate:       testMSCValidDate,
			CurrencyCode:                    "156",
			TransactionAmount:               100,
			DeviceManufacturerIdentifier:    "050110030273",
			PaymentDataType:                 "3DSecure",
			PaymentData: model.PaymentData{
				OnlinePaymentCryptogram: "AOZSYAeX7VKTAAKv5hDuAoABFA==",
				EciIndicator:            "5",
			},
		},
		Chcd:       testChcd,
		Mchntid:    testMchntId,
		TerminalId: testTerminalId,
		CliSN:      mongo.DaySNColl.GetDaySN(testMchntId, testTerminalId),
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
		AcctNum:     testCUPCard,
		IdentType:   "0",
		IdentNum:    testCUPIdentNum,
		PhoneNum:    testCUPPhone,
		AcctType:    "10",
		ValidDate:   testCUPValidDate,
		Cvv2:        testCUPCVV2,
		SendSmsId:   "",
		SmsCode:     "",
		Chcd:        testChcd,
		Mchntid:     testMchntId,
		Terminalid:  testTerminalId,
		CliSN:       mongo.DaySNColl.GetDaySN(testMchntId, testTerminalId),
		SysSN:       mongo.SnColl.GetSysSN(),
	}
	Consume(p)
}
