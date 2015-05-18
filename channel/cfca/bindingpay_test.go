package cfca

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
)

var (
	// merId            = "1426727710113"
	// bindingId        = "1000000000011"
	chanMerId       = "001405"
	chanBingingId   = "cf00fd61d5ef4d924485db88b584897e"
	sysOrderNum     = "aaaaaaaaaaaaaaaaaabb"
	sysOrigOrderNum = "aaaaaaaaaaaaaaaaaabb"
	acctName        = "张三"
	acctNum         = "6222020302062061908"
	identType       = "0"
	identNum        = "350583199009153732"
	phoneNum        = "18205960039"
	acctType        = "10"
	validDate       = ""
	cvv2            = ""
	sendSmsId       = "1000000000009"
	smsCode         = "12353"
	bankId          = "102"
	priKeyPem       = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCvJC9MMGRKmxRBI0KMjDtz2KooIc6XOljHPWhTfAamhV3A5v5y
PiZr4haMDpulU08Y0JxsegwDwfbscQrhG7nvilIqIa+HiI1xkfFxjtNUrMN5hpvO
8HUUfwqzb5EdllQcv/C0xxBkeCECIb86JJry7ty4mNBkN2idbGxldMi90QIDAQAB
AoGATvTIIdfbDss06Vyk/smlb8dohmkfQov6Q/AKHUDXmrCbIIDCiuw70/z73y4i
uviAuxYovrqSugryb4tStUMTogmft4methz1/O/083XHwBNKBPnS2fobYDfBxqkX
tH26woCjrEr/O/wngo6iFp7b5yJlyXapN0x+iOF3CShIhAECQQD2gZ6LLYdxSP8i
aRYAPOh10mF5IHt2dl89eOjNiqVGMlkV5aXNT80jAQr/kWGZfIjscb/xkawSKQKs
ovcn99GRAkEAteL02mBrCLfn2idBwXTdil+yeigReAZmRpqQuAfTRZN4RM+5Dw3q
X0IiCkR3oyiwx89n1eGmz1JTZRxoY1AIQQJAWVbQ5xAxLlWOYiJD3wI0Hb+JpCSp
ml18VwMjHJtLGw3US6NXW/m4Fx+hpM5D2STRWyA+uIZbHpnOZlMJ0Gp4gQJBAK38
66JV5y1Q1r2tHc6UHzQ1tMH7wDIjVQSm6FbSTXxZxAt29Rx8gD0dQvi1ZAg0bV7F
fRtwnqPlqZaoJQcTUMECQQD1Dh+Mu3OMb5AHnrtbk9l1qjM3U81QBKdyF0RY+djo
b3cR9I7+hurpqhJmQ7yuvAWe2xWc+YNTQ48FDJTogXlB
-----END RSA PRIVATE KEY-----`
	txSNUnBinding = "cf00fd61d5ef4d924485db88b584897e"
	settFlag      = "457"
	merOrderNum   = ""
	transAmt      = 1000
)

func TestSendRequest(t *testing.T) {

	req := &BindingRequest{
		Version: "2.0",
		Head: requestHead{
			InstitutionID: chanMerId, //测试ID
			TxCode:        "2502",
		},
		Body: requestBody{
			TxSNBinding: chanBingingId,
		},
		SignCert: priKeyPem,
	}

	response := sendRequest(req)
	if response == nil {
		t.Error("test unsucessful")
	}
}

func TestProcessBindingEnquiry(t *testing.T) {

	be := &model.BindingEnquiry{
		ChanBindingId: chanBingingId,
		ChanMerId:     chanMerId,
		SignCert:      priKeyPem,
	}
	resp := DefaultClient.ProcessBindingEnquiry(be)
	DefaultClient.ProcessBindingEnquiry(be)

	log.Debugf("response message  %s", resp)
}

func TestProcessBindingCreate(t *testing.T) {

	be := &model.BindingCreate{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		BankId:        bankId,
		AcctName:      acctName,
		AcctNum:       acctNum,
		IdentType:     identType,
		IdentNum:      identNum,
		PhoneNum:      phoneNum,
		AcctType:      acctType,
		ValidDate:     validDate,
		Cvv2:          cvv2,
		SendSmsId:     sendSmsId,
		SmsCode:       smsCode,
		SignCert:      priKeyPem,
	}
	resp := DefaultClient.ProcessBindingCreate(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessBindingRemove(t *testing.T) {

	be := &model.BindingRemove{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		TxSNUnBinding: txSNUnBinding,
		SignCert:      priKeyPem,
	}
	resp := DefaultClient.ProcessBindingRemove(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessBindingPayment(t *testing.T) {
	be := &model.BindingPayment{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		SettFlag:      settFlag,
		//需要变化
		SysOrderNum: sysOrderNum,
		TransAmt:    int64(transAmt),
		SignCert:    priKeyPem,
	}
	resp := DefaultClient.ProcessBindingPayment(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessBindingRefund(t *testing.T) {
	be := &model.BindingRefund{
		ChanMerId:       chanMerId,
		SysOrderNum:     sysOrderNum,
		SysOrigOrderNum: sysOrigOrderNum,
		TransAmt:        int64(transAmt),
		SignCert:        priKeyPem,
	}
	resp := DefaultClient.ProcessBindingRefund(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessPaymentEnquiry(t *testing.T) {
	be := &model.OrderEnquiry{
		ChanMerId:   chanMerId,
		SysOrderNum: sysOrderNum,
		SignCert:    priKeyPem,
	}
	resp := DefaultClient.ProcessPaymentEnquiry(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessRefundEnquiry(t *testing.T) {
	be := &model.OrderEnquiry{
		ChanMerId:   chanMerId,
		SysOrderNum: sysOrderNum,
		SignCert:    priKeyPem,
	}
	resp := DefaultClient.ProcessRefundEnquiry(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessTransChecking(t *testing.T) {
	// be := &model.BillingSummary{
	// 	ChanMerId: "001405",
	// 	SettDate:  "2015-03-19",
	// 	SignCert:  priKeyPem,
	// }
	resp := DefaultClient.ProcessTransChecking(chanMerId, "2015-04-13", priKeyPem)
	log.Debugf("response message  %#v", resp)
}
