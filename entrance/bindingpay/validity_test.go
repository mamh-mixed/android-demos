package bindingpay

import (
	"quickpay/model"
	"testing"

	"github.com/omigo/g"
)

func TestBindingCreateRequestValidity(t *testing.T) {
	var (
		bc  model.BindingCreate
		ret *model.BindingReturn
	)
	ret = bindingCreateRequestValidity(bc)
	if ret == nil {
		g.Error("\n", "验证 '报文要素缺失' 失败")
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

	ret = bindingCreateRequestValidity(bc)
	if ret != nil {
		t.Errorf("%s\n", "验证 '报文正确' 失败")
	}

	bc.PhoneNum = "18205960039"
	bc.AcctType = "20"
	bc.ValidDate = "2013"
	ret = bindingCreateRequestValidity(bc)
	if ret == nil {
		t.Errorf("%s\n", "验证 '卡片有效期有误' 失败")
	}

	bc.ValidDate = "2012"
	bc.Cvv2 = "2345"
	ret = bindingCreateRequestValidity(bc)
	if ret == nil {
		t.Errorf("%s\n", "验证 'CVV2有误' 失败")
	}
}

func TestBindingRemoveRequestValidity(t *testing.T) {
	var (
		in  model.BindingRemove
		ret *model.BindingReturn
	)

	ret = bindingRemoveRequestValidity(in)

	if ret == nil {
		t.Error("测试解除绑定关系报文要素缺失失败")
	}
	t.Logf("%+v", ret)

	in.BindingId = "1000000001"
	ret = bindingRemoveRequestValidity(in)
	if ret != nil {
		t.Error("测试解除绑定关系报文要素缺失失败")
	}
	t.Logf("%+v", ret)
}

func TestBindingPaymentRequestValidity(t *testing.T) {
	var ret *model.BindingReturn
	var in = model.BindingPayment{
		SubMerId:    "",
		MerOrderNum: "1000000003",
		TransAmt:    1000,
		BindingId:   "1000000001",
		SendSmsId:   "",
		SmsCode:     "",
	}

	ret = bindingPaymentRequestValidity(in)

	if ret != nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	}

	in.TransAmt = 0
	ret = bindingPaymentRequestValidity(in)
	if ret == nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	}

	in.TransAmt = 1000
	in.SendSmsId = "100100100"
	ret = bindingPaymentRequestValidity(in)
	if ret == nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	}
}

func TestRefundRequestValidity(t *testing.T) {
	var bp = &model.BindingRefund{
		MerOrderNum:  "1000000004", // 商户订单号
		OrigOrderNum: "1000000003", // 原支付订单号
		TransAmt:     10000,
	}
	var ret *model.BindingReturn

	ret = bindingRefundRequestValidity(bp)
	if ret != nil {
		t.Error("测试退款 '报文要素缺失' 失败")
	}

	bp.TransAmt = -10000
	if ret != nil {
		t.Error("测试退款 ‘退款金额有误’ 失败")
	}

}

func TestNoTrackPaymentRequestValidity(t *testing.T) {
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
	ret := noTrackPaymentRequestValidity(in)
	if ret != nil {
		t.Error("测试无卡支付失败，返回信息： %+v", ret)
	}
}
