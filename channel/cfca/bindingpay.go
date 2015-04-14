package cfca

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

var Obj cfca

type cfca struct{}

// 中金交易类型
const (
	version              = "2.0"
	correctCode          = "2000"
	BindingCreateTxCode  = "2501"
	BindingEnquiryTxCode = "2502"
	BindingRemoveTxCode  = "2503"
	BindingPaymentTxCode = "2511"
	BindingRefundTxCode  = "2521"
	PaymentEnquiryTxCode = "2512"
	RefundEnquiryTxCode  = "2522"
	TransCheckingTxCode  = "1810"
)

// ProcessBindingEnquiry 查询绑定关系
func (c *cfca) ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingEnquiryTxCode,
		},
		Body: requestBody{
			TxSNBinding: be.ChanBindingId,
		},
		SignCert: be.SignCert,
	}

	// 向中金发起请求
	log.Debugf("request for cfca param (%+v)", req)
	resp := sendRequest(req)

	// 应答码转换。。。
	ret = transformResp(resp, req.Head.TxCode)

	return
}

// ProcessBindingCreate 建立绑定关系
func (c *cfca) ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn) {
	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingCreateTxCode,
		},
		Body: requestBody{
			TxSNBinding:          be.ChanBindingId,
			BankID:               be.BankId,
			AccountName:          be.AcctName,
			AccountNumber:        be.AcctNum,
			IdentificationType:   be.IdentType,
			IdentificationNumber: be.IdentNum,
			PhoneNumber:          be.PhoneNum,
			CardType:             be.AcctType,
			ValidDate:            be.ValidDate,
			CVN2:                 be.Cvv2,
		},
		SignCert: be.SignCert,
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return
}

// ProcessBindingRemove 解除绑定关系
func (c *cfca) ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingRemoveTxCode,
		},
		Body: requestBody{
			TxSNUnBinding: be.TxSNUnBinding,
			TxSNBinding:   be.ChanBindingId,
		},
		SignCert: be.SignCert,
	}

	// 向中金发起请求
	resp := sendRequest(req)

	// 应答码转换。。。
	ret = transformResp(resp, req.Head.TxCode)

	return
}

// ProcessBindingPayment 快捷支付
func (c *cfca) ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingPaymentTxCode,
		},
		Body: requestBody{
			PaymentNo:      be.ChanOrderNum,
			Amount:         be.TransAmt,
			TxSNBinding:    be.ChanBindingId,
			SettlementFlag: be.SettFlag,
			Remark:         be.Remark,
		},
		SignCert: be.SignCert,
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return

}

// ProcessPaymentEnquiry 快捷支付查询
func (c *cfca) ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        PaymentEnquiryTxCode,
		},
		Body: requestBody{
			PaymentNo: be.ChanOrderNum,
		},
		SignCert: be.SignCert,
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return

}

// ProcessBindingRefund 快捷支付退款
func (c *cfca) ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingRefundTxCode,
		},
		Body: requestBody{
			TxSN:      be.ChanOrderNum,     //退款交易流水号
			PaymentNo: be.ChanOrigOrderNum, //原交易流水号
			Amount:    be.TransAmt,
			Remark:    be.Remark,
		},
		SignCert: be.SignCert,
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return
}

// ProcessRefundEnquiry 快捷支付退款查询
func (c *cfca) ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        RefundEnquiryTxCode,
		},
		Body: requestBody{
			TxSN: be.ChanOrderNum, //退款交易流水号
		},
		SignCert: be.SignCert,
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return
}

// ProcessTransChecking 交易对账，清算
func (c *cfca) ProcessTransChecking(chanMerId, settDate, signCert string) (resp *BindingResponse) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			TxCode: TransCheckingTxCode,
		},
		Body: requestBody{
			InstitutionID: chanMerId,
			Date:          settDate,
		},
		SignCert: signCert,
	}
	//请求
	resp = sendRequest(req)
	return
}
