package bindingpay

import (
	"quickpay/model"
	"testing"

	"github.com/omigo/g"
)

func TestBindingCreateRequestValidity(t *testing.T) {
	request := model.BindingCreateIn{}
	code, msg := bindingCreateRequestValidity(request)
	if code != "200050" {
		g.Error("\n", "验证 '报文要素缺失' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}
	request.BindingId = "1000000000001"
	request.AcctName = "张三"
	request.AcctNum = "6210948000000219"
	request.IdentType = "0"
	request.IdentNum = "36050219880401"
	request.PhoneNum = "15600009909"
	request.AcctType = "20"
	request.ValidDate = "1903"
	request.Cvv2 = "232"
	request.SendSmsId = "1000000000009"
	request.SmsCode = "12353"

	code, msg = bindingCreateRequestValidity(request)
	if code != "00" {
		t.Errorf("%s\n", "验证 '报文正确' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.IdentType = "XXX"
	code, msg = bindingCreateRequestValidity(request)
	if code != "200120" {
		t.Errorf("%s\n", "验证 '证件类型有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.IdentType = "0"
	request.PhoneNum = "wonsikin"
	code, msg = bindingCreateRequestValidity(request)
	if code != "200130" {
		t.Errorf("%s\n", "验证 '手机号有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.PhoneNum = "18205960039"
	request.AcctType = "20"
	request.ValidDate = "2013"
	code, msg = bindingCreateRequestValidity(request)
	if code != "200140" {
		t.Errorf("%s\n", "验证 '卡片有效期有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.ValidDate = "2012"
	request.Cvv2 = "2345"
	code, msg = bindingCreateRequestValidity(request)
	if code != "200150" {
		t.Errorf("%s\n", "验证 'CVV2有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}
}

func TestBindingRemoveRequestValidity(t *testing.T) {
	var (
		in   model.BindingRemoveIn
		code string
		err  error
	)

	code, err = bindingRemoveRequestValidity(in)
	if err == nil {
		t.Error("测试解除绑定关系报文要素缺失失败")
	} else {
		t.Logf("%s", code)
	}

	in.BindingId = "1000000001"
	code, err = bindingRemoveRequestValidity(in)
	if err != nil {
		t.Error("测试解除绑定关系报文要素缺失失败")
	} else {
		t.Logf("%s", code)
	}
}

func TestBindingPaymentRequestValidity(t *testing.T) {
	var (
		code string
		err  error
	)
	in := model.BindingPaymentIn{
		SubMerId:    "",
		MerOrderNum: "1000000003",
		TransAmt:    1000,
		BindingId:   "1000000001",
		SendSmsId:   "",
		SmsCode:     "",
	}

	code, err = bindingPaymentRequestValidity(in)

	if err != nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	} else {
		if code != "00" {
			t.Errorf("测试失败，应该返回‘00’，但是返回 %s", code)
		} else {
			t.Logf("%s", code)
		}
	}

	in.TransAmt = 0
	code, err = bindingPaymentRequestValidity(in)
	if err == nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	} else {
		t.Logf("%s", code)
	}

	in.TransAmt = 1000
	in.SendSmsId = "100100100"
	code, err = bindingPaymentRequestValidity(in)
	if err == nil {
		t.Error("测试绑定支付 '报文要素缺失' 失败")
	} else {
		t.Logf("%s", code)
	}
}

func TestRefundRequestValidity(t *testing.T) {

}
