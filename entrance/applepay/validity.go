package applepay

import (
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

func validateApplePay(ap *model.ApplePay) (ret *model.BindingReturn) {
	// TransType
	if ap.TransType == "" {
		return model.NewBindingReturn("200050", "字段 transType 不能为空")
	}
	if ap.TransType != "SALE" && ap.TransType != "AUTH" {
		return mongo.RespCodeColl.Get("100030")
	}

	// SubMerId
	if ap.SubMerId != "" && !isAlphanumeric(ap.SubMerId) {
		return model.NewBindingReturn("200051", "字段 subMerId 不符合要求")
	}

	// TerminalId
	if ap.TerminalId != "" && !isAlphanumeric(ap.TerminalId) {
		return model.NewBindingReturn("200051", "字段 terminalId 不符合要求")
	}

	// MerOrderNum
	if ap.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}
	if !isAlphanumeric(ap.MerOrderNum) {
		return model.NewBindingReturn("200051", "字段 merOrderNum 不符合要求")
	}

	// TransactionId
	if ap.TransactionId == "" {
		return model.NewBindingReturn("200050", "字段 transactionId 不能为空")
	}
	// TODO 判断transactonID，目前接口需求不确定，先以线下网关的规定为准，只能纯数字，最多20位
	if matched, _ := regexp.MatchString(`^\d{3,20}$`, ap.TransactionId); !matched {
		return mongo.RespCodeColl.Get("200253")
	}

	// 判断主账号ApplicationPrimaryAccountNumber
	if ap.ApplePayData.ApplicationPrimaryAccountNumber == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.applicationPrimaryAccountNumber 不能为空")
	}
	if matched, _ := regexp.MatchString(`^\d{4,}$`, ap.ApplePayData.ApplicationPrimaryAccountNumber); !matched {
		return mongo.RespCodeColl.Get("200110")
	}

	// ApplePayData.ApplicationExpirationDate
	if ap.ApplePayData.ApplicationExpirationDate == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.applicationExpirationDate 不能为空")
	}
	if matched, _ := regexp.MatchString(`^\d{2}(0[1-9]|1[1-2])(0[1-9]|[1-2][0-9]|3[0-1])$`, ap.ApplePayData.ApplicationExpirationDate); !matched {
		return mongo.RespCodeColl.Get("200250")
	}

	// ApplePayData.CurrencyCode
	if ap.ApplePayData.CurrencyCode == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.currencyCode 不能为空")
	}
	if matched, _ := regexp.MatchString(`^\d{3}$`, ap.ApplePayData.CurrencyCode); !matched {
		return mongo.RespCodeColl.Get("200251")
	}

	// ApplePayData.TransactionAmount
	if ap.ApplePayData.TransactionAmount == 0 {
		return model.NewBindingReturn("200050", "字段 applePayData.transactionAmount 不能为空")
	}

	// ApplePayData.PaymentDataType
	if ap.ApplePayData.PaymentDataType == "" {
		return model.NewBindingReturn("200050", "字段 applePayData.paymentDataType 不能为空")
	}
	if ap.ApplePayData.PaymentDataType != "EMV" && ap.ApplePayData.PaymentDataType != "3DSecure" {
		return mongo.RespCodeColl.Get("200252")
	}

	return nil
}

// isAlphabeticOrNumeric 用来判断一个字符串是否是字母或者数字
func isAlphanumeric(str string) (result bool) {
	matched, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, str)
	if matched {
		return true
	}
	return false
}
