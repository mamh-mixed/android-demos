package bindingpay

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"

	"github.com/omigo/log"
)

// 建立绑定关系
func BindingCreateHandle(data []byte, merId string) (ret *model.BindingReturn) {
	bc := new(model.BindingCreate)
	err := json.Unmarshal(data, bc)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	bc.MerId = merId

	m, _ := mongo.MerchantColl.Find(merId)

	// 解密特定字段
	aes := new(tools.AesCBCMode)
	aes.DecodeKey(m.EncryptKey)
	bc.AcctNumDecrypt = aes.Decrypt(bc.AcctNum)
	bc.AcctNameDecrypt = aes.Decrypt(bc.AcctName)
	bc.IdentNumDecrypt = aes.Decrypt(bc.IdentNum)
	bc.PhoneNumDecrypt = aes.Decrypt(bc.PhoneNum)
	if bc.AcctType == "20" {
		bc.ValidDateDecrypt = aes.Decrypt(bc.ValidDate)
		bc.Cvv2Decrypt = aes.Decrypt(bc.Cvv2)
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

// 解除绑定关系
func BindingRemoveHandle(data []byte, merId string) (ret *model.BindingReturn) {
	br := new(model.BindingRemove)
	err := json.Unmarshal(data, br)
	if err != nil {
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

// 查询绑定关系
func BindingEnquiryHandle(data []byte, merId string) (ret *model.BindingReturn) {
	be := new(model.BindingEnquiry)
	err := json.Unmarshal(data, be)
	if err != nil {
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

// BindingPaymentHandle 绑定支付关系
func BindingPaymentHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingPayment)
	err := json.Unmarshal(data, b)
	if err != nil {
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	// 验证请求报文格式
	ret = validateBindingPayment(b)
	if ret != nil {
		return ret
	}
	//  todo 业务处理
	ret = core.ProcessBindingPayment(b)

	return ret
}

// 退款处理
func BindingRefundHandle(data []byte, merId string) (ret *model.BindingReturn) {
	b := new(model.BindingRefund)
	err := json.Unmarshal(data, b)
	if err != nil {
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
		log.Errorf("解析报文错误 :%s", err)
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
		log.Errorf("解析报文错误 :%s", err)
		return mongo.RespCodeColl.Get("200020")
	}
	b.MerId = merId

	ret = validateNoTrackPayment(b)
	if ret != nil {
		return ret
	}

	//  todo 无卡支付暂不开放；业务处理
	return model.NewBindingReturn("000000", "unimplement，暂不支持此类业务")
}
