package core

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"
	"testing"
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
	t.Logf("%+v", ret)
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
	g.Debug("%+v", ret)
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
	g.Debug("%+v,%+v", ret, ret.OrigTransDetail)
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
	g.Debug("%+v", ret)
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
	g.Debug("%+v", ret)
}
