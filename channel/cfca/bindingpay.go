package cfca

import "quickpay/model"

const (
	correctCode = "2000"
)

const (
	version              = "2.0"
	BindingCreateTxCode  = "2501"
	BindingEnquiryTxCode = "2502"
	BindingRemoveTxCode  = "2503"
	BindingPaymentTxCode = "2511"
)

// ProcessBindingEnquiry 查询绑定关系
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 将参数转化为CfcaRequest
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.InstitutionId,
			TxCode:        BindingEnquiryTxCode,
		},
		Body: requestBody{
			TxSNBinding: be.BindingId,
		},
	}

	// 向中金发起请求
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
			InstitutionID: be.InstitutionId,
			TxCode:        BindingCreateTxCode,
		},
		Body: requestBody{
			TxSNBinding:          be.BindingId,
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
			InstitutionID: be.InstitutionId,
			TxCode:        BindingRemoveTxCode,
		},
		Body: requestBody{
			TxSNUnBinding: be.TxSNUnBinding,
			TxSNBinding:   be.BindingId,
		},
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
			InstitutionID: be.InstitutionId,
			TxCode:        BindingPaymentTxCode,
		},
		Body: requestBody{
			PaymentNo:      be.MerOrderNum,
			Amount:         be.TransAmt,
			TxSNBinding:    be.BindingId,
			SettlementFlag: be.SettlementFlag,
			Remark:         be.Remark,
		},
	}
	//请求
	resp := sendRequest(req)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return

}
