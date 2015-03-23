package cfca

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"

	"github.com/omigo/g"
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
	resp := ProcessBindingEnquiry(be)

	g.Debug("response message  %s", resp)
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
	resp := ProcessBindingCreate(be)
	g.Debug("response message  %s", resp)
}

func TestProcessBindingRemove(t *testing.T) {

	be := &model.BindingRemove{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		TxSNUnBinding: txSNUnBinding,
		SignCert:      priKeyPem,
	}
	resp := ProcessBindingRemove(be)
	g.Debug("response message  %s", resp)
}

func TestProcessBindingPayment(t *testing.T) {

	be := &model.BindingPayment{
		ChanMerId:     chanMerId,
		ChanBindingId: chanBingingId,
		SettFlag:      settFlag,
		//需要变化
		ChanOrderNum: chanOrderNum,
		TransAmt:     int64(transAmt),
		SignCert:     priKeyPem,
	}
	resp := ProcessBindingPayment(be)
	g.Debug("response message  %s", resp)
}

func TestProcessBindingRefund(t *testing.T) {

	be := &model.BindingRefund{
		ChanMerId:        chanMerId,
		ChanOrderNum:     chanOrderNum,
		ChanOrigOrderNum: chanOrigOrderNum,
		TransAmt:         int64(transAmt),
		SignCert:         priKeyPem,
	}
	resp := ProcessBindingRefund(be)
	g.Debug("response message  %s", resp)
}

func TestProcessPaymentEnquiry(t *testing.T) {

	be := &model.OrderEnquiry{
		ChanMerId:    chanMerId,
		ChanOrderNum: chanOrderNum,
		SignCert:     priKeyPem,
	}
	resp := ProcessPaymentEnquiry(be)
	g.Debug("response message  %s", resp)
}

func TestProcessRefundEnquiry(t *testing.T) {

	be := &model.OrderEnquiry{
		ChanMerId:    chanMerId,
		ChanOrderNum: chanOrderNum,
		SignCert:     priKeyPem,
	}
	resp := ProcessRefundEnquiry(be)
	g.Debug("response message  %s", resp)
}

func TestProcessTransChecking(t *testing.T) {

	// be := &model.BillingSummary{
	// 	ChanMerId: "001405",
	// 	SettDate:  "2015-03-19",
	// 	SignCert:  priKeyPem,
	// }
	resp := ProcessTransChecking(chanMerId, "2015-03-20", priKeyPem)
	g.Debug("response message  %s", resp)
}
