package core

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"
)

func TestProcessBindingCreate(t *testing.T) {
	bc := &model.BindingCreate{
		MerId:     "1426727710113",
		BindingId: "1000000000011",
		AcctName:  "张三",
		AcctNum:   "6222020302062061908",
		IdentType: "0",
		IdentNum:  "350583199009153732",
		PhoneNum:  "18205960039",
		AcctType:  "10",
		ValidDate: "",
		Cvv2:      "",
		SendSmsId: "1000000000009",
		SmsCode:   "12353",
		BankId:    "102",
	}

	ret := ProcessBindingCreate(bc)

	t.Logf("%+v", ret)

	if ret.RespCode != "000000" {
		t.Errorf("Excpeted success,but get failure: %+v", ret.RespMsg)
	}
}

func TestProcessBindingEnquiry(t *testing.T) {
	be := &model.BindingEnquiry{
		MerId:     "MI12345",
		BindingId: "BI1426664578581",
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
		MerId:     "1426562901844",
		BindingId: "1426562901897",
	}

	ret := ProcessBindingReomve(be)

	t.Logf("%+v", ret)
}

func TestProcessBindingPayment(t *testing.T) {

	be := &model.BindingPayment{
		MerId:       "001405",
		BindingId:   "1000000000011",
		TransAmt:    800,
		MerOrderNum: "20000000010000000",
	}
	ret := ProcessBindingPayment(be)

	if ret.RespCode == "" {
		t.Errorf("process payment but not get a respCode %+v", ret)
	}
	t.Logf("%+v", ret)
}

func TestProcessBindingRefund(t *testing.T) {

	be := &model.BindingRefund{
		MerId:        "001405",
		TransAmt:     800,
		OrigOrderNum: "20000000010000000",
		MerOrderNum:  "300000200000300",
	}
	ret := ProcessBindingRefund(be)

	if ret.RespCode == "" {
		t.Errorf("process payment but not get a respCode %+v", ret)
	}
	g.Debug("%+v", ret)
}

func TestProcessOrderEnquiry(t *testing.T) {

	be := &model.OrderEnquiry{
		MerId:        "001405",
		OrigOrderNum: "20000000002000000",
		ShowOrigInfo: "1",
	}
	ret := ProcessOrderEnquiry(be)

	if ret.RespCode == "" {
		t.Errorf("process order query but not get a respCode %+v", ret)
	}
	g.Debug("%+v", ret)
}

func TestProcessBillingDetails(t *testing.T) {

	be := &model.BillingDetails{
		MerId:    "001405",
		SettDate: "2015-03-22",
	}
	ret := ProcessBillingDetails(be)

	if ret.RespCode == "" {
		t.Errorf("process billing details but not get a respCode %+v", ret)
	}
	g.Debug("%+v", ret)
}

func TestProcessBillingSummary(t *testing.T) {

	be := &model.BillingSummary{
		MerId:    "001405",
		SettDate: "2015-03-21",
	}
	ret := ProcessBillingSummary(be)

	if ret.RespCode == "" {
		t.Errorf("process billing summary but not get a respCode %+v", ret)
	}
	g.Debug("%+v", ret)
}
