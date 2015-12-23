package scanpay

import (
	"strconv"

	"github.com/CardInfoLink/quickpay/model"
)

const (
	veriTime     = "veriTime"
	origOrderNum = "origOrderNum"
	payType      = "payType"
	cardBin      = "cardBin"
)

// validatePurchaseCoupons 验证卡券核销的参数
func validatePurchaseCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)

	}
	if matched, err := validateTxndir(req.Txndir); !matched {
		return err
	}
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	// 验证格式
	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateVeriTime(req); !matched {
		return err
	}
	if req.Txamt != "" {
		if matched, err := validateTxamt(req); !matched {
			return err
		}
	}
	if matched, err := validateTerminalId(req); !matched {
		return err
	}
	if matched, err := validateScanCodeId(req.ScanCodeId); !matched {
		return err
	}
	if matched, err := validateSign(req.Sign); !matched {
		return err
	}
	return
}

// validatePurchaseActCoupons 验证‘刷卡活动券验证’的参数
func validatePurchaseActCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.PayType == "":
		return fieldEmptyError(payType)

	}

	if matched, err := validateTxndir(req.Txndir); !matched {
		return err
	}
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}

	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateVeriTime(req); !matched {
		return err
	}
	if matched, err := validateTxamt(req); !matched {
		return err
	}
	if matched, err := validatePayType(req); !matched {
		return err
	}
	if matched, err := validateTerminalId(req); !matched {
		return err
	}
	if matched, err := validateScanCodeId(req.ScanCodeId); !matched {
		return err
	}
	if matched, err := validateSign(req.Sign); !matched {
		return err
	}
	if matched, err := validateOrigOrderNum(req.OrigOrderNum); !matched {
		return err
	}
	if matched, err := validateCardBin(req.Cardbin); !matched {
		return err
	}

	return
}

// validateQueryPurchaseCoupons 验证‘电子券验证结果查询’的参数
func validateQueryPurchaseCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}
	if matched, err := validateTxndir(req.Txndir); !matched {
		return err
	}
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}

	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateVeriTime(req); !matched {
		return err
	}
	// if req.Txamt != "" {
	// 	if matched, err := validateTxamt(req); !matched {
	// 		return err
	// 	}
	// }
	// if req.PayType != "" {
	// 	if matched, err := validatePayType(req); !matched {
	// 		return err
	// 	}
	// }
	if matched, err := validateTerminalId(req); !matched {
		return err
	}
	if matched, err := validateScanCodeId(req.ScanCodeId); !matched {
		return err
	}
	if matched, err := validateSign(req.Sign); !matched {
		return err
	}
	if matched, err := validateOrigOrderNum(req.OrigOrderNum); !matched {
		return err
	}
	return
}

// validateUndoPurchaseActCoupons 验证‘卡券撤销’的参数
func validateUndoPurchaseActCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}
	if matched, err := validateTxndir(req.Txndir); !matched {
		return err
	}
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}

	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateTerminalId(req); !matched {
		return err
	}
	if matched, err := validateScanCodeId(req.ScanCodeId); !matched {
		return err
	}
	if matched, err := validateSign(req.Sign); !matched {
		return err
	}
	if matched, err := validateOrigOrderNum(req.OrigOrderNum); !matched {
		return err
	}
	return
}

// validatePayType 验证支付方式
func validatePayType(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {
	if len(req.PayType) > 2 {
		return false, fieldFormatError(payType)
	}
	intPayType, err := strconv.Atoi(req.PayType)
	if err != nil {
		return false, fieldFormatError(payType)
	}
	req.IntPayType = intPayType
	return true, nil
}

// validateVeriTime 验证核销次数
func validateVeriTime(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {
	if req.VeriTime != "" {
		if len(req.VeriTime) > 10 {
			return false, fieldFormatError(veriTime)
		}

		intVeriTime, err := strconv.Atoi(req.VeriTime)
		if err != nil {
			return false, fieldFormatError(veriTime)
		}
		if intVeriTime <= 0 {
			return false, fieldFormatError(veriTime)
		}
		req.IntVeriTime = intVeriTime
	} else {
		req.IntVeriTime = 1
	}
	return true, nil
}

// validateChcd 验证渠道代码
func validateChcd(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {

	if req.Chcd != "" {
		if len(req.Chcd) > 5 {
			return false, fieldFormatError(chcd)
		}
		if req.Chcd != "ULIVE" {
			return false, fieldContentError(chcd)
		}
	}
	return true, nil
}

// validateTerminalId 验证终端
func validateTerminalId(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {
	if len(req.Terminalid) > 8 {
		return false, fieldFormatError(terminalid)
	}
	return true, nil
}

// validateScanCodeId 验证扫码号
func validateScanCodeId(scanCodeIdValue string) (bool, *model.ScanPayResponse) {

	if len(scanCodeIdValue) > 32 {
		return false, fieldFormatError(scanCodeId)
	}
	return true, nil
}

func validateSign(signValue string) (bool, *model.ScanPayResponse) {

	if len(signValue) > 128 {
		return false, fieldFormatError(sign)
	}
	return true, nil
}

// validateOrigOrderNum 验证原订单号
func validateOrigOrderNum(no string) (bool, *model.ScanPayResponse) {

	if len(no) > 64 {
		return false, fieldFormatError(origOrderNum)
	}
	// 是否包含中文或其他非法字符
	if len([]rune(no)) != len(no) {
		return false, fieldFormatError(origOrderNum)
	}
	return true, nil
}

func validateCardBin(cardbin string) (bool, *model.ScanPayResponse) {

	if len(cardbin) > 30 {
		return false, fieldFormatError(cardBin)
	}
	return true, nil
}

func validateTxndir(txndirValue string) (bool, *model.ScanPayResponse) {

	if len(txndirValue) > 1 {
		return false, fieldFormatError(txndir)
	}
	return true, nil
}
