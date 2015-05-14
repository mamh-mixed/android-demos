package applepay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
)

func TestValidateApplyPay(t *testing.T) {
	var ap = &model.ApplePay{
		TransType:     "SALE",
		MerOrderNum:   "100000000018",
		TransactionId: "49170302047456006041",
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

func TestIsAlphabeticOrNumeric(t *testing.T) {
	result := isAlphanumeric("abc123")

	if !result {
		t.Error("出错啦")
	}

	result = isAlphanumeric("<script>alert()</script>")
	if result {
		t.Error("出错啦")
	}

}
