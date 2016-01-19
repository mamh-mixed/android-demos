package bindingpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/CardInfoLink/log"
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

func TestValidateNoTrackPayment(t *testing.T) {
	var in = &model.NoTrackPayment{
		TransType:        "SALE",
		SubMerId:         "",
		MerOrderNum:      "Wa1000000008",
		TransAmt:         10000,
		AcctNameDecrypt:  "张三",
		AcctNum:          "6210948000000219",
		AcctNumDecrypt:   "6210948000000219",
		IdentType:        "0",
		IdentNum:         "350583199009153732",
		IdentNumDecrypt:  "350583199009153732",
		PhoneNum:         "18205960039",
		PhoneNumDecrypt:  "18205960039",
		AcctType:         "20",
		ValidDate:        "0919",
		ValidDateDecrypt: "0125",
		Cvv2Decrypt:      "123",
		Cvv2:             "123",
		CurrCode:         "156",
		SendSmsId:        "",
		SmsCode:          "",
	}
	ret := validateNoTrackPayment(in)
	if ret != nil {
		t.Error("测试无卡支付失败，返回信息： %+v", ret)
	}

}

func TestIsChineseOrJapaneseOrAlphanumeric(t *testing.T) {
	result := isChineseOrJapaneseOrAlphanumeric("zhangsan")

	if !result {
		t.Error("哪里出错啦")
	}

	result = isChineseOrJapaneseOrAlphanumeric("姚明")

	if !result {
		t.Error("出错啦")
	}

	result = isChineseOrJapaneseOrAlphanumeric("姚明isBest強い")

	if !result {
		t.Error("出错啦")
	}

	result = isChineseOrJapaneseOrAlphanumeric("alp·json1233+-_")

	if !result {
		t.Error("出错啦")
	}

	result = isAlphanumericOrSpecial("44581312312323+-_aa")

	if !result {
		t.Error("出错啦")
	}
}
