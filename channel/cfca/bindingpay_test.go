package cfca

import (
	"quickpay/model"
	"testing"

	"github.com/omigo/g"
)

func TestSendRequest(t *testing.T) {

	req := &BindingRequest{
		Version: "2.0",
		Head: requestHead{
			InstitutionID: "001405", //测试ID
			TxCode:        "2502",
		},
		Body: requestBody{
			TxSNBinding: "213121231313",
		},
	}

	response := sendRequest(req, priKeyPem)
	if response == nil {
		t.Error("test unsucessful")
	}
}

func TestProcessBindingEnquiry(t *testing.T) {

	be := &model.BindingEnquiry{
		ChanBindingId: "123456789",
		ChanMerId:     "001405",
		SignCert:      priKeyPem,
	}
	resp := ProcessBindingEnquiry(be)

	g.Debug("response message  %s", resp)
}

func TestProcessBindingCreate(t *testing.T) {

	be := &model.BindingCreate{
		ChanMerId:     "001405",
		ChanBindingId: "123456789011",
		BankId:        "102",
		AcctName:      "test",
		AcctNum:       "6222022003008481261",
		IdentType:     "0",
		IdentNum:      "440583199111031012",
		PhoneNum:      "15618103236",
		AcctType:      "10",
		ValidDate:     "",
		Cvv2:          "",
		SendSmsId:     "",
		SmsCode:       "",
		SignCert:      priKeyPem,
	}
	resp := ProcessBindingCreate(be)
	g.Debug("response message  %s", resp)
}

func TestProcessBindingRemove(t *testing.T) {

	be := &model.BindingRemove{
		ChanMerId:     "001405",
		ChanBindingId: "123456789",
		TxSNUnBinding: "3333444",
		SignCert:      priKeyPem,
	}
	resp := ProcessBindingRemove(be)
	g.Debug("response message  %s", resp)
}

func TestProcessBindingPayment(t *testing.T) {

	be := &model.BindingPayment{
		ChanMerId:     "001405",
		ChanBindingId: "123456789011",
		SettFlag:      "475",
		//需要变化
		MerOrderNum: "6222022003008481260",
		TransAmt:    12000,
		SignCert:    priKeyPem,
	}
	resp := ProcessBindingPayment(be)
	g.Debug("response message  %s", resp)
}

func TestProcessBindingRefund(t *testing.T) {

	be := &model.BindingRefund{
		ChanMerId: "001405",
		// BindingId:     "1234567890",
		//需要变化
		MerOrderNum:  "6222022003008481265",
		OrigOrderNum: "6222022003008481264",
		TransAmt:     12000,
		SignCert:     priKeyPem,
	}
	resp := ProcessBindingRefund(be)
	g.Debug("response message  %s", resp)
}
