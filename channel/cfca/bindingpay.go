package cfca

import (
	"github.com/omigo/g"
	"quickpay/model"
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
	}

	// 向中金发起请求
	resp := sendRequest(req, be.SignCert)

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
	}
	//请求
	resp := sendRequest(req, be.SignCert)
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
	}

	// 向中金发起请求
	resp := sendRequest(req, be.SignCert)

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
	}
	//请求
	g.Debug("request for cfca param (%s)", req)
	resp := sendRequest(req, be.SignCert)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return

}

// ProcessPaymentEnquiry 快捷支付查询
func ProcessPaymentEnquiry(be *model.BindingPayment) (ret *model.BindingReturn) {

	//组装参数
	req := &BindingRequest{
		Version: version,
		Head: requestHead{
			InstitutionID: be.ChanMerId,
			TxCode:        PaymentEnquiryTxCode,
		},
		Body: requestBody{
			PaymentNo: be.MerOrderNum,
		},
	}
	//请求
	g.Debug("request for cfca param (%s)", req)
	resp := sendRequest(req, be.SignCert)
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
			TxSN:      be.MerOrderNum,  //退款交易流水号
			PaymentNo: be.OrigOrderNum, //原交易流水号
			Amount:    be.TransAmt,
			Remark:    be.Remark,
		},
	}
	//请求
	resp := sendRequest(req, be.SignCert)
	//应答码转换
	ret = transformResp(resp, req.Head.TxCode)
	return
}
