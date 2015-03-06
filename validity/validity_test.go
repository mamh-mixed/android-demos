package validity

import (
	"github.com/omigo/g"
	"quickpay/domain"
	"testing"
)

func TestBindingCreateRequestValidity(t *testing.T) {
	var request domain.BindingCreateRequest
	code, msg := BindingCreateRequestValidity(request)
	if code != "200050" {
		g.Error("\n", "验证 '报文要素缺失' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}
	request.MerBindingId = "1000000000001"
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

	code, msg = BindingCreateRequestValidity(request)
	if code != "00" {
		g.Error("\n", "验证 '报文正确' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.IdentType = "XXX"
	code, msg = BindingCreateRequestValidity(request)
	if code != "200120" {
		g.Error("\n", "验证 '证件类型有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.IdentType = "0"
	request.PhoneNum = "wonsikin"
	code, msg = BindingCreateRequestValidity(request)
	if code != "200130" {
		g.Error("\n", "验证 '手机号有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.PhoneNum = "18205960039"
	request.AcctType = "20"
	request.ValidDate = "2013"
	code, msg = BindingCreateRequestValidity(request)
	if code != "200140" {
		g.Error("\n", "验证 '卡片有效期有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}

	request.ValidDate = "2012"
	request.Cvv2 = "2345"
	code, msg = BindingCreateRequestValidity(request)
	if code != "200150" {
		g.Error("\n", "验证 'CVV2有误' 失败")
	} else {
		g.Info("%s---%s", code, msg)
	}
}
