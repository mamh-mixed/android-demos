package cil

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// only for test
var (
	applePayMerId  = "123456" // apple pay 测试用商户号
	testChcd       = "00000050"
	testMchntId    = "050310058120002"
	testTerminalId = "00000001"

	// 万事达卡测试数据
	testMSCCard       = "5457210001000019"
	testMSCCVV2       = "300"
	testMSCValidDate  = "1412"
	testMSCTrackdata2 = "5457210001000019=1412101080080748"

	// VISA卡测试数据
	testVISCard       = "4761340000000019"
	testVISCVV2       = "830"
	testVISValidDate  = "1712"
	testVISTrackdata2 = "4761340000000019=171210114991787"

	// 银联卡测试数据
	testCUPCard      = "6225220100740059"
	testCUPCVV2      = "111"
	testCUPValidDate = "1605"
	testCUPPhone     = "13611111111"
	testCUPIdentNum  = "130412"

	// Apple Pay测试数据
	testAPPCard       = "5180841200282463"
	testAPPExpireDate = "180531"
)

// apple pay测试增加了如下路由
// {
//     "merId": "123456",
//     "cardBrand": "VIS",
//     "chanCode": "APT",
//     "chanMerId": "APT123456"
// }

// apple pay测试增加了渠道商户
// {
// 	"chanCode" : "APT",
//     "chanMerId" : "APT123456",
//     "chanMerName" : "Apple Pay测试渠道商户",
// 	"terminalId": "TID123456789",
// 	"insCode" : "99667788"
// }

func init() {
	// mongo.Connect()
	Connect()
}

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
		CliSN:      mongo.SnColl.GetDaySN(testMchntId, testTerminalId),
		SysSN:      mongo.SnColl.GetSysSN(),
	}

	ret := ConsumeByApplePay(ap)
	if ret.RespCode != "000000" {
		t.Errorf("ApplePay error %#v", ret)
	}
}

// 测试无卡直接支付
func TestConsumeNoTrack(t *testing.T) {
	p := &model.NoTrackPayment{
		MerId:       applePayMerId,
		SubMerId:    "SM123456",
		MerOrderNum: strconv.FormatInt(time.Now().UnixNano(), 10),
		TransAmt:    10,
		CurrCode:    "156",
		AcctName:    "Peter",
		AcctNum:     testMSCCard,
		IdentType:   "0",
		IdentNum:    testCUPIdentNum,
		PhoneNum:    testCUPPhone,
		AcctType:    "10",
		ValidDate:   testMSCValidDate,
		Cvv2:        testMSCCVV2,
		SendSmsId:   "",
		SmsCode:     "",
		Chcd:        testChcd,
		Mchntid:     testMchntId,
		TerminalId:  testTerminalId,
		CliSN:       mongo.SnColl.GetDaySN(testMchntId, testTerminalId),
		SysSN:       mongo.SnColl.GetSysSN(),
	}

	ret := Consume(p)
	if ret.RespCode != "000000" {
		t.Errorf("ApplePay error %#v", ret)
	}
}

// 测试冲正
func TestReversalHandle(t *testing.T) {
	// 运行冲正案例的时候，请调小超时时间
	transTimeout = 3 * time.Second // 超时时间
	reversalTimeouts = [...]time.Duration{transTimeout, transTimeout * 1, transTimeout * 8, transTimeout * 50, transTimeout * 1140}

	t.Log("以下报文错误，所以会超时")
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
		TerminalId:  testTerminalId,
		CliSN:       "1" + mongo.SnColl.GetDaySN(testMchntId, testTerminalId),
		SysSN:       mongo.SnColl.GetSysSN(),
	}
	Consume(p)
}
