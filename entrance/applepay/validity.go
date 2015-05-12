package applepay

import (
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

func validateApplePay(ap *model.ApplePay) (ret *model.BindingReturn) {
	if ap.TransType == "" {
		return model.NewBindingReturn("200050", "字段 transType 不能为空")
	}
	if ap.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}
	if ap.TransactionId == "" {
		return model.NewBindingReturn("200050", "字段 transactionId 不能为空")
	}
	// 判断主账号
	if ap.ApplePayData.ApplicationPrimaryAccountNumber == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.applicationPrimaryAccountNumber 不能为空")
	}
	if ap.ApplePayData.ApplicationExpirationDate == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.applicationExpirationDate 不能为空")
	}
	if ap.ApplePayData.CurrencyCode == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.currencyCode 不能为空")
	}
	if ap.ApplePayData.TransactionAmount == 0 {
		return model.NewBindingReturn("200050", "字段 applePayData.transactionAmount 不能为空")
	}
	if ap.ApplePayData.PaymentDataType == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.paymentDataType 不能为空")
	}

	if matched, _ := regexp.MatchString(`^\d{4,}$`, ap.ApplePayData.ApplicationPrimaryAccountNumber); !matched {
		return mongo.RespCodeColl.Get("200110")
	}

	if matched, _ := regexp.MatchString(`^\d{2}(0[1-9]|1[1-2])(0[1-9]|[1-2][0-9]|3[0-1])$`, ap.ApplePayData.ApplicationExpirationDate); !matched {
		return mongo.RespCodeColl.Get("200250")
	}

	if matched, _ := regexp.MatchString(`^\d{3}$`, ap.ApplePayData.CurrencyCode); !matched {
		return mongo.RespCodeColl.Get("200251")
	}

	if ap.ApplePayData.PaymentDataType != "EMV" && ap.ApplePayData.PaymentDataType != "3DSecure" {
		return mongo.RespCodeColl.Get("200252")
	}

	if ap.TransType != "SALE" && ap.TransType != "AUTH" {
		return mongo.RespCodeColl.Get("100030")
	}
	return nil
}
