package bindingpay

import (
	// "github.com/CardInfoLink/log"
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// validateBindingCreate 建立绑定关系的时候验证请求报文
func validateBindingCreate(request *model.BindingCreate) (ret *model.BindingReturn) {
	if request.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if request.AcctNameDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 acctName 不能为空")
	}

	if request.AcctNumDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 acctNum 不能为空")
	}

	if !isAlphanumeric(request.BindingId) {
		return model.NewBindingReturn("200051", "字段 bindingId 格式错误")
	}

	if err := validateAcctName(request.AcctNameDecrypt); err != nil {
		return err
	}

	if err := validateAcctNum(request.AcctNumDecrypt); err != nil {
		return err
	}

	// IdentType
	if request.IdentType != "" {
		if matched, _ := regexp.MatchString(`^([0-9]|X)$`, request.IdentType); !matched {
			return model.NewBindingReturn("200051", "字段 identType 不符合要求")
		}
	}

	// IdentNum
	if request.IdentNum != "" && request.IdentNumDecrypt != "" {
		if !isAlphanumericOrSpecial(request.IdentNumDecrypt) {
			return model.NewBindingReturn("200051", "字段 identNum 不符合要求")
		}
	}

	if request.AcctType != "10" && request.AcctType != "20" {
		return mongo.RespCodeColl.Get("200230")
	}

	if request.AcctType == "20" {
		// 贷记卡
		if request.ValidDateDecrypt == "" {
			return model.NewBindingReturn("200050", "字段 validDate 不能为空")
		}

		if request.Cvv2Decrypt == "" {
			return model.NewBindingReturn("200050", "字段 cvv2 不能为空")
		}

		// 判断格式，需要使用解密后的参数
		if matched, _ := regexp.MatchString(`^\d{2}(0[1-9]|1[1-2])$`, request.ValidDateDecrypt); !matched {
			return mongo.RespCodeColl.Get("200140")
		}

		if matched, _ := regexp.MatchString(`^\d{3}$`, request.Cvv2Decrypt); !matched {
			return mongo.RespCodeColl.Get("200150")
		}
	}

	return nil
}

// validatePaySettlement 验证支付结算接口字段
func validatePaySettlement(in *model.PaySettlement) (ret *model.BindingReturn) {

	if !isAlphanumeric(in.MerOrderNum) {
		return model.NewBindingReturn("200080", "订单号 merOrderNum 格式错误")
	}

	if err := validateSettOrderNum(in.SettOrderNum); err != nil {
		return err
	}

	if err := validateAmt(in.SettAmt); err != nil {
		return err
	}

	if in.AcctNameDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 settAccountName 不能为空")
	}

	if in.AcctNumDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 settAccountNum 不能为空")
	}

	if err := validateAcctName(in.AcctNameDecrypt); err != nil {
		return err
	}

	if err := validateAcctNum(in.AcctNumDecrypt); err != nil {
		return err
	}

	switch in.SettAccountType {
	case "11", "12":
		if in.Province == "" || in.City == "" || in.SettBranchName == "" {
			return model.NewBindingReturn("200050", "分支行信息不完整")
		}
	case "20":
	default:
		return model.NewBindingReturn("200050", "账户类型 settAccountType 取值错误")
	}

	return nil
}

// validateGetCardInfo 验证获取卡片接口字段
func validateGetCardInfo(in *model.CardInfo) (ret *model.BindingReturn) {

	if !isAlphanumericOrSpecial(in.CardNum) {
		return mongo.RespCodeColl.Get("200110")
	}
	return nil
}

// validateBindingRemove 移除绑定关系的时候验证请求报文
func validateBindingRemove(in *model.BindingRemove) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if !isAlphanumeric(in.BindingId) {
		return model.NewBindingReturn("200051", "字段 bindingId 格式错误")
	}

	return nil
}

// validateBindingEnquiry 查询绑定关系的时候验证请求报文
func validateBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	if be.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if !isAlphanumeric(be.BindingId) {
		return model.NewBindingReturn("200051", "字段 bindingId 格式错误")
	}

	return nil
}

// validateBindingPayment 绑定支付的请求报文验证
func validateBindingPayment(in *model.BindingPayment) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if err := validateAmt(in.TransAmt); err != nil {
		return err
	}

	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}

	if !isAlphanumeric(in.BindingId) {
		return model.NewBindingReturn("200051", "字段 bindingId 格式错误")
	}

	if !isAlphanumeric(in.MerOrderNum) {
		return model.NewBindingReturn("200080", "订单号 merOrderNum 格式错误")
	}

	if in.TerminalId != "" && !isAlphanumeric(in.TerminalId) {
		return model.NewBindingReturn("200051", "字段 terminalId 格式错误")
	}

	// 验证短信验证码是否填写
	if in.SendSmsId != "" && in.SmsCode == "" {
		return model.NewBindingReturn("200050", "字段 smsCode 不能为空")
	}

	if in.SettOrderNum != "" {
		if err := validateSettOrderNum(in.SettOrderNum); err != nil {
			return err
		}
	}

	return nil
}

// validateSendBindingPaySMS 交易发送短信接口报文验证
func validateSendBindingPaySMS(in *model.BindingPayment) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if err := validateAmt(in.TransAmt); err != nil {
		return err
	}

	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}

	if !isAlphanumeric(in.BindingId) {
		return model.NewBindingReturn("200051", "字段 bindingId 格式错误")
	}

	if !isAlphanumeric(in.MerOrderNum) {
		return model.NewBindingReturn("200080", "订单号 merOrderNum 格式错误")
	}

	if in.SettOrderNum != "" {
		if err := validateSettOrderNum(in.SettOrderNum); err != nil {
			return err
		}
	}

	return nil
}

// validateBindingPayWithSMS 带验证码交易报文验证
func validateBindingPayWithSMS(in *model.BindingPayment) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}
	if in.SmsCode == "" {
		return model.NewBindingReturn("200050", "字段 smsCode 不能为空")
	}
	if !isAlphanumeric(in.MerOrderNum) {
		return model.NewBindingReturn("200080", "订单号 merOrderNum 格式错误")
	}
	if in.SettOrderNum != "" {
		if err := validateSettOrderNum(in.SettOrderNum); err != nil {
			return err
		}
	}
	// TODO:短信验证码格式验证
	return nil
}

// validateBindingRefund 退款请求报文验证
func validateBindingRefund(in *model.BindingRefund) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}

	if in.OrigOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 origOrderNum 不能为空")
	}

	if !isAlphanumeric(in.OrigOrderNum) {
		return model.NewBindingReturn("200080", "订单号 origOrderNum 格式错误")
	}

	if !isAlphanumeric(in.MerOrderNum) {
		return model.NewBindingReturn("200080", "订单号 merOrderNum 格式错误")
	}

	if in.SettOrderNum != "" {
		if err := validateSettOrderNum(in.SettOrderNum); err != nil {
			return err
		}
	}

	if err := validateAmt(in.TransAmt); err != nil {
		return err
	}

	return nil
}

// validateOrderEnquiry 订单查询报文验证
func validateOrderEnquiry(in *model.OrderEnquiry) (ret *model.BindingReturn) {
	if in.OrigOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 origOrderNum 不能为空")
	}

	if !isAlphanumeric(in.OrigOrderNum) {
		return model.NewBindingReturn("200080", "订单号 origOrderNum 格式错误")
	}

	if matched, _ := regexp.MatchString(`^[1|0]?$`, in.ShowOrigInfo); !matched {
		return model.NewBindingReturn("200050", "字段 showOrigInfo 取值错误")
	}
	return
}

// validateBillingSummary 交易对账汇总验证
func validateBillingSummary(in *model.BillingSummary) (ret *model.BindingReturn) {
	if matched, _ := regexp.MatchString(`^[1-2][0-9][0-9][0-9]-(0[1-9]|1[0-2])-[0-3]{0,1}[0-9]$`, in.SettDate); !matched {
		return model.NewBindingReturn("200200", "日期 settDate 格式错误")
	}
	return
}

// validateBillingSummary 交易对账汇总验证
func validateBillingDetails(in *model.BillingDetails) (ret *model.BindingReturn) {
	if matched, _ := regexp.MatchString(`^[1-2][0-9][0-9][0-9]-(0[1-9]|1[0-2])-[0-3]{0,1}[0-9]$`, in.SettDate); !matched {
		return model.NewBindingReturn("200200", "日期 settDate 格式错误")
	}

	if in.NextOrderNum != "" && !isAlphanumeric(in.NextOrderNum) {
		return model.NewBindingReturn("200080", "订单号 nextOrderNum 格式错误")
	}
	return
}

// validateNoTrackPayment 无卡直接支付请求报文验证
func validateNoTrackPayment(in *model.NoTrackPayment) (ret *model.BindingReturn) {
	// TransType
	if in.TransType == "" {
		return model.NewBindingReturn("200050", "字段 transType 不能为空")
	}
	if matched, _ := regexp.MatchString(`^SALE$|^AUTH$`, in.TransType); !matched {
		return mongo.RespCodeColl.Get("100030")
	}

	// SubMerId
	if in.SubMerId != "" && !isAlphanumeric(in.SubMerId) {
		return model.NewBindingReturn("200051", "字段 subMerId 不符合要求")
	}

	// TerminalId
	if in.TerminalId != "" && !isAlphanumeric(in.TerminalId) {
		return model.NewBindingReturn("200051", "字段 terminalId 不符合要求")
	}

	// MerOrderNum
	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}
	if !isAlphanumeric(in.MerOrderNum) {
		return model.NewBindingReturn("200051", "字段 merOrderNum 不符合要求")
	}

	if err := validateAmt(in.TransAmt); err != nil {
		return err
	}

	// CurrCode
	if in.CurrCode != "" {
		// 判断交易币种格式
		if matched, _ := regexp.MatchString(`^\d{3}$`, in.CurrCode); !matched {
			return mongo.RespCodeColl.Get("200251")
		}
	}

	// AcctName
	// if in.AcctName == "" || in.AcctNameDecrypt == "" {
	// 	return model.NewBindingReturn("200050", "字段 acctName 不能为空")
	// }
	// if !isChineseOrJapaneseOrAlphanumeric(in.AcctNameDecrypt) {
	// 	return mongo.RespCodeColl.Get("200100")
	// }

	// AcctNum
	if in.AcctNum == "" || in.AcctNumDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 acctNum 不能为空")
	}
	if !isAlphanumericOrSpecial(in.AcctNumDecrypt) {
		return mongo.RespCodeColl.Get("200110")
	}

	// IdentType
	if in.IdentType != "" {
		if matched, _ := regexp.MatchString(`^([0-9]|X)$`, in.IdentType); !matched {
			return model.NewBindingReturn("200051", "字段 identType 不符合要求")
		}
	}

	// IdentNum
	if in.IdentNum != "" && in.IdentNumDecrypt != "" {
		if !isAlphanumericOrSpecial(in.IdentNumDecrypt) {
			return model.NewBindingReturn("200051", "字段 identNum 不符合要求")
		}
	}

	// PhoneNum
	if in.PhoneNum != "" && in.PhoneNumDecrypt != "" {
		if matched, _ := regexp.MatchString(`^\d+$`, in.PhoneNumDecrypt); !matched {
			return model.NewBindingReturn("200051", "字段 phoneNum 不符合要求")
		}
	}

	// AcctType
	if in.AcctType == "" {
		return model.NewBindingReturn("200050", "字段 acctType 不能为空")
	}
	if matched, _ := regexp.MatchString(`^10$|^20$`, in.AcctType); !matched {
		return mongo.RespCodeColl.Get("200230")
	}
	if in.AcctType == "20" {
		// ValidDate
		if in.ValidDate == "" || in.ValidDateDecrypt == "" {
			return model.NewBindingReturn("200050", "字段 validDate 不能为空")
		}
		if matched, _ := regexp.MatchString(`^\d{2}(0[1-9]|1[1-2])$`, in.ValidDateDecrypt); !matched {
			return mongo.RespCodeColl.Get("200140")
		}

		// Cvv2
		if in.Cvv2 == "" || in.Cvv2Decrypt == "" {
			return model.NewBindingReturn("200050", "字段 cvv2 不能为空")
		}
		if matched, _ := regexp.MatchString(`^\d{3}$`, in.Cvv2Decrypt); !matched {
			return mongo.RespCodeColl.Get("200150")
		}
	}

	return nil
}

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

func validateAcctNum(acctNum string) *model.BindingReturn {
	if len(acctNum) > 32 {
		return mongo.RespCodeColl.Get("200110")
	}
	if !isAlphanumericOrSpecial(acctNum) {
		return mongo.RespCodeColl.Get("200110")
	}
	return nil
}

func validateAcctName(acctName string) *model.BindingReturn {

	runeAcctName := []rune(acctName)
	if len(runeAcctName) > 64 {
		return mongo.RespCodeColl.Get("200100")
	}
	if !isChineseOrJapaneseOrAlphanumeric(acctName) {
		return mongo.RespCodeColl.Get("200100")
	}
	return nil
}

func validateSettOrderNum(settOrderNum string) *model.BindingReturn {

	if len(settOrderNum) > 16 {
		return model.NewBindingReturn("200080", "结算订单号 settOrderNum 长度过长")
	}
	if !isAlphanumeric(settOrderNum) {
		return model.NewBindingReturn("200080", "结算订单号 settOrderNum 格式错误")
	}
	return nil
}

func validateAmt(amt int64) *model.BindingReturn {
	if amt <= 0 {
		return model.NewBindingReturn("200051", "金额过小")
	}
	if amt > 1000000000 {
		return model.NewBindingReturn("200051", "金额过大")
	}
	return nil
}

// isAlphabeticOrNumeric 用来判断一个字符串是否是字母或者数字
func isAlphanumeric(str string) (result bool) {
	matched, _ := regexp.MatchString(`^[A-Za-z0-9-]{0,32}$`, str)
	if matched {
		return true
	}
	return false
}

// isAlphanumeric 用来判断一个字符串是否是字母或者数字或者特殊字符
func isAlphanumericOrSpecial(str string) (result bool) {
	matched, _ := regexp.MatchString(`^(?i)[a-z0-9_+-\\·]+$`, str)
	if matched {
		return true
	}
	return false
}

// isChineseOrJapaneseOrAlphanumeric 用来判断一个字符串是否只包含汉字，日本字或者字母数字或者特殊字符
func isChineseOrJapaneseOrAlphanumeric(str string) (result bool) {
	matched, _ := regexp.MatchString(`^(?i)(\p{Han}|\p{Hiragana}|[a-z0-9]|[_+-\\·])+$`, str)
	if matched {
		return true
	}
	return false
}
