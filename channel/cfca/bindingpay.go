package cfca

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/g"
)

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
)

// ProcessBindingEnquiry 查询绑定关系
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
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
	g.Debug("request for cfca param (%+v)", req)
	resp := sendRequest(req)

	// 应答码转换。。。
	ret = transformResp(resp, req.Head.TxCode)

	return
}

// ProcessBindingCreate 建立绑定关系
func ProcessBindingCreate(be *model.BindingCreate) (ret *model.BindingReturn) {
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
func ProcessBindingRemove(be *model.BindingRemove) (ret *model.BindingReturn) {
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
func ProcessBindingPayment(be *model.BindingPayment) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        BindingPaymentTxCode,
		},
		Body: requestBody{
			PaymentNo:      be.MerOrderNum,
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
func ProcessPaymentEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        PaymentEnquiryTxCode,
		},
		Body: requestBody{
			PaymentNo: be.OrigOrderNum,
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
func ProcessBindingRefund(be *model.BindingRefund) (ret *model.BindingReturn) {

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
func ProcessRefundEnquiry(be *model.OrderEnquiry) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        RefundEnquiryTxCode,
		},
		Body: requestBody{
			TxSN: be.OrigOrderNum, //退款交易流水号
		},
		SignCert: be.SignCert,
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return
}
