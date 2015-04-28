package core

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"strconv"
	"testing"
	"time"
)

func TestProcessBindingCreate(t *testing.T) {
	bc := &model.BindingCreate{
		MerId:     merId,
		BindingId: bindingId,
		AcctName:  acctName,
		AcctNum:   acctNum,
		IdentType: identType,
		IdentNum:  identNum,
		PhoneNum:  phoneNum,
		AcctType:  acctType,
		ValidDate: validDate,
		Cvv2:      cvv2,
		SendSmsId: sendSmsId,
		SmsCode:   smsCode,
		BankId:    bankId,
	}

	ret := ProcessBindingCreate(bc)

	t.Logf("%+v", ret)

	if ret.RespCode != "000000" {
		t.Errorf("Excpeted success,but get failure: %+v", ret.RespMsg)
	}
}

func TestProcessBindingEnquiry(t *testing.T) {
	be := &model.BindingEnquiry{
		MerId:     merId,
		BindingId: bindingId,
	}

	// br, err := mongo.FindBindingRelation(be.MerId, be.BindingId)
	//
	// if err != nil {
	// 	t.Errorf("'FindBindingRelation' error: %s", err.Error())
	// }
	//
	// br.BindingStatus = "000009"
	// if err = mongo.UpdateBindingRelation(br); err != nil {
	// 	t.Errorf("'UpdateBindingRelation' error: %s", err.Error())
	// }

	ret := ProcessBindingEnquiry(be)

	t.Logf("%+v", ret)
}

func TestProcessBindingRemove(t *testing.T) {
	be := &model.BindingRemove{
		MerId:     removeMerId,
		BindingId: removeBindingId,
	}

	ret := ProcessBindingReomve(be)

	t.Logf("%+v", ret)
}

func TestProcessBindingPayment(t *testing.T) {

	be := &model.BindingPayment{
		MerId:       merId,
		BindingId:   bindingId,
		TransAmt:    int64(transAmt),
		MerOrderNum: merOrderNum,
	}
	ret := ProcessBindingPayment(be)

	if ret.RespCode == "" {
		t.Errorf("process payment but not get a respCode %+v", ret)
	}
	// t.Logf("%+v", ret)
	log.Debugf("%+v", ret)
}

func TestProcessBindingRefund(t *testing.T) {

	be := &model.BindingRefund{
		MerId:        merId,
		TransAmt:     int64(transAmt),
		OrigOrderNum: origOrderNum,
		MerOrderNum:  merOrderNum,
	}
	ret := ProcessBindingRefund(be)

	if ret.RespCode == "" {
		t.Errorf("process payment but not get a respCode %+v", ret)
	}
	log.Debugf("%+v", ret)
}

func TestProcessOrderEnquiry(t *testing.T) {

	be := &model.OrderEnquiry{
		MerId:        merId,
		OrigOrderNum: origOrderNum,
		ShowOrigInfo: showOrigInfo,
	}
	ret := ProcessOrderEnquiry(be)

	if ret.RespCode == "" {
		t.Errorf("process order query but not get a respCode %+v", ret)
	}
	log.Debugf("%+v,%+v", ret, ret.OrigTransDetail)
}

func TestProcessBillingDetails(t *testing.T) {

	be := &model.BillingDetails{
		MerId:    merId,
		SettDate: settDate,
	}
	ret := ProcessBillingDetails(be)

	if ret.RespCode == "" {
		t.Errorf("process billing details but not get a respCode %+v", ret)
	}
	log.Debugf("%+v", ret)
}

func TestProcessBillingSummary(t *testing.T) {

	be := &model.BillingSummary{
		MerId:    merId,
		SettDate: settDate,
	}
	ret := ProcessBillingSummary(be)

	if ret.RespCode == "" {
		t.Errorf("process billing summary but not get a respCode %+v", ret)
	}
	log.Debugf("%+v", ret)
}

func TestProcessNoTrackPayment(t *testing.T) {

	ntp := &model.NoTrackPayment{
		MerId:            testMerID,
		TransType:        "SALE",
		SubMerId:         "SM123456",
		MerOrderNum:      strconv.FormatInt(time.Now().UnixNano(), 10),
		TransAmt:         120,
		CurrCode:         "156",
		AcctNameDecrypt:  "Peter",
		AcctNumDecrypt:   testCUPCard,
		IdentType:        "0",
		IdentNumDecrypt:  testCUPIdentNum,
		PhoneNumDecrypt:  testCUPPhone,
		AcctType:         "10",
		ValidDateDecrypt: testCUPValidDate,
		Cvv2Decrypt:      testCUPCVV2,
	}

	var aes = tools.NewAESCBCEncrypt(testEncryptKey)

	ntp.AcctName = aes.Encrypt(ntp.AcctNameDecrypt)
	ntp.AcctNum = aes.Encrypt(ntp.AcctNumDecrypt)
	ntp.IdentNum = aes.Encrypt(ntp.IdentNumDecrypt)
	ntp.PhoneNum = aes.Encrypt(ntp.PhoneNumDecrypt)
	ntp.ValidDate = aes.Encrypt(ntp.ValidDateDecrypt)
	ntp.Cvv2 = aes.Encrypt(ntp.Cvv2Decrypt)
	if aes.Err != nil {
		panic(aes.Err)
	}

	ret := ProcessNoTrackPayment(ntp)

	if ret == nil {
		t.Error("NoTrackPayment process error")
	}

	t.Logf("%+v", ret)
}

// only for test
var (
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

	testMerID = "APPTEST"

	// 测试用的密钥
	testEncryptKey = "AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8="
)
