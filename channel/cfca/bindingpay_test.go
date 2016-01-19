package cfca

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

var (
	// merId            = "1426727710113"
	// bindingId        = "1000000000011"
	chanMerId = "001405"
	priKeyPem = `-----BEGIN RSA PRIVATE KEY-----
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
	// 	priKeyPem = `-----BEGIN RSA PRIVATE KEY-----
	// MIIEowIBAAKCAQEAs7qZX0NY89AbwbcX+rHNeI7BM5WQP4e284hQMPJe8Sf9h1Ln
	// ALxa57Fa1+Jv15WUVsVhXWLuYDctN2GKw7WFDUsLZR9Kl9wGc2V6oe9lcYL3HBHr
	// 4zidAZ0Ds9rJFQGNWvVlxT7Tbij520EpPbOobTycV/aMwakMys1X1teJofne1iTk
	// bPKv9y/cjCmcSJylMlC4fOqUNywSMe7til1qFGXokhJPykiUKGmm7aij0LgIG+ja
	// ljWiKKNhkyM+zXznuhPWC4n52UaT4ySNqQWG86yPHkTVmA16mDuDfdjGdefTCY3f
	// 0CNiE2UTLYHnv/hcl6OcNRvHVGmAwiin5ZW7ZwIDAQABAoIBAA3AGvfU2TteEqGR
	// dn4xmDy+/z0JFk4l/fNH3IrCIE1m99igmRfPUU0knoShyFjAEudIlphkd0RZNFZm
	// Wg17F8DamvaTnEteBAhUHTvCawpyMKGvJVLvA+QZtSzPck7vHXd2CuE6W0a0T9lc
	// dOFJm3swBe4c6MvwORBTAYNoXhlMOcs5nBdhLNOxC4pDDeDiXeO8bRQHG5xUxp7D
	// A07RNJaFqa6uyRV3nSk218xAklAzUExwzn/LcY0H6v+nuXcq3Lz+MTWNfY8QhuBL
	// umkcrKwlA3on8+ABmcdP9AMZv+4NoqiXfDA6hd+r1cNBl1BZJK2VrQlNyyNqe4Ps
	// mRsFafECgYEA4G9Nsc1iQ+7fjNrnXSkqsT4ecXdV0RDTZm75MrLtqcZG02+8pX2T
	// 71XmjdT3zmq7kxTa9qD7oAuRNjQcz18eKeGo6cM9zqaTyJmqdNj7F934HT3bDIS/
	// /Gyh0OBkgJD9tqq6po19wJkJ9FIfe9RpRKvyCZloG/g06/kucr3B2wUCgYEAzQGt
	// Dz3GvHRBPNYxmEWFas83CQPwEhzUfJvETZaRYIk7j9uLsgwiSEPBpG0eyBWKEFsJ
	// p+JWojRHKkUNHhpQPwbvJoAeQ0TEqsXSL3BoFWMFcMu2ktrcQmMP4NGliqCFL45n
	// JatHEWKUBugUYJatxwJYSxQNTY7FbvtqewVjgHsCgYEApBA5V5SwFMD2Il2TbAK/
	// 9rlWP9Pgo+gM4YCWIn2yRr1u4Pl5ifB2yCqfU2cvj3FulWJVfpzH1IMgL+OAfAco
	// Ya1YcSoMcJhMyAOtG6XSR+w3iAjDrC8OuVQgJjUiwuk6zuGXeFFOcBBvum6eHUN1
	// gIHBeUrvVCLpbObHEZGtuJECgYBs1gvzgH+Gw01zJ9/ykE4Rc2srbRzB8O5aLTQd
	// YOdTUef+KrdSUiDNLrOaQJhL7yt6HWrV51LJEGoLpdcd+ShLHbpPPUtTuSmT5Cv/
	// JXUMjaJwzKXj9y9iS0c9uu9g1nF+2uIl3HWBZE1kEUfoM3aUpckMKtwZJcfpcK6K
	// G0VFCwKBgG+FznQbNueYUtyJx/eD/gn3bBI9k55ra0DmBhgnorPzGdc+/LMr5CS+
	// TK8p8uw4sJiIAWCbJV3+f4QJ/ZTpKj2QnrAXDWRwUFMzmPV0e+Cro/AEHvx6v1ck
	// WTc7WnkJXcXk2xfJiKz+O4S23s55ULKLc/uWVisupBaM/hQWKOiM
	// -----END RSA PRIVATE KEY-----`

	// 	chanMerId = "001583"
	// 	priKeyPem = `-----BEGIN RSA PRIVATE KEY-----
	// MIIEowIBAAKCAQEAs7qZX0NY89AbwbcX+rHNeI7BM5WQP4e284hQMPJe8Sf9h1Ln
	// ALxa57Fa1+Jv15WUVsVhXWLuYDctN2GKw7WFDUsLZR9Kl9wGc2V6oe9lcYL3HBHr
	// 4zidAZ0Ds9rJFQGNWvVlxT7Tbij520EpPbOobTycV/aMwakMys1X1teJofne1iTk
	// bPKv9y/cjCmcSJylMlC4fOqUNywSMe7til1qFGXokhJPykiUKGmm7aij0LgIG+ja
	// ljWiKKNhkyM+zXznuhPWC4n52UaT4ySNqQWG86yPHkTVmA16mDuDfdjGdefTCY3f
	// 0CNiE2UTLYHnv/hcl6OcNRvHVGmAwiin5ZW7ZwIDAQABAoIBAA3AGvfU2TteEqGR
	// dn4xmDy+/z0JFk4l/fNH3IrCIE1m99igmRfPUU0knoShyFjAEudIlphkd0RZNFZm
	// Wg17F8DamvaTnEteBAhUHTvCawpyMKGvJVLvA+QZtSzPck7vHXd2CuE6W0a0T9lc
	// dOFJm3swBe4c6MvwORBTAYNoXhlMOcs5nBdhLNOxC4pDDeDiXeO8bRQHG5xUxp7D
	// A07RNJaFqa6uyRV3nSk218xAklAzUExwzn/LcY0H6v+nuXcq3Lz+MTWNfY8QhuBL
	// umkcrKwlA3on8+ABmcdP9AMZv+4NoqiXfDA6hd+r1cNBl1BZJK2VrQlNyyNqe4Ps
	// mRsFafECgYEA4G9Nsc1iQ+7fjNrnXSkqsT4ecXdV0RDTZm75MrLtqcZG02+8pX2T
	// 71XmjdT3zmq7kxTa9qD7oAuRNjQcz18eKeGo6cM9zqaTyJmqdNj7F934HT3bDIS/
	// /Gyh0OBkgJD9tqq6po19wJkJ9FIfe9RpRKvyCZloG/g06/kucr3B2wUCgYEAzQGt
	// Dz3GvHRBPNYxmEWFas83CQPwEhzUfJvETZaRYIk7j9uLsgwiSEPBpG0eyBWKEFsJ
	// p+JWojRHKkUNHhpQPwbvJoAeQ0TEqsXSL3BoFWMFcMu2ktrcQmMP4NGliqCFL45n
	// JatHEWKUBugUYJatxwJYSxQNTY7FbvtqewVjgHsCgYEApBA5V5SwFMD2Il2TbAK/
	// 9rlWP9Pgo+gM4YCWIn2yRr1u4Pl5ifB2yCqfU2cvj3FulWJVfpzH1IMgL+OAfAco
	// Ya1YcSoMcJhMyAOtG6XSR+w3iAjDrC8OuVQgJjUiwuk6zuGXeFFOcBBvum6eHUN1
	// gIHBeUrvVCLpbObHEZGtuJECgYBs1gvzgH+Gw01zJ9/ykE4Rc2srbRzB8O5aLTQd
	// YOdTUef+KrdSUiDNLrOaQJhL7yt6HWrV51LJEGoLpdcd+ShLHbpPPUtTuSmT5Cv/
	// JXUMjaJwzKXj9y9iS0c9uu9g1nF+2uIl3HWBZE1kEUfoM3aUpckMKtwZJcfpcK6K
	// G0VFCwKBgG+FznQbNueYUtyJx/eD/gn3bBI9k55ra0DmBhgnorPzGdc+/LMr5CS+
	// TK8p8uw4sJiIAWCbJV3+f4QJ/ZTpKj2QnrAXDWRwUFMzmPV0e+Cro/AEHvx6v1ck
	// WTc7WnkJXcXk2xfJiKz+O4S23s55ULKLc/uWVisupBaM/hQWKOiM
	// -----END RSA PRIVATE KEY-----
	// `

	chanBingingId   = "e169a3bd64ab455045b0129e1a18d53d"
	sysOrderNum     = util.SerialNumber()
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
		PrivateKey: priKeyPem,
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
		PrivateKey:    priKeyPem,
	}
	resp := DefaultClient.ProcessBindingEnquiry(be)

	log.Debugf("response message  %s", resp)
}

func TestProcessBindingCreate(t *testing.T) {

	be := &model.BindingCreate{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		BankCode:      bankId,
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
		PrivateKey:    priKeyPem,
	}
	resp := DefaultClient.ProcessBindingCreate(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessBindingRemove(t *testing.T) {

	be := &model.BindingRemove{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		TxSNUnBinding: txSNUnBinding,
		PrivateKey:    priKeyPem,
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
		PrivateKey:  priKeyPem,
	}
	resp := DefaultClient.ProcessBindingPayment(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessSendBindingPaySMS(t *testing.T) {

	be := &model.BindingPayment{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		SettFlag:      settFlag,
		SysOrderNum:   sysOrderNum,
		TransAmt:      int64(transAmt),
		PrivateKey:    priKeyPem,
	}
	resp := DefaultClient.ProcessSendBindingPaySMS(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessPaymentWithSMS(t *testing.T) {
	be := &model.BindingPayment{
		ChanMerId: chanMerId,
		// ChanBindingId: chanBingingId,
		// SettFlag:      settFlag,
		SysOrderNum: "72bd6786041d47917a8d06ee71dfa761",
		// TransAmt:      int64(transAmt),
		PrivateKey: priKeyPem,
		SmsCode:    "123456",
	}
	resp := DefaultClient.ProcessPaymentWithSMS(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessBindingRefund(t *testing.T) {
	be := &model.BindingRefund{
		ChanMerId:       chanMerId,
		SysOrderNum:     sysOrderNum,
		SysOrigOrderNum: sysOrigOrderNum,
		TransAmt:        int64(transAmt),
		PrivateKey:      priKeyPem,
	}
	resp := DefaultClient.ProcessBindingRefund(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessPaymentEnquiry(t *testing.T) {
	be := &model.OrderEnquiry{
		ChanMerId:   "001583",
		SysOrderNum: "729290073920001LibTW1n120151013",
		PrivateKey:  priKeyPem,
		Mode:        model.MarketMode,
	}
	resp := DefaultClient.ProcessPaymentEnquiry(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessRefundEnquiry(t *testing.T) {
	be := &model.OrderEnquiry{
		ChanMerId:   chanMerId,
		SysOrderNum: sysOrderNum,
		PrivateKey:  priKeyPem,
	}
	resp := DefaultClient.ProcessRefundEnquiry(be)
	log.Debugf("response message  %s", resp)
}

func TestProcessTransChecking(t *testing.T) {
	// be := &model.BillingSummary{
	// 	ChanMerId: "001405",
	// 	SettDate:  "2015-03-19",
	// 	PrivateKey:  priKeyPem,
	// }
	resp := DefaultClient.ProcessTransChecking(chanMerId, "2015-04-13", priKeyPem)
	log.Debugf("response message  %#v", resp)
}
