package cfca

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

// DefaultClient 默认 CFCA 绑定支付客户端
var DefaultClient CFCABindingPay

// CFCABindingPay CFCA 绑定支付
type CFCABindingPay struct{}

// 中金交易类型
const (
	version     = "2.0"
	correctCode = "2000"

	BindingCreate       = "2501"
	BindingEnquiry      = "2502"
	BindingRemove       = "2503"
	TransChecking       = "1810"
	MarketPaySettlement = "1341"

	MerModePay    = "2511" // 商户模式快捷支付
	MarketModePay = "1371" // 市场模式快捷支付

	MerModeRefund    = "2521" // 商户模式快捷支付退款
	MarketModeRefund = "1373" // 市场模式快捷支付退款

	MerModePayEnquiry    = "2512" // 商户模式快捷支付查询
	MarketModePayEnquiry = "1372" // 市场模式快捷支付查询

	MerModeRefundEnquiry    = "2522" // 商户模式退款查询
	MarketModeRefundEnquiry = "1374" // 市场模式退款查询

	MerModePayWithSMS    = "2542" // 商户模式短信支付
	MarketModePayWithSMS = "1376" // 市场模式短信支付

	MerModeSendSMS    = "2541" // 商户模式发送验证短信
	MarketModeSendSMS = "1375" // 市场模式发送短信验证
)

// ProcessPaySettlement 绑定支付结算
func (c *CFCABindingPay) ProcessPaySettlement(be *model.PaySettlement) (ret *model.BindingReturn) {
	// 组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			TxCode: MarketPaySettlement,
		},
		Body: requestBody{
			InstitutionID: be.ChanMerId,
			SerialNumber:  be.SysOrderNum,
			OrderNo:       be.SettOrderNum,
			Amount:        be.SettAmt,
			AccountType:   be.SettAccountType,
		},
		PrivateKey: be.PrivateKey,
	}

	if be.SettAccountType != "20" {
		bankInfo := &bankAccount{
			BankID:        be.BankCode,
			AccountName:   be.AcctNameDecrypt,
			AccountNumber: be.AcctNumDecrypt,
			BranchName:    be.SettBranchName,
			Province:      be.Province,
			City:          be.City,
		}
		req.Body.BankAccount = bankInfo
	} else {
		req.Body.PaymentAccountName = be.AcctNameDecrypt
		req.Body.PaymentAccountNumber = be.AcctNumDecrypt
	}

	return transformResp(sendRequest(req), req.Head.TxCode)
}

// ProcessPaymentWithSMS 快捷支付(验证并绑定)
func (c *CFCABindingPay) ProcessPaymentWithSMS(be *model.BindingPayment) (ret *model.BindingReturn) {
	// 组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
		},
		Body: requestBody{
			PaymentNo:         be.SysOrderNum,
			SMSValidationCode: be.SmsCode,
		},
		PrivateKey: be.PrivateKey,
	}

	// 判断交易模式
	switch be.Mode {
	case model.MarketMode:
		req.Head.TxCode = MarketModePayWithSMS
		req.Body.OrderNo = be.SettOrderNum
	case model.MerMode:
		req.Head.TxCode = MerModePayWithSMS
	default:
		log.Errorf("unsupport mode %s", be.Mode)
		return mongo.RespCodeColl.Get("000001")
	}

	return transformResp(sendRequest(req), req.Head.TxCode)
}

// ProcessSendBindingPaySMS 快捷支付(发送验证短信)
func (c *CFCABindingPay) ProcessSendBindingPaySMS(be *model.BindingPayment) (ret *model.BindingReturn) {

	// 判断交易模式
	switch be.Mode {
	case model.MarketMode:
		be.SettFlag = ""
		return quickPayment(be, MarketModeSendSMS)
	case model.MerMode:
		be.SettOrderNum = ""
		return quickPayment(be, MerModeSendSMS)
	default:
		log.Errorf("unsupport mode %s", be.Mode)
		return mongo.RespCodeColl.Get("000001")
	}
}

// ProcessBindingCreate 建立绑定关系
func (c *CFCABindingPay) ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn) {
	// 组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingCreate,
		},
		Body: requestBody{
			TxSNBinding:          be.ChanBindingId,
			BankID:               be.BankCode,
			AccountName:          be.AcctNameDecrypt,
			AccountNumber:        be.AcctNumDecrypt,
			IdentificationType:   be.IdentType,
			IdentificationNumber: be.IdentNumDecrypt,
			PhoneNumber:          be.PhoneNumDecrypt,
			CardType:             be.AcctType,
			ValidDate:            be.ValidDateDecrypt,
			CVN2:                 be.Cvv2Decrypt,
		},
		PrivateKey: be.PrivateKey,
	}

	// 请求
	resp := sendRequest(req)
	// 应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return ret
}

// ProcessBindingEnquiry 查询绑定关系
func (c *CFCABindingPay) ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingEnquiry,
		},
		Body: requestBody{
			TxSNBinding: be.ChanBindingId,
		},
		PrivateKey: be.PrivateKey,
	}

	// 向中金发起请求
	resp := sendRequest(req)

	// 应答码转换
	ret = transformResp(resp, req.Head.TxCode)

	return ret
}

// ProcessBindingRemove 解除绑定关系
func (c *CFCABindingPay) ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingRemove,
		},
		Body: requestBody{
			TxSNUnBinding: be.TxSNUnBinding,
			TxSNBinding:   be.ChanBindingId,
		},
		PrivateKey: be.PrivateKey,
	}

	// 向中金发起请求
	resp := sendRequest(req)

	// 应答码转换
	ret = transformResp(resp, req.Head.TxCode)

	return ret
}

// ProcessBindingPayment 快捷支付
func (c *CFCABindingPay) ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {

	// 判断交易模式
	switch be.Mode {
	case model.MarketMode:
		be.SettFlag = ""
		return quickPayment(be, MarketModePay)
	case model.MerMode:
		be.SettOrderNum = ""
		return quickPayment(be, MerModePay)
	default:
		log.Errorf("unsupport mode %s", be.Mode)
		return mongo.RespCodeColl.Get("000001")
	}
}

// quickPayment 快捷支付，根据业务代码的不同，走不同接口
// 2511、2541接口。
func quickPayment(be *model.BindingPayment, txCode string) (ret *model.BindingReturn) {
	// 组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        txCode,
		},
		Body: requestBody{
			OrderNo:        be.SettOrderNum,
			PaymentNo:      be.SysOrderNum,
			Amount:         be.TransAmt,
			TxSNBinding:    be.ChanBindingId,
			SettlementFlag: be.SettFlag,
			Remark:         be.Remark,
		},
		PrivateKey: be.PrivateKey,
	}

	return transformResp(sendRequest(req), req.Head.TxCode)
}

// ProcessPaymentEnquiry 快捷支付查询
func (c *CFCABindingPay) ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	// 组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
		},
		Body: requestBody{
			PaymentNo: be.SysOrderNum,
		},
		PrivateKey: be.PrivateKey,
	}

	// 判断交易模式
	switch be.Mode {
	case model.MarketMode:
		req.Head.TxCode = MarketModePayEnquiry
	case model.MerMode:
		req.Head.TxCode = MerModePayEnquiry
	default:
		log.Errorf("unsupport mode %d", be.Mode)
		return mongo.RespCodeColl.Get("000001")
	}

	return transformResp(sendRequest(req), req.Head.TxCode)
}

// ProcessBindingRefund 快捷支付退款
func (c *CFCABindingPay) ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {
	// 组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
		},
		Body: requestBody{
			TxSN:      be.SysOrderNum,     //退款交易流水号
			PaymentNo: be.SysOrigOrderNum, //原交易流水号
			Amount:    be.TransAmt,
			Remark:    be.Remark,
		},
		PrivateKey: be.PrivateKey,
	}

	// 判断交易模式
	switch be.Mode {
	case model.MarketMode:
		req.Head.TxCode = MarketModeRefund
		req.Body.OrderNo = be.SettOrderNum
	case model.MerMode:
		req.Head.TxCode = MerModeRefund
	default:
		log.Errorf("unsupport mode %s", be.Mode)
		return mongo.RespCodeColl.Get("000001")
	}

	// 应答码转换
	return transformResp(sendRequest(req), req.Head.TxCode)
}

// ProcessRefundEnquiry 快捷支付退款查询
func (c *CFCABindingPay) ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {
	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
		},
		Body: requestBody{
			TxSN: be.SysOrderNum, //退款交易流水号
		},
		PrivateKey: be.PrivateKey,
	}

	// 判断交易模式
	switch be.Mode {
	case model.MarketMode:
		req.Head.TxCode = MarketModeRefund
	case model.MerMode:
		req.Head.TxCode = MerModeRefund
	default:
		log.Errorf("unsupport mode %s", be.Mode)
		return mongo.RespCodeColl.Get("000001")
	}

	return transformResp(sendRequest(req), req.Head.TxCode)
}

// ProcessTransChecking 交易对账，清算
func (c *CFCABindingPay) ProcessTransChecking(chanMerId, settDate, PrivateKey string) (resp *BindingResponse) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			TxCode: TransChecking,
		},
		Body: requestBody{
			InstitutionID: chanMerId,
			Date:          settDate,
		},
		PrivateKey: PrivateKey,
	}

	// 请求
	resp = sendRequest(req)

	return resp
}
