package core

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/channel/cil"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"testing"
	"time"
)

func init() {
	mongo.Connect()
	BuildTree()
	cil.Connect()
}

func TestProcessApplePay(t *testing.T) {

	ap := &model.ApplePay{
		MerId:         applePayMerId,
		TransType:     "SALE",
		MerOrderNum:   fmt.Sprintf("%d", time.Now().UnixNano()),
		TransactionId: fmt.Sprintf("%020d", time.Now().UnixNano()),
		ApplePayData: model.ApplePayData{
			ApplicationPrimaryAccountNumber: testAPPCard,
			ApplicationExpirationDate:       testAPPExpireDate,
			CurrencyCode:                    "156",
			TransactionAmount:               120,
			DeviceManufacturerIdentifier:    "040010030273",
			PaymentDataType:                 "3DSecure",
			PaymentData: model.PaymentData{
				OnlinePaymentCryptogram: "AcqhpcYAIdfgEP3QIUGgMAACAAA",
				EciIndicator:            "5",
			},
		},
	}

	ret := ProcessApplePay(ap)

	if ret == nil {
		t.Error("Apple pay process error")
	}

	t.Logf("%+v", ret)
}
