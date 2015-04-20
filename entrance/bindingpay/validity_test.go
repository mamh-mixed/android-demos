package bindingpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
)

func validateTestBindingCreate(t *testing.T) {
	var (
		bc  *model.BindingCreate
		ret *model.BindingReturn
	)
	ret = validateBindingCreate(bc)
	if ret == nil {
		log.Errorf("\n", "验证 '报文要素缺失' 失败")
	}

	bc.BindingId = "1000000000001"
	bc.AcctName = "张三"
	bc.AcctNum = "6210948000000219"
	bc.IdentType = "0"
	bc.IdentNum = "36050219880401"
	bc.PhoneNum = "15600009909"
	bc.AcctType = "20"
	bc.ValidDate = "1903"
	bc.Cvv2 = "232"
	bc.SendSmsId = "1000000000009"
	bc.SmsCode = "12353"

	ret = validateBindingCreate(bc)
	if ret != nil {
		t.Errorf("%s\n", "验证 '报文正确' 失败")
	}

	bc.PhoneNum = "18205960039"
	bc.AcctType = "20"
	bc.ValidDate = "2013"
	ret = validateBindingCreate(bc)
	if ret == nil {
		t.Errorf("%s\n", "验证 '卡片有效期有误' 失败")
	}

	bc.ValidDate = "2012"
	bc.Cvv2 = "2345"
	ret = validateBindingCreate(bc)
	if ret == nil {
		t.Errorf("%s\n", "验证 'CVV2有误' 失败")
	}
}

func validateTestBindingRemove(t *testing.T) {
	var (
		in  *model.BindingRemove
		ret *model.BindingReturn
	)

	ret = validateBindingRemove(in)

	if ret == nil {
		t.Error("测试解除绑定关系报文要素缺失失败")
	}
	t.Logf("%+v", ret)

	in.BindingId = "1000000001"
	ret = validateBindingRemove(in)
	if ret != nil {
		t.Error("测试解除绑定关系报文要素缺失失败")
	}
	t.Logf("%+v", ret)
}

func validateTestBindingPayment(t *testing.T) {
	var ret *model.BindingReturn
	var in = &model.BindingPayment{
		SubMerId:    "",
		MerOrderNum: "1000000003",
		TransAmt:    1000,
		BindingId:   "1000000001",
		SendSmsId:   "",
		SmsCode:     "",
	}

	ret = validateBindingPayment(in)

	if ret != nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	}

	in.TransAmt = 0
	ret = validateBindingPayment(in)
	if ret == nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	}

	in.TransAmt = 1000
	in.SendSmsId = "100100100"
	ret = validateBindingPayment(in)
	if ret == nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	}
}

func validateTestRefund(t *testing.T) {
	var bp = &model.BindingRefund{
		MerOrderNum:  "1000000004", // 商户订单号
		OrigOrderNum: "1000000003", // 原支付订单号
		TransAmt:     10000,
	}
	var ret *model.BindingReturn

	ret = validateBindingRefund(bp)
	if ret != nil {
		t.Error("测试退款 '报文要素缺失' 失败")
	}

	bp.TransAmt = -10000
	if ret != nil {
		t.Error("测试退款 ‘退款金额有误’ 失败")
	}

}

func TestValidateOrderEnquiry(t *testing.T) {
	oe := &model.OrderEnquiry{
		OrigOrderNum: "200000000",
		ShowOrigInfo: "3",
	}

	ret := validateOrderEnquiry(oe)
	if ret != nil {
		t.Errorf("测试订单查询失败 %+v", ret)
	}
}

func TestValidateBillingSummary(t *testing.T) {
	oe := &model.BillingSummary{
		SettDate: "2015-13-03",
	}

	ret := validateBillingSummary(oe)
	if ret != nil {
		t.Errorf("测试对账汇总报文验证失败 %+v", ret)
	}
}

func TestValidateBillingDetails(t *testing.T) {
	oe := &model.BillingDetails{
		SettDate:     "2015-01-03",
		NextOrderNum: "2222222222222000000",
	}
	ret := validateBillingDetails(oe)
	if ret != nil {
		t.Errorf("测试交易明细报文失败 %+v", ret)
	}
}

func validateTestNoTrackPayment(t *testing.T) {
	var in = &model.NoTrackPayment{
		SubMerId:    "",
		MerOrderNum: "1000000008",
		TransAmt:    10000,
		AcctName:    "张三",
		AcctNum:     "6210948000000219",
		IdentType:   "",
		IdentNum:    "",
		PhoneNum:    "",
		AcctType:    "20",
		ValidDate:   "",
		Cvv2:        "",
		SendSmsId:   "",
		SmsCode:     "",
	}
	ret := validateNoTrackPayment(in)
	if ret != nil {
		t.Error("测试无卡支付失败，返回信息： %+v", ret)
	}
}

func TestValidateApplyPay(t *testing.T) {
	var ap = &model.ApplePay{
		TransType:     "SALE",
		MerOrderNum:   "100000000018",
		TransactionId: "49170302b04f74b56b0060f33e11a135134e48e8af80a50cefea6c079353b419",
		ApplePayData: model.ApplePayData{
			ApplicationPrimaryAccountNumber: "4097900050058723",
			ApplicationExpirationDate:       "200228",
			CurrencyCode:                    "840",
			TransactionAmount:               120,
			DeviceManufacturerIdentifier:    "040010030273",
			PaymentDataType:                 "3DSecure",
			PaymentData: model.PaymentData{
				OnlinePaymentCryptogram: "AcqhpcYAIdfgEP3QIUGgMAACAAA",
				EciIndicator:            "5",
			},
		},
	}

	ret := validateApplePay(ap)

	if ret != nil {
		t.Errorf("验证apple pay 数据失败", ret)
	}

	ap.TransType = "wsj"
	ret = validateApplePay(ap)
	if ret == nil {
		t.Error("验证apple pay 数据失败")
	}
	t.Logf("%+v\n", ret)

	ap.ApplePayData.ApplicationPrimaryAccountNumber = "123"
	ret = validateApplePay(ap)
	if ret == nil {
		t.Error("验证apple pay 数据失败")
	}

	t.Logf("%+v\n", ret)
}
