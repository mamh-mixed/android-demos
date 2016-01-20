package bindingpay

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/log"
)

var sysKey = goconf.Config.App.EncryptKey

// BindingCreateHandle 建立绑定关系
func BindingCreateHandle(data []byte, merId string) (ret *model.BindingReturn) {
	bc := new(model.BindingCreate)
	err := json.Unmarshal(data, bc)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	bc.MerId = merId

	m, _ := mongo.MerchantColl.Find(merId)

	// 解密特定字段
	aes := security.NewAESCBCEncrypt(m.EncryptKey, sysKey)
	bc.AcctNumDecrypt, bc.AcctNum = aes.DcyAndUseSysKeyEcy(bc.AcctNum, sysKey)
	bc.AcctNameDecrypt, bc.AcctName = aes.DcyAndUseSysKeyEcy(bc.AcctName, sysKey)
	bc.IdentNumDecrypt, bc.IdentNum = aes.DcyAndUseSysKeyEcy(bc.IdentNum, sysKey)
	bc.PhoneNumDecrypt, bc.PhoneNum = aes.DcyAndUseSysKeyEcy(bc.PhoneNum, sysKey)
	if bc.AcctType == "20" {
		bc.ValidDateDecrypt, bc.ValidDate = aes.DcyAndUseSysKeyEcy(bc.ValidDate, sysKey)
		bc.Cvv2Decrypt, bc.Cvv2 = aes.DcyAndUseSysKeyEcy(bc.Cvv2, sysKey)
	}
	// 报文解密错误
	if aes.Err != nil {
		log.Errorf("decrypt fail : merId=%s, request=%+v, err=%s", merId, bc, aes.Err)
		return mongo.RespCodeColl.Get("200021")
	}

	log.Debugf("after decrypt field: acctNum=%s, acctName=%s, phoneNum=%s, identNum=%s, validDate=%s, cvv2=%s",
		bc.AcctNumDecrypt, bc.AcctNameDecrypt, bc.PhoneNumDecrypt, bc.IdentNumDecrypt, bc.ValidDateDecrypt, bc.Cvv2Decrypt)

	// 验证请求报文是否完整，格式是否正确
	ret = validateBindingCreate(bc)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBindingCreate(bc)

	return ret
}

// BindingPaymentSettlementHandle 支付结算
func BindingPaymentSettlementHandle(data []byte, merId string) (ret *model.BindingReturn) {

	r := new(model.PaySettlement)
	err := json.Unmarshal(data, r)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	r.MerId = merId

	m, _ := mongo.MerchantColl.Find(merId)
	aes := security.NewAESCBCEncrypt(m.EncryptKey, sysKey)
	r.AcctNameDecrypt, r.SettAccountName = aes.DcyAndUseSysKeyEcy(r.SettAccountName, sysKey)
	r.AcctNumDecrypt, r.SettAccountNum = aes.DcyAndUseSysKeyEcy(r.SettAccountNum, sysKey)
	if aes.Err != nil {
		log.Errorf("decrypt fail : merId=%s, request=%+v, err=%s", merId, r, aes.Err)
		return mongo.RespCodeColl.Get("200021")
	}

	log.Debugf("after decrypt: settAccountName=%s, settAccountNum=%s", r.AcctNameDecrypt, r.AcctNumDecrypt)

	// 验证
	ret = validatePaySettlement(r)
	if ret != nil {
		return ret
	}

	return core.ProcessPaySettlement(r)
}

// GetCardInfoHandle 获取卡片信息
func GetCardInfoHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.CardInfo)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	ret = validateGetCardInfo(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessGetCardInfo(b)
	return ret
}

// BindingRemoveHandle 解除绑定关系
func BindingRemoveHandle(data []byte, merId string) (ret *model.BindingReturn) {
	br := new(model.BindingRemove)
	err := json.Unmarshal(data, br)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	br.MerId = merId

	ret = validateBindingRemove(br)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBindingReomve(br)
	return ret
}

// BindingEnquiryHandle 查询绑定关系
func BindingEnquiryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	be := new(model.BindingEnquiry)
	err := json.Unmarshal(data, be)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	be.MerId = merId

	// 验证请求报文格式
	ret = validateBindingEnquiry(be)
	if ret != nil {
		return ret
	}

	ret = core.ProcessBindingEnquiry(be)

	return ret
}

// BindingPaymentHandle 绑定支付
func BindingPaymentHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBindingPayment(b)
	if ret != nil {
		return ret
	}
	// 直接支付
	ret = core.ProcessBindingPayment(b, false)

	return ret
}

// SendBindingPaySMS 绑定支付发送短信验证码
func SendBindingPaySMS(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateSendBindingPaySMS(b)
	if ret != nil {
		return ret
	}
	// 发送支付短信验证码
	ret = core.ProcessBindingPayment(b, true)

	return ret
}

// BindingPayWithSMS 带验证码的支付
func BindingPayWithSMS(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBindingPayWithSMS(b)
	if ret != nil {
		return ret
	}
	// 带验证码的支付
	ret = core.ProcessPaymentWithSMS(b)

	return ret
}

// BindingRefundHandle 退款处理
func BindingRefundHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingRefund)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBindingRefund(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBindingRefund(b)

	return ret
}

// BillingSummaryHandle 交易对账汇总
func BillingSummaryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BillingSummary)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}

	b.MerId = merId
	// 验证请求报文格式
	ret = validateBillingSummary(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBillingSummary(b)

	// mock return
	return ret
}

// BillingDetailsHandle 交易对账明细
func BillingDetailsHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BillingDetails)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBillingDetails(b)
	if ret != nil {
		return ret
	}
	// 业务处理
	ret = core.ProcessBillingDetails(b)
	// mock return
	return ret
}

// OrderEnquiryHandle 查询订单状态
func OrderEnquiryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.OrderEnquiry)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateOrderEnquiry(b)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	ret = core.ProcessOrderEnquiry(b)

	return ret
}

// NoTrackPaymentHandle 无卡直接支付的处理
func NoTrackPaymentHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.NoTrackPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	m, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		return mongo.RespCodeColl.Get("300030")
	}

	aes := security.NewAESCBCEncrypt(m.EncryptKey, sysKey)
	b.AcctNumDecrypt, b.AcctNum = aes.DcyAndUseSysKeyEcy(b.AcctNum, sysKey)
	// b.AcctNameDecrypt, b.AcctName = aes.DcyAndUseSysKeyEcy(b.AcctName)
	if b.IdentNum != "" {
		b.IdentNumDecrypt, b.IdentNum = aes.DcyAndUseSysKeyEcy(b.IdentNum, sysKey)
	}
	if b.PhoneNum != "" {
		b.PhoneNumDecrypt, b.PhoneNum = aes.DcyAndUseSysKeyEcy(b.PhoneNum, sysKey)
	}

	if b.AcctType == "20" {
		b.ValidDateDecrypt, b.ValidDate = aes.DcyAndUseSysKeyEcy(b.ValidDate, sysKey)
		b.Cvv2Decrypt, b.Cvv2 = aes.DcyAndUseSysKeyEcy(b.Cvv2, sysKey)
	}

	// 报文解密错误
	if aes.Err != nil {
		log.Errorf("decrypt fail : merId=%s, request=%+v, err=%s", merId, b, aes.Err)
		return mongo.RespCodeColl.Get("200021")
	}

	log.Debugf("after decrypt field: acctNum=%s, acctName=%s, phoneNum=%s, identNum=%s, validDate=%s, cvv2=%s",
		b.AcctNumDecrypt, b.AcctNameDecrypt, b.PhoneNumDecrypt, b.IdentNumDecrypt, b.ValidDateDecrypt, b.Cvv2Decrypt)

	ret = validateNoTrackPayment(b)
	if ret != nil {
		return ret
	}
	log.Debugf("请求对象： %+v；校验结果：%+v", b, ret)
	ret = core.ProcessNoTrackPayment(b)

	return ret
}

// ApplePayHandle Apple Pay 支付入口
func ApplePayHandle(data []byte, merId string) (ret *model.BindingReturn) {
	ap := new(model.ApplePay)

	if err := json.Unmarshal(data, ap); err != nil {
		log.Errorf("can't unmarshal `%s` to json: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}

	ap.MerId = merId

	if ret = validateApplePay(ap); ret != nil {
		return ret
	}

	ret = core.ProcessApplePay(ap)
	return ret
}
